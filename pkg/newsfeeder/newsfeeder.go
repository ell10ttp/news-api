package newsfeeder

import (
	"news-api/pkg/models"
	"sort"

	"github.com/mmcdole/gofeed"
)

func GetFeed(source models.Source) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(source.URL)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func GetFeedByCategory(source models.Source, category models.Category) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	url := source.Categories[category.String()]
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func SortFeedByPublished(feed gofeed.Feed) gofeed.Feed {
	sort.SliceStable(feed.Items, func(i, j int) bool {
		return feed.Items[i].PublishedParsed.After(*feed.Items[j].PublishedParsed)
	})
	return feed
}
