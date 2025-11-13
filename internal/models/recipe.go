package models

// Main recipe model
type Recipe struct {
	RecipeID     int                `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string             `gorm:"unique" json:"name"`
	Difficulty   int                `json:"difficulty"`
	Description  *string            `json:"description,omitempty"` //optional
	Ingredients  []RecipeIngredient `gorm:"foreignKey:RecipeID" json:"ingredients,omitempty"`
	Instructions []Instruction      `gorm:"foreignKey:RecipeID" json:"instructions,omitempty"`
	UserID       string             `gorm:"type:varchar(32);not null" json:"user_id"`
}
