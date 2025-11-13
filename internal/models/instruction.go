package models

// Model for Instructions / Method
type Instruction struct {
	InstructionID int     `gorm:"primaryKey;autoIncrement" json:"id"`
	RecipeID      int     `gorm:"not null;index" json:"recipe_id"`
	StepNumber    int     `gorm:"not null;check:step_number>0" json:"stepNumber"`
	StepText      string  `json:"stepText"`
	Duration      *int    `json:"stepTime,omitempty"` //optional
	Notes         *string `json:"notes,omitempty"`    //optional
}
