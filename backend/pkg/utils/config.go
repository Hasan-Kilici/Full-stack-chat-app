package utils

import (
	"github.com/go-ini/ini"
	"log"
)

// ConfigList represents the configuration settings.
type ConfigList struct {
	Port     			string
	Database 			string
	RedisAddr 			string
	SMTPHost 			string
	SMTPPort 			string
	SMTPPassword 		string
	SMTPFrom 			string
}

var Config ConfigList

// LoadConfig loads the configuration from the specified file path.
func LoadConfig(path string) ConfigList {
	cfg, err := ini.Load(path)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	
	Config := ConfigList{
		Port:     			cfg.Section("server").Key("port").MustString("8080"),
		Database: 			cfg.Section("db").Key("dsn").String(),
		RedisAddr: 			cfg.Section("db").Key("redisAddr").String(),
		SMTPHost: 			cfg.Section("smtp").Key("host").String(),
		SMTPPort: 			cfg.Section("smtp").Key("port").String(),
		SMTPPassword: 		cfg.Section("smtp").Key("password").String(),
		SMTPFrom: 			cfg.Section("smtp").Key("from").String(),
	}
	
	return Config
}