// Package gdo provides the garage door webserver
// this file organizes the server part
package gdo

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/djerman3/gdo/data"
	"github.com/djerman3/gdo/staticdata"
	"github.com/gorilla/mux"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	door              *gdoDoor
	homepageTpl       *template.Template
	googleOauthConfig *oauth2.Config
)

// Webserver is the target object for webserver functions
type Webserver struct {
	tlsSrv   *http.Server
	srv      *http.Server
	cachedir string
}

//NewWebserver gets a server ready to go
func NewWebserver(cfg *Config) (srv *Webserver, err error) {
	// open the door
	if door == nil {
		door = &gdoDoor{}
		door.Init(cfg)
	}
	// prepare google oauth2
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  cfg.Oauth2.Redirect,
		ClientID:     cfg.Oauth2.ClientID,
		ClientSecret: cfg.Oauth2.ClientSecret,
		Scopes:       cfg.Oauth2.Scopes,
		Endpoint:     google.Endpoint,
	}
	// prepare router
	r := newRouter()
	// allocate server
	s := &Webserver{cachedir: cfg.Server.Cachedir}
	var m *autocert.Manager
	if cfg.Testing {
		s.srv = &http.Server{
			Handler:        r,
			Addr:           fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port),
			ReadTimeout:    time.Duration(cfg.Server.Timeout.Read) * time.Second,
			WriteTimeout:   time.Duration(cfg.Server.Timeout.Write) * time.Second,
			IdleTimeout:    time.Duration(cfg.Server.Timeout.Idle) * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	} else {
		hostPolicy := func(ctx context.Context, host string) error {
			// Note: change to your real domain
			allowedHost := cfg.Server.Host
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
		}
		m = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(s.cachedir),
		}
		s.tlsSrv = &http.Server{
			Handler:        r,
			Addr:           ":443", // fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.TLSPort),
			ReadTimeout:    time.Duration(cfg.Server.Timeout.Read) * time.Second,
			WriteTimeout:   time.Duration(cfg.Server.Timeout.Write) * time.Second,
			IdleTimeout:    time.Duration(cfg.Server.Timeout.Idle) * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		s.tlsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		s.srv = &http.Server{
			Handler:        m.HTTPHandler(r),
			Addr:           ":80", //fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port),
			ReadTimeout:    time.Duration(cfg.Server.Timeout.Read) * time.Second,
			WriteTimeout:   time.Duration(cfg.Server.Timeout.Write) * time.Second,
			IdleTimeout:    time.Duration(cfg.Server.Timeout.Idle) * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	}
	return s, nil
}

//ListenAndServe to implement http.Server
func (s *Webserver) ListenAndServe() error {
	if s.tlsSrv != nil {
		go func() {
			err := s.tlsSrv.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
			}
		}()
	}
	return s.srv.ListenAndServe()
}

//NewGdoRouter produces the router with configured handlers
func newRouter() *mux.Router {
	r := mux.NewRouter()
	// handle home
	r.HandleFunc("/", HomeHandler)
	// handle api routes
	r.HandleFunc("/api", http.NotFound)
	r.HandleFunc("/api/click", ClickHandler)
	// OauthGoogle
	r.HandleFunc("/auth/login", oauthGoogleLogin)
	r.HandleFunc("/auth/logout", oauthLogout)
	r.HandleFunc("/auth/callback", oauthGoogleCallback)
	r.HandleFunc("/auth/redirect", oauthRedirect)

	// handle static routes
	r.Handle("/static", http.FileServer(staticdata.AssetFile()))

	// Load templates
	data, err := data.Asset("templates/index.html")
	if err != nil {
		log.Fatalf("Messed up Asset %v", err)
	}
	homepageTpl = template.Must(template.New("homepage_view").Parse(string(data)))
	return r
}

// Render a template, or server error.
func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		log.Printf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}

//HomeHandler handles the homepage and anything matching "/"
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	//push(w, "/static/style.css")
	//push(w, "/static/navigation_bar.css")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	//get content data
	state, asserted, err := door.ReadPin()
	if err != nil {
		log.Printf("Failed to get pi %v\n", err)
		http.Error(w, "Error fetching states!", http.StatusInternalServerError)
		return
	}
	//get auth status
	userName := GetLoggedInUser(w, r)
	login := false
	if userName == "" {
		login = true
	}

	if !asserted {
		state = "Probably " + state
	}
	fullData := map[string]interface{}{
		"time":     time.Now().Format(time.UnixDate),
		"state":    state,
		"username": userName,
		"login":    login,
	}

	render(w, r, homepageTpl, "homepage_view", fullData)
}

//ClickHandler handles the homepage and anything matching "/api/{owner}/{device}/{state}"
// Owner = cj or sj
// Device = ipad, laptop, or all
// State = lock or unlock
func ClickHandler(w http.ResponseWriter, r *http.Request) {
	//push(w, "/static/style.css")
	//push(w, "/static/navigation_bar.css")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if GetLoggedInUser(w, r) == "" {
		http.Error(w, "You aren't logged in!	", http.StatusUnauthorized)
		return
	}
	err := door.DoClick()
	if err != nil {
		log.Printf("Failed to click %v\n", err)
		http.Error(w, "Error trying to click, check the door!", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
