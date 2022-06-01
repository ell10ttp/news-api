package models

import "fmt"

type Category int

// TODO extract this to a config map for additional categories to be added
// alternatively put this as values for retrieval in param store/s3/database
// and have aws/db integration into api to allow client admin route to expand
const (
	UK Category = iota + 1
	World
	Business
	Technology
	Entertainment
	Politics
)

func (c Category) String() string {
	categories := [...]string{"uk", "world", "business", "technology", "entertainment", "politics"}
	if c < UK || c > Politics {
		return fmt.Sprintf("Category(%d)", int(c))
	}
	return categories[c-1]
}

func (c Category) IsValid() bool {
	switch c {
	case UK, World, Business, Technology, Politics:
		return true
	}
	return false
}
