package main

import (
	"crypto/rsa"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// WebHooks is a list
var WebHooks = []string{
	"http://sub1.domain.com/wafer_hook",
	"http://127.0.0.1:8080/api/wafer_hook",
}

// User info
type User struct {
	ID     int64    `json:"id"`
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
}

var signingKey *rsa.PrivateKey

func init() {
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

func getRedirectURL(r *http.Request) (redirect string) {
	if param, ok := r.URL.Query()["redirect"]; ok && len(param) > 0 {
		redirect = param[0]
	} else {
		redirect = r.Referer()
	}
	return
}

func main() {

	tmpl := template.Must(template.ParseGlob("*.tmpl"))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		// only accept GET requests
		if r.Method != "GET" {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		// render template
		if err := tmpl.ExecuteTemplate(w, "login.html.tmpl", struct {
			Redirect string
		}{getRedirectURL(r)}); err != nil {
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

		log.Print(redirect)

		if err := tmpl.ExecuteTemplate(w, "postlogin.html.tmpl", struct {
			JWT      string
			WebHooks []string
			Redirect string
		}{
			JWT:      token,
			WebHooks: WebHooks,
			Redirect: redirect,
		}); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

	})

	http.HandleFunc("/api/wafer_hook", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		if params, ok := r.URL.Query()["jwt"]; ok {
			http.SetCookie(w, &http.Cookie{
				Name:  "jwt",
				Value: params[0],
				Path:  "/",
			})
		}

		w.WriteHeader(200)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
