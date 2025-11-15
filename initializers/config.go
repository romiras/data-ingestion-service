package initializers

import "github.com/spf13/viper"

func NewConfig(env string) *viper.Viper {
	config := viper.New()
	config.AddConfigPath("./config")
	config.SetConfigName(env)
	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}
	config.AutomaticEnv()
	return config
}
