package pulsarL

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/flyerxp/lib/logger"
	"github.com/flyerxp/lib/middleware/pulsarL"
	"strconv"
	"testing"
	"time"
)

func TestProd(T *testing.T) {
	//return
	time.Sleep(time.Second * 1)
	fmt.Println("开始发10000条消息")
	t := time.Now()
	for i := 0; i < 6; i++ {
		_ = pulsarL.Producer(&pulsarL.OutMessage{
			Topic:      10101001,
			Content:    map[string]string{"a": "b", "10101001": "==============" + strconv.Itoa(i) + "=================="},
			Properties: map[string]string{"prop": "prop"},
			Delay:      0,
		}, context.Background())
		_ = pulsarL.Producer(&pulsarL.OutMessage{
			TopicStr:   "test",
			Content:    map[string]string{"a": "b", "test": "==============test=================="},
			Properties: map[string]string{"prop": "prop"},
			Delay:      0,
		}, context.Background())
	}
	fmt.Println(time.Since(t).Milliseconds(), "总耗时！")
	pulsarL.Flush()
	logger.WriteLine()
}
func TestConsum(T *testing.T) {
	//return
	topics := pulsarL.GetStringTopics([]int{10101001})
	topics = append(topics, "test")
	c := pulsarL.NewConsumer(topics, "testConsume", pulsarL.WithDlq(pulsar.DLQPolicy{
		MaxDeliveries:    5,
		DeadLetterTopic:  "dead_test",
		RetryLetterTopic: "retry_test",
	}))

	count := 0
	c.Consumer(func(message *pulsar.ConsumerMessage, message2 *pulsarL.PulsarMessage) bool {
		//c.Stop()
		fmt.Println(message2.String())
		fmt.Println(message.Properties())
		fmt.Println(message2.Topic)
		count++
		if count == 9 {
			pulsarL.RetryAfter(message, time.Second*10, map[string]string{"aaa": "abdddddddddddddddddddddddddddddddddddddddddcd"})
			return true
		}
		return true
	})
}
