package words

import (
	"github.com/reiver/go-porterstemmer"
)

type Stemmer interface {
	Stem(word string) string
}

type porterStemmer struct{}

func NewStemmer() Stemmer {
	return &porterStemmer{}
}

func (ps *porterStemmer) Stem(word string) string {
	return porterstemmer.StemString(word)
}
