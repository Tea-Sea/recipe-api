package models

// Unit model
type Unit struct {
	UnitID int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Label  string `gorm:"type:varchar(32);not null" json:"label"`
}
