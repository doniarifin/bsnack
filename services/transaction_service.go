package services

import (
	"bsnack/database/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TransactionService struct {
	DB     *gorm.DB
	redis  *redis.Client
	prdSrv *ProductService
	cstSrv *CustomerService
}

func NewTransactionService(db *gorm.DB, redis *redis.Client, prd *ProductService, cst *CustomerService) *TransactionService {
	return &TransactionService{
		DB:     db,
		redis:  redis,
		prdSrv: prd,
		cstSrv: cst,
	}
}

func (s *TransactionService) CreateTransaction(payload *models.Transaction) error {
	prod, err := s.prdSrv.CheckStock(payload.ProductID)
	if err != nil {
		return err
	}

	if payload.Quantity > prod.Stock {
		return fmt.Errorf("product stock not enough: requested %d, available %d", payload.Quantity, prod.Stock)
	}
	if !payload.IsNewCustomer && payload.CustomerID == "" {
		return fmt.Errorf("erorr: if this not new customer, please select existing customer_id")
	}

	prod.Stock = prod.Stock - payload.Quantity
	err = s.prdSrv.UpdateProduct(prod)
	if err != nil {
		return err
	}

	//check customer
	var cst models.Customer
	if payload.IsNewCustomer {
		cst.Name = payload.CustomerName
		cst.Points = payload.Quantity * int(prod.Price) / 1000
		err := s.cstSrv.CreateCustomer(&cst)
		if err != nil {
			return err
		}
	} else {
		cst.ID = payload.CustomerID
		cst.Name = payload.CustomerName
		cst.Points = cst.Points + (payload.Quantity * int(prod.Price) / 1000)
		err := s.cstSrv.UpdateCustomer(&cst)
		if err != nil {
			return err
		}
	}

	payload.ID = uuid.New().String()
	payload.ProductName = prod.Name
	payload.TotalPrice = float64(payload.Quantity) * prod.Price

	payload.TransactionAt = time.Now()
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()

	s.redis.Del(ctx, "trxs:all")
	return s.DB.Create(payload).Error
}

func (s *TransactionService) UpdateTransaction(payload *models.Transaction) error {
	payload.UpdatedAt = time.Now()

	s.redis.Del(ctx, "trxs:all")
	s.redis.Del(ctx, fmt.Sprintf("trx:%s", payload.ID))
	return s.DB.Save(payload).Error
}

func (s *TransactionService) GetAll() ([]models.Transaction, error) {
	val, err := s.redis.Get(ctx, "trxs:all").Result()
	if err == nil {
		var cached []models.Transaction
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return cached, nil
		}
	}

	var data []models.Transaction
	err = s.DB.Order("created_at DESC").Find(&data).Error

	jsonData, _ := json.Marshal(data)
	s.redis.Set(ctx, "trxs:all", jsonData, 30*time.Second)

	return data, err
}

func (s *TransactionService) GetByID(id string) (*models.Transaction, error) {
	key := fmt.Sprintf("trx:%s", id)
	val, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		var cached models.Transaction
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return &cached, nil
		}
	}

	var t models.Transaction
	err = s.DB.Where("id = ?", id).First(&t).Error
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	jsonData, _ := json.Marshal(t)
	s.redis.Set(ctx, key, jsonData, 30*time.Second)

	return &t, nil
}
