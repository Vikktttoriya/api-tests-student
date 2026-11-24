package favorites

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/apiErrors"
	"api-tests-template/internal/fixtures"
	"api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	commonModels "api-tests-template/internal/models"

	base "api-tests-template/tests"
)

const (
	// nonExistentAdvertisementId валидный, но несуществующий UUID. Нулевой UUID не подходит:
	// сервер считает его пустым значением и отвечает 400 вместо 404
	nonExistentAdvertisementId = "00000000-0000-4000-8000-000000000001"
	invalidAdvertisementId     = "not-a-uuid"
	photoPath                  = "testdata/photo.png"
)

// TestSuite seller — пользователь из переменных окружения, владелец объявлений;
// buyer — новый пользователь, чьё избранное проверяем. Своё объявление в избранное добавить нельзя,
// поэтому для сценариев нужны два разных пользователя
type TestSuite struct {
	suite.Suite
	sellerData models.LoginOkResponse
	buyerData  models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	s.sellerData = base.LoginTestUser(s.T())

	base.Precondition("Регистрация нового покупателя с пустым избранным")
	s.buyerData = fixtures.NewUser(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Positive: добавленное объявление появляется в избранном, удалённое из избранного — пропадает
func (s *TestSuite) TestAddAndRemoveFavorite() {
	advertisement := s.newSellerAdvertisement()

	s.addToFavoritesAndCheck(advertisement.Id)

	s.Run("Покупатель удаляет объявление из избранного", func() {
		base.Step("DELETE /favorites/%s", advertisement.Id)
		favorites.RemoveFavorite(s.T(), s.buyerData.Token, advertisement.Id)
	})

	s.Run("Объявление пропало из списка избранного", func() {
		base.Step("GET /favorites")
		list := favorites.GetFavorites(s.T(), s.buyerData.Token)
		assertions.NotContainsFavorite(s.T(), list.Items, advertisement.Id)
	})
}

// Positive: объявление, удалённое продавцом, пропадает из избранного покупателя
func (s *TestSuite) TestFavoriteDisappearsAfterAdvertisementDeleted() {
	advertisement := s.newSellerAdvertisement()

	s.addToFavoritesAndCheck(advertisement.Id)

	s.Run("Продавец удаляет объявление", func() {
		base.Step("DELETE /advertisement?id=%s", advertisement.Id)
		advertisements.DeleteAdvertisement(s.T(), s.sellerData.Token, advertisement.Id)
	})

	s.Run("Объявление пропало из избранного покупателя", func() {
		base.Step("GET /favorites")
		list := favorites.GetFavorites(s.T(), s.buyerData.Token)
		assertions.NotContainsFavorite(s.T(), list.Items, advertisement.Id)
	})
}

// Negative: добавление и удаление из избранного несуществующего объявления, добавление объявления с некорректным id
func (s *TestSuite) TestFavoritesWithUnknownAdvertisement() {
	addFavorite := func(advertisementId string, expectedStatusCode int) commonModels.ErrorResponse {
		return favorites.AddFavoriteWithError(s.T(), s.buyerData.Token, advertisementId, expectedStatusCode)
	}
	removeFavorite := func(advertisementId string, expectedStatusCode int) commonModels.ErrorResponse {
		return favorites.RemoveFavoriteWithError(s.T(), s.buyerData.Token, advertisementId, expectedStatusCode)
	}

	cases := []struct {
		name               string
		action             func(advertisementId string, expectedStatusCode int) commonModels.ErrorResponse
		advertisementId    string
		expectedStatusCode int
		expectedError      string
		expectedMessage    string
	}{
		{"добавление несуществующего объявления", addFavorite, nonExistentAdvertisementId,
			http.StatusNotFound, apiErrors.ErrorNotFound, apiErrors.MessageAdvertisementNotFound},
		{"удаление несуществующего объявления", removeFavorite, nonExistentAdvertisementId,
			http.StatusNotFound, apiErrors.ErrorNotFound, apiErrors.MessageFavoriteNotFound},
		{"добавление объявления с некорректным id", addFavorite, invalidAdvertisementId,
			http.StatusBadRequest, apiErrors.ErrorInvalidRequest,
			fmt.Sprintf(apiErrors.MessageInvalidUUIDLengthFormat, len(invalidAdvertisementId))},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			base.Step("%s: ad_id=%s, ожидаем %d", tc.name, tc.advertisementId, tc.expectedStatusCode)
			errorResponse := tc.action(tc.advertisementId, tc.expectedStatusCode)
			assertions.ErrorResponseEquals(s.T(), errorResponse, tc.expectedError, tc.expectedMessage)
		})
	}
}

// Negative: ручки избранного без авторизации возвращают 401
func (s *TestSuite) TestFavoritesUnauthorized() {
	operations := []struct {
		name   string
		action func(token string) commonModels.ErrorResponse
	}{
		{"GET /favorites", func(token string) commonModels.ErrorResponse {
			return favorites.GetFavoritesWithError(s.T(), token, http.StatusUnauthorized)
		}},
		{"POST /favorites", func(token string) commonModels.ErrorResponse {
			return favorites.AddFavoriteWithError(s.T(), token, nonExistentAdvertisementId, http.StatusUnauthorized)
		}},
		{"DELETE /favorites/{ad_id}", func(token string) commonModels.ErrorResponse {
			return favorites.RemoveFavoriteWithError(s.T(), token, nonExistentAdvertisementId, http.StatusUnauthorized)
		}},
	}

	for _, operation := range operations {
		for _, tc := range fixtures.UnauthorizedCases() {
			s.Run(operation.name+" "+tc.Name, func() {
				base.Step("%s %s, ожидаем 401", operation.name, tc.Name)
				errorResponse := operation.action(tc.Token)
				assertions.ErrorResponseEquals(s.T(), errorResponse, apiErrors.ErrorUnauthorized, tc.ExpectedMessage)
			})
		}
	}
}

