package models

import (
    "time"
    "github.com/google/uuid"
)

type Order struct {
    BaseModel
    UserID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
    BusinessID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"business_id"`
    OrderNumber     string     `gorm:"unique;not null;size:50" json:"order_number"`
    Status          string     `gorm:"default:'pending';size:20" json:"status"`
    TotalAmount     float64    `gorm:"not null;type:decimal(10,2)" json:"total_amount"`
    DeliveryFee     float64    `gorm:"default:0;type:decimal(10,2)" json:"delivery_fee"`
    PaymentMethod   string     `gorm:"default:'mpesa';size:20" json:"payment_method"`
    PaymentStatus   string     `gorm:"default:'pending';size:20" json:"payment_status"`
    DeliveryAddress string     `gorm:"type:text;not null" json:"delivery_address"`
    MpesaReceipt    string     `gorm:"size:100" json:"mpesa_receipt"`
    MpesaRequestID  string     `gorm:"size:100" json:"mpesa_request_id"`
    MpesaResultCode string     `gorm:"size:10" json:"mpesa_result_code"`
    PaidAt          *time.Time `json:"paid_at"`
    CompletedAt     *time.Time `json:"completed_at"`
    CancelledAt     *time.Time `json:"cancelled_at"`
    CancellationReason string  `gorm:"type:text" json:"cancellation_reason"`
    Notes           string     `gorm:"type:text" json:"notes"`
    User            User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Business        Business   `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
    Items           []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// Order status constants
const (
    OrderStatusPending    = "pending"
    OrderStatusProcessing = "processing"
    OrderStatusShipped    = "shipped"
    OrderStatusCompleted  = "completed"
    OrderStatusCancelled  = "cancelled"
)

// Payment status constants
const (
    PaymentStatusPending   = "pending"
    PaymentStatusPaid      = "paid"
    PaymentStatusFailed    = "failed"
    PaymentStatusRefunded  = "refunded"
)

// Payment method constants
const (
    PaymentMethodMpesa = "mpesa"
    PaymentMethodCash  = "cash"
    PaymentMethodCard  = "card"
)

func (o *Order) IsValidStatus(status string) bool {
    validStatuses := []string{
        OrderStatusPending,
        OrderStatusProcessing,
        OrderStatusShipped,
        OrderStatusCompleted,
        OrderStatusCancelled,
    }
    for _, s := range validStatuses {
        if s == status {
            return true
        }
    }
    return false
}

func (o *Order) CanTransitionTo(newStatus string) bool {
    transitions := map[string][]string{
        OrderStatusPending:    {OrderStatusProcessing, OrderStatusCancelled},
        OrderStatusProcessing: {OrderStatusShipped, OrderStatusCancelled},
        OrderStatusShipped:    {OrderStatusCompleted, OrderStatusCancelled},
        OrderStatusCompleted:  {},
        OrderStatusCancelled:  {},
    }

    allowed, ok := transitions[o.Status]
    if !ok {
        return false
    }

    for _, s := range allowed {
        if s == newStatus {
            return true
        }
    }
    return false
}