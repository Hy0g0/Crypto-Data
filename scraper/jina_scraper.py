import requests
import re

url = "https://r.jina.ai/https://www.livecoinwatch.com/"
response = requests.get(url)
data = response.text

crypto_pattern = r"https:\/\/www\.livecoinwatch\.com\/price\/(\w+-\w+)"
price_pattern = r"\$(\d+[\d,]*\.?\d*)"

# Split the data into lines for easier processing
lines = data.splitlines()

crypto_prices = {}

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

for crypto, price in crypto_prices.items():
    print(f"{crypto}: {price}")