package models

import "time"

type Transaction struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Amount      float64   `json:"amount"`
}
