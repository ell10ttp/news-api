package application

import (
	"encoding/json"
	"fmt"
	"net/http"

	"news-api/pkg/constants"
	"news-api/pkg/logger"
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
	list := constants.GetSourceList()
	response, _ := json.Marshal(list)

	w.WriteHeader(http.StatusOK)
	code, err := w.Write(response)
	if err != nil {
		logger.Debug("GENERIC", fmt.Sprintf("failed to write error to responseWriter. int code: %d", code), "", 1)
		logger.Error("GENERIC", err, 1)
	}
}
