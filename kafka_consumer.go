package main

/*

BUILD:
- Install / set up Go
- `go get github.com/Shopify/sarama`
- `github.com/wvanbergen/kafka/consumergroup`
- `go build kafka_consumer.go`

INSTALL:
Use gpscp to install the resulting `kafka_consumer` executable into ~gpadmin/ on all segment hosts

EXAMPLE:
-- Here, the table name, pos_trans, happens to be the same as the Kafka topic name.
CREATE EXTERNAL WEB TABLE pos_trans_kafka
(LIKE pos_trans)
EXECUTE '/tmp/kafka_consumer -zookeeper 10.0.2.2:2181 -topic pos_trans 2>>$HOME/`printf "kafka_consumer_%02d.log" $GP_SEGMENT_ID`'
ON ALL FORMAT 'TEXT' (DELIMITER '|' NULL '');

TODO:
	- Make the `zookeeper` values accessible via an environment variable, so you
		don't need to pass via command line

*/

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
)

const (
	waitMilliseconds  = 250                   // Maximum time to wait for data, ms; exits if no data appears within this window
	consumerGroupName = "GPDB_Consumer_Group" // Arbitrary, but constant
)

var (
	topic    = flag.String("topic", "", "Name of the topic to consume")
	zkString = flag.String("zookeeper", "localhost:2181", "Comma separated list of ZK host:port values")
	verbose  = flag.Bool("verbose", false, "Whether to turn on sarama logging")

	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func main() {

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	flag.Parse()

	if *verbose {
		sarama.Logger = logger
	}

	zkPeers := strings.Split(*zkString, ",")

	config := consumergroup.NewConfig()

	consumer, consumerErr := consumergroup.JoinConsumerGroup(
		consumerGroupName,
		[]string{*topic},
		zkPeers,
		//nil)
		config)

	if consumerErr != nil {
		log.Fatalln(consumerErr)
	}

	defer func() {
		if closeErr := consumer.Close(); closeErr != nil {
			panic(closeErr)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Get signal for finish
	doneCh := make(chan struct{})
	tStart := time.Now()
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages(): // Messages appears to return one message here
				fmt.Println(string(msg.Value))
				consumer.CommitUpto(msg)
			case <-signals:
				doneCh <- struct{}{}
			default:
				elapsed := time.Since(tStart)
				if elapsed > waitMilliseconds*time.Millisecond {
					signals <- os.Interrupt
				}
			}
		}
	}()

	<-doneCh
}
