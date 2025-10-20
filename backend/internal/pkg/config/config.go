package config

import (
	"context"
	"fmt"

	"github.com/ai-companion/backend/internal/pkg/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	//Database DatabaseConfig `mapstructure:"database"`
	//Redis RedisConfig `mapstructure:"redis"`
	LLM LLMConfig `mapstructure:"llm"`
	//TTS      TTSConfig      `mapstructure:"tts"`
	//ASR      ASRConfig      `mapstructure:"asr"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LLMConfig struct {
	Provider string `mapstructure:"provider"`
	Model    string `mapstructure:"model"`
	BaseUrl  string `mapstructure:"baseUrl"`
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)

	// 环境变量覆盖
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Error("Config file not found, using defaults: %v", err)
		fmt.Println("Config file not found")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		logger.Errorf("Unable to decode config: %v", err)
		fmt.Println("Unable to decode config")
	}

	return &config
}

/*
GetString 从配置文件中获取字符串值
  - 参数:
    @param configPath: 配置项的路径
  - 返回:
    @param string: 配置项的值
*/
func GetString(configPath string) string {
	configValue := viper.GetString(configPath)
	if configValue == "" {
		if !viper.IsSet(configPath) {
			logger.Error(context.Background(), "config not found ", configPath)
		}
	}
	return configValue
}
