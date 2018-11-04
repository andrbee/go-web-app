package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
	"html/template"
)

var templates *template.Template
var client *redis.Client
var store = sessions.NewCookieStore([]byte("Qwe12345678"))

func main() {

	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		})

	templates = template.Must(template.ParseGlob("templates/*.html"))

	r := mux.NewRouter();	

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static",http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/", indexPostHandler).Methods("POST")

	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")

	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")

	http.Handle("/",r)
	http.ListenAndServe(":8000",nil)

}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionId")
	_, ok := session.Values["user"]
	if !ok {
		http.Redirect(w, r, "/login", 302)
		return
	}

	comments, err := client.LRange("comments",0,10).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)		
		w.Write([]byte("Internal Server Error"))
		return
	}
	templates.ExecuteTemplate(w, "index.html",comments)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment")
	if comment != "" {
		client.LPush("comments", comment)	
	}	
	http.Redirect(w, r, "/",302)
}


func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html",nil)	
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {	
	r.ParseForm()

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	hashPass, err := client.Get("user:" + username).Bytes()
	if err != nil {
		templates.ExecuteTemplate(w, "login.html", "User is not found")
		return
	}

	err = bcrypt.CompareHashAndPassword(hashPass, []byte (password))

	if err != nil {
		templates.ExecuteTemplate(w, "login.html", "User hasn't access")
		return
	}

	session, _ := store.Get(r, "sessionId")
	session.Values["user"] = username
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "register.html",nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm();

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	error := ""
	if username == "" || password == "" {
		error = "Login Invalid !"
	}

	if error != "" {
		templates.ExecuteTemplate(w, "register.html", error)	
		return
	}
	cost := bcrypt.DefaultCost
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	} 

	client.Set("user:" + username, hashPass, 0)
	http.Redirect(w, r, "/login", 302)	
}