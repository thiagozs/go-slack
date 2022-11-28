package fuzzy

import (
	fzf "go.deanishe.net/fuzzy"
)

type Result struct {
	Match   bool
	Query   string
	Score   float64
	SortKey string
}

type FzfSearcher struct {
	terms []string
}

func NewFzfSearcher(terms []string) *FzfSearcher {
	return &FzfSearcher{
		terms: terms,
	}
}

func (f *FzfSearcher) SetTerms(terms []string) {
	f.terms = terms
}

func (f *FzfSearcher) GetTerms() []string {
	return f.terms
}

func (f *FzfSearcher) Search(term string) []Result {
	rr := []Result{}
	result := fzf.SortStrings(f.terms, term)

	for _, r := range result {
		rr = append(rr, Result{
			Match:   r.Match,
			Query:   r.Query,
			Score:   r.Score,
			SortKey: r.SortKey,
		})
	}

	return rr
}
