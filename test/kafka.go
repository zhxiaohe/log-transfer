package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

var (
	brokers  = ""
	version  = "2.2.1"
	group    = ""
	topics   = "dev_access_test"
	assignor = ""
	oldest   = true
	verbose  = false
)

type outputComponent struct {
	brokers  []string
	topic    string
	producer sarama.AsyncProducer
}

func KAFInit() *outputComponent {
	fmt.Println(brokers)
	/*构造配置*/

	version, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	// 初始化配置
	config := sarama.NewConfig()
	config.Version = version
	config.Producer.Flush.Frequency = 200 * time.Millisecond // 200毫秒刷新
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForLocal      // 只等待leader ack , sarama.WaitForAll等待所有同步副本ack消息
	config.Producer.Partitioner = sarama.NewHashPartitioner //通过msg中的key生成hash选择分区, 默认随机分配

	// sarama.NewManualPartitioner 返回一个手动选择分区的分割器,也就是获取msg中指定的`partition`
	// sarama.NewRandomPartitioner 通过随机函数随机获取一个分区号
	// sarama.NewRoundRobinPartitioner 环形选择,也就是在所有分区中循环选择一个
	// sarama.NewHashPartitioner 通过msg中的key生成hash值,选择分区

	// config.Producer.RequiredAcks = sarama.WaitForAll // 同步

	//producer.Input() <- &sarama.ProducerMessage
	// 创建上下文
	// ctx, cancel := context.WithCancel(context.Background())

	// 初始化客户端
	brokerList := strings.Split(brokers, ",")
	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatal("NewSyncProducer err:", err)
	}

	//关闭
	// defer producer.Close()
	go func() {
		// for err := range producer.Errors() {
		// 	log.Println("Failed to write access log entry:", err)
		// }
		select {
		case err := <-producer.Errors():
			log.Println("Failed to write access log entry:", err)
		case msg := <-producer.Successes():
			b, _ := msg.Value.Encode()
			log.Println("Success:", b)
		}
	}()

	return &outputComponent{producer: producer}
}

func (o *outputComponent) Start(event chan map[string]interface{}) {
	producer := o.producer

	for m := range event {
		key := fmt.Sprintf("%s", m["key"])
		value := fmt.Sprintf("%s", m["value"])
		msg := &sarama.ProducerMessage{Topic: topics, Key: sarama.StringEncoder(key), Value: sarama.StringEncoder(value)}
		producer.Input() <- msg
	}

}

func (o *outputComponent) write() {
	producer := o.producer
	limit := 10
	for i := 0; i < limit; i++ {
		fmt.Println("i:", i)
		str := strconv.Itoa(int(time.Now().UnixNano()))
		msg := &sarama.ProducerMessage{Topic: topics, Key: nil, Value: sarama.StringEncoder(str)}
		producer.Input() <- msg
	}
}

func (o *outputComponent) close() {
	if err := o.producer.Close(); err != nil {
		log.Panicln("kafka close err: ", err)
	}
}
