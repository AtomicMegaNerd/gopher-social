package main

import (
	"time"

	"github.com/atomicmeganerd/gopher-social/internal/env"
)

const version = "0.1.0"

type config struct {
	addr        string
	db          dbConfig
	cache       cacheConfig
	env         string
	version     string
	apiURL      string
	mail        mailConfig
	frontendURL string
	auth        authConfig
}

func NewConfig() config {
	return config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "http://localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			// postgres://user:password@host:port/dbname?sslmode=disable
			addr:         env.GetString("DATABASE_URL", ""), // no default, must be set
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 20),
			minIdleConns: env.GetInt("DB_MIN_IDLE_CONNS", 5),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		cache: cacheConfig{
			addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PASSWORD", ""),
			db:       env.GetInt("REDIS_DB", 0),
			enabled:  env.GetBool("REDIS_ENABLE", false),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
		version: env.GetString("VERSION", "0.1.1"),
		auth: authConfig{
			basic: basicAuthConfig{
				username: env.GetString("BASIC_USERNAME", ""),
				password: env.GetString("BASIC_PASSWORD", ""),
			},
			jwtToken: jwtTokenConfig{
				secret:    env.GetString("JWT_SECRET", ""),
				tokenHost: env.GetString("JWT_TOKEN_HOST", ""),
				expiry:    time.Hour * 24 * 3, // 3 days
			},
		},
	}
}

type authConfig struct {
	basic    basicAuthConfig
	jwtToken jwtTokenConfig
}

// NOTE: We are only using basic auth here as part of the course to learn how to set that
// up with Go + chi. Obviously in most cases this would not be a best practice.
type basicAuthConfig struct {
	username string
	password string // WARNING: Sensitive secret, do not expose
}

type jwtTokenConfig struct {
	secret    string // WARNING: Sensitive secret, do not expose
	tokenHost string
	expiry    time.Duration
}

type mailConfig struct {
	sendGrid  sendGridConfig
	fromEmail string
	exp       time.Duration
}

type sendGridConfig struct {
	apiKey string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	minIdleConns int
	maxIdleTime  string
}

type cacheConfig struct {
	addr     string
	password string // WARNING: Sensitive secret, do not expose
	db       int
	enabled  bool
}
