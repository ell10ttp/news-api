package application

import "news-api/pkg/sourceapi"

type Model struct {
	SourceAPI sourceapi.ISourceAPI
}

func Init() *Model {
	app := Model{}
	return &app
}

func (app *Model) SetSourceAPI(sourceAPI sourceapi.ISourceAPI) {
	app.SourceAPI = sourceAPI
}
