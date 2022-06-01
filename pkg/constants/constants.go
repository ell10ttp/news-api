package constants

import (
	"news-api/pkg/models"

	"github.com/biter777/countries"
)

func GetSourceList() (list models.SourceList) {
	list.Sources = []models.Source{
		{
			ID:          "1",
			Name:        "BBC News UK",
			Description: "Trusted World and UK news as well as local and regional perspectives.",
			URL:         "http://feeds.bbci.co.uk/news/uk/rss.xml",
			Category:    "National",
			Language:    countries.UnitedKingdom.Alpha3(),
			Country:     countries.UnitedKingdom.Info().Name,
		},
		{
			ID:          "2",
			Name:        "BBC News Technology",
			Description: "Breaking news and analysis on computing, the web, blogs, games, gadgets, social media, broadband and more.",
			URL:         "http://feeds.bbci.co.uk/news/technology/rss.xml",
			Category:    "Technology",
			Language:    countries.UnitedKingdom.Alpha3(),
			Country:     countries.UnitedKingdom.Info().Name,
		},
		{
			ID:          "3",
			Name:        "Sky News UK",
			Description: "Expert comment and analysis on the latest UK news, with headlines from England, Scotland, Northern Ireland and Wales.",
			URL:         "https://feeds.skynews.com/feeds/rss/uk.xml",
			Category:    "National",
			Language:    countries.UnitedKingdom.Alpha3(),
			Country:     countries.UnitedKingdom.Info().Name,
		},
		{
			ID:          "4",
			Name:        "Sky News Technology",
			Description: "Provides you with all the latest tech and gadget news, game reviews, Internet and web news across the globe.",
			URL:         "https://feeds.skynews.com/feeds/rss/uk.xml",
			Category:    "Technology",
			Language:    countries.UnitedKingdom.Alpha3(),
			Country:     countries.UnitedKingdom.Info().Name,
		},
	}

	return list
}
