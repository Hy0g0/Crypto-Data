import json
from confluent_kafka import Consumer, KafkaException
import clickhouse_connect

# Configuration for the Kafka Consumer
kafka_conf = {
    'bootstrap.servers': 'kafka.kafka.svc.cluster.local:9094',
    'group.id': 'crypto-consumer-group',
    'auto.offset.reset': 'earliest',
}

# Initialize the Kafka Consumer
consumer = Consumer(kafka_conf)
consumer.subscribe(['crypto-prices'])

# ClickHouse connection setup
client = clickhouse_connect.get_client(host='clickhouse-service.clickhouse.svc.cluster.local', port=8123)

# Function to insert data into ClickHouse
def insert_into_clickhouse(crypto, price, site):
    insert_query = """
    INSERT INTO crypto_prices (crypto, price, site)
    VALUES (%(crypto)s, %(price)s, %(site)s)
    """
    client.command(insert_query, {'crypto': crypto, 'price': price, 'site': site})

# Function to consume messages from Kafka
def consume_from_kafka():
    print("Starting Kafka Consumer...")

    try:
        while True:
            # Poll Kafka for messages
            msg = consumer.poll(timeout=1.0)  # Poll for 1 second
            
            if msg is None:  # No message
                continue
            
            if msg.error():  # Handle any errors
                raise KafkaException(msg.error())
            
            # Process the message (assuming it's a JSON message)
            message = msg.value().decode('utf-8')
            print(f"Received message: {message}")
            
            try:
                # Parse the JSON message
                data = json.loads(message)
                crypto = data.get("crypto")
                price = data.get("price")
                site = data.get("site")
                
                # Insert data into ClickHouse
                insert_into_clickhouse(crypto, price, site)
                print(f"Inserted into ClickHouse: {crypto}, {price}, {site}")
            
            except json.JSONDecodeError as e:
                print(f"Failed to decode JSON: {e}")
    
    finally:
        # Close the consumer when done
        consumer.close()

# Run the consumer
consume_from_kafka()
