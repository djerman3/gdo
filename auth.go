// Package gdo provides the garage door webserver
// this file organizes the server side authentication part
package gdo

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// RandomString produces a base64 string (websafe?) from n random bytes
func RandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("RandomString: Not Enough Random %v", err)
	}
	s := base64.RawStdEncoding.EncodeToString(b)
	return s
}

// Session state string for each request (note, for each? When dose this change)
var oauthStateString string
var lastOauthStateString string
var oauthStateStringTime time.Time

// extra?  Here's a "go home func"
func oauthRedirect(w http.ResponseWriter, r *http.Request) {

	HomeHandler(w, r)
}

// OAuthGoogleLogin does the google login redirect.
func oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}
func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

type GoogleUserProfile struct {
	Email    string `json:"email"`
	Verified bool   `json:"verified_email"`
	ID       string `json:"id"`
	Picture  string `json:"picture"`
}

type OauthLogin struct {
	Cookie      http.Cookie
	NextHandler http.Handler
}

func NewLoginHandler(next http.Handler) OauthLogin {
	h := OauthLogin{
		Cookie: http.Cookie{
			Name:   "access_token",
			Value:  "",
			Domain: "batcave.jerman.info",
			Path:   "/",
			MaxAge: 30000,
		},
		NextHandler: next,
	}
	return h
}
func NewLogoutHandler(next http.Handler) OauthLogin {
	h := OauthLogin{
		Cookie: http.Cookie{
			Name:   "access_token",
			Value:  "",
			Domain: "batcave.jerman.info",
			Path:   "/",
			MaxAge: -1,
		},
		NextHandler: next,
	}
	return h
}

//Make OauthLogin a http.Handler
func (l OauthLogin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("access_token")
	if err != nil {
		if err.Error() != "http: named cookie not present" {
			log.Printf("Error getting auth cookie:%v\n", err)
		}
		log.Printf("Error getting auth cookie:%v\n", err)
		c = &l.Cookie
		r.AddCookie(c)
	}
	log.Printf("Loin auth cookie:%#v\n", *c)

	//set cookie
	if l.Cookie.MaxAge < c.MaxAge {
		c.MaxAge = l.Cookie.MaxAge
	}
	http.SetCookie(w, c)
	r.AddCookie(c)
	l.NextHandler.ServeHTTP(w, r)
}

// GoogleLogin sends us here with a form to catch the userinfo
func oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("Auth Content: %#v\n", string(data))

	u := GoogleUserProfile{}
	json.Unmarshal(data, &u)
	okUser := false
	for _, au := range cfg.Oauth2.Access {
		if au.GoogleID == u.ID {
			okUser = true
		}
		log.Printf("Compare:%v : %v\n", au.GoogleID, u.ID)
	}
	if okUser {
		cookie := http.Cookie{
			Name:    "access_token",
			Domain:  "batcave.jerman.info",
			Value:   u.Email,
			Path:    "/",
			Expires: time.Now().AddDate(0, 0, 1),
		}

		http.SetCookie(w, &cookie)

		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	oauthLogout(w, r)
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func oauthLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		log.Println(err.Error())
	} else {
		cookie = &http.Cookie{
			Name:   "access_token",
			Value:  "",
			Domain: "batcave.jerman.info",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
		log.Println("logout")
	}
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
func getUserDataFromGoogle(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}

//GetLoggedInUser returns the logged in user or cleans up if youve expired
func GetLoggedInUser(w http.ResponseWriter, r *http.Request) string {
	var s string
	c, err := r.Cookie("access_token")
	if err != nil {
		log.Printf("Error foo getting auth cookie:%v\n", err)
	} else {
		log.Printf("cookie:%#v\n", *c)

		if err != nil {
			log.Printf("Error bar getting auth cookie:%v\n", err)
		} else {
			s = c.Value
		}
		//http.SetCookie(w, c)
	}
	return s
}
