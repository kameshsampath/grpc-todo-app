package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func GetEnv() string {
	return getEnvironmentValue("ENV")
}

func GetPort() int {
	port, _ := strconv.Atoi(getEnvironmentValue("PORT"))
	return port
}

func GetBrokers() []string {
	brokers := getEnvironmentValue("BROKERS")
	return strings.Split(brokers, ",")
}

func GetTopics() []string {
	topics := getEnvironmentValue("TOPICS")
	return strings.Split(topics, ",")
}

func GetDefaultProducerTopic() string {
	return GetTopics()[0]
}

func getEnvironmentValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("%s environment variable is missing.", key)
	}
	return os.Getenv(key)
}
