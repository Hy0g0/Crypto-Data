package epitech.mpl10.data;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import net.datafaker.Faker;

import java.time.Instant;
import java.util.*;

public class DataGenerator {

    private static final Random random = new Random();
    private static final Faker faker = new Faker();
    private static final ObjectMapper objectMapper = new ObjectMapper();

    private DataGenerator() {}

    public static Collection<String> generateLolMatchEvents(final int numberOfEvents) {
        int counter = 0;
        final List<String> lolMatchEvents = new ArrayList<>();

        while (counter++ < numberOfEvents) {
            try {
                String jsonContent = objectMapper.writeValueAsString(generateRandomLolEvent());
                lolMatchEvents.add(jsonContent);
            } catch (JsonProcessingException e) {
                e.printStackTrace();
            }
        }

        return lolMatchEvents;
    }

    private static Map<String, Object> generateRandomLolEvent() {
        Map<String, Object> content = new HashMap<>();

        content.put("id", UUID.randomUUID().toString());
        content.put("match", generateMatch());
        content.put("type", "kill_feed");
        content.put("payload", generatePayload());
        content.put("ts", Instant.now().toEpochMilli());
        content.put("ingame_timestamp", random.nextInt(1000));
        content.put("game", generateGame());

        return content;
    }

    private static Map<String, Object> generateMatch() {
        Map<String, Object> match = new HashMap<>();
        match.put("id", faker.number().numberBetween(1, 1000000));
        return match;
    }

    private static Map<String, Object> generatePayload() {
        Map<String, Object> payload = new HashMap<>();
        payload.put("killed", generatePlayer("red"));
        payload.put("type", "player");
        payload.put("killer", generatePlayer("blue"));
        payload.put("assists", Collections.singletonList(generateAssist("blue")));
        return payload;
    }

    private static Map<String, Object> generatePlayer(String team) {
        Map<String, Object> player = new HashMap<>();
        player.put("type", "player");
        player.put("object", generatePlayerObject(team));
        return player;
    }

    private static Map<String, Object> generatePlayerObject(String team) {
        Map<String, Object> playerObject = new HashMap<>();
        playerObject.put("id", faker.number().numberBetween(1, 100000));
        playerObject.put("name", faker.name().firstName());
        playerObject.put("role", team.equals("red") ? "adc" : "sup");
        playerObject.put("team", team);
        playerObject.put("champion", generateChampion());
        return playerObject;
    }

    private static Map<String, Object> generateChampion() {
        Map<String, Object> champion = new HashMap<>();
        champion.put("id", faker.number().numberBetween(1, 10000));
        champion.put("name", faker.leagueOfLegends().champion());
        return champion;
    }

    private static Map<String, Object> generateAssist(String team) {
        Map<String, Object> assist = new HashMap<>();
        assist.put("id", faker.number().numberBetween(1, 100000));
        assist.put("name", faker.name().firstName());
        assist.put("role", "sup");
        assist.put("team", team);
        assist.put("champion", generateChampion());
        return assist;
    }

    private static Map<String, Object> generateGame() {
        Map<String, Object> game = new HashMap<>();
        game.put("id", faker.number().numberBetween(1, 1000000));
        return game;
    }

    public static void main(String[] args) {
        try {
            String jsonContent = objectMapper.writerWithDefaultPrettyPrinter().writeValueAsString(generateRandomLolEvent());
            System.out.println(jsonContent);
        } catch (JsonProcessingException e) {
            e.printStackTrace();
        }
    }
}