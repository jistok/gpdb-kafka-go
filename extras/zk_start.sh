#!/bin/bash

. ./kafka_env.sh
nohup $kafka_dir/bin/zookeeper-server-start.sh $kafka_dir/config/zookeeper.properties >> zk.log 2>&1 </dev/null &

