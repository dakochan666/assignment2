package controllers

import (
	"assignment2/helpers"
	"assignment2/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	db *gorm.DB
}

type ReqOrder struct {
	CustomerName string    `json:"customer_name"`
	Items        []ReqItem `json:"items"`
}

type ResOrder struct {
	ID           uint      `json:"id"`
	CustomerName string    `json:"customer_name"`
	OrderedAt    time.Time `json:"ordered_at"`
	Items        []ResItem `json:"items,omitempty"`
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{
		db: db,
	}
}

func (o *OrderController) CreateOrder(ctx *gin.Context) {
	var newReqOrder ReqOrder

	err := ctx.ShouldBindJSON(&newReqOrder)
	if err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	newOrder := models.Order{
		CustomerName: newReqOrder.CustomerName,
	}

	err = o.db.Create(&newOrder).Error
	if err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	var newResItem []ResItem

	for _, item := range newReqOrder.Items {
		newItem := models.Item{
			ItemCode:    item.ItemCode,
			Description: item.Description,
			Quantity:    item.Quantity,
			OrderID:     newOrder.ID,
		}
		err = o.db.Create(&newItem).Error
		if err != nil {
			helpers.BadRequestResponse(ctx, err.Error())
			return
		}
		newResItem = append(newResItem, ResItem{
			ID:          newItem.ID,
			ItemCode:    newItem.ItemCode,
			Description: newItem.Description,
			Quantity:    newItem.Quantity,
		})
	}

	newResOrder := ResOrder{
		ID:           newOrder.ID,
		CustomerName: newOrder.CustomerName,
		OrderedAt:    *newOrder.OrderedAt,
		Items:        newResItem,
	}

	helpers.WriteJsonResponse(ctx, http.StatusCreated, gin.H{
		"success": true,
		"data":    newResOrder,
	})
}

func (o *OrderController) GetOrders(ctx *gin.Context) {
	limit := ctx.Query("limit")
	limitInt := 10

	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err == nil {
			limitInt = l
		}
	}

	var orders []models.Order
	var total int64

	err := o.db.Limit(limitInt).Preload("Items").Find(&orders).Count(&total).Error
	if err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
		"query": map[string]interface{}{
			"limit": limitInt,
			"total": total,
		},
	})
}

func (o *OrderController) UpdateOrder(ctx *gin.Context) {
	orderId := ctx.Param("orderId")
	var newReqOrder ReqOrder
	var order models.Order

	err := ctx.ShouldBindJSON(&newReqOrder)
	if err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	newOrder := models.Order{
		CustomerName: newReqOrder.CustomerName,
	}

	err = o.db.Preload("Items").First(&order, orderId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "Order data not found")
			return
		}
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	err = o.db.Model(&order).Updates(newOrder).Error
	if err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	for _, item := range newReqOrder.Items {
		newItem := models.Item{
			ItemCode:    item.ItemCode,
			Description: item.Description,
			Quantity:    item.Quantity,
		}
		err = o.db.First(&newItem, "order_id = ?", orderId).Error
		if err != nil {
			helpers.BadRequestResponse(ctx, err.Error())
			return
		}

		err = o.db.Model(&newItem).Updates(newItem).Error
		if err != nil {
			helpers.BadRequestResponse(ctx, err.Error())
			return
		}
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"success": true,
		"result":  fmt.Sprintf("OrderID %s has been succesfully updated", orderId),
	})
}

func (o *OrderController) DeleteOrder(ctx *gin.Context) {
	orderId := ctx.Param("orderId")
	var order models.Order
	var items []models.Item

	err := o.db.First(&order, orderId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "Order data not found")
			return
		}
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	err = o.db.Delete(&order).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err.Error())
			return
		}
		helpers.BadRequestResponse(ctx, err.Error())
		return
	} else {
		o.db.Where("order_id = ?", order.ID).Delete(&items)
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("OrderID %d has been successfully deleted", order.ID),
	})
}
