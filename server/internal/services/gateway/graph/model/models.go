package model

import (
	"fmt"
	"time"

	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// MapMenuItem converts sqlc.MenuItem to the appropriate GraphQL MenuItem type
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

// TimestampToTime converts pgtype.Timestamptz to *time.Time
func TimestampToTime(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	return &ts.Time
}
