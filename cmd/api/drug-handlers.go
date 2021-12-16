package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/borbert/budgeting_app_golang/models"
)

func getOpenFDA(drug models.Drug) models.Drug {

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
	client := &http.Client{}
	key := ""
	theUrl := ""
	log.Println(theUrl + key + "&query=" + url.QueryEscape(drug.Ndc9))

	req, err := http.NewRequest("GET", theUrl+key+"&query="+url.QueryEscape(drug.Ndc9), nil)
	if err != nil {
		log.Println(err)
		return drug
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return drug
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return drug
	}

	var responseObject OpenFDA

	json.Unmarshal(bodyBytes, &responseObject)

	if len(responseObject.Results) > 0 {
		return drug
	} else {

		return drug
	}
}
