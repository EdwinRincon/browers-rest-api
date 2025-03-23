package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/EdwinRincon/browersfc-api/helper"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Date time.Time

const dateFormat = "2006-01-02"

func (d *Date) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d))
}

func (d Date) Value() (driver.Value, error) {
	return time.Time(d).Format(dateFormat), nil
}

func (d *Date) Scan(v interface{}) error {
	switch value := v.(type) {
	case time.Time:
		*d = Date(value)
		return nil
	case []byte:
		// Parse the string date from the byte slice
		t, err := time.Parse(dateFormat, string(value))
		if err != nil {
			return fmt.Errorf("failed to parse date from []byte: %w", err)
		}
		*d = Date(t)
		return nil
	case string:
		// Handle string values
		t, err := time.Parse(dateFormat, value)
		if err != nil {
			return fmt.Errorf("failed to parse date from string: %w", err)
		}
		*d = Date(t)
		return nil
	case nil:
		// Handle NULL values
		return nil
	default:
		return fmt.Errorf("can't scan %T into Date", v)
	}
}

type User struct {
	ID                  string         `gorm:"type:char(36);primaryKey" json:"id"`
	Name                string         `gorm:"type:varchar(35);not null" json:"name" binding:"required"`
	LastName            string         `gorm:"type:varchar(35);not null" json:"last_name" binding:"required"`
	Username            string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"username" binding:"required"`
	IsActive            bool           `gorm:"default:true" json:"is_active"`
	Birthdate           Date           `json:"birthdate" example:"1990-01-01"`
	ImgProfile          string         `gorm:"type:varchar(255)" json:"img_profile,omitempty"`
	ImgBanner           string         `gorm:"type:varchar(255)" json:"img_banner,omitempty"`
	Password            string         `gorm:"type:varchar(60);not null" json:"-" binding:"required"`
	FailedLoginAttempts uint8          `gorm:"default:0" json:"-"`
	RoleID              uint8          `json:"role_id" binding:"required"`
	Role                Role           `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	CreatedAt           time.Time      `json:"created_at,omitempty"`
	UpdatedAt           time.Time      `json:"updated_at,omitempty"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate UUID: %w", err)
	}

	hash, err := helper.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.ID = id.String()
	user.Password = hash
	return nil
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserMin struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UserResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	LastName   string `json:"last_name"`
	Username   string `json:"username"`
	IsActive   bool   `json:"is_active"`
	Birthdate  Date   `json:"birthdate"`
	ImgProfile string `json:"img_profile,omitempty"`
	ImgBanner  string `json:"img_banner,omitempty"`
	RoleName   string `json:"role_name"`
}

type UserUpdate struct {
	Name       *string `json:"name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	Username   *string `json:"username,omitempty"`
	Birthdate  *Date   `json:"birthdate,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
	ImgProfile *string `json:"img_profile,omitempty"`
	ImgBanner  *string `json:"img_banner,omitempty"`
}
