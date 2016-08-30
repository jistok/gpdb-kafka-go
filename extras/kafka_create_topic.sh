#!/bin/bash

topic="chicago_crimes"
partition_count=2

./kafka_2.11-0.10.0.0/bin/kafka-topics.sh --create --topic $topic --replication-factor 1 \
  --partitions $partition_count --zookeeper localhost:2181

