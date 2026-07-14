package repository

import (
    "errors"
    "time"
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type OrderRepository struct {
    db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
    return &OrderRepository{db: db}
}

// CreateOrder creates a new order
func (r *OrderRepository) CreateOrder(order *models.Order) error {
    return r.db.Create(order).Error
}

// CreateOrderItems creates multiple order items
func (r *OrderRepository) CreateOrderItems(items []models.OrderItem) error {
    return r.db.Create(&items).Error
}

// GetOrderByID gets an order by ID
func (r *OrderRepository) GetOrderByID(id uuid.UUID) (*models.Order, error) {
    var order models.Order
    err := r.db.Preload("Items").Preload("Items.Product").
        Preload("User").Preload("Business").
        Where("id = ?", id).
        First(&order).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &order, nil
}

// GetOrderByOrderNumber gets an order by order number
func (r *OrderRepository) GetOrderByOrderNumber(orderNumber string) (*models.Order, error) {
    var order models.Order
    err := r.db.Preload("Items").Preload("User").Preload("Business").
        Where("order_number = ?", orderNumber).
        First(&order).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &order, nil
}

// GetOrdersByUserID gets all orders for a user
func (r *OrderRepository) GetOrdersByUserID(userID uuid.UUID, status string, limit, offset int) ([]models.Order, int64, error) {
    var orders []models.Order
    var total int64

    db := r.db.Model(&models.Order{}).
        Preload("Items").
        Preload("Business").
        Where("user_id = ?", userID)

    if status != "" && status != "all" {
        db = db.Where("status = ?", status)
    }

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    err := db.Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&orders).Error

    return orders, total, err
}

// GetOrdersByBusinessID gets all orders for a business
func (r *OrderRepository) GetOrdersByBusinessID(businessID uuid.UUID, status string, limit, offset int) ([]models.Order, int64, error) {
    var orders []models.Order
    var total int64

    db := r.db.Model(&models.Order{}).
        Preload("Items").
        Preload("User").
        Where("business_id = ?", businessID)

    if status != "" && status != "all" {
        db = db.Where("status = ?", status)
    }

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    err := db.Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&orders).Error

    return orders, total, err
}

// UpdateOrderStatus updates an order's status
func (r *OrderRepository) UpdateOrderStatus(id uuid.UUID, status string) error {
    updates := map[string]interface{}{
        "status": status,
    }

    // Set completion timestamps based on status
    if status == models.OrderStatusCompleted {
        now := time.Now()
        updates["completed_at"] = &now
    }
    if status == models.OrderStatusCancelled {
        now := time.Now()
        updates["cancelled_at"] = &now
    }

    return r.db.Model(&models.Order{}).
        Where("id = ?", id).
        Updates(updates).Error
}

// UpdatePaymentStatus updates an order's payment status
func (r *OrderRepository) UpdatePaymentStatus(id uuid.UUID, paymentStatus string, mpesaReceipt, mpesaRequestID, mpesaResultCode string) error {
    updates := map[string]interface{}{
        "payment_status": paymentStatus,
    }

    if mpesaReceipt != "" {
        updates["mpesa_receipt"] = mpesaReceipt
    }
    if mpesaRequestID != "" {
        updates["mpesa_request_id"] = mpesaRequestID
    }
    if mpesaResultCode != "" {
        updates["mpesa_result_code"] = mpesaResultCode
    }

    if paymentStatus == models.PaymentStatusPaid {
        now := time.Now()
        updates["paid_at"] = &now
    }

    return r.db.Model(&models.Order{}).
        Where("id = ?", id).
        Updates(updates).Error
}

// UpdateOrder updates an order
func (r *OrderRepository) UpdateOrder(order *models.Order) error {
    return r.db.Save(order).Error
}

// GetOrderStats gets order statistics for a business
func (r *OrderRepository) GetOrderStats(businessID uuid.UUID) (map[string]interface{}, error) {
    var totalOrders int64
    var pendingOrders int64
    var processingOrders int64
    var completedOrders int64
    var cancelledOrders int64
    var totalRevenue float64

    // Total orders
    if err := r.db.Model(&models.Order{}).
        Where("business_id = ?", businessID).
        Count(&totalOrders).Error; err != nil {
        return nil, err
    }

    // Pending orders
    if err := r.db.Model(&models.Order{}).
        Where("business_id = ? AND status = ?", businessID, models.OrderStatusPending).
        Count(&pendingOrders).Error; err != nil {
        return nil, err
    }

    // Processing orders
    if err := r.db.Model(&models.Order{}).
        Where("business_id = ? AND status = ?", businessID, models.OrderStatusProcessing).
        Count(&processingOrders).Error; err != nil {
        return nil, err
    }

    // Completed orders
    if err := r.db.Model(&models.Order{}).
        Where("business_id = ? AND status = ?", businessID, models.OrderStatusCompleted).
        Count(&completedOrders).Error; err != nil {
        return nil, err
    }

    // Cancelled orders
    if err := r.db.Model(&models.Order{}).
        Where("business_id = ? AND status = ?", businessID, models.OrderStatusCancelled).
        Count(&cancelledOrders).Error; err != nil {
        return nil, err
    }

    // Total revenue (completed orders only)
    if err := r.db.Model(&models.Order{}).
        Where("business_id = ? AND status = ?", businessID, models.OrderStatusCompleted).
        Select("COALESCE(SUM(total_amount), 0)").
        Scan(&totalRevenue).Error; err != nil {
        return nil, err
    }

    stats := map[string]interface{}{
        "total_orders":     totalOrders,
        "pending_orders":   pendingOrders,
        "processing_orders": processingOrders,
        "completed_orders": completedOrders,
        "cancelled_orders": cancelledOrders,
        "total_revenue":    totalRevenue,
    }

    return stats, nil
}

// SearchOrders searches orders by customer name or order number
func (r *OrderRepository) SearchOrders(businessID uuid.UUID, query string, status string, limit, offset int) ([]models.Order, int64, error) {
    var orders []models.Order
    var total int64

    db := r.db.Model(&models.Order{}).
        Preload("Items").
        Preload("User").
        Where("business_id = ?", businessID)

    if query != "" {
        db = db.Joins("JOIN users ON orders.user_id = users.id").
            Where("orders.order_number ILIKE ? OR users.name ILIKE ?", "%"+query+"%", "%"+query+"%")
    }

    if status != "" && status != "all" {
        db = db.Where("orders.status = ?", status)
    }

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    err := db.Order("orders.created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&orders).Error

    return orders, total, err
}

// GetRecentOrders gets recent orders for a business
func (r *OrderRepository) GetRecentOrders(businessID uuid.UUID, limit int) ([]models.Order, error) {
    var orders []models.Order
    err := r.db.Preload("Items").Preload("User").
        Where("business_id = ?", businessID).
        Order("created_at DESC").
        Limit(limit).
        Find(&orders).Error
    return orders, err
}

// GetOrdersByDateRange gets orders within a date range
func (r *OrderRepository) GetOrdersByDateRange(businessID uuid.UUID, startDate, endDate time.Time) ([]models.Order, error) {
    var orders []models.Order
    err := r.db.Preload("Items").Preload("User").
        Where("business_id = ? AND created_at BETWEEN ? AND ?", businessID, startDate, endDate).
        Order("created_at DESC").
        Find(&orders).Error
    return orders, err
}