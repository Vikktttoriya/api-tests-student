package advertisementCreate

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/apiErrors"
	"api-tests-template/internal/fixtures"
	"api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

const (
	maxTitleLength        = 50
	swaggerMaxTitleLength = 255
	maxDescriptionLength  = 500
	maxPrice              = 1000000
	maxQuantity           = 100
	notNumberValue        = "abc"

	photo1Path           = "testdata/photo_1.png"
	photo2Path           = "testdata/photo_2.png"
	photo3Path           = "testdata/photo_3.png"
	photo4Path           = "testdata/photo_4.png"
	unsupportedPhotoPath = "testdata/unsupported.bmp"
	notPhotoPath         = "testdata/unsupported.txt"
	// invalidPhotoNumber номер некорректного файла в форме: он передаётся единственным
	invalidPhotoNumber = 1
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	s.loginData = base.LoginTestUser(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Positive: объявление со всеми параметрами доступно по id, находится поиском по названию, его фото загружены
func (s *TestSuite) TestCreateAdvertisementWithAllParams() {
	t := s.T()

	request := fixtures.AdvertisementRequest(photo1Path, photo2Path, photo3Path)
	base.Precondition(fmt.Sprintf("Готовим запрос со всеми параметрами и %d фотографиями", len(request.PhotoPaths)))

	base.Step("POST /advertisement: %q", request.Title)
	created := fixtures.CreateAdvertisement(t, s.loginData.Token, request)

	s.Run("Ответ на создание содержит переданные параметры", func() {
		assertions.AdvertisementMatchesRequest(s.T(), created, request)
		assertions.AdvertisementBelongsToUser(s.T(), created, s.loginData.User.Id)
	})

	s.Run("GET /advertisement возвращает те же параметры", func() {
		base.Step("GET /advertisement?id=%s", created.Id)
		advertisement := advertisements.GetAdvertisement(s.T(), created.Id)
		assertions.AdvertisementMatchesRequest(s.T(), advertisement, request)
		assertions.AdvertisementBelongsToUser(s.T(), advertisement, s.loginData.User.Id)
		assertions.PhotosMatch(s.T(), advertisement.Photos, created.Photos)
	})

	s.Run("Объявление находится поиском по полному названию", func() {
		base.Step("GET /advertisements?search=%s", request.Title)
		list := advertisements.SearchAdvertisementsUntilFound(s.T(), s.loginData.Token, request.Title)
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, request.Title)
	})

	s.Run("Фотографии доступны через GET /advertisements/{id}/photos", func() {
		base.Step("GET /advertisements/%s/photos", created.Id)
		photos := advertisements.GetAdvertisementPhotos(s.T(), s.loginData.Token, created.Id)
		assertions.PhotosMatch(s.T(), photos, created.Photos)
	})
}

// Negative: создание объявления без обязательных параметров или с некорректными параметрами возвращает 400.
// Каждый кейс портит одно поле валидного запроса
func (s *TestSuite) TestCreateAdvertisementValidation() {
	cases := []struct {
		name            string
		modify          func(request *advModels.CreateAdvertisementRequest)
		expectedMessage string
	}{
		{"пустое название",
			func(r *advModels.CreateAdvertisementRequest) { r.Title = "" },
			apiErrors.MessageTitleRequired},
		{"название длиннее 50 символов",
			func(r *advModels.CreateAdvertisementRequest) { r.Title = strings.Repeat("a", maxTitleLength+1) },
			fmt.Sprintf(apiErrors.MessageTitleTooLongFormat, maxTitleLength, maxTitleLength+1)},
		{"название длиннее 255 символов",
			func(r *advModels.CreateAdvertisementRequest) { r.Title = strings.Repeat("a", swaggerMaxTitleLength+1) },
			fmt.Sprintf(apiErrors.MessageTitleTooLongFormat, maxTitleLength, swaggerMaxTitleLength+1)},
		{"пустое описание",
			func(r *advModels.CreateAdvertisementRequest) { r.Description = "" },
			apiErrors.MessageDescriptionRequired},
		{"описание длиннее 500 символов",
			func(r *advModels.CreateAdvertisementRequest) {
				r.Description = strings.Repeat("a", maxDescriptionLength+1)
			},
			apiErrors.MessageDescriptionLength},
		{"отрицательная цена",
			func(r *advModels.CreateAdvertisementRequest) { r.Price = strconv.Itoa(-1) },
			apiErrors.MessagePriceOutOfRange},
		{"цена больше 1 000 000",
			func(r *advModels.CreateAdvertisementRequest) { r.Price = strconv.Itoa(maxPrice + 1) },
			apiErrors.MessagePriceOutOfRange},
		{"цена не число",
			func(r *advModels.CreateAdvertisementRequest) { r.Price = notNumberValue },
			fmt.Sprintf(apiErrors.MessagePriceNotNumberFormat, notNumberValue)},
		{"отрицательное количество",
			func(r *advModels.CreateAdvertisementRequest) { r.Quantity = strconv.Itoa(-1) },
			apiErrors.MessageQuantityOutOfRange},
		{"нулевое количество",
			func(r *advModels.CreateAdvertisementRequest) { r.Quantity = strconv.Itoa(0) },
			apiErrors.MessageQuantityOutOfRange},
		{"количество больше 100",
			func(r *advModels.CreateAdvertisementRequest) { r.Quantity = strconv.Itoa(maxQuantity + 1) },
			apiErrors.MessageQuantityOutOfRange},
		{"количество не число",
			func(r *advModels.CreateAdvertisementRequest) { r.Quantity = notNumberValue },
			apiErrors.MessageQuantityNotNumber},
		{"без фотографий",
			func(r *advModels.CreateAdvertisementRequest) { r.PhotoPaths = nil },
			apiErrors.MessagePhotoRequired},
		{"больше 3 фотографий",
			func(r *advModels.CreateAdvertisementRequest) {
				r.PhotoPaths = []string{photo1Path, photo2Path, photo3Path, photo4Path}
			},
			apiErrors.MessageTooManyPhotos},
		{"фото неподдерживаемого формата (BMP)",
			func(r *advModels.CreateAdvertisementRequest) { r.PhotoPaths = []string{unsupportedPhotoPath} },
			fmt.Sprintf(apiErrors.MessagePhotoTypeNotAllowedFormat, invalidPhotoNumber, filepath.Ext(unsupportedPhotoPath))},
		{"файл не является изображением (TXT)",
			func(r *advModels.CreateAdvertisementRequest) { r.PhotoPaths = []string{notPhotoPath} },
			fmt.Sprintf(apiErrors.MessagePhotoTypeNotAllowedFormat, invalidPhotoNumber, filepath.Ext(notPhotoPath))},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			request := fixtures.AdvertisementRequest(photo1Path)
			tc.modify(&request)

			base.Step("POST /advertisement: %s, ожидаем 400", tc.name)
			errorResponse := advertisements.CreateAdvertisementWithError(s.T(), s.loginData.Token, request, http.StatusBadRequest)
			assertions.ErrorResponseEquals(s.T(), errorResponse, apiErrors.ErrorInvalidRequest, tc.expectedMessage)
		})
	}
}

// Negative: создание объявления без авторизации возвращает 401
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	request := fixtures.AdvertisementRequest(photo1Path)

	for _, tc := range fixtures.UnauthorizedCases() {
		s.Run(tc.Name, func() {
			base.Step("POST /advertisement %s, ожидаем 401", tc.Name)
			errorResponse := advertisements.CreateAdvertisementWithError(s.T(), tc.Token, request, http.StatusUnauthorized)
			assertions.ErrorResponseEquals(s.T(), errorResponse, apiErrors.ErrorUnauthorized, tc.ExpectedMessage)
		})
	}
}
