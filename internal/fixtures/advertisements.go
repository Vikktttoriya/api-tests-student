package fixtures

import (
	"fmt"
	"strconv"
	"testing"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

const (
	advertisementTitlePrefix = "Кот"
	advertisementDescription = "Хороший кот, милый, пушистый"
	advertisementPrice       = 15000
	advertisementQuantity    = 2
	uniqueTitleSuffixLength  = 8
)

// AdvertisementRequest возвращает валидный запрос на создание объявления с уникальным названием и переданными фото
func AdvertisementRequest(photoPaths ...string) models.CreateAdvertisementRequest {
	return models.CreateAdvertisementRequest{
		Title:       fmt.Sprintf("%s %s", advertisementTitlePrefix, utils.RandomString(uniqueTitleSuffixLength)),
		Description: advertisementDescription,
		Price:       strconv.Itoa(advertisementPrice),
		Quantity:    strconv.Itoa(advertisementQuantity),
		PhotoPaths:  photoPaths,
	}
}

// CreateAdvertisement создаёт объявление по запросу и удаляет его после теста
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	created := advertisements.CreateAdvertisement(t, token, request)
	t.Cleanup(func() {
		advertisements.DeleteAdvertisementIfExists(t, token, created.Id)
	})

	return created
}

// NewAdvertisement создаёт объявление с одной фотографией и удаляет его после теста
func NewAdvertisement(t *testing.T, token, photoPath string) models.Advertisement {
	return CreateAdvertisement(t, token, AdvertisementRequest(photoPath))
}
