package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(env_path string) {
	wd, err := GetProjectRootPath()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Here's the working directory where environment file is to be loaded ==> ", wd)
	err = godotenv.Load(wd + env_path)
	if err != nil {
		log.Println(err)
		log.Fatalf("Error loading specified environment file")
	}
}

func GoDotEnvVariable(key string) string {
	return os.Getenv(key)
}
