from kafka import KafkaProducer
import requests
import json
import time
from datetime import datetime

# Configuration
KAFKA_BOOTSTRAP_SERVERS = 'localhost:9092'
KAFKA_TOPIC = 'crypto_prices'
COINGECKO_API_URL = 'https://api.coingecko.com/api/v3'
CRYPTO_IDS = ['bitcoin', 'ethereum', 'ripple'] 

# Initialisation du producteur Kafka
producer = KafkaProducer(
    bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS,
    value_serializer=lambda x: json.dumps(x).encode('utf-8')
)

def get_crypto_data():
    """Récupère les données de prix pour les cryptomonnaies spécifiées"""
    try:
        response = requests.get(
            f"{COINGECKO_API_URL}/simple/price",
            params={
                'ids': ','.join(CRYPTO_IDS),
                'vs_currencies': 'usd,eur',
                'include_24hr_change': 'true',
                'include_market_cap': 'true'
            }
        )
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Erreur lors de la requête API: {e}")
        return None

def main():
    while True:
        try:
            # Récupération des données
            crypto_data = get_crypto_data()
            
            if crypto_data:
                # Ajout du timestamp
                message = {
                    'timestamp': datetime.now().isoformat(),
                    'data': crypto_data
                }
                
                # Envoi des données à Kafka
                producer.send(KAFKA_TOPIC, value=message)
                print(f"Données envoyées à Kafka: {message}")
            
            # Attente de 1 minute avant la prochaine requête
            # Pour respecter les limites de l'API CoinGecko
            time.sleep(60)
            
        except Exception as e:
            print(f"Erreur: {e}")
            time.sleep(10)

if __name__ == "__main__":
    try:
        main()
    finally:
        producer.close() 