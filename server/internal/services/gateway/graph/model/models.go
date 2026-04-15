package model

import (
	"fmt"
	"time"

	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// Custom Order with internal UserID for DataLoader resolution.
// Overrides the gqlgen-generated Order so we can carry the DB foreign key.
type Order struct {
	ID     string    `json:"id"`
	UserID int32     `json:"-"`
	User   *User     `json:"user"`
	Items      []*OrderItem `json:"items"`
	Ticket     *Ticket     `json:"ticket,omitempty"`
	TotalPrice float64     `json:"totalPrice"`
	Status     OrderStatus `json:"status"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
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
