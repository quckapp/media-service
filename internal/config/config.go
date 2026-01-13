package config

import "os"

type Config struct {
	Port          string
	MongoURI      string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	JWTSecret     string
	AWSRegion     string
	AWSAccessKey  string
	AWSSecretKey  string
	S3Bucket      string
	KafkaBrokers  string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "5001"),
		MongoURI:      getEnv("MONGODB_URI", "mongodb://localhost:27017/quckchat_media"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		AWSRegion:     getEnv("AWS_REGION", "ap-south-1"),
		AWSAccessKey:  getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretKey:  getEnv("AWS_SECRET_ACCESS_KEY", ""),
		S3Bucket:      getEnv("AWS_S3_BUCKET", "quckchat-media"),
		KafkaBrokers:  getEnv("KAFKA_BROKERS", "localhost:9092"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
