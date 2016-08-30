#!/bin/bash

./kafka_2.11-0.10.0.0/bin/kafka-run-class.sh kafka.tools.ImportZkOffsets --input-file zk_offsets.txt --zkconnect localhost:2181

