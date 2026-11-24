package apiErrors

const (
	ErrorUnauthorized   = "unauthorized"
	ErrorInvalidRequest = "invalid_request"
	ErrorNotFound       = "not_found"
	ErrorForbidden      = "forbidden"
	ErrorAlreadyExists  = "already_exists"
	ErrorLimitReached   = "limit_reached"
)

const (
	MessageAuthorizationRequired = "Authorization header required"
	MessageAuthorizationInvalid  = "Invalid or expired token"
	MessageAdvertisementNotFound = "Advertisement not found"
)

// Сообщения валидации POST /advertisement. Константы с суффиксом Format подставляются через fmt.Sprintf
const (
	MessageTitleRequired        = "title is required"
	MessageTitleTooLongFormat   = "title must not exceed %d characters provided %d"
	MessageDescriptionRequired  = "description is required"
	MessageDescriptionLength    = "description must be between 1 and 500 characters"
	MessagePriceOutOfRange      = "price must be between 0 and 1000000 inclusive"
	MessagePriceNotNumberFormat = "cannot convert string value %s to float"
	MessageQuantityOutOfRange   = "quantity must be between 1 and 100 inclusive"
	MessageQuantityNotNumber    = "cannot convert string value to int"
	MessagePhotoRequired        = "at least 1 photo is required"
	MessageTooManyPhotos        = "maximum 3 photos allowed"
	// MessagePhotoTypeNotAllowedFormat номер фото в форме и его расширение в кавычках, например: photo 1 ... ".bmp"
	MessagePhotoTypeNotAllowedFormat = "photo %d: file type not allowed: extension %q"
)

// Сообщения ручек избранного. Константы с суффиксом Format подставляются через fmt.Sprintf
const (
	MessageFavoriteNotFound            = "Favorite entry not found"
	MessageInvalidUUIDLengthFormat     = "invalid UUID length: %d"
	MessageAlreadyInFavorites          = "Advertisement is already in favorites"
	MessageOwnAdvertisementToFavorites = "Cannot add own advertisement to favorites"
	MessageFavoritesLimitReachedFormat = "Favorites limit of %d reached"
)
