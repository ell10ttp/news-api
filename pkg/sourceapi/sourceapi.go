package sourceapi

import (
	"errors"
	"net/url"
	"news-api/pkg/models"

	"github.com/biter777/countries"
)

type ISourceAPI interface {
	GetSourceList() models.SourceList
	GetSource(sourceId int) (models.Source, error)
	CreateSource(sourceMap map[string]interface{}) (models.Source, error)
}

type SourceAPI struct{}

// NewSourceAPI create pointer to SourceAPI
func NewSourceAPI() *SourceAPI {
	api := SourceAPI{}
	return &api
}

var (
	sourceList = models.SourceList{
		Sources: []models.Source{
			{
				ID:          1,
				Name:        "BBC News",
				Description: "Trusted World and UK news as well as local and regional perspectives.",
				URL:         "http://feeds.bbci.co.uk/news/uk/rss.xml",
				Language:    countries.UnitedKingdom.Alpha3(),
				Country:     countries.UnitedKingdom.Info().Name,
			},
			{
				ID:          2,
				Name:        "Sky News",
				Description: "Expert comment and analysis on the latest UK news, with headlines from England, Scotland, Northern Ireland and Wales.",
				URL:         "https://feeds.skynews.com/feeds/rss/uk.xml",
				Language:    countries.UnitedKingdom.Alpha3(),
				Country:     countries.UnitedKingdom.Info().Name,
			},
		},
	}
)

func (s *SourceAPI) GetSourceList() (list models.SourceList) {
	return sourceList
}

func addSource(newSource models.Source) {
	sourceList.Sources = append(sourceList.Sources, newSource)
}

func nextSourceId() int {
	id := sourceList.Sources[len(sourceList.Sources)-1].ID + 1
	return id
}

func (s *SourceAPI) GetSource(sourceId int) (models.Source, error) {
	for _, src := range sourceList.Sources {
		if src.ID == sourceId {
			return src, nil
		}
	}
	return models.Source{}, errors.New("source id not found")
}

func (s *SourceAPI) CreateSource(sourceMap map[string]interface{}) (models.Source, error) {
	newSource := models.Source{
		ID:          nextSourceId(),
		Name:        sourceMap["Name"].(string),
		Description: sourceMap["Description"].(string),
		URL:         sourceMap["Url"].(string),
		Language:    sourceMap["Language"].(string),
		Country:     sourceMap["Country"].(string),
	}
	if newSource.Name == "" {
		return models.Source{}, errors.New("invalid source: \"Name\" was empty")
	}
	if _, err := url.ParseRequestURI(newSource.URL); err != nil {
		return models.Source{}, errors.New("invalid source: \"Url\" failed to parse")
	}

	addSource(newSource)
	return newSource, nil
}
