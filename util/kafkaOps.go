package util

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

func producer() {

	config := sarama.NewConfig()
	// WaitForAll waits for all in-sync replicas to commit before responding
	config.Producer.RequiredAcks = sarama.WaitForAll
	// NewRandomPartitioner returns a Partitioner which chooses a random partition each time
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		fmt.Println("producer close err! ", err)
		return
	}
	defer client.Close()

	// set message content
	msg := &sarama.ProducerMessage{}
	msg.Topic = "ProvinceCity"
	msg.Key = sarama.StringEncoder("江苏")
	var valueString = strings.Join([]string{"南京", "苏州"}, ",")
	msg.Value = sarama.StringEncoder(valueString)

	// send message
	pid, offset, err := client.SendMessage(msg)

	if err != nil {
		fmt.Println("send message failed! ", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}

func consumer() {
	var wg sync.WaitGroup
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Println("close consumer err! ", err)
		}
	}()

	fmt.Println("consumer connnect kafka success...")
	partitions, err := consumer.Partitions("revolution")
	if err != nil {
		fmt.Println("get partitions failed ", err)
		return
	}

	for _, p := range partitions {
		partitionConsumer, err := consumer.ConsumePartition("ProvinceCity", p, sarama.OffsetOldest)
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
		go func() {
			//TODO:
			for m := range partitionConsumer.Messages() {
				fmt.Printf("key: %s, text: %s, offset: %d\n", string(m.Key), string(m.Value), m.Offset)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func RunKafkaSample() {
	producer()
	fmt.Println("-----------+++----------------------")
	consumer()
}
