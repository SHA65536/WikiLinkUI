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
type Searchinfo struct {
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
	Searchinfo Searchinfo `json:"searchinfo"`
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
	Pageid               int       `json:"pageid"`
	Ns                   int       `json:"ns"`
	Title                string    `json:"title"`
	Contentmodel         string    `json:"contentmodel"`
	Pagelanguage         string    `json:"pagelanguage"`
	Pagelanguagehtmlcode string    `json:"pagelanguagehtmlcode"`
	Pagelanguagedir      string    `json:"pagelanguagedir"`
	Touched              time.Time `json:"touched"`
	Lastrevid            int       `json:"lastrevid"`
	Length               int       `json:"length"`
	Extract              string    `json:"extract"`
}
type RandomQuery struct {
	Pages map[string]RandomPage `json:"pages"`
}

// Search Response from the search route
type SearchResponse struct {
	Error  string          `json:"error,omitempty"`
	Result []SearchArticle `json:"result,omitempty"`
}

type SearchArticle struct {
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	Pageid  int    `json:"pageid"`
}

// Result Response from result route
type ResultResponse struct {
	Error        string   `json:"error,omitempty"`
	ResultIds    []uint32 `json:"ids,omitempty"`
	ResultTitles []string `json:"titles,omitempty"`
}
