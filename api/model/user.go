package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Date time.Time

const dateFormat = "2006-01-02"

func (d *Date) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}

	t, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
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
	ID         string         `gorm:"type:char(36);primaryKey" json:"id"`
	Name       string         `gorm:"type:varchar(35);not null" json:"name" binding:"required,min=2,max=35"`
	LastName   string         `gorm:"type:varchar(35);not null" json:"last_name" binding:"required,min=2,max=35"`
	Username   string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"username" binding:"required,safe_email,allowed_domain"`
	Birthdate  Date           `json:"birthdate" example:"1990-01-01"`
	ImgProfile string         `gorm:"type:varchar(255)" json:"img_profile,omitempty" binding:"omitempty,url"`
	ImgBanner  string         `gorm:"type:varchar(255)" json:"img_banner,omitempty" binding:"omitempty,url"`
	RoleID     uint8          `json:"role_id" binding:"required,min=1"`
	Role       *Role          `gorm:"foreignKey:RoleID" json:"role,omitempty" binding:"-"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (u *User) BeforeCreate(tx *gorm.DB) error {

	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
