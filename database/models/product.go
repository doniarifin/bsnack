package models

import "time"

type Product struct {
	ID             string     `json:"id" gorm:"primaryKey"`
	Name           string     `json:"name"`
	Type           string     `json:"type"`
	Flavor         string     `json:"flavor"`
	Size           string     `json:"size"`
	Price          float64    `json:"price"`
	Stock          int        `json:"stock"`
	ProductionDate time.Time  `json:"production_date"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
