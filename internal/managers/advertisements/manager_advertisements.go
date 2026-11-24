package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements/models"
	commonModels "api-tests-template/internal/models"
	"api-tests-template/internal/utils"
)

const (
	searchTimeout      = 5 * time.Second
	searchPollInterval = 500 * time.Millisecond
)

func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	advertisementsRawResponse := advertisements.HttpGetAdvertisements(t, token, text)
	defer advertisementsRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(advertisementsRawResponse.Body)
	require.NoError(t, errResponse)

	advertisementsList := models.AdvertisementsList{}
	err := json.Unmarshal(respArray, &advertisementsList)
	require.NoError(t, err)

	return advertisementsList
}

// SearchAdvertisementsUntilFound повторяет поиск, пока выдача не станет непустой или не истечёт таймаут,
// и возвращает последнюю выдачу. Поисковая выдача может обновиться не сразу после создания объявления
func SearchAdvertisementsUntilFound(t *testing.T, token, text string) models.AdvertisementsList {
	var list models.AdvertisementsList
	utils.WaitUntil(searchTimeout, searchPollInterval, func() bool {
		list = GetAdvertisementsList(t, token, text)
		return len(list.Items) > 0
	})

	return list
}

// CreateAdvertisement создаёт объявление и возвращает его, ожидает 201
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	return response.ParseWithStatus[models.Advertisement](t,
		advertisements.HttpPostAdvertisement(t, token, request), http.StatusCreated)
}

// CreateAdvertisementWithError отправляет запрос на создание объявления, который должен завершиться ошибкой,
// и возвращает тело ошибки. Если объявление всё же создалось, оно удаляется после теста, чтобы не копить мусор
func CreateAdvertisementWithError(t *testing.T, token string, request models.CreateAdvertisementRequest,
	expectedStatusCode int) commonModels.ErrorResponse {
	rawResponse := advertisements.HttpPostAdvertisement(t, token, request)
	if rawResponse.StatusCode == http.StatusCreated {
		created := response.ParseWithStatus[models.Advertisement](t, rawResponse, http.StatusCreated)
		t.Cleanup(func() {
			DeleteAdvertisementIfExists(t, token, created.Id)
		})
		require.Failf(t, "Объявление создано, хотя ожидалась ошибка",
			"Ожидался статус %d, создано объявление %s (будет удалено после теста)", expectedStatusCode, created.Id)
	}

	return response.ParseWithStatus[commonModels.ErrorResponse](t, rawResponse, expectedStatusCode)
}

// GetAdvertisement возвращает объявление по id, ожидает 200
func GetAdvertisement(t *testing.T, advertisementId string) models.Advertisement {
	return response.ParseWithStatus[models.Advertisement](t,
		advertisements.HttpGetAdvertisement(t, advertisementId), http.StatusOK)
}

// DeleteAdvertisement удаляет объявление, ожидает 204
func DeleteAdvertisement(t *testing.T, token, advertisementId string) {
	response.ReadBodyWithStatus(t, advertisements.HttpDeleteAdvertisement(t, token, advertisementId), http.StatusNoContent)
}

// DeleteAdvertisementIfExists удаляет объявление, если оно ещё существует. Используется для очистки данных после теста.
// Существование проверяется через GET: повторный DELETE уже удалённого объявления возвращает 403, а не 404
func DeleteAdvertisementIfExists(t *testing.T, token, advertisementId string) {
	getResponse := advertisements.HttpGetAdvertisement(t, advertisementId)
	defer getResponse.Body.Close()

	if getResponse.StatusCode == http.StatusNotFound {
		return
	}
	DeleteAdvertisement(t, token, advertisementId)
}

// GetAdvertisementPhotos возвращает фотографии объявления, ожидает 200
func GetAdvertisementPhotos(t *testing.T, token, advertisementId string) []models.Photo {
	return response.ParseWithStatus[[]models.Photo](t,
		advertisements.HttpGetAdvertisementPhotos(t, token, advertisementId), http.StatusOK)
}
