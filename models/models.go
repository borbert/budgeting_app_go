package models

import (
	"database/sql"
	"time"
)

//Models is the wrapper for database
type Models struct {
	DB DBModel
}

//New Models returns models with db pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

//Type for budgeting account items
type BudgetAcct struct {
	ID            int            `json:"id"`
	UserID        int            `json:"user_id"`
	Item          string         `json:"item"`
	Description   string         `json:"description"`
	Amount        float64        `json:"amount"`
	BudgetingType string         `json:"budgeting_type"`
	Biweekly      bool           `json:"biweekly"`
	ApplyDefault  bool           `json:"apply_default_amount"`
	DefaultAmt    float64        `json:"default_amt"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	TerminatedAt  time.Time      `json:"terminated_at"`
	Tags          map[int]string `json:"tags"`
}

//Type for tags
type Tag struct {
	TagID       int       `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

//Type for tags
type BudgetTag struct {
	ID          int    `json:"-"`
	TagID       int    `json:"-"`
	ItemID      int    `json:"-"`
	Description string `json:"description"`
}

type User struct {
	ID       int    `json:"_"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type UserPreferences struct {
	ID                      int    `json:"id"`
	User_id                 string `json:"user_id"`
	Budget_items_sort_order []int  `json:"budget_items_sort_order"`
}
