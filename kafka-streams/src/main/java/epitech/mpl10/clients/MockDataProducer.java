package epitech.mpl10.clients;

import epitech.mpl10.data.DataGenerator;
import org.apache.kafka.clients.producer.*;
import org.apache.kafka.common.serialization.StringSerializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.Random;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class MockDataProducer implements AutoCloseable {
    private static final Logger LOG = LoggerFactory.getLogger(MockDataProducer.class);

    private ExecutorService executorService = Executors.newFixedThreadPool(1);
    private static final String LOL_TOPIC = "league-of-legends";
    private static final String NULL_KEY = null;
    private volatile boolean keepRunning = true;
    private final Random random = new Random();

    @Override
    public void close() {
        LOG.info("Shutting down data generation");
        keepRunning = false;
        if (executorService != null) {
            executorService.shutdownNow();
            executorService = null;
        }
    }

    public void produceLolMatchEventsData() {
        produceLolMatchEventsData(false);
    }

    public void produceLolMatchEventsDataSchemaRegistry() {
        produceLolMatchEventsData(true);
    }

    private void produceLolMatchEventsData(boolean produceSchemaRegistry) {
        // Runnable generateTask = () -> {
        LOG.info("Creating task for generating mock purchase transactions");
        final Map<String, Object> configs = producerConfigs();
        final Callback callback = (metadata, exception) -> {
            if (exception != null) {
                System.out.printf("Producing records encountered error %s %n", exception);
            } else {
                System.out.printf("Record produced - offset - %d timestamp - %d %n", metadata.offset(), metadata.timestamp());
            }

        };

        if (produceSchemaRegistry) {
            // configs.put("schema.registry.url", "http://localhost:8081");
            // configs.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, KafkaProtobufSerializer.class);
        } else {
            // configs.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, ProtoSerializer.class);
        }

        try (Producer<String, String> producer = new KafkaProducer<>(configs)) {
            while (keepRunning) {
                Collection<String> lolMatchEvents = DataGenerator.generateLolMatchEvents(20);

                var producerRecords = lolMatchEvents.stream()
                        .map(r -> new ProducerRecord<>(LOL_TOPIC, "1111", r))
                        .toList();
                producerRecords.forEach((pr -> producer.send(pr, callback)));
            }
        }
        LOG.info("Done generating data");
        // };
        // executorService.submit(generateTask);
    }

    private static Map<String, Object> producerConfigs() {
        Map<String, Object> producerConfigs = new HashMap<>();
        producerConfigs.put("bootstrap.servers", "localhost:29092,localhost:39092,localhost:49092");
        producerConfigs.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, StringSerializer.class);
        producerConfigs.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, StringSerializer.class);
        producerConfigs.put("acks", "all");
        return producerConfigs;
    }
}
