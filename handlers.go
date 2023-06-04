package wikilinkui

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type SourceStruct struct {
	Error  string
	Result []string
}

// SourceRoute is the index of the applicatoin
func (h *UIHandler) SourceRoute(w http.ResponseWriter, r *http.Request) {
	var res SourceStruct
	var err error

	// Set up log
	sTime := time.Now()
	log := h.Logger.Logger.With().Str("ip", r.RemoteAddr).Str("path", r.URL.Path).Str("route", "index").Logger()

	// Get query params
	randSearch := r.URL.Query().Get("random")
	srcSearch := r.URL.Query().Get("srcSearch")

	if randSearch == "true" {
		// Random articles
		res.Result, err = h.Random()
		if err != nil {
			res.Error = err.Error()
			ReqLog(log, w, r, sTime, res.Error, zerolog.InfoLevel)
		} else if len(res.Result) == 0 {
			res.Error = "No matches found from random"
			ReqLog(log, w, r, sTime, res.Error, zerolog.InfoLevel)
		}
	} else if srcSearch != "" {
		// Search articles
		res.Result, err = h.Search(srcSearch)
		if err != nil {
			res.Error = err.Error()
			ReqLog(log, w, r, sTime, res.Error, zerolog.InfoLevel)
		} else if len(res.Result) == 0 {
			res.Error = fmt.Sprintf("No matches found for: %s", srcSearch)
			ReqLog(log, w, r, sTime, res.Error, zerolog.InfoLevel)
		}
	}

	h.Templates["source"].Execute(w, res)
	if res.Error == "" {
		ReqLog(log, w, r, sTime, "success", zerolog.InfoLevel)
	}
}

type DestStruct struct {
	SrcFinal string
	Error    string
	Result   []string
}

func (h *UIHandler) DestRoute(w http.ResponseWriter, r *http.Request) {
	var res DestStruct
	var err error

	// Set up log
	sTime := time.Now()
	log := h.Logger.Logger.With().Str("ip", r.RemoteAddr).Str("path", r.URL.Path).Str("route", "dest").Logger()

	// Get query params
	randSearch := r.URL.Query().Get("random")
	res.SrcFinal = r.URL.Query().Get("srcFinal")
	dstSearch := r.URL.Query().Get("dstSearch")

	if randSearch == "true" {
		// Random articles
		res.Result, err = h.Random()
		if err != nil {
			res.Error = err.Error()
			ReqLog(log, w, r, sTime, res.Error, zerolog.InfoLevel)
		} else if len(res.Result) == 0 {
			res.Error = "No matches found from random"
			ReqLog(log, w, r, sTime, res.Error, zerolog.InfoLevel)
		}
	} else if dstSearch != "" {
		// Search articles
		res.Result, err = h.Search(dstSearch)
		if err != nil {
			res.Error = err.Error()
			ReqLog(log, w, r, sTime, res.Error, zerolog.InfoLevel)
		} else if len(res.Result) == 0 {
			res.Error = fmt.Sprintf("No matches found for: %s", dstSearch)
			ReqLog(log, w, r, sTime, res.Error, zerolog.InfoLevel)
		}
	}

	h.Templates["destination"].Execute(w, res)
	if res.Error == "" {
		ReqLog(log, w, r, sTime, "success", zerolog.InfoLevel)
	}
}

type FinalStruct struct {
	Error        string   `json:"error,omitempty"`
	ResultIds    []uint32 `json:"ids,omitempty"`
	ResultTitles []string `json:"titles,omitempty"`
}

// FinalRoute searches the path between two articles
func (h *UIHandler) FinalRoute(w http.ResponseWriter, r *http.Request) {
	var res *FinalStruct
	var err error

	// Set up log
	sTime := time.Now()
	log := h.Logger.Logger.With().Str("ip", r.RemoteAddr).Str("path", r.URL.Path).Str("route", "final").Logger()

	// Get query params
	src := r.URL.Query().Get("src")
	dst := r.URL.Query().Get("dst")
	if src != "" && dst != "" {
		res, err = h.Path(src, dst)
		if err != nil {
			res = &FinalStruct{Error: err.Error()}
			ReqLog(log, w, r, sTime, res.Error, zerolog.InfoLevel)
		} else if len(res.ResultTitles) == 0 {
			res = &FinalStruct{Error: "Path not found!"}
			ReqLog(log, w, r, sTime, res.Error, zerolog.InfoLevel)
		}
	} else {
		res = &FinalStruct{Error: "invalid parameters"}
		ReqLog(log, w, r, sTime, res.Error, zerolog.InfoLevel)
	}

	h.Templates["final"].Execute(w, *res)
	if res.Error == "" {
		ReqLog(log, w, r, sTime, "success", zerolog.InfoLevel)
	}
}
