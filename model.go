package main

import "github.com/dgrijalva/jwt-go"

// Application configuration
type Application struct {
	Name          string `json:"name"`
	RedirectURL   string `json:"redirect_url"`
	ForceRedirect bool   `json:"force_redirect"`
	WebHookURL    string `json:"webhook_url"`
}

// Applications configuration
type Applications []*Application

// Config is the configuration
type Config struct {
	ListenAddr   string       `json:"listen_addr"`
	Key          string       `json:"key"`
	Applications Applications `json:"applications"`
}

// Lookup Applications by name
func (a Applications) Lookup(name string) (*Application, bool) {
	for _, app := range a {
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

// JWT converts the user to a JSON Web Token
func (u *User) JWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = map[string]interface{}{
		"id":     u.ID,
		"name":   u.Name,
		"groups": u.Groups,
	}
	return token.SignedString(signingKey)
}
