#!/bin/bash

topic="chicago_crimes"

./kafka_2.11-0.10.0.0/bin/kafka-topics.sh --delete --topic $topic --zookeeper localhost:2181

