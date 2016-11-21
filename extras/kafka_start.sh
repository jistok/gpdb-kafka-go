#!/bin/bash

. ./kafka_env.sh
nohup $kafka_dir/bin/kafka-server-start.sh $kafka_dir/config/server.properties >> kafka.log 2>&1 </dev/null &

