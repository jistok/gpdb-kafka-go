#!/bin/bash

topic="chicago_crimes"

. ./kafka_env.sh
$kafka_dir/bin/kafka-topics.sh --delete --topic $topic --zookeeper localhost:2181

