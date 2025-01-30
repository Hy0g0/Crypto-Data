package epitech.mpl10;

import com.launchdarkly.eventsource.EventSource;
import com.launchdarkly.eventsource.background.BackgroundEventHandler;
import com.launchdarkly.eventsource.background.BackgroundEventSource;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.common.serialization.StringSerializer;

import java.net.URI;
import java.util.Properties;
import java.util.concurrent.TimeUnit;

public class WikimediaChangesProducer {

    private static final KafkaProducer<String, String> kafkaProducer;
    private static final String WIKIMEDDIA_RECENTCHANGE_TOPIC = "wikimedia.recentchange";
    private static final String SOURCE_URL = "https://stream.wikimedia.org/v2/stream/recentchange";

    static {
        String bootstrapServers = "localhost:29092,localhost:39092,localhost:49092";

        Properties properties = new Properties();
        properties.setProperty(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, bootstrapServers);
        properties.setProperty(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        properties.setProperty(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());

        kafkaProducer = new KafkaProducer<>(properties);
    }

    public static void main(String[] args) throws InterruptedException {
        BackgroundEventHandler eventHandler = new WikimediaChangeHandler(kafkaProducer, WIKIMEDDIA_RECENTCHANGE_TOPIC);
        EventSource.Builder eventSourceBuilder = new EventSource.Builder(URI.create(SOURCE_URL));
        BackgroundEventSource.Builder backgroundEventSourceBuilder = new BackgroundEventSource.Builder(eventHandler, eventSourceBuilder);
        BackgroundEventSource backgroundEventSource = backgroundEventSourceBuilder.build();

        // Start the Producer in another thread
        backgroundEventSource.start();

        TimeUnit.MINUTES.sleep(20);
    }
}
