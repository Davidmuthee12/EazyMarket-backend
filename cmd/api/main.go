package main

import (
	"time"

	_ "github.com/Davidmuthee12/eazymarket/docs"
	"github.com/Davidmuthee12/eazymarket/internals/auth"
	"github.com/Davidmuthee12/eazymarket/internals/db"
	"github.com/Davidmuthee12/eazymarket/internals/env"
	"github.com/Davidmuthee12/eazymarket/internals/mailer"
	ratelimiter "github.com/Davidmuthee12/eazymarket/internals/rateLimiter"
	"github.com/Davidmuthee12/eazymarket/internals/store"
	cache "github.com/Davidmuthee12/eazymarket/internals/store/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const version = "1.8.0"

//	@title	EazyMarket APP API

//	@description	This API for the Eazymarket app.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
	cfg := config{
		addr:   env.GetString("ADDR", "8000"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:          env.GetString("DB_ADDR", "postgres://adminpassword@localhost/eazymarket?sslmode=disable"),
			maxOpenConn:   env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIddleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIddleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		auth: authConfig{
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3, //3 days
				iss:    "eazymarket",
			},
		},
		mail: mailConfig{
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			fromEmail: env.GetString("FROM_EMAIL", ""),
			exp:       time.Hour * 24 * 3,
		},
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5174"),
		rootDomain:  env.GetString("ROOT_DOMAIN", "eazymarket.com"),
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PASSWORD", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		rateLimiter: ratelimiter.Config{
			RequestPerTimeFrame: env.GetInt("RATELIMITER_REQUEST_COUNT", 20),
			TimeFrame:           time.Second * 5,
			Enabled:             env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}

	// logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	logger.Infow("Starting EazyMarket API", "version", version)

	// DB
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConn,
		cfg.db.maxIddleConns,
		cfg.db.maxIddleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Database connection pool established")

	// Cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("Redis cache connection pool established")

		defer rdb.Close()
	}

	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	store := store.NewStorage(db)

	jwtAuthentticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	var emailClient mailer.Client
	if cfg.mail.sendGrid.apiKey != "" && cfg.mail.fromEmail != "" {
		emailClient = mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	} else {
		logger.Warn("Mailer not configured: missing SENDGRID_API_KEY or FROM_EMAIL, welcome emails will be skipped")
	}

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        emailClient,
		authenticator: jwtAuthentticator,
		ratelimiter:   rateLimiter,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
