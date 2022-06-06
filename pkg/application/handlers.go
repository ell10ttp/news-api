package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"news-api/pkg/logger"
	"news-api/pkg/models"
	"news-api/pkg/newsfeeder"

	"github.com/gorilla/mux"
	"github.com/mmcdole/gofeed"
)

func (app *Model) ping(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal("ping successful")

	w.WriteHeader(http.StatusOK)
	code, err := w.Write(response)
	if err != nil {
		logger.Debug("GENERIC", fmt.Sprintf("failed to write error to responseWriter. int code: %d", code), "", 1)
		logger.Error("GENERIC", err, 1)
	}
}

func (app *Model) getSourceList(w http.ResponseWriter, r *http.Request) {
	list := app.SourceAPI.GetSourceList()
	response, _ := json.Marshal(list)

	w.WriteHeader(http.StatusOK)
	code, err := w.Write(response)
	if err != nil {
		logger.Debug("GENERIC", fmt.Sprintf("failed to write error to responseWriter. int code: %d", code), "", 1)
		logger.Error("GENERIC", err, 1)
	}
}

func (app *Model) postSource(w http.ResponseWriter, r *http.Request) {
	sourceMap, err := mapStringInterface(w, r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err.Error())
		return
	}

	source, err := app.SourceAPI.CreateSource(sourceMap)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := struct {
		Action     string `json:"action"`
		Successful bool   `json:"Successful"`
		SourceID   int    `json:"sourceId"`
	}{
		Action:     "create source",
		Successful: true,
		SourceID:   source.ID,
	}

	formattedResponse, err := json.Marshal(response)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	code, err := w.Write(formattedResponse)
	if err != nil {
		logger.Debug("GENERIC", fmt.Sprintf("failed to write error to responseWriter. int code: %d", code), "", 1)
		logger.Error("GENERIC", err, 1)
	}
}

func (app *Model) getSource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sourceId := vars["sourceId"]
	intId, err := strconv.Atoi(sourceId)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "did not provide integer id for source")
		return
	}

	source, err := app.SourceAPI.GetSource(intId)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err.Error())
		return
	}
	if source.ID == 0 {
		app.clientError(w, http.StatusBadRequest, "source id not found")
		return
	}

	response, _ := json.Marshal(source)
	w.WriteHeader(http.StatusOK)
	code, err := w.Write(response)
	if err != nil {
		logger.Debug("GENERIC", fmt.Sprintf("failed to write error to responseWriter. int code: %d", code), "", 1)
		logger.Error("GENERIC", err, 1)
	}
}

func (app *Model) getSourceCategories(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sourceId := vars["sourceId"]
	intId, err := strconv.Atoi(sourceId)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "did not provide integer id for source")
		return
	}

	source, err := app.SourceAPI.GetSource(intId)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err.Error())
		return
	}
	if source.ID == 0 {
		app.clientError(w, http.StatusBadRequest, "source id not found")
		return
	}

	catKeys := make([]string, 0, len(source.Categories))
	for k := range source.Categories {
		catKeys = append(catKeys, k)
	}
	sort.Strings(catKeys)

	response := struct {
		Action             string   `json:"action"`
		Successful         bool     `json:"successful"`
		NumberOfCategories int      `json:"numberOfCategories"`
		Categories         []string `json:"categories"`
	}{
		Action:             "retrieve available categories",
		Successful:         true,
		NumberOfCategories: len(source.Categories),
		Categories:         catKeys,
	}

	formattedResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	code, err := w.Write(formattedResponse)
	if err != nil {
		logger.Debug("GENERIC", fmt.Sprintf("failed to write error to responseWriter. int code: %d", code), "", 1)
		logger.Error("GENERIC", err, 1)
	}
}

func (app *Model) getFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sourceId := vars["sourceId"]
	intId, err := strconv.Atoi(sourceId)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "did not provide integer id for source")
		return
	}

	fmt.Println("sourceID: ", sourceId)

	source, err := app.SourceAPI.GetSource(intId)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "source id not found")
		return
	}
	if source.ID == 0 {
		app.clientError(w, http.StatusBadRequest, "source id not found")
		return
	}

	var feed gofeed.Feed

	categoryStr := r.URL.Query().Get("category")
	if categoryStr == "" {
		feedPtr, err := newsfeeder.GetFeed(source.URL)
		if err != nil {
			app.clientError(w, http.StatusBadRequest, err.Error())
			return
		}
		if feedPtr != nil {
			feed = *feedPtr
		}
	} else {
		category := models.StrToCategory(categoryStr)
		fmt.Println("category is ", categoryStr, " plus ", category)
		if category > models.UK || category < models.Politics {
			// valid category iota between low and high
			fmt.Println("source is category", source.Categories[category.String()])
			feedPtr, err := newsfeeder.GetFeed(source.Categories[category.String()])
			if err != nil {
				app.clientError(w, http.StatusBadRequest, err.Error())
				return
			}
			if feedPtr != nil {
				feed = *feedPtr
			}
		}
	}

	// todo: accept multiple queries, add advanced filtering through url param brackets etc
	sortStr := r.URL.Query().Get("sort")
	if strings.EqualFold(sortStr, "true") {
		feed = newsfeeder.SortFeedByPublished(feed)
	}

	fmt.Println(feed)
	response, _ := json.Marshal(feed.Items)
	w.WriteHeader(http.StatusOK)
	code, err := w.Write(response)
	if err != nil {
		logger.Debug("GENERIC", fmt.Sprintf("failed to write error to responseWriter. int code: %d", code), "", 1)
		logger.Error("GENERIC", err, 1)
	}
}

// Utility functions - not handlers of routes

// MapStringInterface decode http request to map[string]interface
func mapStringInterface(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	b, err := decodeRequestBody(r, w)
	if err != nil {
		return nil, fmt.Errorf("failed to decode map string interface")
	}

	bMap := *b
	result, ok := bMap.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse map")
	}

	return result, nil
}

// decodeRequestBody decode http request
func decodeRequestBody(r *http.Request, w http.ResponseWriter) (*interface{}, error) {
	var tmp interface{}
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	return &tmp, err
}
