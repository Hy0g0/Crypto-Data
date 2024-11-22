from pycoingecko import CoinGeckoAPI
from kafka import KafkaProducer
import json
import time
import logging
from typing import Dict, Any

class CoinGeckoScraper:
    def __init__(self, kafka_bootstrap_servers: str = 'localhost:9092'):
        self.cg = CoinGeckoAPI()
        self.producer = KafkaProducer(
            bootstrap_servers=kafka_bootstrap_servers,
            value_serializer=lambda x: json.dumps(x).encode('utf-8')
        )
        self.logger = logging.getLogger(__name__)

    def get_crypto_data(self) -> Dict[str, Any]:
        try:
            # Récupère les top 100 cryptos par capitalisation
            return self.cg.get_coins_markets(
                vs_currency='usd',
                order='market_cap_desc',
                per_page=100,
                sparkline=False
            )
        except Exception as e:
            self.logger.error(f"Erreur lors de la récupération des données: {e}")
            return {}

    def send_to_kafka(self, data: Dict[str, Any], topic: str = 'crypto_data'):
        try:
            self.producer.send(topic, value=data)
            self.producer.flush()
        except Exception as e:
            self.logger.error(f"Erreur lors de l'envoi vers Kafka: {e}")

    def run(self, interval: int = 60):
        while True:
            data = self.get_crypto_data()
            if data:
                self.send_to_kafka(data)
                self.logger.info(f"Données envoyées pour {len(data)} cryptomonnaies")
            time.sleep(interval)

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    scraper = CoinGeckoScraper()
    scraper.run() 