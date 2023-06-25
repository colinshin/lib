package pulsarL

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/flyerxp/lib/logger"
	json2 "github.com/flyerxp/lib/utils/json"
	cmap "github.com/orcaman/concurrent-map/v2"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type PulsarMessage struct {
	ProducerTime int64             `json:"producer_time"`
	From         string            `json:"from"`
	Topic        string            `json:"topic"`
	RequestId    string            `json:"request_id"`
	Content      json.RawMessage   `json:"content"`
	DelayTime    time.Duration     `json:"delay_time"`
	Properties   map[string]string `json:"properties"`
}

type OutMessage struct {
	Topic      int               `json:"topic"`
	TopicStr   string            `json:"topic_str"`
	Content    any               `json:"content"`
	Properties map[string]string `json:"properties"`
	Delay      int               `json:"delay"`
	Key        string            `json:"key"`
	Track      string            `json:"track"`
}

var producerQue *pulsarProducer

type pulsarProducer struct {
	//Pool     *ants.Pool
	Que      cmap.ConcurrentMap[string, pulsar.Producer]
	isInitEd bool
	Sending  int32
	Wg       sync.WaitGroup
}

func init() {
	initProducer()
}
func initProducer() {
	producerQue = new(pulsarProducer)
	producerQue.isInitEd = true
	producerQue.Que = cmap.New[pulsar.Producer]()
}

func Producer(o *OutMessage, ctx context.Context) error {
	objTopic, ok := getTopic(o.Topic)
	if !ok {
		panic(errors.New(fmt.Sprintf("%d no find message %s", o.Topic, o.Content)))
	}
	pMessage := getPulsarMessage(o, objTopic, ctx)
	pClient, e := GetEngine(objTopic.Cluster, ctx)

	if e != nil {
		logger.AddError(zap.Error(e), zap.Any("message", o))
	}
	codeStr := o.TopicStr
	if o.Topic > 0 {
		codeStr = strconv.Itoa(objTopic.Code)
	}
	var p pulsar.Producer
	p, ok = producerQue.Que.Get(codeStr)
	if !ok {
		//官方此处存在性能问题,协程下直接卡死
		//，目前无解,NewRequestID
		p, e = pClient.CurrPulsar.CreateProducer(pulsar.ProducerOptions{
			Topic:              codeStr,
			ProducerAccessMode: pulsar.ProducerAccessModeShared,
			//DisableBlockIfQueueFull: true,
			BatchingMaxSize:                 1048576, //1M
			SendTimeout:                     time.Second * 5,
			BatchingMaxPublishDelay:         2 * time.Second,
			BatchingMaxMessages:             100,
			PartitionsAutoDiscoveryInterval: time.Second * 86400 * 5,
		})
		producerQue.Que.Set(codeStr, p)
	}
	if e != nil {
		logger.AddError(zap.Error(e))
		panic(e)
	}
	atomic.AddInt32(&producerQue.Sending, 1)
	producerQue.Wg.Add(1)
	p.SendAsync(ctx, pMessage, func(id pulsar.MessageID, message *pulsar.ProducerMessage, err error) {
		if err != nil {
			logger.AddError(zap.Error(err), zap.Any(codeStr, pMessage.Payload))
		}
		atomic.AddInt32(&producerQue.Sending, -1)
		producerQue.Wg.Done()
	})

	return nil
}
func getPulsarMessage(o *OutMessage, objTopic *TopicS, ctx context.Context) *pulsar.ProducerMessage {
	if o.Properties == nil {
		o.Properties = map[string]string{}
	}
	//全局唯一id
	if rId, ok := ctx.Value("GlobalRequestId").(string); ok {
		o.Properties["GlobalRequestId"] = rId
	} else {
		o.Properties["GlobalRequestId"] = ""
	}
	BContent, err := json2.Encode(o.Content)
	if err != nil {
		logger.AddError(zap.Error(err))
		panic(err)
	}

	codeStr := o.TopicStr
	if o.Topic > 0 {
		codeStr = strconv.Itoa(objTopic.Code)
	}
	payload := PulsarMessage{
		ProducerTime: time.Now().UnixMilli(),
		From:         o.Track,
		Topic:        codeStr,
		RequestId:    o.Properties["GlobalRequestId"],
		Content:      BContent,
		DelayTime:    time.Duration(objTopic.Delay),
		Properties:   o.Properties,
	}
	if o.Delay > 0 {
		payload.DelayTime = time.Duration(o.Delay)
	}
	payloadB, errJ := json2.Encode(payload)
	if errJ != nil {
		logger.AddError(zap.Error(errJ), zap.Any("fail payload", payload))
		panic(errors.New("json Fail"))
	}
	return &pulsar.ProducerMessage{
		Payload:      payloadB,
		Properties:   payload.Properties,
		DeliverAfter: time.Second * payload.DelayTime,
		Key:          o.Key,
	}
}
func producerPreInit(t []string) {
	var p pulsar.Producer
	for _, codeStr := range t {
		objTopic, ok := getTopic(codeStr)
		pClient, _ := GetEngine(objTopic.Cluster, context.Background())
		p, ok = producerQue.Que.Get(codeStr)
		if !ok {
			//官方此处存在性能问题,协程下直接卡死
			//，目前无解,NewRequestID
			p, _ = pClient.CurrPulsar.CreateProducer(pulsar.ProducerOptions{
				Topic:              codeStr,
				ProducerAccessMode: pulsar.ProducerAccessModeShared,
				//DisableBlockIfQueueFull: true,
				BatchingMaxSize:                 1048576, //1M
				SendTimeout:                     time.Second * 5,
				BatchingMaxPublishDelay:         2 * time.Second,
				BatchingMaxMessages:             100,
				PartitionsAutoDiscoveryInterval: time.Second * 86400 * 5,
			})
			producerQue.Que.Set(codeStr, p)
		}
	}
}
func Flush() {
	if producerQue != nil {
		producerQue.Wg.Wait()
		for _, v := range producerQue.Que.Items() {
			e := v.Flush()
			if e != nil {
				logger.AddError(zap.Error(e))
			}
		}
	}
}