package utils

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

func SearchIndex(index bleve.Index, terms ...string) (*bleve.SearchResult, error) {
	var queries []*query.MatchQuery
	for _, term := range terms {
		queries = append(queries, bleve.NewMatchQuery(term))
	}
	boolQuery := bleve.NewBooleanQuery()
	for _, matchQuery := range queries {
		boolQuery.AddMust(matchQuery)
	}

	searchRequest := bleve.NewSearchRequest(boolQuery)
	res, err := index.Search(searchRequest)
	return res, err
}
