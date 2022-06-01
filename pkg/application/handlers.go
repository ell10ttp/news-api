package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"news-api/pkg/logger"
	"news-api/pkg/newsfeeder"

	"github.com/gorilla/mux"
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

	type Response struct {
		Action     string `json:"action"`
		Successful bool   `json:"Successful"`
		SourceID   int    `json:"sourceId"`
	}
	response := Response{
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

func (app *Model) parseFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sourceId := vars["sourceId"]
	intId, err := strconv.Atoi(sourceId)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "did not provide integer id for source")
		return
	}

	source, err := app.SourceAPI.GetSource(intId)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "source id not found")
		return
	}
	if source.ID == 0 {
		app.clientError(w, http.StatusBadRequest, "source id not found")
		return
	}

	newsfeeder.ParseFeed(source)

	response, _ := json.Marshal(source)
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
