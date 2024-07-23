package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Date datatypes.Date

const dateFormat = "2006-01-02"

func (d *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return err
	}
	*d = Date(datatypes.Date(t))
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format(dateFormat))
}

func (d Date) Value() (driver.Value, error) {
	return time.Time(d).Format(dateFormat), nil
}

func (d *Date) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*d = Date(datatypes.Date(v))
		return nil
	case []byte:
		// Convert byte slice to string
		strValue := string(v)
		// Parse string as time.Time
		parsedTime, err := time.Parse(dateFormat, strValue)
		if err != nil {
			return err
		}
		*d = Date(datatypes.Date(parsedTime))
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into Date", value)
	}
}

type Users struct {
	ID                  string `gorm:"type:char(36);primaryKey" json:"id"`
	Name                string `gorm:"type:varchar(35);not null" json:"name" binding:"required"`
	LastName            string `gorm:"type:varchar(35);not null" json:"lastname" binding:"required"`
	Username            string `gorm:"type:varchar(15);not null;unique" json:"username" binding:"required"`
	IsActive            string `gorm:"type:char(1)" json:"is_active" binding:"oneof=S N"`
	Birthdate           Date   `json:"birthdate" binding:"required"`
	ImgProfile          string `gorm:"type:varchar(255)" json:"img_profile"`
	ImgBanner           string `gorm:"type:varchar(255)" json:"img_banner"`
	Password            string `gorm:"type:varchar(60);not null" json:"password" binding:"required"`
	FailedLoginAttempts uint8  `gorm:"default:0" json:"failed_login_attempts"`
	RolesID             uint8  `json:"roles_id" binding:"required"`
	Roles               Roles
	gorm.Model
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserMin struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// UUID v4 Generator
func generateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func (user *Users) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := generateUUID()
	if err != nil {
		return err
	}

	hash, err := helper.HashPassword(user.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return
	}

	user.ID = id
	user.Password = hash
	return
}

type UsersResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	LastName   string `json:"last_name"`
	Username   string `json:"username"`
	IsActive   string `json:"is_active"`
	Birthdate  Date   `json:"birthdate"`
	ImgProfile string `json:"img_profile"`
	ImgBanner  string `json:"img_banner"`
	RoleName   string `json:"role_name"`
}
