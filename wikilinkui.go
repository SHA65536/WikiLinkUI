package wikilinkui

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const WikiSearchEndpoint = `https://he.wikipedia.org/w/api.php?action=query&list=search&srnamespace=0&srlimit=5&prop=info&utf8=&format=json&origin=*&srsearch=`
const WikiRandomEndpoint = "https://he.wikipedia.org/w/api.php?action=query&generator=random&grnnamespace=0&grnlimit=1&prop=info|extracts&exlimit=1&explaintext=true&exsentences=1&utf8=&format=json&origin=*"

//go:embed static/indexstyles.css
var indexstyles []byte

//go:embed static/resultstyles.css
var resultstyles []byte

//go:embed static/main.js
var mainjs []byte

//go:embed static/index.html
var indexhtml []byte

//go:embed static/results.html
var resultshtml string

type UIHandler struct {
	Locale      string
	LinkAPI     string
	Client      *http.Client
	Router      *chi.Mux
	ResultTempl *template.Template
}

func MakeUIHandler(locale string, api_url string) (*UIHandler, error) {
	var ui = &UIHandler{
		Locale:      locale,
		LinkAPI:     api_url,
		Client:      http.DefaultClient,
		ResultTempl: template.Must(template.New("results").Parse(resultshtml)),
	}

	ui.Router = chi.NewRouter()
	ui.Router.Use(middleware.Logger)
	// Main route
	ui.Router.Get("/", ui.MainRoute)
	// Search route
	ui.Router.Get("/search", ui.SearchRoute)
	// Random route
	ui.Router.Get("/random", ui.RandomRoute)
	// Result route
	ui.Router.Get("/result", ui.ResultRoute)

	// Static files
	ui.Router.Get("/main.js", func(w http.ResponseWriter, r *http.Request) {
		w.Write(mainjs)
	})
	ui.Router.Get("/indexstyles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Write(indexstyles)
	})
	ui.Router.Get("/resultstyles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Write(resultstyles)
	})
	return ui, nil
}

func (u *UIHandler) Serve(addr string) error {
	return http.ListenAndServe(addr, u.Router)
}

// Main webpage route
func (u *UIHandler) MainRoute(w http.ResponseWriter, r *http.Request) {
	w.Write(indexhtml)
}

// Search for articles
func (u *UIHandler) SearchRoute(w http.ResponseWriter, r *http.Request) {
	var wikires = &SearchResult{}
	var qres = &SearchResponse{}
	query := r.URL.Query().Get("q")
	if query == "" {
		qres.Error = "must have query parameter!"
		render.JSON(w, r, qres)
		return
	}
	resp, err := u.Client.Get(WikiSearchEndpoint + url.QueryEscape(query))
	if err != nil {
		qres.Error = "failed reaching to wikipedia!"
		render.JSON(w, r, qres)
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(wikires); err != nil {
		qres.Error = "failed decoding response from wikipedia!"
		render.JSON(w, r, qres)
		return
	}

	for _, s := range wikires.Query.Search {
		qres.Result = append(qres.Result, SearchArticle{
			Title:   s.Title,
			Snippet: s.Snippet,
			Pageid:  s.Pageid,
		})
	}

	render.JSON(w, r, qres)
}

// Random Articles
func (u *UIHandler) RandomRoute(w http.ResponseWriter, r *http.Request) {
	var res = &SearchResponse{}
	var links = make([]SearchArticle, 10)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		i := i
		go func() {
			defer wg.Done()
			u.GetRandom(links, i)
		}()
	}

	wg.Wait()

	res.Result = links

	render.JSON(w, r, res)
}

// Results for path
func (u *UIHandler) ResultRoute(w http.ResponseWriter, r *http.Request) {
	var res = &ResultResponse{}
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	// Must have both params
	if end == "" || start == "" {
		res.Error = "must have start and end parameters!"
		render.JSON(w, r, res)
		return
	}
	resp, err := u.Client.Get(fmt.Sprintf("http://%s/search?start=%s&end=%s", u.LinkAPI, url.QueryEscape(start), url.QueryEscape(end)))
	if err != nil {
		res.Error = "failed getting response from link api!"
		render.JSON(w, r, res)
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		res.Error = "failed decoding response from link api!"
		render.JSON(w, r, res)
		return
	}

	u.ResultTempl.Execute(w, res)
}

func (u *UIHandler) GetRandom(out []SearchArticle, idx int) {
	var res = &RandomResult{}
	var page string
	resp, err := u.Client.Get(WikiRandomEndpoint)
	if err != nil {
		out[idx] = SearchArticle{
			Title:   "Error getting random!",
			Snippet: "Error getting random!",
			Pageid:  -1,
		}
		return
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		out[idx] = SearchArticle{
			Title:   "Error getting random!",
			Snippet: "Error getting random!",
			Pageid:  -1,
		}
		return
	}
	for page = range res.Query.Pages {
	}
	out[idx] = SearchArticle{
		Title:   res.Query.Pages[page].Title,
		Snippet: res.Query.Pages[page].Extract,
		Pageid:  res.Query.Pages[page].Pageid,
	}
}