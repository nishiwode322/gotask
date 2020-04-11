package main

import (
	"database/sql"
	"fmt"

	"./util"
	"github.com/Shopify/sarama"
)

func task1() {
	msg := util.FillMessageStruct("ProvinceCity", "江苏", util.CityList{"南京", "苏州"})

	config := sarama.NewConfig()
	// WaitForAll waits for all in-sync replicas to commit before responding
	config.Producer.RequiredAcks = sarama.WaitForAll
	// NewRandomPartitioner returns a Partitioner which chooses a random partition each time
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	ko := util.KafkaOps{BrokerAddress: []string{"localhost:9092"}}
	ko.SyncSendMessage(config, []*sarama.ProducerMessage{msg})

	fmt.Println("consumer start")
	bufferSize := 100
	dataChannel := make(chan *sarama.ConsumerMessage, bufferSize)
	go ko.ReceiveMessage("ProvinceCity", nil, dataChannel)

	for item := range dataChannel {
		fmt.Printf("key: %s, text: %s, offset: %d\n", string(item.Key), string(item.Value), item.Offset)

		value := util.CityList{}
		value.DecodeString(string(item.Value))

		/* sync mysql */
		db, _ := sql.Open("mysql", "root:123456@(127.0.0.1:3306)/provincecity")
		defer db.Close()
		err := db.Ping()
		if err != nil {
			fmt.Println("database connect error! ", err)
			return
		}
		sql := fmt.Sprintf("insert into province values (0,'%s')", string(item.Key))
		_, err = db.Exec(sql)
		if err != nil {
			fmt.Println("insert sql err! ", err, sql)
		}
		/*
				batch inserts

			keyValueFormatStrings := []string{}
			for _, x := range value {
				tmp := fmt.Sprintf("(0,'%s')", x)
				keyValueFormatStrings = append(keyValueFormatStrings, tmp)
			}
			valueSql := "insert into province values " + strings.Join(keyValueFormatStrings, ",")
			result, err := db.Exec(valueSql)
			n, _ := result.RowsAffected()
			fmt.Println("affected row number is:", n)
		*/

		/*sync redis*/

	}
}
