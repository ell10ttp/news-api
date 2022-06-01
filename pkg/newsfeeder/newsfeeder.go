package newsfeeder

import (
	"fmt"
	"news-api/pkg/models"

	"github.com/mmcdole/gofeed"
)

func GetFeed(source models.Source) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(source.URL)
	fmt.Println(feed)
}
