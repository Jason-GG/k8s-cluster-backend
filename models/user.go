package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}
