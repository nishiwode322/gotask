package util

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

type kafkaops struct {
	BrokerAddress []string
}

func (k kafkaops) SyncSendMessage(config *sarama.Config, messages []*sarama.ProducerMessage) {
	client, err := sarama.NewSyncProducer(k.BrokerAddress, config)
	if err != nil {
		fmt.Println("producer close err! ", err)
		return
	}
	defer client.Close()

	// send message
	for index, msg := range messages {
		fmt.Printf("send message %v: %v,%v,%v\n", index, msg.Topic, msg.Key, msg.Value)
		pid, offset, err := client.SendMessage(msg)

		if err != nil {
			fmt.Println("send message failed! ", err)
		} else {
			fmt.Printf("message's pid:%v offset:%v\n", pid, offset)
		}
	}
}

func (k kafkaops) ReceiveMessage(topic string, config *sarama.Config, data chan<- *sarama.ConsumerMessage) {
	var wg sync.WaitGroup
	consumer, err := sarama.NewConsumer(k.BrokerAddress, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Println("close consumer err! ", err)
		}
	}()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println("get partitions failed ", err)
		return
	}

	for _, p := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, p, sarama.OffsetOldest)
		if err != nil {
			fmt.Println("partitionConsumer err! ", err)
			continue
		}
		defer func() {
			if err := partitionConsumer.Close(); err != nil {
				fmt.Println("partitionConsumer close err! ", err)
			}
		}()
		wg.Add(1)
		go func(c chan<- *sarama.ConsumerMessage) {
			// send received message to channel
			for m := range partitionConsumer.Messages() {
				fmt.Printf("key: %s, text: %s, offset: %d\n", string(m.Key), string(m.Value), m.Offset)
				c <- m
			}
			wg.Done()
		}(data)
	}
	defer close(data)
	wg.Wait()
}

func RunKafkaSample() {
	//broker address: "localhost:9092"
	// set message content
	msg := &sarama.ProducerMessage{}
	msg.Topic = "ProvinceCity"
	msg.Key = sarama.StringEncoder("江苏")
	var valueString = strings.Join([]string{"南京", "苏州"}, ",")
	msg.Value = sarama.StringEncoder(valueString)

	config := sarama.NewConfig()
	// WaitForAll waits for all in-sync replicas to commit before responding
	config.Producer.RequiredAcks = sarama.WaitForAll
	// NewRandomPartitioner returns a Partitioner which chooses a random partition each time
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	ko := kafkaops{[]string{"localhost:9092"}}
	ko.SyncSendMessage(config, []*sarama.ProducerMessage{msg})

	fmt.Println("consumer start")
	bufferSize := 100
	dataChannel := make(chan *sarama.ConsumerMessage, bufferSize)
	go ko.ReceiveMessage("ProvinceCity", nil, dataChannel)

	for item := range dataChannel {
		fmt.Printf("key: %s, text: %s, offset: %d\n", string(item.Key), string(item.Value), item.Offset)
	}
}
