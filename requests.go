package wikilinkui

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const WikiSearchEndpoint = `https://he.wikipedia.org/w/api.php?action=query&list=search&srnamespace=0&srlimit=8&utf8=&format=json&srsearch=`
const WikiRandomEndpoint = `https://he.wikipedia.org/w/api.php?action=query&generator=random&grnnamespace=0&grnlimit=8&utf8=&format=json`

// Search searches wikipedia for articles
// consults the redis cache for quick retrieval
func (h *UIHandler) Search(query string) ([]string, error) {
	var res SearchResult
	var titles []string

	// Searching for cached result
	val, err := h.Redis.GetValue(url.QueryEscape(query))
	if err == nil {
		if err := json.Unmarshal([]byte(val), &titles); err == nil {
			return titles, nil
		}
	} else {
		h.Logger.Logger.Debug().Msgf("redis get err: %v", err)
	}

	// Searching wikipedia
	resp, err := h.Client.Get(WikiSearchEndpoint + url.QueryEscape(query))
	if err != nil {
		return nil, fmt.Errorf("failed reaching to wikipedia ")
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed decoding response from wikipedia! ")
	}

	titles = make([]string, len(res.Query.Search))
	for i := range titles {
		titles[i] = res.Query.Search[i].Title
	}

	// Updating cache
	data, _ := json.Marshal(titles)
	if err := h.Redis.PutValue(url.QueryEscape(query), string(data)); err != nil {
		h.Logger.Logger.Debug().Msgf("redis set err: %v", err)
	}

	return titles, nil
}

// Path searches LinkAPI for a path between articles
// consults the redis cache for quick retrieval
func (h *UIHandler) Path(src, dst string) (*FinalStruct, error) {
	var res FinalStruct
	query := fmt.Sprintf("start=%s&end=%s", url.QueryEscape(src), url.QueryEscape(dst))

	// Searching for cached result
	val, err := h.Redis.GetValue(query)
	if err == nil {
		// Returning cached result
		if err := json.Unmarshal([]byte(val), &res); err == nil {
			return &res, nil
		}
	} else {
		h.Logger.Logger.Debug().Msgf("redis get err: %v", err)
	}

	// Searching LinkAPI
	resp, err := h.Client.Get(fmt.Sprintf("http://%s/search?%s", h.LinkAPI, query))
	if err != nil {
		return nil, fmt.Errorf("failed reaching to linkapi")
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed decoding response from linkapi")
	}

	// Updating cache
	data, _ := json.Marshal(res)
	if err := h.Redis.PutValue(query, string(data)); err != nil {
		h.Logger.Logger.Debug().Msgf("redis set err: %v", err)
	}

	return &res, nil
}

// Random gets random articles from wikipedia
// does not consult cache
func (h *UIHandler) Random() ([]string, error) {
	var res RandomResult
	resp, err := h.Client.Get(WikiRandomEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed reaching to wikipedia ")
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed decoding response from wikipedia! ")
	}

	var titles = make([]string, 0, len(res.Query.Pages))
	for _, v := range res.Query.Pages {
		titles = append(titles, v.Title)
	}

	return titles, nil
}
