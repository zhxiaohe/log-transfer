package output

import (
	"fmt"
	"log"
	C "log_transfer/config"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

var (
	brokers = "10.139.2.222:9092, 10.139.7.53:9092, 10.139.5.249:9092"
	// brokers  = "10.139.2.222:9092"
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
	config   *sarama.Config
}

func KafkaInit() *outputComponent {
	fmt.Println(brokers)
	/*构造配置*/

	version, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}
	brokers := C.C.Out.Kafka.Brokers
	// 初始化配置
	config := sarama.NewConfig()
	config.Version = version
	config.Producer.Flush.Frequency = 200 * time.Millisecond // 200毫秒刷新
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForLocal // 只等待leader ack , sarama.WaitForAll等待所有同步副本ack消息
	// config.Producer.Partitioner = sarama.NewHashPartitioner //通过msg中的key生成hash选择分区, 默认随机分配
	config.Producer.Partitioner = sarama.NewRandomPartitioner //通过msg中的key生成hash选择分区, 默认随机分配

	brokerList := strings.Split(brokers, ",")
	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatal("NewSyncProducer err:", err)
	}

	return &outputComponent{config: config, producer: producer}
}

func (o *outputComponent) Start(event chan map[string]interface{}) error {
	// producer := o.producer
	go func() {
		for {
			select {
			case err := <-o.producer.Errors():
				log.Println("Failed to write access log entry:", err)
			case <-o.producer.Successes():
				// b, _ := msg.Value.Encode()
				// log.Println("Success:", b)
			}
		}
	}()

	for m := range event {
		key := fmt.Sprintf("%s", m["key"])
		value := fmt.Sprintf("%s", m["value"])
		msg := &sarama.ProducerMessage{Topic: topics, Key: sarama.StringEncoder(key), Value: sarama.StringEncoder(value)}
		// msg := &sarama.ProducerMessage{Topic: topics, Key: nil, Value: sarama.StringEncoder(value)}
		o.producer.Input() <- msg
		// fmt.Println(key)
	}
	defer o.producer.AsyncClose()
	return nil

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
	producer := o.producer
	if err := producer.Close(); err != nil {
		log.Panicln("kafka close err: ", err)
	}
}
