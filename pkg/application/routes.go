package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"news-api/pkg/logger"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *Model) Routes() http.Handler {

	router := mux.NewRouter()

	standardMiddleware := alice.New(
		logRequest,
	)

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Request-Headers", "Access-Control-Allow-Origin", "Access-Control-Request-Method", "Connection", "Host", "Origin", "User-Agent", "Referer", "Cache-Control", "X-header", "X-Forwarded-Tls-Client-Cert-Info"}),
		handlers.AllowedOrigins([]string{"*"}), // TODO: restrict to specific urls of client app
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
	)

	router.Use(cors)
	router.Use(panicRecovery)

	router.HandleFunc("/ping", app.ping).Methods(http.MethodGet)

	source := router.PathPrefix("/source").Subrouter()
	source.HandleFunc("", app.getSourceList).Methods(http.MethodGet)
	source.HandleFunc("", app.postSource).Methods(http.MethodPost)
	source.HandleFunc("/{sourceId:[0-9]*}", app.getSource).Methods(http.MethodGet)
	source.HandleFunc("/{sourceId:[0-9]*}/categories", app.getSourceCategories).Methods(http.MethodGet)
	source.HandleFunc("/{sourceId:[0-9]*}/feed", app.getFeed).Methods(http.MethodGet)
	// source.HandleFunc("/{sourceId:[0-9]*}/feed/params", app.getCategoryFeed).Methods(http.MethodGet)

	return standardMiddleware.Then(router)
}

// logRequest - output logs to stdout e.g. INFO	2020/09/22 19:56:28 192.168.0.1:8443 - HTTP/1.1 GET /api/user/X00001
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpointInfo := fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		logger.Info("GENERIC", endpointInfo, 1)
		next.ServeHTTP(w, r)
	})
}

func panicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				logger.Error("SERVER", fmt.Errorf("server panicked. recovering from error %v\n %s", err, buf), 7)
				code, err := w.Write([]byte(`{"error":"server got panic"}`)) // todo: informative error message
				if err != nil {
					logger.Debug("GENERIC", fmt.Sprintf("failed to write error to responseWriter. int code: %d", code), "", 1)
					logger.Error("GENERIC", err, 1)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Model) clientError(w http.ResponseWriter, status int, message string) {
	responseStruct := struct {
		Reason string
	}{
		Reason: message,
	}

	response, err := json.Marshal(responseStruct)
	if err != nil {
		logger.Error("GENERIC", err, 1)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	code, err := w.Write(response)
	if err != nil {
		logger.Debug("GENERIC", fmt.Sprintf("failed to write error to responseWriter. int code: %d", code), "", 1)
		logger.Error("GENERIC", err, 1)
	}
}
