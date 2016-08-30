#!/bin/bash

kafka_dir="./kafka_2.11-0.10.0.0"
nohup $kafka_dir/bin/zookeeper-server-start.sh $kafka_dir/config/zookeeper.properties >> zk.log 2>&1 </dev/null &

