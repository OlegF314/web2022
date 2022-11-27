package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/namsral/flag"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/go-redis/redis/v8"
)

var (
	// kafka
	kafkaBrokerUrl     string
	kafkaVerbose       bool
	kafkaTopic         string
	kafkaConsumerGroup string
	kafkaClientId      string
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0, 
    })
	flag.StringVar(&kafkaBrokerUrl, "kafka-brokers", "localhost:19092,localhost:29092", "Kafka brokers in comma separated value")
	flag.BoolVar(&kafkaVerbose, "kafka-verbose", true, "Kafka verbose logging")
	flag.StringVar(&kafkaTopic, "kafka-topic", "foo", "Kafka topic. Only one topic per worker.")
	flag.StringVar(&kafkaConsumerGroup, "kafka-consumer-group", "consumer-group", "Kafka consumer group")
	flag.StringVar(&kafkaClientId, "kafka-client-id", "my-client-id", "Kafka client id")

	flag.Parse()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	brokers := strings.Split(kafkaBrokerUrl, ",")

	config := kafka.ReaderConfig{
		Brokers:         brokers,
		GroupID:         kafkaClientId,
		Topic:           kafkaTopic,
		MinBytes:        10e3,            
		MaxBytes:        10e6,            
		MaxWait:         1 * time.Second, 
		ReadLagInterval: -1,
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	for {
		_, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Error().Msgf("error while receiving message: %s", err.Error())
			continue
		}

		err = rdb.Set(ctx, "task", "in process", 0).Err()
    		if err != nil {
        		panic(err)
        	}
		time.Sleep(3)
		err = rdb.Set(ctx, "task", "done", 0).Err()
    		if err != nil {
        		panic(err)
    		}
	}
}