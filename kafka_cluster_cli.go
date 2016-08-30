package main

/*

- Adapted from https://github.com/bsm/sarama-cluster/blob/master/cmd/sarama-cluster-cli/main.go

- This will store offsets in the __consumer_offsets topic, within Kafka (it doesn't
use Zookeeper for this).

- GPDB External Table Example
DROP EXTERNAL TABLE IF EXISTS crimes_kafka;
CREATE EXTERNAL WEB TABLE crimes_kafka
(LIKE crimes)
EXECUTE '$HOME/kafka_cluster_cli -brokers 172.16.1.1:9092 -group GPDBKafka -topics chicago_crimes 2>>$HOME/`printf "kafka_consumer_%02d.log" $GP_SEGMENT_ID`'
ON ALL FORMAT 'CSV' (DELIMITER ',' NULL '')
LOG ERRORS INTO err_crimes SEGMENT REJECT LIMIT 1 PERCENT;

*/

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

var (
	groupID    = flag.String("group", "", "REQUIRED: The shared consumer group name")
	brokerList = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The comma separated list of brokers in the Kafka cluster")
	topicList  = flag.String("topics", "", "REQUIRED: The comma separated list of topics to consume")
	offset     = flag.String("offset", "newest", "The offset to start with. Can be `oldest`, `newest`")
	verbose    = flag.Bool("verbose", false, "Whether to turn on sarama logging")

	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func main() {
	flag.Parse()

	if *groupID == "" {
		printUsageErrorAndExit("You have to provide a -group name.")
	} else if *brokerList == "" {
		printUsageErrorAndExit("You have to provide -brokers as a comma-separated list, or set the KAFKA_PEERS environment variable.")
	} else if *topicList == "" {
		printUsageErrorAndExit("You have to provide -topics as a comma-separated list.")
	}

	// Init config
	config := cluster.NewConfig()
	if *verbose {
		sarama.Logger = logger
	} else {
		config.Consumer.Return.Errors = true
		config.Group.Return.Notifications = true
	}

	switch *offset {
	case "oldest":
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	case "newest":
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	default:
		printUsageErrorAndExit("-offset should be `oldest` or `newest`")
	}

	// Init consumer, consume errors & messages
	consumer, err := cluster.NewConsumer(strings.Split(*brokerList, ","), *groupID, strings.Split(*topicList, ","), config)
	if err != nil {
		printErrorAndExit(69, "Failed to start consumer: %s", err)
	}

	go func() {
		for err := range consumer.Errors() {
			logger.Printf("Error: %s\n", err.Error())
		}
	}()

	go func() {
		for note := range consumer.Notifications() {
			logger.Printf("Rebalanced: %+v\n", note)
		}
	}()

	// TODO: embed a timeout here, so that, if there hasn't been any message consumed
	// for some period of time, we close this process down.
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(500 * time.Millisecond)
		timeout <- true
	}()

	go func() {
		for {
			select {
			case msg := <-consumer.Messages():
				//fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
				fmt.Printf("%s\n", msg.Value)
				consumer.MarkOffset(msg, "")
			case <-timeout:
				syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
				return // Must exit this function
			}
		}
	}()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-wait

	if err := consumer.Close(); err != nil {
		logger.Println("Failed to close consumer: ", err)
	}
}

func printErrorAndExit(code int, format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", fmt.Sprintf(format, values...))
	fmt.Fprintln(os.Stderr)
	os.Exit(code)
}

func printUsageErrorAndExit(format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", fmt.Sprintf(format, values...))
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Available command line options:")
	flag.PrintDefaults()
	os.Exit(64)
}
