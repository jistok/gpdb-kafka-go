#!/bin/bash

. ./kafka_env.sh
$kafka_dir/bin/kafka-run-class.sh kafka.tools.ImportZkOffsets --input-file zk_offsets.txt --zkconnect localhost:2181

