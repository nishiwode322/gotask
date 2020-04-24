package task

import (
	"fmt"

	"gotask/util"

	"github.com/Shopify/sarama"
	"github.com/garyburd/redigo/redis"
)

func Task1() {
	// parse province and city information
	provinceCityMap, err := util.ParseProvinceAndCity("http://www.hotelaah.com/dijishi.html")
	if err != nil {
		fmt.Println("parse province city err ", err)
		return
	}
	fmt.Println("parse province city success")

	//set kafka producer messages struct
	messageList := []*sarama.ProducerMessage{}
	for province := range provinceCityMap {
		tmp := provinceCityMap[province]
		msg := util.FillMessageStruct("ProvinceCity", province, &tmp)
		messageList = append(messageList, msg)
	}

	//set kafka producer config
	config := sarama.NewConfig()
	//WaitForAll waits for all in-sync replicas to commit before responding
	config.Producer.RequiredAcks = sarama.WaitForAll
	//NewRandomPartitioner returns a Partitioner which chooses a random partition each time
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	//init and send producer messages to kafka
	ko := util.KafkaOps{BrokerAddress: []string{"localhost:9092"}}
	ko.SyncSendMessage(config, messageList)

	//start kafka consumer
	fmt.Println("----------------------------------------------")
	bufferSize := 100
	dataChannel := make(chan *sarama.ConsumerMessage, bufferSize)
	go ko.ReceiveMessage("ProvinceCity", nil, dataChannel)

	//init mysql
	db := util.MysqlOps{UserName: "root", PassWord: "123456", IP: "127.0.0.1", Port: "3306", DataBase: "provincecity"}
	err = db.Open()
	if err != nil {
		fmt.Println("database instance create err!")
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("database connect err! ", err)
		return
	}

	//init redis
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("connect redis err! ", err)
		return
	}
	defer conn.Close()

	// number of records in the province table
	var provinceCount int64
	provinceCount = 0
	//handle messages from kafka
	for item := range dataChannel {

		value := util.CityList{}
		value.DecodeString(string(item.Value))

		//sync mysql
		sql := fmt.Sprintf("insert into province values (0,'%s')", string(item.Key))
		_, err = db.Exec(sql)
		if err != nil {
			fmt.Println("insert province sql err! ", err, sql)
		}

		//number of records add 1
		provinceCount++

		valueList := [][]interface{}{}
		for _, v := range value {
			tmp := []interface{}{0, v, provinceCount}
			valueList = append(valueList, tmp)
		}
		_, err = db.BatchInserts("city", []string{"city_id", "city_name", "province_id"}, valueList)
		if err != nil {
			fmt.Println("insert city sql err! ", err)
		}

		//change to GBK encode
		for i, v := range value {
			t, _ := util.DecodeUTF(v)
			value[i] = string(t)
		}
		//sync redis
		for _, v := range value {
			reply, err := conn.Do("zadd", string(item.Key), 0, v)
			fmt.Printf("reply=%#v,err=%v\n", reply, err)
		}
	}
}
