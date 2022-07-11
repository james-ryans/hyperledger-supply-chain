package service

import (
	"github.com/meneketehe/hehe/app/model"
	repository "github.com/meneketehe/hehe/app/repository/organization"
)

func GetMe() (*model.Organization, error) {
	return repository.GetMe("28899ce0-b79d-497b-ae78-fd3b896e0429")
}
