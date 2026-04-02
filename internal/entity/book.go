package entity

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	ID         uint     `gorm:"primaryKey" json:"id"`
	Title      string   `gorm:"not null" json:"title" validate:"required"`
	Author     string   `gorm:"not null" json:"author" validate:"required"`
	CategoryID uint     `gorm:"not null" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (b *Book) TableName() string {
	return "books"
}
