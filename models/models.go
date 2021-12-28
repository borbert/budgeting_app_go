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
	ID       int    `json:"-"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type UserPreferences struct {
	ID                      int    `json:"-"`
	User_id                 string `json:"user_id"`
	Budget_items_sort_order []int  `json:"preferences"`
}

type Drug struct {
	Ndc11       string `json:"ndc11"`
	Ndc10       string `json:"ndc10"`
	Ndc9        string `json:"ndc9"`
	Description string `json:"description"`
	Favorite    bool   `json:"favorite"`
}

type OpenFDA struct {
	Meta struct {
		Disclaimer  string `json:"disclaimer"`
		Terms       string `json:"terms"`
		License     string `json:"license"`
		LastUpdated string `json:"last_updated"`
		Results     struct {
			Skip  int `json:"skip"`
			Limit int `json:"limit"`
			Total int `json:"total"`
		} `json:"results"`
	} `json:"meta"`
	Results []struct {
		ProductNdc        string `json:"product_ndc"`
		GenericName       string `json:"generic_name"`
		LabelerName       string `json:"labeler_name"`
		BrandName         string `json:"brand_name"`
		ActiveIngredients []struct {
			Name     string `json:"name"`
			Strength string `json:"strength"`
		} `json:"active_ingredients"`
		Finished  bool `json:"finished"`
		Packaging []struct {
			PackageNdc         string `json:"package_ndc"`
			Description        string `json:"description"`
			MarketingStartDate string `json:"marketing_start_date"`
			Sample             bool   `json:"sample"`
		} `json:"packaging"`
		ListingExpirationDate string `json:"listing_expiration_date"`
		Openfda               struct {
			ManufacturerName []string `json:"manufacturer_name"`
			Rxcui            []string `json:"rxcui"`
			SplSetID         []string `json:"spl_set_id"`
			Nui              []string `json:"nui"`
			PharmClassMoa    []string `json:"pharm_class_moa"`
			PharmClassEpc    []string `json:"pharm_class_epc"`
			Unii             []string `json:"unii"`
		} `json:"openfda"`
		MarketingCategory  string   `json:"marketing_category"`
		DosageForm         string   `json:"dosage_form"`
		SplID              string   `json:"spl_id"`
		ProductType        string   `json:"product_type"`
		Route              []string `json:"route"`
		MarketingStartDate string   `json:"marketing_start_date"`
		ProductID          string   `json:"product_id"`
		ApplicationNumber  string   `json:"application_number"`
		BrandNameBase      string   `json:"brand_name_base"`
		PharmClass         []string `json:"pharm_class"`
	} `json:"results"`
}
