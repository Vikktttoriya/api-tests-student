package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/favorites/models"
	commonModels "api-tests-template/internal/models"
)

// MaxFavoritesPerUser лимит избранного на одного пользователя (описан в swagger ручки POST /favorites)
const MaxFavoritesPerUser = 50

// GetFavorites возвращает первую страницу избранного пользователя и общее количество записей, ожидает 200
func GetFavorites(t *testing.T, token string) models.FavoritesList {
	return response.ParseWithStatus[models.FavoritesList](t, favorites.HttpGetFavorites(t, token), http.StatusOK)
}

// GetFavoritesWithError запрашивает избранное с ожидаемой ошибкой и возвращает тело ошибки
func GetFavoritesWithError(t *testing.T, token string, expectedStatusCode int) commonModels.ErrorResponse {
	return response.ParseWithStatus[commonModels.ErrorResponse](t, favorites.HttpGetFavorites(t, token), expectedStatusCode)
}

// AddFavorite добавляет объявление в избранное и возвращает созданную запись, ожидает 201
func AddFavorite(t *testing.T, token, advertisementId string) models.FavoriteItem {
	request := models.AddFavoriteRequest{AdvertisementId: advertisementId}
	return response.ParseWithStatus[models.FavoriteItem](t,
		favorites.HttpPostFavorite(t, token, request), http.StatusCreated)
}

// AddFavoriteWithError отправляет запрос на добавление в избранное, который должен завершиться ошибкой,
// и возвращает тело ошибки. Если объявление всё же добавилось, оно убирается из избранного после теста
func AddFavoriteWithError(t *testing.T, token, advertisementId string, expectedStatusCode int) commonModels.ErrorResponse {
	request := models.AddFavoriteRequest{AdvertisementId: advertisementId}
	rawResponse := favorites.HttpPostFavorite(t, token, request)
	if rawResponse.StatusCode == http.StatusCreated {
		response.ReadBodyWithStatus(t, rawResponse, http.StatusCreated)
		t.Cleanup(func() {
			RemoveFavoriteIfExists(t, token, advertisementId)
		})
		require.Failf(t, "Объявление добавлено в избранное, хотя ожидалась ошибка",
			"Ожидался статус %d для объявления %s", expectedStatusCode, advertisementId)
	}

	return response.ParseWithStatus[commonModels.ErrorResponse](t, rawResponse, expectedStatusCode)
}

// AddFavoritesBulk добавляет пачку объявлений в избранное и возвращает результат, ожидает 200
func AddFavoritesBulk(t *testing.T, token string, advertisementIds []string) models.BulkAddFavoritesResult {
	request := models.BulkAddFavoritesRequest{AdvertisementIds: advertisementIds}
	return response.ParseWithStatus[models.BulkAddFavoritesResult](t,
		favorites.HttpPostFavoritesBulk(t, token, request), http.StatusOK)
}

// RemoveFavorite удаляет объявление из избранного, ожидает 204
func RemoveFavorite(t *testing.T, token, advertisementId string) {
	response.ReadBodyWithStatus(t, favorites.HttpDeleteFavorite(t, token, advertisementId), http.StatusNoContent)
}

// RemoveFavoriteWithError отправляет запрос на удаление из избранного с ожидаемой ошибкой и возвращает тело ошибки
func RemoveFavoriteWithError(t *testing.T, token, advertisementId string, expectedStatusCode int) commonModels.ErrorResponse {
	return response.ParseWithStatus[commonModels.ErrorResponse](t,
		favorites.HttpDeleteFavorite(t, token, advertisementId), expectedStatusCode)
}

// RemoveFavoriteIfExists убирает объявление из избранного, если оно там есть. Используется для очистки данных после теста:
// 404 означает, что объявления в избранном уже нет
func RemoveFavoriteIfExists(t *testing.T, token, advertisementId string) {
	rawResponse := favorites.HttpDeleteFavorite(t, token, advertisementId)
	defer rawResponse.Body.Close()

	require.Containsf(t, []int{http.StatusNoContent, http.StatusNotFound}, rawResponse.StatusCode,
		"Не удалось убрать объявление %s из избранного после теста", advertisementId)
}
