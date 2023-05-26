package wikilinkui

import "time"

// Search Result from Wikipedia
type SearchResult struct {
	Batchcomplete string         `json:"batchcomplete"`
	Continue      SearchContinue `json:"continue"`
	Query         SearchQuery    `json:"query"`
}
type SearchContinue struct {
	Sroffset int    `json:"sroffset"`
	Continue string `json:"continue"`
}
type SearchInfo struct {
	Totalhits int `json:"totalhits"`
}
type Search struct {
	Ns        int       `json:"ns"`
	Title     string    `json:"title"`
	Pageid    int       `json:"pageid"`
	Size      int       `json:"size"`
	Wordcount int       `json:"wordcount"`
	Snippet   string    `json:"snippet"`
	Timestamp time.Time `json:"timestamp"`
}
type SearchQuery struct {
	SearchInfo SearchInfo `json:"searchinfo"`
	Search     []Search   `json:"search"`
}

// Random results from wikipedia
type RandomResult struct {
	Batchcomplete string         `json:"batchcomplete"`
	Continue      RandomContinue `json:"continue"`
	Query         RandomQuery    `json:"query"`
}
type RandomContinue struct {
	Grncontinue string `json:"grncontinue"`
	Continue    string `json:"continue"`
}
type RandomPage struct {
	Pageid int    `json:"pageid"`
	Ns     int    `json:"ns"`
	Title  string `json:"title"`
}
type RandomQuery struct {
	Pages map[string]RandomPage `json:"pages"`
}
