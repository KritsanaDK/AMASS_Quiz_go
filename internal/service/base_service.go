package service

import (
	"amass/internal/repository"
	"amass/internal/utils"
	"context"
)

type BaseService struct {
	Ctx   context.Context
	Debug bool
}

type AllRepository struct {
	IUserRepository repository.IUserRepository
}

type AllService struct {
	IUserService  IUserService
	IUtilsService utils.IUtilsService
}
