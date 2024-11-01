package model

import (
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,max=100"`
}

type RegisterUser struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" binding:"required,max=100"`
	Email       string    `json:"email" binding:"required,email,max=100"`
	Password    string    `json:"password,omitempty" binding:"required,max=100"`
	Bio         string    `json:"bio" binding:"max=500"`
	Gender      string    `json:"gender" binding:"oneof=MALE FEMALE OTHERS"`
	Preference  string    `json:"preference" binding:"oneof=MALE FEMALE OTHERS"`
	DateOfBirth time.Time `json:"age" binding:"required" time_format:"2006-01-02"`
	Images      []string  `json:"images" binding:"required,min=1,max=5"`
}

type AuthUser struct {
	User
	AuthToken           string    `json:"token"`
	AuthTokenExpires    time.Time `json:"token_expires"`
	RefreshToken        string    `json:"refresh_token"`
	RefreshTokenExpires time.Time `json:"refresh_token_expires"`
}

type User struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Email            string     `json:"email"`
	Password         string     `gorm:"->:false;<-:create" json:"password,omitempty"`
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
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	URL       string     `json:"url"`
	CreatedAt time.Time  `gorm:"<-:create" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (u *User) CheckPassword(password string) error {
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
	for _, image := range images {
		u.Images = append(u.Images, Image{
			ID:        ulid.Make().String(),
			URL:       image,
			UserID:    u.ID,
			CreatedAt: time.Now(),
		})
	}
}

func (ru *RegisterUser) ToUserModel() *User {
	age := calculateAge(ru.DateOfBirth)

	return &User{
		ID:         ru.ID,
		Name:       ru.Name,
		Email:      ru.Email,
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
