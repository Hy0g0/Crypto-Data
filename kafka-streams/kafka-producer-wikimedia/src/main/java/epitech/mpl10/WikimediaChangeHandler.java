package epitech.mpl10;

import com.launchdarkly.eventsource.MessageEvent;
import com.launchdarkly.eventsource.background.BackgroundEventHandler;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.json.JSONObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class WikimediaChangeHandler implements BackgroundEventHandler {

    private final KafkaProducer<String, String> kafkaProducer;
    private final String topic;
    private static final Logger log = LoggerFactory.getLogger(WikimediaChangeHandler.class.getSimpleName());

    public WikimediaChangeHandler(KafkaProducer<String, String> kafkaProducer, String topic) {
        this.kafkaProducer = kafkaProducer;
        this.topic = topic;
    }

    @Override
    public void onOpen() {

    }

    @Override
    public void onClosed() {
        this.kafkaProducer.close();
    }

    @Override
    public void onMessage(String s, MessageEvent messageEvent) {
        String message = messageEvent.getData();
        log.info("Message received: {}", message);

        // Parse the JSON message
        JSONObject jsonObject = new JSONObject(message);
        String domain = jsonObject.getJSONObject("meta").getString("domain");

        // Async
        this.kafkaProducer.send(new ProducerRecord<>(this.topic, domain, message));
    }

    @Override
    public void onComment(String s) {

    }

    @Override
    public void onError(Throwable throwable) {

    }
}
