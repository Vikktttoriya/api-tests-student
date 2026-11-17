package models

type AddFavoriteRequest struct {
	AdvertisementId string `json:"advertisement_id"`
}

type BulkAddFavoritesRequest struct {
	AdvertisementIds []string `json:"advertisement_ids"`
}
