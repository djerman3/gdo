// Package gdo provides the garage door webserver
// this file organizes the server part
package gdo

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/djerman3/gdo/data"
	"github.com/djerman3/gdo/staticdata"
	"github.com/gorilla/mux"

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
	tlsSrv  *http.Server
	srv     *http.Server
	Testing bool
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
	s := &Webserver{Testing: cfg.Testing}
	if false { //cfg.Testing {
		s.srv = &http.Server{
			Handler:        r,
			Addr:           fmt.Sprintf("%s:%d", cfg.Server.Addr, 8050),
			ReadTimeout:    time.Duration(cfg.Server.Timeout.Read) * time.Second,
			WriteTimeout:   time.Duration(cfg.Server.Timeout.Write) * time.Second,
			IdleTimeout:    time.Duration(cfg.Server.Timeout.Idle) * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	} else {
		cer, err := tls.LoadX509KeyPair(cfg.Server.TLS.CertFile, cfg.Server.TLS.KeyFile)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		tc := &tls.Config{Certificates: []tls.Certificate{cer}}

		s.tlsSrv = &http.Server{
			TLSConfig:      tc,
			Handler:        r,
			Addr:           cfg.Server.TLSAddr,
			ReadTimeout:    time.Duration(cfg.Server.Timeout.Read) * time.Second,
			WriteTimeout:   time.Duration(cfg.Server.Timeout.Write) * time.Second,
			IdleTimeout:    time.Duration(cfg.Server.Timeout.Idle) * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		m := mux.NewRouter()
		m.HandleFunc("/", redirectToHTTPS)
		s.srv = &http.Server{
			Handler:        m,
			Addr:           cfg.Server.Addr,
			ReadTimeout:    time.Duration(cfg.Server.Timeout.Read) * time.Second,
			WriteTimeout:   time.Duration(cfg.Server.Timeout.Write) * time.Second,
			IdleTimeout:    time.Duration(cfg.Server.Timeout.Idle) * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	}
	return s, nil
}

func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		// see comments below and consider the codes 308, 302, or 301
		http.StatusTemporaryRedirect)
}

//ListenAndServe to implement http.Server
func (s *Webserver) ListenAndServe() error {
	if s.tlsSrv != nil {
		go func() {
			log.Printf("listening TLS on %s\n", s.tlsSrv.Addr)
			err := s.tlsSrv.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
			}
		}()
	}
	log.Printf("listening HTTP on %s\n", s.srv.Addr)

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
