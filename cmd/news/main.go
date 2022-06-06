package main

import (
	"fmt"
	"net/http"
	"os"

	"news-api/pkg/application"
	"news-api/pkg/logger"
	"news-api/pkg/sourceapi"
)

func main() {

	// Configure application instance and check auth server status
	app := application.Init()
	app.SetSourceAPI(sourceapi.NewSourceAPI())

	fmt.Println(os.Getenv("SERVER_PORT"))
	serverPort := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	srv := &http.Server{
		Addr:    serverPort,
		Handler: app.Routes(),
	}

	logger.Info("SETUP", fmt.Sprintf("server listening on port %s", serverPort), 1)
	logger.Info("SETUP", "api ready", 1)

	err := srv.ListenAndServe()
	panic(err)
}
