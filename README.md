# gpdb-kafka-go

This is a simple Kafka consumer to enable Greenplum Database (GPDB) to load data from
a Kafka topic using GPDB's _external web table_ capability.  This client is written
in Go.

## What is Kafka?
_Kafka is a distributed, partitioned, replicated commit log service. It provides the functionality of a messaging system, but with a unique design._ See the [Kafka docs](http://kafka.apache.org/documentation.html#introduction) for a nice introduction.

## What is GPDB?
GPDB is an open source massively parallel processing (MPP) relational database.  Using its
[_external tables_](http://gpdb.docs.pivotal.io/4370/ref_guide/sql_commands/CREATE_EXTERNAL_TABLE.html)
feature, it is able to load massive amounts of data very efficiently by exploiting its parallel
architecture.  This data can be in files on ETL hosts, within HDFS, or in Amazon S3.  This Kafka
consumer will utilize GPDB _external web tables_, since they are able to run an executable program,
in parallel, to ingest data.

## Why Go?
* I am trying to learn Go.
* Go builds a single, statically linked binary, which has all its dependencies built in.

## What's Required to Get This Running?
* A Running Kafka installation.  I followed the [Quick Start](http://kafka.apache.org/07/quickstart.html).
* A [Go installation](https://golang.org/doc/install)
* An installation of GPDB.  The [GPDB Sandbox VM](https://network.pivotal.io/products/pivotal-gpdb#/releases/1683/file_groups/411) would work just fine for trying this out.
* The [Go source file](./kafka_consumer.go) for the Kafka consumer

## Running Through an Example
1. Resolve the dependencies: `go get github.com/wvanbergen/kafka/consumergroup github.com/Shopify/sarama`
1. Build the executable: `go build kafka_consumer.go`
1. Install the resulting executable, kafka_consumer, into the ``$HOME` directory of the gpadmin user on each of your GPDB segment hosts (or, just onto the single host if you are using the GPDB Sandbox VM).
