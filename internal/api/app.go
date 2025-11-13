package api

import (
	"log"
	"recipe-api/internal/repository"
)

type App struct {
	Repo   *repository.App
	Logger *log.Logger
}
