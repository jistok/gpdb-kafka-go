#!/bin/bash

. ./kafka_env.sh
$kafka_dir/bin/kafka-topics.sh --list --zookeeper $zk_ip_port

