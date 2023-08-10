package config

import (
	"fmt"

	"github.com/ignitedotdev/auth-ms/internal/utils"
)

type configModel struct {
	Environment         string
	DevClientOrigin     string
	StagingClientOrigin string
	ServicePort         string
	PostgresURI         string
	JwtSecret           string
	WebClientDomain     string
	ServiceDomain       string
	APIKey              string
}

func (c *configModel) loadfromEnv() {
	c.Environment = utils.GoDotEnvVariable("Environment")
	c.DevClientOrigin = utils.GoDotEnvVariable("DevCrossOrigin")
	c.StagingClientOrigin = utils.GoDotEnvVariable("StagingCrossOrigin")
	c.ServicePort = utils.GoDotEnvVariable("ServicePort")
	c.PostgresURI = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s,sslmode=%s",
		utils.GoDotEnvVariable("DBHost"), utils.GoDotEnvVariable("DBUser"), utils.GoDotEnvVariable("DBPwd"),
		utils.GoDotEnvVariable("DBName"), utils.GoDotEnvVariable("DBPort"), utils.GoDotEnvVariable("SSLMode"))

	c.JwtSecret = utils.GoDotEnvVariable("JwtSecret")
	c.WebClientDomain = utils.GoDotEnvVariable("WebClientDomain")
	c.ServiceDomain = utils.GoDotEnvVariable("ServerDomain")
	c.APIKey = utils.GoDotEnvVariable("APIKey")
}

func LoadConfig(env_path string) {
	utils.LoadEnv(env_path)
	Config.loadfromEnv()
}

var Config *configModel = new(configModel)
