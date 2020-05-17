package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

var redirectURI = "http://localhost:8000/auth/razer/callback"
var clientID = "<your client id>"
var clientSecret = "<your client secret>"

// 1) Your user goes to http://yourdomain:8000/auth/razer/login
// 2) User will need to login user Razer ID
// 3) User will be directed to your callback path http://yourdomain:8000/auth/razer/callback for further processing
// 4) Your endpoint will process the callback with the code provided by Razer ID to fetch user details
// 5) Your post authorization flow

func RedirectToRazer(w http.ResponseWriter, r *http.Request) {
	// This method redirects user to the Razer ID login page.
	if (*r).Method == "OPTIONS" {
		return
	}
	// The scopes used here are (openid,email,profile). For more scopes, please refer to the Razre ID official documentation.
	u := "https://oauth2.razer.com/authorize_openid?response_type=code&l=en&scope=openid+email+profile&client_id=" + clientID + "&state=login&redirect_uri=" + redirectURI
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	// This url is redirected by Razer ID. Razer ID will give you a 302 for your browser to redirect.
	if (*r).Method == "OPTIONS" {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	code := r.URL.Query()["code"][0]
	fmt.Println(code)
	selection := "https://oauth2.razer.com/token"
	client := http.Client{}
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	req, _ := http.NewRequest("POST", selection, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)

	}
	if resp.StatusCode == http.StatusBadRequest {

	}
	if resp.Body != nil {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		// For debugging purpose
		fmt.Println(result["access_token"])
		fmt.Println(result["expires_in"])
		fmt.Println(result["id_token"])
		fmt.Println(result["scope"])
		fmt.Println(result["token_type"])

		selection := "https://oauth2.razer.com/userinfo"
		client := http.Client{}
		ureq, _ := http.NewRequest("GET", selection, nil)
		ureq.Header.Set("Authorization", "Bearer "+result["access_token"].(string))
		uresp, err := client.Do(ureq)
		if err != nil {
			panic(err)

		}
		if uresp.StatusCode == http.StatusBadRequest {

		}
		var uresult map[string]interface{}
		json.NewDecoder(uresp.Body).Decode(&uresult)
		// You might want to generate your JWT token here after assertion.
		u := "http://localhost/your/path/after/assertion"
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
		return

	}

}

func ServeService() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/auth/razer/login", RedirectToRazer)
	router.HandleFunc("/auth/razer/callback", Callback)
	return router
}

func main() {
	http.ListenAndServe(":8000", ServeService())
}
