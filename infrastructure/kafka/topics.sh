#!/bin/bash
# topics.sh - Initializes Kafka topic channels.

echo "--> Provisioning Kafka event topics..."

# Wait for Kafka to boot
sleep 5

# Create core transaction/event streams
kafka-topics.sh --create --if-not-exists --bootstrap-server localhost:9092 --replication-factor 1 --partitions 3 --topic checkout.submitted
kafka-topics.sh --create --if-not-exists --bootstrap-server localhost:9092 --replication-factor 1 --partitions 3 --topic order.created
kafka-topics.sh --create --if-not-exists --bootstrap-server localhost:9092 --replication-factor 1 --partitions 3 --topic inventory.deducted

echo "Kafka topics set up successfully."
