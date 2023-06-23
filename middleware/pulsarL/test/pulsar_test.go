package pulsarL

import (
	"context"
	"fmt"
	"github.com/flyerxp/lib/middleware/pulsarL"
	"strconv"
	"testing"
	"time"
)

func TestProd(T *testing.T) {
	time.Sleep(time.Second)
	for i := 0; i <= 10; i++ {
		t := time.Now()
		_ = pulsarL.Producer(&pulsarL.OutMessage{
			Topic:      10101001,
			Content:    map[string]string{"a": "b", "c": "==============" + strconv.Itoa(i) + "=================="},
			Properties: map[string]string{},
			IsSync:     false,
			Delay:      0,
		}, context.Background())
		fmt.Println(time.Since(t).Milliseconds(), "总耗时")
	}
	pulsarL.Flush()
	pulsarL.Reset()
	//logger.WriteLine()
}
func TestConsum(T *testing.T) {
	/*
		r, _ := pulsarL.GetEngine("pubPulsar", context.Background())
		p := r.GetPulsar()
		consumer, err := p.Subscribe(pulsar.ConsumerOptions{
			Topic:            10101001,
			SubscriptionName: "my-sub",
			Type:             pulsar.Shared,
		})
		fmt.Println(err)
		defer consumer.Close()
		for i := 0; i <= 1; i++ {
			msg, err := consumer.Receive(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Received message msgId: %#v -- content: '%s'\n",
				msg.ID(), string(msg.Payload()))
		}*/

}
func TestConsumRead(T *testing.T) {
	/*	r, _ := pulsarL.GetEngine("pubPulsar", context.Background())
		p := r.GetPulsar()
		reader, err := p.CreateReader(pulsar.ReaderOptions{
			Topic:          "test",
			StartMessageID: pulsar.EarliestMessageID(),
		})
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()

		for reader.HasNext() {
			msg, err := reader.Next(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Received message msgId: %#v -- content: '%s'\n",
				msg.ID(), string(msg.Payload()))
		}*/

}
