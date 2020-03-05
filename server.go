package challenge

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

const (
	JWT_KEY		=	"test"
)

type App struct {
	name 		string
	wd			string
	siteInfo 	siteInfo
}
type siteInfo struct {
	Name 	  string
	User      string
	Content   interface{}
	Flag	  string
}


//func main(){
//	http.HandleFunc("/", handler)
//	http.HandleFunc("/test", handlerTest)
//	http.HandleFunc("/vue", handlerVue)
//	log.Fatal(http.ListenAndServe(":8080", nil))
//}

func NewApp(name string) *App{
	wd, _ := os.Getwd()
	app := &App{
		name:name,
		wd:wd[:len(wd)-4],
		siteInfo:siteInfo{
			Name: "MyChallenge",
		},
	}
	return app
}

func (app *App) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/", app.handleLogin())
	m.HandleFunc("/user", app.handleUser())
	m.HandleFunc("/logout", app.handleLogout())

	m.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir(app.wd + "/files/public"))))

	return m
}

func (app *App) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "jwt", MaxAge: -1})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *App) handleLogin() http.HandlerFunc {
	get := app.handleLoginGET()
	post := app.handleLoginPOST()

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			get(w, r)
			return

		case http.MethodPost:
			post(w, r)
			return
		}

		http.NotFound(w, r)
	}
}

func (app *App) handleLoginGET() http.HandlerFunc {
	tmpl, err := template.ParseFiles(
		app.wd+"/files/private/base.tmpl.html",
		app.wd+"/files/private/navbar.tmpl.html",
		app.wd+"/files/private/index.tmpl.html",
	)
	if err != nil {
		log.Println("error index tmpl: ", err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, app.siteInfo); err != nil {
			log.Println("template err login: ", err)
		}
	}
}


func (app *App) handleLoginPOST() http.HandlerFunc {
	tmpl, err := template.ParseFiles(
		app.wd+"/files/private/base.tmpl.html",
		app.wd+"/files/private/navbar.tmpl.html",
		app.wd+"/files/private/index.tmpl.html",
	)
	if err != nil {
		log.Println("error index tmpl: ", err)
	}

	type loginData struct {
		User       string
		LoginError string
	}

	readParams := func(r *http.Request) (loginData, error) {

		data := loginData{
			User: r.PostFormValue("user"),
		}

		if data.User == "" {
			return data, fmt.Errorf("User cannot be empty")
		}

		if data.User == "admin"{
			return data, fmt.Errorf("Cannot enter as admin")
		}

		app.siteInfo.User = data.User
		return data, nil
	}

	displayErr := func(w http.ResponseWriter, params loginData, err error) {
		tmplData := app.siteInfo
		params.LoginError = err.Error()
		tmplData.Content = params
		if err := tmpl.Execute(w, tmplData); err != nil {
			log.Println("template err login: ", err)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {

		params, err := readParams(r)
		if err != nil {
			displayErr(w, params, err)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": params.User,
		})
		tokenString, _ := token.SignedString([]byte(JWT_KEY))
		cookie := http.Cookie{
			Name:    "jwt",
			Value:   tokenString,
			MaxAge: int(time.Hour.Seconds()),
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/user", http.StatusSeeOther)
	}
}

func (app *App) handleUser() http.HandlerFunc {
	tmpl, err := template.ParseFiles(
		app.wd+"/files/private/base.tmpl.html",
		app.wd+"/files/private/navbar.tmpl.html",
		app.wd+"/files/private/user.tmpl.html",
	)
	if err != nil {
		log.Println("error login tmpl: ", err)
	}

	type cookieData struct {
		Cookie       	string
		CookieError 	string
		User			interface{}
	}

	readCookie := func(r *http.Request) (cookieData, error) {
		tokenString, err := r.Cookie("jwt");
		if err != nil {
			return  cookieData{}, fmt.Errorf("error retrieving token")
		}
		data := cookieData{
			Cookie:      tokenString.Value,
		}
		return data, nil
	}

	displayErr := func(w http.ResponseWriter, params cookieData, err error) {
		fmt.Println("error display")
		tmplData := app.siteInfo
		params.CookieError = err.Error()
		tmplData.Content = params
		if err := tmpl.Execute(w, tmplData); err != nil {
			log.Println("template err user display error: ", err)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {

		params, err := readCookie(r)
		if err != nil {
			displayErr(w, params, err)
			return
		}

		token, err := jwt.Parse(params.Cookie, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(JWT_KEY), nil
		})

		if err != nil {
			displayErr(w, params, err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims);
		if !ok || !token.Valid {
			displayErr(w, params, errors.New("Error while getting the user from the token"))
			return
		}

		tmplData := app.siteInfo
		params.User = claims["user"]
		tmplData.Content = params
		if params.User == "admin" {
			tmplData.Flag = "this_is_the_flag"
		}

		if err := tmpl.Execute(w, tmplData); err != nil {
			log.Println("template err user: ", err)
		}
	}
}