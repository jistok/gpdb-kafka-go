#!/bin/bash

. ./kafka_env.sh
$kafka_dir/bin/kafka-topics.sh --list --zookeeper 104.198.249.198:2181

