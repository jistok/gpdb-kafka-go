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

var (
	brokerList  = flag.String("brokers", "localhost:9092", "The comma separated list of brokers in the Kafka cluster")
	topic       = flag.String("topic", "", "The topic to produce to")
	partitioner = flag.String("partitioner", "random", "The partitioning scheme to use. Can be `hash`, or `random`")
	verbose     = flag.Bool("verbose", false, "Whether to turn on sarama logging")

	logger = log.New(os.Stderr, "", log.LstdFlags)
)

// Usage -- override Usage within the flag package
var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n  %s reads from stdin and, loads each line into Kafka\n\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	if len(os.Args) == 1 {
		Usage()
		return
	}

	flag.Parse()

	if *verbose {
		sarama.Logger = logger
	}

	var partitionerConstructor sarama.PartitionerConstructor
	switch *partitioner {
	case "hash":
		partitionerConstructor = sarama.NewHashPartitioner
	case "random":
		partitionerConstructor = sarama.NewRandomPartitioner
	default:
		log.Fatalf("Partitioner %s not supported.", *partitioner)
	}

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
		if nLines%2500 == 0 {
			time.Sleep(5 * time.Second)
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
