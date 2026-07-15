package api

import (
	"log"
	"recipe-api/internal/middleware"
	"recipe-api/internal/repository"
)

type App struct {
	Repo        *repository.App
	Logger      *log.Logger
	RateLimiter *middleware.RateLimiter
}
