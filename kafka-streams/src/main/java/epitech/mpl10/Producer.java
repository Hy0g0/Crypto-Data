package epitech.mpl10;

import epitech.mpl10.clients.MockDataProducer;

public class Producer {
    public static void main(String[] args) throws InterruptedException {
        try (MockDataProducer mockDataProducer = new MockDataProducer()) {
            mockDataProducer.produceLolMatchEventsData();
//            Thread.sleep(10000);
        }
    }
}
