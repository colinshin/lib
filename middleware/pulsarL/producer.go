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
	"runtime"
	"strconv"
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
	IsSync     bool              `json:"is_sync"`
	Delay      int               `json:"delay"`
}

var producerQue = cmap.New[pulsar.Producer]() //cmap.ConcurrentMap[int, pulsar.Producer]

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
	start := time.Now()
	producer, e := pClient.CurrPulsar.CreateProducer(pulsar.ProducerOptions{
		Topic:                   codeStr,
		ProducerAccessMode:      pulsar.ProducerAccessModeShared,
		DisableBlockIfQueueFull: true,
		BatchingMaxSize:         1048576, //1M
		SendTimeout:             time.Second * 5,
	})
	fmt.Println("发送耗时=========================", time.Since(start).Milliseconds())
	if e != nil {
		logger.AddError(zap.Error(e))
		panic(e)
	}

	if o.IsSync {
		_, errSend := producer.Send(ctx, pMessage)
		if errSend != nil {
			logger.AddError(zap.Error(errSend))
			return errSend
		}
	} else {

		producer.SendAsync(ctx, pMessage, func(id pulsar.MessageID, message *pulsar.ProducerMessage, err error) {
			if err != nil {
				logger.AddError(zap.String("MessageId", id.String()), zap.String("Message:", string(message.Payload)))
			} else {
				logger.AddNotice(zap.String("MessageId", id.String()))
			}
		})

	}
	count := producerQue.Count()

	producerQue.Set(strconv.Itoa(count), producer)
	if count > 90 {
		Flush()
	} else if count > 50 {
		go Flush()
	}

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
	track := ""
	for i := 1; i <= 5; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if ok {
			track += fmt.Sprintf("fun%s:%s(%d)", runtime.FuncForPC(pc).Name(), file, line)
		} else {
			break
		}
	}
	codeStr := o.TopicStr
	if o.Topic > 0 {
		codeStr = strconv.Itoa(objTopic.Code)
	}
	payload := PulsarMessage{
		ProducerTime: time.Now().UnixMilli(),
		From:         track,
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
	}
}
func Flush() {
	for i, v := range producerQue.Items() {
		v.Close()
		producerQue.Remove(i)
	}
}
