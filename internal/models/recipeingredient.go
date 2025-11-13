package models

// Linker model for Recipe, Ingredient and Unit
type RecipeIngredient struct {
	RecipeIngredientID int         `gorm:"primaryKey;autoIncrement" json:"id"`
	RecipeID           int         `gorm:"not null;index" json:"recipe_id"`
	IngredientID       int         `gorm:"not null;index" json:"ingredient_id"`
	UnitID             *int        `gorm:"index" json:"unit_id,omitempty"`            //optional
	Amount             *float32    `gorm:"type:numeric(4,2)" json:"amount,omitempty"` //optional
	Ingredient         *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient"`
	Unit               *Unit       `gorm:"foreignKey:UnitID;references:UnitID" json:"unit,omitempty"`
}
