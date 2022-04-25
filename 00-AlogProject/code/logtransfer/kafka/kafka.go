package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"logtransfer/es"
)
// 初始化kafka连接的一个client

type KafkaClient struct {
	client sarama.Consumer
	addrs []string
	topic string
}

var (
	kafkaClient *KafkaClient
)

// init初始化client
func Init(addrs []string, topic string) (err error) {
	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	kafkaClient = &KafkaClient{
		client: consumer,
		addrs: addrs,
		topic: topic,
	}
	return
}

func Run() {
	partitionList, err := kafkaClient.client.Partitions(kafkaClient.topic) // topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println("分区: ", partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := kafkaClient.client.ConsumePartition(kafkaClient.topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
				log_data := es.LogData{
					Topic: kafkaClient.topic,
					Data:  string(msg.Value),
				}
				es.SendToChan(log_data)   // 函数调函数 一个函数的执行时间和另一个函数相关，应该通过channel进行性能优化
			}
		}(pc)
	}
	defer kafkaClient.client.Close()
	select {}
}