import requests
import re
import time
from confluent_kafka import Producer

# Callback function for delivery confirmation
def delivery_report(err, msg):
    if err is not None:
        print(f"Message delivery failed: {err}")
    else:
        print(f"Message delivered to {msg.topic()} [{msg.partition()}]")


# Configuration for the Kafka Producer
conf = {
    'bootstrap.servers': 'redpanda.redpanda.svc.cluster.local:9094',
    'group.id': 'consumer',
    'session.timeout.ms': 6000,
    'auto.offset.reset': 'earliest',
}

# Initialize the producer once, outside the loop
producer = Producer(conf)

url = "https://r.jina.ai/https://www.livecoinwatch.com/"
crypto_pattern = r"https:\/\/www\.livecoinwatch\.com\/price\/(\w+-\w+)"
price_pattern = r"\$(\d+[\d,]*\.?\d*)"

def scrape_and_send():
    print("Scraping data...")
    try:
        response = requests.get(url, timeout=10)
        response.raise_for_status()  # Raise an error for bad status codes
    except requests.RequestException as e:
        print(f"Failed to fetch data: {e}")
        return

    data = response.text
    lines = data.splitlines()

    crypto_prices = {}

    # Extract crypto names and prices from the webpage
    for i in range(len(lines)):
        match = re.search(crypto_pattern, lines[i])
        if match:
            crypto_name = match.group(1)
            
            # Jump 1 line and then check for the price
            if i + 2 < len(lines):
                price_match = re.search(price_pattern, lines[i + 2])
                if price_match:
                    price = price_match.group(0)
                    crypto_prices[crypto_name] = price

    # Send crypto prices to Redpanda topic
    for crypto, price in crypto_prices.items():
        try:
            print(f"Sending: {crypto}: {price}")
            producer.produce(
                'crypto-prices', 
                key="livecoin",
                value=price,
                callback=delivery_report
            )
        except Exception as e:
            print(f"Failed to produce message: {e}")

    # Wait for all messages to be delivered
    producer.flush()

# Loop to continuously scrape and send data every second
while True:
    scrape_and_send()
    time.sleep(1)

