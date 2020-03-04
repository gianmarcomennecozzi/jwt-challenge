package challenge

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

const (
	JWT_KEY		=	"haaukins"
)



//func addCookie(w http.ResponseWriter, name string, value string) {
//	expire := time.Now().AddDate(0, 0, 1)
//	cookie := http.Cookie{
//		Name:    name,
//		Value:   value,
//		Expires: expire,
//	}
//	http.SetCookie(w, &cookie)
//}
//
//func handler(w http.ResponseWriter, req *http.Request) {
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"foo": "bar",
//		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
//	})
//
//	// Sign and get the complete encoded token as a string using the secret
//	tokenString, _ := token.SignedString([]byte(JWT_KEY))
//	fmt.Println("browser= "+tokenString)
//	addCookie(w, "menne", tokenString)
//	http.Redirect(w, req, "/test", http.StatusFound)
//}
//
//func handlerTest(w http.ResponseWriter, req *http.Request) {
//	tokenString, err := req.Cookie("menne");
//	if err != nil {
//		fmt.Errorf("error retrieving token")
//	}
//
//	token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
//		// Don't forget to validate the alg is what you expect:
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
//		}
//
//		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
//		return []byte(JWT_KEY), nil
//	})
//
//	claims, ok := token.Claims.(jwt.MapClaims);
//	if !ok || !token.Valid {
//		fmt.Printf("errror")
//	} else {
//		fmt.Fprint(w, claims["foo"])
//	}
//}


//func main(){
//	http.HandleFunc("/", handler)
//	http.HandleFunc("/test", handlerTest)
//	http.HandleFunc("/vue", handlerVue)
//	log.Fatal(http.ListenAndServe(":8080", nil))
//}

type App struct {
	name string
}
//func main(){
//	http.HandleFunc("/", handler)
//	http.HandleFunc("/test", handlerTest)
//	http.HandleFunc("/vue", handlerVue)
//	log.Fatal(http.ListenAndServe(":8080", nil))
//}

func NewApp(name string) *App{
	app := &App{name:name}
	return app
}

func (app *App) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/", app.handleIndex())

	return m
}

func (app *App) handleIndex() http.HandlerFunc  {
	tmpl, err := template.ParseFiles(
		"/home/gian/go/src/jwt-challenge/files/private/base.tmpl.html",
	)
	if err != nil {
		log.Println("error index tmpl: ", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			fmt.Println("error here")
			return
		}

		data := "test"
		if err := tmpl.Execute(w, data); err != nil {
			fmt.Println("dwadwa")
			log.Println("template err index: ", err)
		}
	}
}