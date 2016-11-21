#!/bin/bash

. ./kafka_env.sh
$kafka_dir/bin/kafka-run-class.sh kafka.tools.ExportZkOffsets --group GPDB_Consumer_Group --output-file zk_offsets.txt --zkconnect localhost:2181

