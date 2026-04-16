package model

import (
	"fmt"
	"os"
	"time"

	"github.com/1kyryll/go-grpc/internal/common/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// Custom Order with internal UserID for DataLoader resolution.
// Overrides the gqlgen-generated Order so we can carry the DB foreign key.
type Order struct {
	ID         string       `json:"id"`
	UserID     int32        `json:"-"`
	User       *User        `json:"user"`
	Items      []*OrderItem `json:"items"`
	Ticket     *Ticket      `json:"ticket,omitempty"`
	TotalPrice float64      `json:"totalPrice"`
	Status     OrderStatus  `json:"status"`
	CreatedAt  time.Time    `json:"createdAt"`
	UpdatedAt  time.Time    `json:"updatedAt"`
}

func (Order) IsSearchResult() {}

// Custom OrderItem with internal MenuItemID for DataLoader resolution.
type OrderItem struct {
	ID                  string   `json:"id"`
	MenuItemID          int32    `json:"-"`
	MenuItem            MenuItem `json:"menuItem"`
	Quantity            int      `json:"quantity"`
	SpecialInstructions *string  `json:"specialInstructions,omitempty"`
	Subtotal            float64  `json:"subtotal"`
}

// Custom Ticket with internal OrderID for nested resolution.
type Ticket struct {
	ID        string       `json:"id"`
	OrderID   int32        `json:"-"`
	Order     *Order       `json:"order"`
	Status    TicketStatus `json:"status"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

// MapMenuItem converts sqlc.MenuItem to the appropriate GraphQL MenuItem type.
func MapMenuItem(item sqlc.MenuItem) MenuItem {
	id := fmt.Sprintf("%d", item.ID)

	var desc *string
	if item.Description.Valid {
		desc = &item.Description.String
	}

	price, _ := item.Price.Float64Value()
	priceFloat := price.Float64

	isAvailable := item.IsAvailable.Valid && item.IsAvailable.Bool

	if item.Category == "DRINK" {
		isAlcoholic := item.IsAlcoholic.Valid && item.IsAlcoholic.Bool
		return DrinkItem{
			ID:          id,
			Name:        item.Name,
			Description: desc,
			Price:       priceFloat,
			Category:    MenuCategoryDrink,
			IsAvailable: isAvailable,
			IsAlcoholic: isAlcoholic,
		}
	}

	var allergens []string
	if item.ContainsAllergens != nil {
		allergens = item.ContainsAllergens
	}

	var category MenuCategory
	switch item.Category {
	case "APPETIZER":
		category = MenuCategoryAppetizer
	case "MAIN":
		category = MenuCategoryMain
	case "DESSERT":
		category = MenuCategoryDessert
	default:
		category = MenuCategory(item.Category)
	}

	return FoodItem{
		ID:                id,
		Name:              item.Name,
		Description:       desc,
		Price:             priceFloat,
		Category:          category,
		IsAvailable:       isAvailable,
		ContainsAllergens: allergens,
	}
}

// TimestampToTime converts pgtype.Timestamptz to time.Time.
func TimestampToTime(ts pgtype.Timestamptz) time.Time {
	if !ts.Valid {
		return time.Time{}
	}
	return ts.Time
}

// TextToStringPtr converts pgtype.Text to *string.
func TextToStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

// StringToText converts a Go string to pgtype.Text.
func StringToText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

type User struct {
	ID             string     `json:"id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	Phone          *string    `json:"phone,omitempty"`
	Role           string     `json:"role"`
	HashedPassword string     `json:"-"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `json:"-"`
}

func (User) IsSearchResult() {}

func (u *User) HashPassword(password string) error {
	bytePassword := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.HashedPassword = string(passwordHash)

	return nil
}

func (u *User) GenToken() (*AuthToken, error) {
	expiredAt := time.Now().Add(time.Hour * 24 * 7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ExpiresAt": expiredAt.Unix(),
		"Id":        u.ID,
		"IssuedAt":  time.Now().Unix(),
		"Issuer":    "Restaurant",
	})

	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return &AuthToken{
		AccessToken: accessToken,
		ExpiredAt:   expiredAt,
	}, nil
}

func (u *User) ComparePassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.HashedPassword)

	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}
