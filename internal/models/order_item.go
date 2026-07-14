package models

import (
    "github.com/google/uuid"
)

type OrderItem struct {
    BaseModel
    OrderID     uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
    ProductID   uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
    ProductName string    `gorm:"not null;size:100" json:"product_name"`
    Price       float64   `gorm:"not null;type:decimal(10,2)" json:"price"`
    Quantity    int       `gorm:"not null" json:"quantity"`
    Subtotal    float64   `gorm:"not null;type:decimal(10,2)" json:"subtotal"`
    Product     Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
    Order       Order     `gorm:"foreignKey:OrderID" json:"-"`
}