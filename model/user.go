package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

var SupportedGender = map[string]string{
	"MALE":   "MALE",
	"FEMALE": "FEMALE",
}

var SupportedPreference = map[string]string{
	"MALE":   "MALE",
	"FEMALE": "FEMALE",
	"BOTH":   "BOTH",
}

var GenderMale = SupportedGender["MALE"]
var GenderFemale = SupportedGender["FEMALE"]
var PreferenceMale = SupportedPreference["MALE"]
var PreferenceFemale = SupportedPreference["FEMALE"]
var PreferenceBoth = SupportedPreference["BOTH"]

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,max=100"`
}

type RegisterUser struct {
	Name        string   `json:"name" binding:"required,max=100"`
	Email       string   `json:"email" binding:"required,email,max=100"`
	Password    string   `json:"password,omitempty" binding:"required,max=100"`
	Bio         string   `json:"bio" binding:"max=500"`
	Gender      string   `json:"gender" binding:"oneof=MALE FEMALE"`
	Preference  string   `json:"preference" binding:"oneof=MALE FEMALE BOTH"`
	DateOfBirth string   `json:"date_of_birth" binding:"required" time_format:"2006-01-02"`
	Images      []string `json:"images" binding:"required,min=1,max=5"`
}

type AuthUser struct {
	User
	AuthToken    string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Email            string     `json:"email"`
	Password         string     `gorm:"<-:create" json:"-"`
	Bio              string     `json:"bio"`
	Gender           string     `json:"gender"`
	Preference       string     `json:"preference"`
	Age              int        `json:"age"`
	Images           []Image    `json:"images"`
	StripeCustomerID string     `json:"-"`
	CreatedAt        time.Time  `gorm:"<-:create" json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

type Image struct {
	ID        int        `json:"id"`
	UserID    string     `json:"user_id"`
	URL       string     `json:"url"`
	IsPrimary bool       `json:"is_primary"`
	CreatedAt time.Time  `gorm:"<-:create" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (u *User) CheckPassword(password string) error {
	fmt.Println(u.Password, password)
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	return nil
}

func (u *User) NewImageFromRequest(images []string) {
	u.Images = []Image{}
	for i, image := range images {
		u.Images = append(u.Images, Image{
			URL:       image,
			UserID:    u.ID,
			CreatedAt: time.Now(),
			IsPrimary: i == 0,
		})
	}
}

func (ru *RegisterUser) ToUserModel() *User {
	dob, err := time.Parse("2006-01-02", ru.DateOfBirth)
	if err != nil {
		return &User{}
	}

	age := calculateAge(dob)

	return &User{
		ID:         ulid.Make().String(),
		Name:       ru.Name,
		Email:      strings.ToLower(ru.Email),
		Password:   ru.Password,
		Bio:        ru.Bio,
		Gender:     ru.Gender,
		Preference: ru.Preference,
		Age:        age,
	}
}

func calculateAge(birthDate time.Time) int {
	today := time.Now()
	age := today.Year() - birthDate.Year()

	// Adjust age if birthday hasn't occurred this year
	if today.YearDay() < birthDate.YearDay() {
		age--
	}

	return age
}
