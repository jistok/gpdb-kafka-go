#!/bin/bash

nohup ./kafka_2.11-0.10.0.0/bin/kafka-server-start.sh ./kafka_2.11-0.10.0.0/config/server.properties >> kafka.log 2>&1 </dev/null &

