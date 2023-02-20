package wikilinkui

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

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
	Redis       *redis.Client
	Client      *http.Client
	Router      *chi.Mux
	ResultTempl *template.Template
	Logger      zerolog.Logger
}

func MakeUIHandler(locale string, api_url string, redis_addr string, logLevel zerolog.Level, writer io.Writer) (*UIHandler, error) {
	var ui = &UIHandler{
		Locale:      locale,
		LinkAPI:     api_url,
		Client:      http.DefaultClient,
		ResultTempl: template.Must(template.New("results").Parse(resultshtml)),
		Redis: redis.NewClient(&redis.Options{
			Addr:     redis_addr,
			Password: "",
			DB:       0,
		}),
	}

	// Creating logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logwriter := io.MultiWriter(os.Stdout, writer)
	ui.Logger = zerolog.New(logwriter).With().Str("service", "linkui").Timestamp().Logger().Level(logLevel)

	ui.Router = chi.NewRouter()
	// Main route
	ui.Router.Get("/", ui.MainRoute)
	// Search route
	ui.Router.Get("/search", ui.SearchRoute)
	// Random route
	ui.Router.Get("/random", ui.RandomRoute)
	// Result route
	ui.Router.Get("/result", ui.ResultRoute)
	// Health route
	ui.Router.Get("/health", ui.HealthRoute)

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

	ui.Logger.Debug().Msg("created router")

	return ui, nil
}

func (u *UIHandler) Serve(addr string) error {
	u.Logger.Info().Msgf("serving linkui on %s", addr)
	return http.ListenAndServe(addr, u.Router)
}

// Main webpage route
func (u *UIHandler) MainRoute(w http.ResponseWriter, r *http.Request) {
	sTime := time.Now()
	log := u.Logger.With().Str("ip", r.RemoteAddr).Str("path", r.URL.Path).Str("route", "index").Logger()
	w.Write(indexhtml)
	ReqLog(log, w, r, sTime, "success", zerolog.InfoLevel)
}

// Search for articles
func (u *UIHandler) SearchRoute(w http.ResponseWriter, r *http.Request) {
	sTime := time.Now()
	log := u.Logger.With().Str("ip", r.RemoteAddr).Str("path", r.URL.Path).Str("route", "search").Logger()
	var qres = &SearchResponse{}
	query := r.URL.Query().Get("q")
	if query == "" {
		qres.Error = "must have query parameter!"
		render.JSON(w, r, qres)
		ReqLog(log, w, r, sTime, "invalid parameters", zerolog.InfoLevel)
		return
	}

	wikires, err := u.WikiSearch(query)
	if err != nil {
		qres.Error = err.Error()
		render.JSON(w, r, qres)
		ReqLog(log, w, r, sTime, qres.Error, zerolog.WarnLevel)
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
	ReqLog(log, w, r, sTime, "success", zerolog.InfoLevel)
}

// Random Articles
func (u *UIHandler) RandomRoute(w http.ResponseWriter, r *http.Request) {
	sTime := time.Now()
	log := u.Logger.With().Str("ip", r.RemoteAddr).Str("path", r.URL.Path).Str("route", "random").Logger()

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
	ReqLog(log, w, r, sTime, "success", zerolog.InfoLevel)
}

// Results for path
func (u *UIHandler) ResultRoute(w http.ResponseWriter, r *http.Request) {
	sTime := time.Now()
	log := u.Logger.With().Str("ip", r.RemoteAddr).Str("path", r.URL.Path).Str("route", "result").Logger()

	var res = &ResultResponse{}
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	// Must have both params
	if end == "" || start == "" {
		res.Error = "must have start and end parameters!"
		render.JSON(w, r, res)
		ReqLog(log, w, r, sTime, "invalid parameters", zerolog.InfoLevel)
		return
	}
	res, err := u.PathSearch(start, end)
	if err != nil {
		res.Error = err.Error()
		render.JSON(w, r, res)
		ReqLog(log, w, r, sTime, res.Error, zerolog.WarnLevel)
		return
	}

	u.ResultTempl.Execute(w, res)
	ReqLog(log, w, r, sTime, "success", zerolog.InfoLevel)
}

// HealthRoute for health checking purposes
func (u *UIHandler) HealthRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
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

func timeToMs(t time.Duration) string {
	return fmt.Sprintf("%dms", t/time.Millisecond)
}

func ReqLog(Logger zerolog.Logger, w http.ResponseWriter, r *http.Request, s time.Time, msg string, level zerolog.Level) {
	Logger.WithLevel(level).Str("took", timeToMs(time.Since(s))).Msg(msg)
}
