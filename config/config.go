package config

import (
	"log/slog"
	"os"
)

type LocalStoreConfig struct {
	Directory string
}

func NewLocalStoreConfig() *LocalStoreConfig {
	return &LocalStoreConfig{
		Directory: os.Getenv("LOCAL_STORE_DIRECTORY"),
	}
}

// RMQConfig is the configuration for rabbitmq senders
type RMQConfig struct {
	Host     string
	User     string
	Password string
	Port     string
	L        *slog.Logger
}

func NewRMQConfig(L *slog.Logger) *RMQConfig {
	return &RMQConfig{
		Host:     os.Getenv("RABBITMQ_HOST"),
		User:     os.Getenv("RABBITMQ_USER"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
		Port:     os.Getenv("RABBITMQ_PORT"),
		L:        L,
	}
}

func THUMBNAIL_QUEUE() string {
	return "thumbnail_jobs"
}

func METADATA_QUEUE() string {
	return "metadata_jobs"
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	Port     string
	DBName   string
}

func NewPostgresConfig() *DBConfig {
	return &DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}
}

type ResolutionConfig struct {
	CollageWidth  int
	CollageHeight int
	SectionWidth  int
	SectionHeight int
	XSections     int
	YSections     int
}

func NewResolutionConfig() *ResolutionConfig {
	return &ResolutionConfig{
		CollageWidth:  8000,
		CollageHeight: 6000,
		SectionWidth:  80,
		SectionHeight: 60,
		XSections:     100,
		YSections:     100,
	}
}

type ServerConfig struct {
	Host string
	Port string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Host: os.Getenv("SERVER_HOST"),
		Port: os.Getenv("SERVER_PORT"),
	}
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func NewRabbitMQConfig() *RabbitMQConfig {
	return &RabbitMQConfig{
		Host:     os.Getenv("RABBITMQ_HOST"),
		Port:     os.Getenv("RABBITMQ_PORT"),
		User:     os.Getenv("RABBITMQ_USER"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
	}
}