// Negative: повторное добавление того же объявления в избранное возвращает 409
func (s *TestSuite) TestAddFavoriteTwice() {
	advertisement := s.newSellerAdvertisement()

	s.Run("Покупатель добавляет объявление в избранное", func() {
		base.Step("POST /favorites: объявление %s", advertisement.Id)
		favorites.AddFavorite(s.T(), s.buyerData.Token, advertisement.Id)
	})

	s.Run("Повторное добавление того же объявления возвращает 409", func() {
		base.Step("POST /favorites повторно: объявление %s", advertisement.Id)
		errorResponse := favorites.AddFavoriteWithError(s.T(), s.buyerData.Token, advertisement.Id, http.StatusConflict)
		assertions.ErrorResponseEquals(s.T(), errorResponse, apiErrors.ErrorAlreadyExists, apiErrors.MessageAlreadyInFavorites)
	})
}

// Negative: добавление своего объявления в избранное возвращает 403
func (s *TestSuite) TestAddOwnAdvertisementToFavorites() {
	advertisement := s.newSellerAdvertisement()

	s.Run("Продавец добавляет своё объявление в избранное", func() {
		base.Step("POST /favorites от владельца: объявление %s", advertisement.Id)
		errorResponse := favorites.AddFavoriteWithError(s.T(), s.sellerData.Token, advertisement.Id, http.StatusForbidden)
		assertions.ErrorResponseEquals(s.T(), errorResponse, apiErrors.ErrorForbidden, apiErrors.MessageOwnAdvertisementToFavorites)
	})
}

// Negative: когда в избранном уже MaxFavoritesPerUser объявлений, добавить ещё одно нельзя
func (s *TestSuite) TestAddFavoriteOverLimit() {
	t := s.T()

	base.Precondition(fmt.Sprintf("Продавец создаёт %d объявлений для заполнения избранного до лимита", favorites.MaxFavoritesPerUser))
	advertisementIds := make([]string, 0, favorites.MaxFavoritesPerUser)
	for range favorites.MaxFavoritesPerUser {
		advertisementIds = append(advertisementIds, fixtures.NewAdvertisement(t, s.sellerData.Token, photoPath).Id)
	}

	base.Precondition("Покупатель добавляет созданные объявления в избранное одним запросом")
	bulkResult := favorites.AddFavoritesBulk(t, s.buyerData.Token, advertisementIds)

	extraAdvertisement := s.newSellerAdvertisement()

	s.Run("В избранном ровно лимит объявлений", func() {
		base.Step("GET /favorites")
		list := favorites.GetFavorites(s.T(), s.buyerData.Token)
		require.Equalf(s.T(), favorites.MaxFavoritesPerUser, list.Total,
			"Избранное не заполнено до лимита: добавлено %d, не найдено %v", len(bulkResult.Added), bulkResult.NotFound)
	})

	s.Run("Добавление сверх лимита возвращает 400", func() {
		base.Step("POST /favorites сверх лимита: объявление %s", extraAdvertisement.Id)
		errorResponse := favorites.AddFavoriteWithError(s.T(), s.buyerData.Token, extraAdvertisement.Id, http.StatusBadRequest)
		assertions.ErrorResponseEquals(s.T(), errorResponse, apiErrors.ErrorLimitReached,
			fmt.Sprintf(apiErrors.MessageFavoritesLimitReachedFormat, favorites.MaxFavoritesPerUser))
	})
}

// newSellerAdvertisement создаёт объявление продавца на время теста и убирает его из избранного покупателя после теста
func (s *TestSuite) newSellerAdvertisement() advModels.Advertisement {
	t := s.T()

	base.Precondition("Продавец создаёт объявление")
	advertisement := fixtures.NewAdvertisement(t, s.sellerData.Token, photoPath)
	t.Cleanup(func() {
		favorites.RemoveFavoriteIfExists(t, s.buyerData.Token, advertisement.Id)
	})

	return advertisement
}

// addToFavoritesAndCheck добавляет объявление в избранное покупателя и проверяет, что оно появилось в списке
func (s *TestSuite) addToFavoritesAndCheck(advertisementId string) {
	s.Run("Покупатель добавляет объявление в избранное", func() {
		base.Step("POST /favorites: объявление %s", advertisementId)
		favorite := favorites.AddFavorite(s.T(), s.buyerData.Token, advertisementId)
		require.Equal(s.T(), advertisementId, favorite.Advertisement.Id, "В ответе на добавление другое объявление")
	})

	s.Run("Объявление появилось в списке избранного", func() {
		base.Step("GET /favorites")
		list := favorites.GetFavorites(s.T(), s.buyerData.Token)
		assertions.ContainsFavorite(s.T(), list.Items, advertisementId)
	})
}
