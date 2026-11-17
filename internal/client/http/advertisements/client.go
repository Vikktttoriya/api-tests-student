package advertisements

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/helpers/multipart"
	"api-tests-template/internal/managers/advertisements/models"
)

func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// HttpGetOldestAdvertisements выполняет GET-запрос для получения limit самых старых объявлений каталога
func HttpGetOldestAdvertisements(t *testing.T, limit int) *http.Response {
	return apiRunner.GetRunner().Create().
		Get(path.AdvertisementsPath).
		Query("sort_by", "created_at").
		Query("sort_order", "asc").
		Query("limit", strconv.Itoa(limit)).
		Expect(t).
		End().Response
}

// HttpPostAdvertisement выполняет POST-запрос (multipart/form-data) для создания объявления с фотографиями
func HttpPostAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) *http.Response {
	fields := map[string]string{
		"title":       request.Title,
		"description": request.Description,
	}
	if request.Price != "" {
		fields["price"] = request.Price
	}
	if request.Quantity != "" {
		fields["quantity"] = request.Quantity
	}
	body, contentType := multipart.BuildForm(t, fields, "photos", request.PhotoPaths)

	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		ContentType(contentType).
		Body(body).
		Expect(t).
		End().Response
}

// HttpGetAdvertisement выполняет GET-запрос для получения объявления по id
func HttpGetAdvertisement(t *testing.T, advertisementId string) *http.Response {
	return apiRunner.GetRunner().Create().
		Get(path.AdvertisementPath).
		Query("id", advertisementId).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement выполняет DELETE-запрос для удаления объявления по id
func HttpDeleteAdvertisement(t *testing.T, token, advertisementId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", advertisementId).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos выполняет GET-запрос для получения фотографий объявления
func HttpGetAdvertisementPhotos(t *testing.T, token, advertisementId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(fmt.Sprintf(path.AdvertisementPhotosPath, advertisementId)).
		Expect(t).
		End().Response
}
