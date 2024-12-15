package models

// User struct represents a user in the database.
type User struct {
	ID    string `gorm:"primaryKey"`
	Name  string
	Email string
}

// Product struct represents a product in the database.
type Product struct {
	ID                      string   `json:"id"`
	UserID                  string   `json:"user_id"`
	ProductName             string   `json:"product_name"`
	ProductDescription      string   `json:"product_description"`
	ProductImages           []string `json:"product_images"`
	CompressedProductImages []string `json:"compressed_product_images"`
	ProductPrice            float64  `json:"product_price"`
}
