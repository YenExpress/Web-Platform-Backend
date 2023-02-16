package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"

	cogman_config "github.com/Joker666/cogman/config"
	"github.com/joho/godotenv"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// create function to validate based on enums
func Enum(
	fl validator.FieldLevel,
) bool {

	enumString := fl.Param()
	value := fl.Field().String()
	enumSlice := strings.Split(enumString, "_")
	for _, v := range enumSlice {
		if value == v {
			return true
		}
	}
	return false
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func RegisterValidation() bool {
	validate = validator.New()
	validate.RegisterValidation("Enum", Enum)
	return true
}

var (
	TaskMasterCfg *cogman_config.Config = &cogman_config.Config{
		ConnectionTimeout: time.Minute * 10,
		RequestTimeout:    time.Second * 5,

		AmqpURI:  goDotEnvVariable("AmqpURI"),
		RedisURI: goDotEnvVariable("RedisURI"),
		MongoURI: goDotEnvVariable("MongoURI"),

		RedisTTL: time.Hour * 24 * 7,  // optional. default value 1 week
		MongoTTL: time.Hour * 24 * 30, // optional. default value 1 month

		HighPriorityQueueCount: 2,
		LowPriorityQueueCount:  4}
	CustomValidatorsActive bool   = RegisterValidation()
	DatabaseURI            string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=verify-full",
		goDotEnvVariable("DBHost"), goDotEnvVariable("DBUser"), goDotEnvVariable("DBPwd"), goDotEnvVariable("DBName"))

	JwtSecret          string = goDotEnvVariable("JwtSecret")
	WebClientDomain    string = goDotEnvVariable("WebClientDomain")
	ServerDomain       string = goDotEnvVariable("ServerDomain")
	PatientRedisClient        = redis.NewClient(&redis.Options{
		Addr:     goDotEnvVariable("RedisHost") + ":" + goDotEnvVariable("RedisPort"),
		Password: goDotEnvVariable("RedisPwd"),
		DB:       0,
	})
)
