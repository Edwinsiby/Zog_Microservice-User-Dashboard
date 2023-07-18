package utils

import "github.com/spf13/viper"

type Config struct {
	DB   string `mapstructure:"DB"`
	DSN  string `mapstructure:"DSN"`
	KEY1 string `mapstructure:"KEY1"`
	KEY2 string `mapstructure:"KEY2"`
	KEY3 string `mapstructure:"KEY3"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
