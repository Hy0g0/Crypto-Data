import requests
import re
import time
import json
from confluent_kafka import Producer

url = "https://r.jina.ai/https://www.livecoinwatch.com/"
response = requests.get(url)
data = response.text
crypto_pattern = r"https:\/\/www\.livecoinwatch\.com\/price\/(\w+)"
price_pattern = r"\|\s\d+\s\|\s\[.*?\]\(.*?\)\s\|\s([\$\d,.]+)"
# Split the data into lines for easier processing
lines = data.splitlines()
crypto_prices = {}

def delivery_report(err, msg):
    if err is not None:
        print(f"Message delivery failed: {err}")
    else:
        print(f"Message delivered to {msg.topic()} [{msg.partition()}]")


# Configuration for the Kafka Producer
conf = {
    'bootstrap.servers': 'redpanda.redpanda.svc.cluster.local:9094'
}

# Initialize the producer once, outside the loop
producer = Producer(conf)

def scrape_and_send():
    try:
        response = requests.get(url, timeout=10).text  # Fetch the response text
    except requests.RequestException as e:
        print(f"Failed to fetch data: {e}")
        return

    x = re.findall(crypto_pattern, response)  # Find all matches for crypto symbols
    x = x[9:]
    
    y = re.findall(price_pattern, response)   # Find all matches for prices

    if len(x) != len(y):
        print(f"Warning: Mismatch between crypto symbols and prices (Symbols: {len(x)}, Prices: {len(y)})")
        return

    # Create a dictionary to store crypto-price pairs
    crypto_prices = {}
    for i in range(len(x)):
        crypto_prices[x[i]] = y[i]

    # Send each crypto and its price as a Kafka message
    for crypto, price in crypto_prices.items():
        try:
            # Prepare the message as a JSON object
            message = json.dumps({"crypto": crypto, "price": price, "site": "livecoin"})
            print(f"Sending: {crypto}: {price}")
            # Produce the message to Kafka
            producer.produce(
                'crypto-prices', 
                value=message,   # JSON message
                callback=delivery_report  # Callback function for delivery confirmation

            )
        except Exception as e:
            print(f"Failed to produce message: {e}")

    # Wait for all messages to be delivered
    producer.flush()


# Loop to continuously scrape and send data every second
while True:
    scrape_and_send()
    time.sleep(1)

