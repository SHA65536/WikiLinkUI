package wikilinkui

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const WikiSearchEndpoint = `https://he.wikipedia.org/w/api.php?action=query&list=search&srnamespace=0&srlimit=5&prop=info&utf8=&format=json&origin=*&srsearch=`
const WikiRandomEndpoint = "https://he.wikipedia.org/w/api.php?action=query&generator=random&grnnamespace=0&grnlimit=1&prop=info|extracts&exlimit=1&explaintext=true&exsentences=1&utf8=&format=json&origin=*"

func (u *UIHandler) WikiSearch(query string) (*SearchResult, error) {
	var res = &SearchResult{}

	// Searching for cached result
	if val, err := u.Redis.GetValue(url.QueryEscape(query)); err == nil {
		if err := json.Unmarshal([]byte(val), res); err == nil {
			return res, nil
		}
	}

	// Requesting search from wikipedia if no cache
	resp, err := u.Client.Get(WikiSearchEndpoint + url.QueryEscape(query))
	if err != nil {
		return nil, fmt.Errorf("failed reaching to wikipedia! ")
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return nil, fmt.Errorf("failed decoding response from wikipedia! ")
	}

	// Updating cache
	data, _ := json.Marshal(res)
	u.Redis.PutValue(url.QueryEscape(query), string(data))
	return res, nil
}

func (u *UIHandler) PathSearch(src, dst string) (*ResultResponse, error) {
	var res = &ResultResponse{}
	query := fmt.Sprintf("start=%s&end=%s", url.QueryEscape(src), url.QueryEscape(dst))

	// Searching for cached result
	if val, err := u.Redis.GetValue(query); err == nil {
		// Returning cached result
		if err := json.Unmarshal([]byte(val), res); err == nil {
			return res, nil
		}
	}

	// Requesting search from linkapi if no cache
	resp, err := u.Client.Get(fmt.Sprintf("http://%s/search?%s", u.LinkAPI, query))
	if err != nil {
		return res, fmt.Errorf("failed getting response from link api! ")
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return res, fmt.Errorf("failed decoding response from link api! ")
	}

	// Updating cache
	data, _ := json.Marshal(res)
	u.Redis.PutValue(url.QueryEscape(query), string(data))

	return res, nil
}
