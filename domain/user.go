package domain

import (
	"regexp"
	"time"
)

// User represents a user in the domain layer.
// This is a pure business entity without infrastructure concerns.
type User struct {
	ID         string
	Name       string
	LastName   string
	Username   string
	Birthdate  *time.Time
	ImgProfile string
	ImgBanner  string
	RoleID     uint64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// IsValid performs comprehensive domain validation for the user.
func (u *User) IsValid() bool {
	return u.isValidName() &&
		u.isValidLastName() &&
		u.isValidEmail() &&
		u.isValidBirthdate() &&
		u.RoleID > 0
}

// isValidEmail checks if the username follows a valid email format.
func (u *User) isValidEmail() bool {
	if u.Username == "" {
		return false
	}
	// Simple email regex for basic validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(u.Username)
}

// isValidName checks if the name contains only valid characters.
func (u *User) isValidName() bool {
	if len(u.Name) < 2 || len(u.Name) > 35 {
		return false
	}
	// Check for valid name characters (letters, spaces, hyphens, apostrophes)
	nameRegex := regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s'-]+$`)
	return nameRegex.MatchString(u.Name)
}

// isValidLastName checks if the last name contains only valid characters.
func (u *User) isValidLastName() bool {
	if len(u.LastName) < 2 || len(u.LastName) > 35 {
		return false
	}
	// Check for valid name characters (letters, spaces, hyphens, apostrophes)
	nameRegex := regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s'-]+$`)
	return nameRegex.MatchString(u.LastName)
}

// isValidBirthdate checks if the birthdate is reasonable.
func (u *User) isValidBirthdate() bool {
	if u.Birthdate == nil {
		return true // Birthdate is optional
	}

	now := time.Now()
	// Check if birthdate is not in the future and age is reasonable (not older than 120 years)
	age := now.Year() - u.Birthdate.Year()
	if now.YearDay() < u.Birthdate.YearDay() {
		age--
	}

	return !u.Birthdate.After(now) && age >= 0 && age <= 120
}
