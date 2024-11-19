package config

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Configs struct {
	App    App    `mapstructure:"app"`
	Secret Secret `mapstructure:"secrets"`
}

type App struct {
	Env string `mapstructure:"env"`
}

type Secret struct {
	GcsCredential map[string]interface{} `mapstructure:"gcs-credential"`
	Datamock      string                 `mapstructure:"datamock"`
}

// var MainConfig *Configs

func InitConfig(ctx context.Context) *Configs {
	configPath, ok := os.LookupEnv("API_CONFIG_PATH")
	if configPath == "" || !ok {
		log.Print("API_CONFIG_PATH not found, using default config")
		configPath = "./config"
	}

	configName, ok := os.LookupEnv("API_CONFIG_NAME")
	if configName == "" || !ok {
		log.Print("API_CONFIG_NAME not found, using default config")
		configName = "config"
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	// viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Print("API_CONFIG_NAME not found, using default config")
		log.Panic(err)
	}

	// viper.AutomaticEnv()
	viper.AutomaticEnv()

	if err := viper.MergeConfig(strings.NewReader(viper.GetString("configs"))); err != nil {
		log.Panic(err.Error())
	}

	// fmt.Println(viper.Get("app.env"))
	// fmt.Println(viper.Get("secrets.gcs-credential"))
	// fmt.Println("configs :", viper.GetString("configs"))
	GetSecretValue()

	var mainConfigs *Configs
	if err := viper.Unmarshal(&mainConfigs); err != nil {
		log.Panic(err.Error())
	}

	return mainConfigs
}

func GetSecretValue() {
	for _, value := range os.Environ() {
		pair := strings.SplitN(value, "=", 2)
		if strings.Contains(pair[0], "SECRET_") {
			keys := strings.Replace(pair[0], "SECRET_", "secrets.", -1)
			keys = strings.Replace(keys, "_", "-", -1)
			newKey := strings.Trim(keys, " ")
			newValue := strings.Trim(pair[1], " ")

			jsonFlag := json.Valid([]byte(newValue))
			if !jsonFlag {
				newValue = strings.Replace(newValue, `\n`, "\n", -1)
			}

			viper.Set(newKey, newValue)
		}
	}
}
