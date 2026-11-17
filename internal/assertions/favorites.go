package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites/models"
)

// ContainsFavorite проверяет, что объявление есть в списке избранного
func ContainsFavorite(t *testing.T, items []models.FavoriteItem, advertisementId string) {
	t.Helper()

	require.Containsf(t, favoriteAdvertisementIds(items), advertisementId,
		"Объявление %s не найдено в избранном", advertisementId)
	logCheck("объявление %s есть в избранном", advertisementId)
}

// NotContainsFavorite проверяет, что объявления нет в списке избранного
func NotContainsFavorite(t *testing.T, items []models.FavoriteItem, advertisementId string) {
	t.Helper()

	require.NotContainsf(t, favoriteAdvertisementIds(items), advertisementId,
		"Объявление %s всё ещё есть в избранном", advertisementId)
	logCheck("объявления %s нет в избранном", advertisementId)
}

// favoriteAdvertisementIds возвращает id объявлений из списка избранного
func favoriteAdvertisementIds(items []models.FavoriteItem) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.Advertisement.Id)
	}

	return ids
}
