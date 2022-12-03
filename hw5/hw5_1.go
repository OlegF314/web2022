package main
import (
  "fmt"
  "net/http"
  "context"
  "github.com/go-redis/redis/v8"
  "time"
  "github.com/segmentio/kafka-go"
  "github.com/segmentio/kafka-go/snappy"
)


var ctx = context.Background()
var writer *kafka.Writer

func Configure(kafkaBrokerUrls []string, clientId string, topic string) (w *kafka.Writer, err error) {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientId,
	}

	config := kafka.WriterConfig{
		Brokers:          kafkaBrokerUrls,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	w = kafka.NewWriter(config)
	writer = w
	return w, nil
}

func Push(parent context.Context, key, value []byte) (err error) {
  message := kafka.Message{
    Key:   key,
    Value: value,
    Time:  time.Now(),
  }
  return writer.WriteMessages(parent, message)
}

func main() {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", 
        DB:       0,  
    })
    http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
        Push(ctx, nil, []byte("task"))
	fmt.Println("task got")
    })
    http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
        val, err := rdb.Get(ctx, "task").Result()
        if err != nil {
            panic(err)
        }
        fmt.Println("task", val)
    })
    http.ListenAndServe(":80", nil)
}

