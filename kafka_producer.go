package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

// Alter these as you like
const (
	tSleepSeconds = 5    // Time to pause, in seconds, between batches
	batchSize     = 5000 // Number of rows per batch
)

var (
	brokerList = flag.String("brokers", "localhost:9092", "The comma separated list of brokers in the Kafka cluster")
	topic      = flag.String("topic", "", "The topic to produce to")
	logger     = log.New(os.Stderr, "", log.LstdFlags)
)

// Usage -- overrides Usage within the flag package
var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n  %s reads from stdin and loads each line into Kafka\n\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
}

func main() {

	if len(os.Args) == 1 {
		Usage()
		return
	}

	flag.Parse()

	partitionerConstructor := sarama.NewRandomPartitioner

	var keyEncoder, valueEncoder sarama.Encoder

	config := sarama.NewConfig()
	config.Producer.Partitioner = partitionerConstructor
	producer, err := sarama.NewSyncProducer(strings.Split(*brokerList, ","), config)
	if err != nil {
		logger.Fatalln("FAILED to open the producer:", err)
	}
	defer producer.Close()

	// Loop over stdin here
	nLines := 1
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// Batch it
		if nLines%batchSize == 0 {
			time.Sleep(tSleepSeconds * time.Second)
		}
		line := scanner.Text()
		valueEncoder = sarama.StringEncoder(line)
		_, _, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: *topic,
			Key:   keyEncoder,
			Value: valueEncoder,
		})
		if err != nil {
			logger.Println("FAILED to produce message:", err)
		}
		nLines++
	}

}
