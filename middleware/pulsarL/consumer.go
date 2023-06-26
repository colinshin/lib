package pulsarL

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/flyerxp/lib/app"
	"github.com/flyerxp/lib/logger"
	json2 "github.com/flyerxp/lib/utils/json"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type Consumer struct {
	Topics            map[string][]string
	Name              string
	Options           *Options
	ConsumerContainer map[string]pulsar.Consumer
	ConnContainer     map[string]*PulsarClient
	IsStop            bool
}
type Options struct {
	Name int
	Dlq  *pulsar.DLQPolicy
}
type Option func(opts *Options)

func NewConsumer(s []string, subName string, f ...Option) *Consumer {
	c := new(Consumer)
	c.Topics = map[string][]string{}
	for _, v := range s {
		t, _ := getTopic(v)
		if _, ok := c.Topics[t.Cluster]; ok {
			c.Topics[t.Cluster] = append(c.Topics[t.Cluster], t.CodeStr)
		} else {
			c.Topics[t.Cluster] = []string{t.CodeStr}
		}
	}
	c.Name = subName
	c.Options = loadOptions(f...)
	return c
}

func GetStringTopics(i []int) []string {
	var s = make([]string, len(i), len(i))
	for ii, v := range i {
		s[ii] = strconv.Itoa(v)
	}
	return s
}
func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

// 死信和失败保存的topic
func WithDlq(policy pulsar.DLQPolicy) Option {
	return func(opts *Options) {
		var p = policy
		opts.Dlq = &p
	}
}
func (c *Consumer) Consumer(F func(*pulsar.ConsumerMessage, *PulsarMessage) bool) {
	var MessageChannel = make(chan pulsar.ConsumerMessage)
	c.ConnContainer = make(map[string]*PulsarClient)
	c.ConsumerContainer = make(map[string]pulsar.Consumer)
	var err error
	ackGroup := pulsar.AckGroupingOptions{MaxSize: 1000, MaxTime: 100 * time.Millisecond}
	dlq := new(pulsar.DLQPolicy)
	if c.Options.Dlq == nil {
		dlq.MaxDeliveries = 20
	} else {
		dlq = c.Options.Dlq
	}
	if dlq.DeadLetterTopic == "" {
		dlq.DeadLetterTopic = "dead_letter_topic"
	}
	if dlq.RetryLetterTopic == "" {
		dlq.RetryLetterTopic = "retry_letter_topic"
	}
	for cluster, gTopics := range c.Topics {
		c.ConnContainer[cluster], err = GetEngine(cluster, context.Background())
		if err != nil {
			logger.AddError(zap.Error(err))
			logger.WriteErr()
		}
		c.ConsumerContainer[cluster], err = c.ConnContainer[cluster].CurrPulsar.Subscribe(pulsar.ConsumerOptions{
			Topics:                 gTopics,
			SubscriptionName:       c.Name,
			Name:                   c.Name,
			Type:                   pulsar.Shared,
			MessageChannel:         MessageChannel,
			AutoAckIncompleteChunk: true,
			AckGroupingOptions:     &ackGroup,
			DLQ:                    dlq,
			RetryEnable:            true,
		})
		if err != nil {
			logger.AddError(zap.Error(err))
			logger.WriteErr()
		}
	}
	defer func() {
		for i := range c.ConsumerContainer {
			c.ConsumerContainer[i].Close()
		}
		app.Shutdown(context.Background())
	}()
	for i := range c.ConsumerContainer {
		c.ConsumerContainer[i].Subscription()
	}
	for !c.IsStop {
		select {
		case cm, ok := <-MessageChannel:
			if ok {
				gProductChan := new(PulsarMessage)
				err = json2.Decode(cm.Payload(), gProductChan)
				if err != nil {
					logger.AddError(zap.Error(err))
					logger.WriteErr()
				}
				if prop, okRetry := cm.Message.Properties()["RECONSUMETIMES"]; okRetry {
					gProductChan.Properties["RECONSUMETIMES"] = prop
				}
				_ = F(&cm, gProductChan)
				err = cm.Consumer.Ack(cm.Message)
				if err != nil {
					logger.AddError(zap.Error(err))
				}
				logger.WriteLine()
			} else {
				logger.AddError(zap.String("consumer", " chan no ok"))
				logger.WriteErr()
			}
		case <-time.After(time.Second * 30):
			fmt.Println("30秒没有消息")
		}
	}
}
func RetryAfter(message *pulsar.ConsumerMessage, t time.Duration, m map[string]string) {
	message.Consumer.ReconsumeLaterWithCustomProperties(message, m, time.Second*t)
}
func (c *Consumer) Stop() {
	c.IsStop = true
}
