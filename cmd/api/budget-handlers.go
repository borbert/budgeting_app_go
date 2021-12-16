package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/borbert/budgeting_app_golang/models"
	"github.com/julienschmidt/httprouter"
)

type jsonResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (app *application) getOneBudgetItem(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Println(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	bdgt, err := app.models.DB.Get(id)

	if err != nil {
		app.logger.Println(err)
	}

	err = app.writeJSON(w, http.StatusOK, bdgt, "budgetItem")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllBudgetItems(w http.ResponseWriter, r *http.Request) {
	budgetItems, err := app.models.DB.All()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, budgetItems, "budget_items")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := app.models.DB.AllTags()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, tags, "tags")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllBudgetItemsbyTag(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	tagID, err := strconv.Atoi(params.ByName("tag"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	budgetItems, err := app.models.DB.All(tagID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, budgetItems, "budgetItems")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}

type BudgetAcctPayload struct {
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

func (app *application) editBudgetItem(w http.ResponseWriter, r *http.Request) {
	var payload BudgetAcctPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	var budgetitem models.BudgetAcct

	if payload.ID != 0 {
		id := payload.ID
		m, _ := app.models.DB.Get(id)
		budgetitem = *m
		budgetitem.UpdatedAt = time.Now()
	}

	budgetitem.ID = payload.ID
	budgetitem.UserID = payload.UserID
	budgetitem.Item = payload.Item
	budgetitem.Description = payload.Description
	budgetitem.Amount = payload.Amount
	budgetitem.BudgetingType = payload.BudgetingType
	budgetitem.Biweekly = payload.Biweekly
	budgetitem.ApplyDefault = payload.ApplyDefault
	budgetitem.DefaultAmt = payload.DefaultAmt
	budgetitem.CreatedAt = payload.CreatedAt
	budgetitem.UpdatedAt = time.Now()
	budgetitem.TerminatedAt = payload.TerminatedAt
	budgetitem.Tags = payload.Tags

	if budgetitem.ID == 0 {
		err = app.models.DB.InsertBudgetItem(budgetitem)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		err = app.models.DB.UpdateBudgetItem(budgetitem)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	ok := jsonResp{
		OK: true,
	}
	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}

func (app *application) deleteBudgetItem(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.models.DB.DeleteBudgetItem(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	ok := jsonResp{
		OK: true,
	}
	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.DB.AllUsers()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, users, "users")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}
func (app *application) getOneUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	email, err := strconv.Atoi(params.ByName("email"))
	if err != nil {
		app.logger.Println(errors.New("invalid email parameter"))
		app.errorJSON(w, err)
		return
	}

	user, err := app.models.DB.Get(email)

	if err != nil {
		app.logger.Println(err)
	}

	err = app.writeJSON(w, http.StatusOK, user, "user")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}
