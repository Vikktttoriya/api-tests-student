package models

type AdvertisementsList struct {
	Items      []Advertisement `json:"items,omitempty"`
	Total      int             `json:"total,omitempty"`
	HasNext    bool            `json:"has_next,omitempty"`
	NextCursor string          `json:"next_cursor,omitempty"`
}

type Advertisement struct {
	Id          string  `json:"id,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Quantity    int     `json:"quantity,omitempty"`
	Photos      []Photo `json:"photos,omitempty"`
	UserId      string  `json:"user_id,omitempty"`
	User        User    `json:"user,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}

type Photo struct {
	Id              string `json:"id,omitempty"`
	AdvertisementId string `json:"advertisement_id,omitempty"`
	Url             string `json:"url,omitempty"`
	SortOrder       int    `json:"sort_order,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
}

type User struct {
	Id        string `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	BirthDate string `json:"birth_date,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}
