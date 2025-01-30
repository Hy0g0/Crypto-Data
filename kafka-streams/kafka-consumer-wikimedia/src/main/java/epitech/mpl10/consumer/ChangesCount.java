package epitech.mpl10.consumer;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.*;
import org.apache.kafka.streams.kstream.Consumed;
import org.apache.kafka.streams.kstream.Grouped;
import org.apache.kafka.streams.kstream.Produced;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.FileInputStream;
import java.io.IOException;
import java.time.Duration;
import java.util.Properties;
import java.util.concurrent.CountDownLatch;

public class ChangesCount {
    private static final Logger LOG = LoggerFactory.getLogger(ChangesCount.class);
    private static final String CHANGES_COUNT_TOPIC = "wikimedia.recentchange.changescount";

    private static Properties getProperties() {
        Properties streamsProps = new Properties();
        try (FileInputStream fis = new FileInputStream("src/main/resources/streams.properties")) {
            streamsProps.load(fis);
        } catch (IOException e) {
            LOG.error("Error while loading the streams.properties file. Make sure you have it properly configured.");
            throw new RuntimeException(e);
        }

        streamsProps.put(StreamsConfig.APPLICATION_ID_CONFIG, "wikimedia-recentchange");
        streamsProps.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");
        streamsProps.put(StreamsConfig.COMMIT_INTERVAL_MS_CONFIG, 200);

        return streamsProps;
    }

    public Topology topology(Properties properties) {
        StreamsBuilder streamsBuilder = new StreamsBuilder();

        streamsBuilder
                .stream(
                        properties.getProperty("input.topic"),
                        Consumed.with(Serdes.String(), Serdes.String())
                )

                .groupBy((key, value) -> key, Grouped.with(Serdes.String(), Serdes.String()))

                .count()

                .toStream()

                .to(CHANGES_COUNT_TOPIC, Produced.with(Serdes.String(), Serdes.Long()));

        return streamsBuilder.build(properties);
    }

    public static void main(String[] args) {
        ChangesCount consumer = new ChangesCount();
        Properties streamsProps = getProperties();
        Topology topology = consumer.topology(streamsProps);

        try (KafkaStreams kafkaStreams = new KafkaStreams(topology, streamsProps)) {
            final CountDownLatch shutdownLatch = new CountDownLatch(1);

            Runtime.getRuntime().addShutdownHook(new Thread(() -> {
                kafkaStreams.close(Duration.ofSeconds(2));
                shutdownLatch.countDown();
            }));
            try {
                kafkaStreams.start();
                shutdownLatch.await();
            } catch (Throwable e) {
                System.exit(1);
            }
        }
        System.exit(0);
    }
}