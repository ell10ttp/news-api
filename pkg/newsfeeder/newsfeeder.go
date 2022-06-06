package newsfeeder

import (
	"sort"

	"github.com/mmcdole/gofeed"
)

func GetFeed(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
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
