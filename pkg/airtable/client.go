package airtable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Base struct {
	Id     string
	ApiKey string
	http.Client
}

type Table struct {
	Name string
	Base
}

// airtable marshaling / unmarshaling structs

type ATRecord struct {
	Id     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

type ATList struct {
	Records []ATRecord `json:"records"`
	Offset  string     `json:"offset"`
}

func (tbl *Table) BaseUri() string {
	return fmt.Sprintf(
		"https://api.airtable.com/v0/%s/%s",
		tbl.Base.Id,
		url.PathEscape(tbl.Name),
	)
}

func (tbl *Table) List() string {
	defaultResponse := "[]"

	req, err := http.NewRequest(
		"GET",
		tbl.BaseUri(),
		strings.NewReader(""),
	)
	if err != nil {
		log.Println(err)
		return defaultResponse
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tbl.Base.ApiKey))
	req.Header.Set("Access-Control-Allow-Origin", "*")
	resp, err := tbl.Base.Client.Do(req)
	if err != nil {
		log.Println(err)
		return defaultResponse
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return defaultResponse
	}

	log.Printf("Get list %s/%s\n", tbl.Base.Id, tbl.Name)

	// convert into proper format
	var data ATList
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		return defaultResponse
	}

	return data.String()
}

func (tbl *Table) View(viewName string) string {
	defaultResponse := "[]"

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s?view=%s", tbl.BaseUri(), url.QueryEscape(viewName)),
		strings.NewReader(""),
	)
	if err != nil {
		log.Println(err)
		return defaultResponse
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tbl.Base.ApiKey))
	req.Header.Set("Access-Control-Allow-Origin", "*")
	resp, err := tbl.Base.Client.Do(req)
	if err != nil {
		log.Println(err)
		return defaultResponse
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return defaultResponse
	}

	log.Printf("Get view %s/%s :: %s\n", tbl.Base.Id, tbl.Name, viewName)

	// convert into proper format
	var data ATList
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		return defaultResponse
	}

	return data.String()
}

func (tbl *Table) Get(recordId string) string {
	defaultResponse := "{}"

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%s", tbl.BaseUri(), recordId),
		strings.NewReader(""),
	)
	if err != nil {
		log.Println(err)
		return defaultResponse
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tbl.Base.ApiKey))
	req.Header.Set("Access-Control-Allow-Origin", "*")
	resp, err := tbl.Base.Client.Do(req)
	if err != nil {
		log.Println(err)
		return defaultResponse
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return defaultResponse
	}

	log.Printf("Get rec %s/%s/%s\n", tbl.Base.Id, tbl.Name, recordId)

	// convert into proper format
	var data ATRecord
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		return defaultResponse
	}

	return data.String()
}

func (rec *ATRecord) String() string {
	flattened := make(map[string]interface{})
	flattened["id"] = rec.Id
	for k, v := range rec.Fields {
		flattened[k] = v
	}

	resp, err := json.Marshal(flattened)
	if err != nil {
		log.Println(err)
		return "{}"
	}

	return string(resp)
}

func (lst *ATList) String() string {
	records := make([]string, len(lst.Records))
	for i, r := range lst.Records {
		records[i] = r.String()
	}
	return "[" + strings.Join(records, ",") + "]"
}
