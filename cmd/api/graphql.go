package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/borbert/budgeting_app_golang/models"
	"github.com/graphql-go/graphql"
)

var budgetItems []*models.BudgetAcct

//graphql schema definition
var fields = graphql.Fields{
	"budgetItem": &graphql.Field{
		Type:        budgetAcctType,
		Description: "Get all budget accounts by id and user id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, ok := p.Args["id"].(int)
			//find budget items by id
			if ok {
				for _, budgetItem := range budgetItems {
					if budgetItem.ID == id {
						return budgetItem, nil
					}
				}
			}

			return nil, nil
		},
	},
	"list": &graphql.Field{
		Type:        graphql.NewList(budgetAcctType),
		Description: "Get all budget account items",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return budgetItems, nil
		},
	},
	"search": &graphql.Field{
		Type:        graphql.NewList(budgetAcctType),
		Description: "Search budget items by description",
		Args: graphql.FieldConfigArgument{
			"descriptionContains": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var theList []*models.BudgetAcct
			search, ok := params.Args["descriptionContains"].(string)
			if ok {
				for _, currentBudgetItem := range budgetItems {
					if strings.Contains(currentBudgetItem.Description, search) {
						log.Println("Found one")
						theList = append(theList, currentBudgetItem)
					}
				}
			}
			return theList, nil
		},
	},
}

//define the budetItem fields available for query
var budgetAcctType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "budgetItem",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"user_id": &graphql.Field{
				Type: graphql.Int,
			},
			"item": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"amount": &graphql.Field{
				Type: graphql.Float,
			},
			"budgeting_type": &graphql.Field{
				Type: graphql.Boolean,
			},
			"biweekly": &graphql.Field{
				Type: graphql.Boolean,
			},
			"apply_default_amount": &graphql.Field{
				Type: graphql.Boolean,
			},
			"default_amt": &graphql.Field{
				Type: graphql.Float,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"terminated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

func (app *application) budgetItemsGraphQL(w http.ResponseWriter, r *http.Request) {
	budgetItems, _ = app.models.DB.All()

	q, _ := io.ReadAll(r.Body)
	query := string(q)

	log.Println(query)

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		app.errorJSON(w, errors.New("failed to create schema"))
		log.Println(err)
		return
	}

	params := graphql.Params{Schema: schema, RequestString: query}
	resp := graphql.Do(params)
	if len(resp.Errors) > 0 {
		app.errorJSON(w, errors.New(fmt.Sprintf("failed: %+v", resp.Errors)))
	}
	j, _ := json.MarshalIndent(resp, "", " ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
