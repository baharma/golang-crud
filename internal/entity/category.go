package entity

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"size:255;not null;uniqueIndex" json:"name" validate:"required"`
	Books []Book `gorm:"foreignKey:CategoryID" json:"books,omitempty"`
}

func (c *Category) TableName() string {
	return "categories"
}
