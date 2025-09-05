package config

import (
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

func getEnv(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return v
}

// getEnvAny returns the first set environment value among keys, or defaultValue.
func getEnvAny(defaultValue string, keys ...string) string {
	for _, k := range keys {
		if v, ok := os.LookupEnv(k); ok {
			return v
		}
	}

	return defaultValue
}

func AppAddr() string {
	// Prefer explicit APP_ADDR (e.g., ":8080").
	// Fallback to PORT (common on PaaS like NeoShowcase/Heroku) if provided.
	if v, ok := os.LookupEnv("APP_ADDR"); ok && v != "" {
		return v
	}
	if p, ok := os.LookupEnv("PORT"); ok && p != "" {
		return ":" + p
	}
	return ":8080"
}

func MySQL() *mysql.Config {
	c := mysql.NewConfig()

	// Prefer NS_MARIADB_* (PaaS-provided) first, then DB_* as fallback.
	c.User = getEnvAny("root", "NS_MARIADB_USER", "DB_USER")
	c.Passwd = getEnvAny("pass", "NS_MARIADB_PASSWORD", "DB_PASS")
	c.Net = getEnv("DB_NET", "tcp")
	c.Addr = fmt.Sprintf(
		"%s:%s",
		getEnvAny("localhost", "NS_MARIADB_HOSTNAME", "DB_HOST"),
		getEnvAny("3306", "NS_MARIADB_PORT", "DB_PORT"),
	)
	c.DBName = getEnvAny("app", "NS_MARIADB_DATABASE", "DB_NAME")
	c.Collation = "utf8mb4_general_ci"
	c.AllowNativePasswords = true
	c.ParseTime = true // DATETIME型をtime.Timeに変換

	return c
}
