package main

import (
	"fio/internal/config"
	"fio/internal/enrich"
	"fio/internal/enrich/api"
	"fio/internal/middleware"
	"fio/internal/person/delivery"
	"fio/internal/person/repo"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

const (
	LoggerPrefixName = "prefix"
)

func main() {

	// Настройка логгера
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05.999999999Z07:00",
		FullTimestamp:   true,
	})

	// Настройка конфига
	cfg := config.Config{}

	// Если не смогли прочитать файл конфига - прочитаем environment
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		err = cleanenv.ReadEnv(&cfg)
		if err != nil {
			log.Fatal(err)
		}
	}

	db, err := sqlx.Connect("postgres", cfg.DbSource)
	if err != nil {
		logger.Fatal(err)
	}

	personRepo := &repo.PostgresRepo{
		Db:              db,
		PaginationLimit: cfg.PaginationLimit,
		Logger:          logger.WithField(LoggerPrefixName, "postgres"),
	}

	apiService := &api.Service{
		Client: http.DefaultClient,
	}

	enrichService := &enrich.Service{
		Logger:             logger.WithField(LoggerPrefixName, "enrich"),
		NationalityService: apiService,
		AgeService:         apiService,
		GenderService:      apiService,
	}

	personHandler := delivery.PersonHandler{
		Logger:        logger.WithField(LoggerPrefixName, "personHandler"),
		EnrichService: enrichService,
		PersonRepo:    personRepo,

		// Разрешенные фильтры, которые принимает хэндлер
		WhiteListFilters: []string{"age", "name", "gender", "nationality"},
	}

	// Настроим веб сервер
	gin.SetMode(cfg.GinMode)
	r := gin.New()

	r.Use(
		middleware.LoggerMiddleware(logger.WithField(LoggerPrefixName, "gin")),
		gin.Recovery(),
	)

	r.GET("/person", personHandler.GetListHandler)
	r.POST("/person", personHandler.CreateHandler)
	r.PUT("/person/:id", personHandler.UpdateHandler)
	r.DELETE("/person/:id", personHandler.DeleteHandler)

	logger.Infof("server is running. Port: %s", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}

}
