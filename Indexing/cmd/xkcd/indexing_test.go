package main

import (
	"testing"
	"xkcd-searcher/pkg/database"
)

var normComics, err = database.LoadComics("../../database.json")
var index = database.BuildIndex(normComics)
var NormalizedQuery = []string{"question",
	"follow"}

func BenchmarkSearchWithIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IndexingSearchComics(NormalizedQuery, index, normComics)
	}
}

func BenchmarkSearchWithoutIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = InefficientSearchComics(NormalizedQuery, normComics)
	}
}
