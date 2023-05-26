package wikilinkui

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

//go:embed content
var contentFiles embed.FS

type UIHandler struct {
	Locale    string
	LinkAPI   string
	Redis     *RedisHandler
	Client    *http.Client
	Router    *chi.Mux
	Templates map[string]*template.Template
	Logger    zerolog.Logger
}

// ("heb", apiAddr, rAddr, rRole, vAddr, vRegion, vRole, level, logf)
func MakeUIHandler(locale, apiAddr, rAddr, vAddr, vRole string, logLevel zerolog.Level, writer io.Writer) (*UIHandler, error) {
	var ui = &UIHandler{
		Locale:  locale,
		LinkAPI: apiAddr,
		Client:  http.DefaultClient,
	}

	// Loading Templates
	ui.Templates = map[string]*template.Template{
		"source": template.Must(template.ParseFS(contentFiles,
			"content/templates/source.html",
			"content/templates/rules.html",
			"content/templates/common.html")),
		"destination": template.Must(template.ParseFS(contentFiles,
			"content/templates/destination.html",
			"content/templates/rules.html",
			"content/templates/common.html")),
		"final": template.Must(template.ParseFS(contentFiles,
			"content/templates/final.html",
			"content/templates/rules.html",
			"content/templates/common.html")),
	}

	// Creating logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logwriter := io.MultiWriter(os.Stdout, writer)
	ui.Logger = zerolog.New(logwriter).With().Str("service", "linkui").Timestamp().Logger().Level(logLevel)

	// Creating redis handler
	redis, err := MakeRedisHandler(rAddr, vAddr, vRole, ui.Logger)
	if err != nil {
		return nil, err
	}
	ui.Redis = redis

	// Setting up router
	ui.Router = chi.NewRouter()

	// Static files
	staticContent, _ := fs.Sub(contentFiles, "content/static")
	fs := http.FileServer(http.FS(staticContent))
	ui.Router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// App Routes
	ui.Router.Get("/", ui.SourceRoute)
	ui.Router.Get("/dest", ui.DestRoute)
	ui.Router.Get("/final", ui.FinalRoute)
	ui.Router.Get("/health", ui.HealthRoute)

	ui.Logger.Debug().Msg("created router")

	return ui, nil
}

func (u *UIHandler) Serve(addr string) error {
	u.Logger.Info().Msgf("serving linkui v2 on %s", addr)
	return http.ListenAndServe(addr, u.Router)
}

// HealthRoute for health checking purposes
func (u *UIHandler) HealthRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func timeToMs(t time.Duration) string {
	return fmt.Sprintf("%dms", t/time.Millisecond)
}

func ReqLog(Logger zerolog.Logger, w http.ResponseWriter, r *http.Request, s time.Time, msg string, level zerolog.Level) {
	Logger.WithLevel(level).Str("took", timeToMs(time.Since(s))).Msg(msg)
}
