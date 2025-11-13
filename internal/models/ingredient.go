package models

// Single Ingredient
type Ingredient struct {
	IngredientID int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Label        string `gorm:"type:varchar(32);not null" json:"label"`
}
