package config

import (
	"fmt"
	"os"
	"strings"

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

// ========== traQ OAuth ==========
func TraqOAuthClientID() string {
	return getEnv("TRAQ_OAUTH_CLIENT_ID", "")
}

func TraqOAuthClientSecret() string {
	return getEnv("TRAQ_OAUTH_CLIENT_SECRET", "")
}

func TraqOAuthRedirectURI() string {
	// 例: http://localhost:8080/api/auth/callback
	return getEnv("TRAQ_OAUTH_REDIRECT_URI", "")
}

// 開発/本番のベースURL
func ServerBaseURL() string {
	// 例: http://localhost:8080
	return getEnv("SERVER_BASE_URL", "")
}

func FrontendBaseURL() string {
	// 例: http://localhost:5173
	return getEnv("FRONTEND_BASE_URL", "")
}

// CookieSecure returns whether cookies should be set with Secure=true.
// Priority:
// 1) If COOKIE_SECURE env is set to one of: "1", "true", "TRUE", "yes", "on" -> true
//    or "0", "false", "FALSE", "no", "off" -> false
// 2) Otherwise, infer from ServerBaseURL scheme (https => true)
// 3) Default to false
func CookieSecure() bool {
	if v, ok := os.LookupEnv("COOKIE_SECURE"); ok {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "1", "true", "yes", "on":
			return true
		case "0", "false", "no", "off":
			return false
		}
	}
	if base := ServerBaseURL(); strings.HasPrefix(strings.ToLower(base), "https://") {
		return true
	}
	return false
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
