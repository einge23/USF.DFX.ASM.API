package services

import (
	"gin-api/models"
	"gin-api/util"
)

func GetSettings() (models.Settings, error) {
	return util.Settings, nil
}