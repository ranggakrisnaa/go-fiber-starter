#!/bin/bash

echo "Starting Kafka setup..."

# Start Kafka containers
docker-compose -f docker-compose.kafka.yml up -d

echo "Waiting for Kafka to start..."
sleep 10

# Create topics
echo "Creating topics..."
docker exec kafka kafka-topics --create --topic user-events --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --if-not-exists
docker exec kafka kafka-topics --create --topic auth-events --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --if-not-exists

# List topics
echo "Available topics:"
docker exec kafka kafka-topics --list --bootstrap-server localhost:9092

echo "Kafka setup complete!"
echo "Kafka UI: http://localhost:8080"
echo "Kafka Broker: localhost:9092"