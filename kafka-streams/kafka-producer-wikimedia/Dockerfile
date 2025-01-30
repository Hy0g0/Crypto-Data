FROM gradle:7.6.1-jdk17 AS build
WORKDIR /app
COPY . .
RUN ./gradlew
ENTRYPOINT ["./gradlew", "runProducer"]