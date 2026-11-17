package assertions

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	titles := make([]string, 0, len(items))
	for _, item := range items {
		titles = append(titles, item.Title)
	}

	require.Containsf(t, titles, title, "Объявление с названием %q не найдено в выдаче", title)
	logCheck("объявление %q найдено в выдаче", title)
}

// AdvertisementMatchesRequest проверяет, что объявление содержит все параметры из запроса на создание:
// название, описание, цену, количество и число фотографий
func AdvertisementMatchesRequest(t *testing.T, actual models.Advertisement, expected models.CreateAdvertisementRequest) {
	t.Helper()

	expectedPrice, err := strconv.ParseFloat(expected.Price, 64)
	require.NoErrorf(t, err, "Цена в запросе должна быть числом: %q", expected.Price)
	expectedQuantity, err := strconv.Atoi(expected.Quantity)
	require.NoErrorf(t, err, "Количество в запросе должно быть числом: %q", expected.Quantity)

	require.Equal(t, expected.Title, actual.Title, "Название объявления не совпадает с переданным")
	require.Equal(t, expected.Description, actual.Description, "Описание объявления не совпадает с переданным")
	require.Equal(t, expectedPrice, actual.Price, "Цена объявления не совпадает с переданной")
	require.Equal(t, expectedQuantity, actual.Quantity, "Количество объявления не совпадает с переданным")
	require.Len(t, actual.Photos, len(expected.PhotoPaths), "Количество фотографий не совпадает с загруженным")
	logCheck("объявление %s содержит переданные параметры: %q, цена %s, количество %s, фото %d",
		actual.Id, actual.Title, expected.Price, expected.Quantity, len(actual.Photos))
}

// AdvertisementBelongsToUser проверяет, что объявление принадлежит пользователю
func AdvertisementBelongsToUser(t *testing.T, actual models.Advertisement, userId string) {
	t.Helper()

	require.Equal(t, userId, actual.UserId, "Объявление принадлежит другому пользователю")
	logCheck("объявление %s принадлежит пользователю %s", actual.Id, userId)
}

// PhotosMatch проверяет, что фотографии совпадают с ожидаемыми по id и ссылке на файл
func PhotosMatch(t *testing.T, actual, expected []models.Photo) {
	t.Helper()

	require.ElementsMatch(t, photoUrlsById(expected), photoUrlsById(actual), "Фотографии объявления не совпадают с загруженными")
	for _, photo := range actual {
		require.NotEmptyf(t, photo.Url, "У фотографии %s пустая ссылка на файл", photo.Id)
	}
	logCheck("найдено фотографий: %d, id и ссылки совпадают с загруженными", len(actual))
}

// photoUrlsById возвращает пары «id: ссылка» для сравнения фотографий без учёта порядка
func photoUrlsById(photos []models.Photo) []string {
	pairs := make([]string, 0, len(photos))
	for _, photo := range photos {
		pairs = append(pairs, photo.Id+": "+photo.Url)
	}

	return pairs
}
