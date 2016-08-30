#!/bin/bash

./kafka_2.11-0.10.0.0/bin/kafka-run-class.sh kafka.tools.ExportZkOffsets --group GPDB_Consumer_Group --output-file zk_offsets.txt --zkconnect localhost:2181

