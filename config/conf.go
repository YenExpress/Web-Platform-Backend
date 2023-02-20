package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"

	cogman_config "github.com/Joker666/cogman/config"

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

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func loadLocalTestEnv() {
	err := godotenv.Load("../../.env.test.local")
	if err != nil {
		log.Fatalf("Error loading .env.test.local file")
	}
}

func goDotEnvVariable(key string) string {
	loadEnv()
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

		RedisTTL: time.Hour * 24 * 7,

		HighPriorityQueueCount: 2,
		LowPriorityQueueCount:  4}
	CustomValidatorsActive bool   = RegisterValidation()
	DatabaseURI            string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=verify-full",
		goDotEnvVariable("DBHost"), goDotEnvVariable("DBUser"), goDotEnvVariable("DBPwd"), goDotEnvVariable("DBName"))

	JwtSecret          string = goDotEnvVariable("JwtSecret")
	WebClientDomain    string = goDotEnvVariable("WebClientDomain")
	ServerDomain       string = goDotEnvVariable("ServerDomain")
	PatientAPIKey      string = goDotEnvVariable("PatientAPIKey")
	PatientRedisClient        = redis.NewClient(&redis.Options{
		Addr:     goDotEnvVariable("RedisHost") + ":" + goDotEnvVariable("RedisPort"),
		Password: goDotEnvVariable("RedisPwd"),
		DB:       0,
	})
)
