package main

import (
	"crypto/rsa"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// Application configuration
type Application struct {
	Name        string `json:"name"`
	RedirectURL string `json:"redirect_url"`
	WebHookURL  string `json:"webhook_url"`
}

// Applications configuration
type Applications []*Application

// Lookup Applications by name
func (a Applications) Lookup(name string) (*Application, bool) {
	for _, app := range applications {
		if app.Name == name {
			return app, true
		}
	}
	return nil, false
}

// WebHooks gets all the hook urls
func (a Applications) WebHooks() (hooks []string) {
	for _, app := range a {
		hooks = append(hooks, app.WebHookURL)
	}
	return
}

// User info
type User struct {
	ID     int64    `json:"id"`
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
}

var (
	signingKey   *rsa.PrivateKey
	applications Applications
)

func init() {

	// read config file
	f, err := os.Open("applications.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(f).Decode(&applications); err != nil {
		log.Fatal(err)
	}

	// read key file
	data, err := ioutil.ReadFile("privkey.pem")
	if err != nil {
		log.Fatal(err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		log.Fatal(err)
	}
	signingKey = key
}

// JWT converts the user to a JSON Web Token
func (u *User) JWT() (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = map[string]interface{}{
		"id":     u.ID,
		"name":   u.Name,
		"groups": u.Groups,
	}
	return token.SignedString(signingKey)
}

// Authenticate a user
func Authenticate(username string, password string) (*User, error) {
	user := &User{0, "Ilia Choly", []string{"dev", "admin"}}
	return user, nil
}

func main() {

	tmpl := template.Must(template.ParseGlob("*.tmpl"))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		// only accept GET requests
		if r.Method != "GET" {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var (
			redirect = r.URL.Query().Get("redirect")
			appname  = r.URL.Query().Get("appname")
		)

		// render template
		data := struct {
			Redirect string
			AppName  string
		}{redirect, appname}

		if err := tmpl.ExecuteTemplate(w, "login.html.tmpl", data); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

	})

	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {

		// only accept POST requests
		if r.Method != "POST" {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		// authenticate
		var (
			username = r.PostFormValue("username")
			password = r.PostFormValue("password")
			redirect = r.PostFormValue("redirect")
			appname  = r.PostFormValue("appname")
		)
		user, err := Authenticate(username, password)
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		// generate JSON web token
		token, err := user.JWT()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		var hooks []string
		if app, ok := applications.Lookup(appname); ok {
			// If an application is specified, only invoke its hook
			if redirect == "" {
				redirect = app.RedirectURL
			}
			hooks = []string{app.WebHookURL}
		} else {
			// if no application is specified, invoke all the hooks
			hooks = applications.WebHooks()
		}

		data := struct {
			JWT      string
			WebHooks []string
			Redirect string
		}{token, hooks, redirect}

		if err := tmpl.ExecuteTemplate(w, "postlogin.html.tmpl", &data); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

	})

	http.HandleFunc("/api/wafer_hook", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "jwt",
			Value: r.URL.Query().Get("jwt"),
			Path:  "/",
		})

		w.WriteHeader(200)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
