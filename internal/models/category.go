package models

type Category struct {
    BaseModel
    Name        string    `gorm:"unique;not null;size:50" json:"name"`
    Slug        string    `gorm:"unique;not null;size:50" json:"slug"`
    Icon        string    `json:"icon"`
    Description string    `json:"description"`
    IsActive    bool      `gorm:"default:true" json:"is_active"`
    Products    []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}