package models

import advModels "api-tests-template/internal/managers/advertisements/models"

type FavoritesList struct {
	Items      []FavoriteItem `json:"items,omitempty"`
	Page       int            `json:"page,omitempty"`
	PageSize   int            `json:"page_size,omitempty"`
	Total      int            `json:"total,omitempty"`
	TotalPages int            `json:"total_pages,omitempty"`
}

type FavoriteItem struct {
	FavoriteId    string                  `json:"favorite_id,omitempty"`
	AddedAt       string                  `json:"added_at,omitempty"`
	Advertisement advModels.Advertisement `json:"advertisement,omitempty"`
}

type BulkAddFavoritesResult struct {
	Added        []string `json:"added,omitempty"`
	Skipped      []string `json:"skipped,omitempty"`
	NotFound     []string `json:"not_found,omitempty"`
	LimitReached bool     `json:"limit_reached,omitempty"`
}
