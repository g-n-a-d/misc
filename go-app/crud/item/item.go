package item

type Item struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Price float32 `json:"price"`
	DiscountPercentage float32 `json:"discountPercentage"`
	Rating float32 `json:"rating"`
	Stock float32 `json:"stock"`
	Brand string `json:"brand"`
	Category string `json:"category"`
	Thumbnail string `json:"thumbnail"`
}