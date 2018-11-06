package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"../middleware"
	"../sessions"
	"../models"
	"../utils"
)

func NewRouter() *mux.Router {

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static",http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/", middleware.AuthRequired(indexGetHandler)).Methods("GET")
	r.HandleFunc("/", middleware.AuthRequired(indexPostHandler)).Methods("POST")

	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")

	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")

	return r
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	comments, err := models.Client.LRange("comments",0,10).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)		
		w.Write([]byte("Internal Server Error"))
		return
	}
	utils.ExecuteTemplate(w, "index.html",comments)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment")
	if comment != "" {
		models.Client.LPush("comments", comment)	
	}	
	http.Redirect(w, r, "/",302)
}


func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "login.html",nil)	
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {	
	r.ParseForm()

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	hashPass, err := models.Client.Get("user:" + username).Bytes()
	if err != nil {
		utils.ExecuteTemplate(w, "login.html", "User is not found")
		return
	}

	err = bcrypt.CompareHashAndPassword(hashPass, []byte (password))

	if err != nil {
		utils.ExecuteTemplate(w, "login.html", "User hasn't access")
		return
	}

	session, _ := sessions.Store.Get(r, "sessionId")
	session.Values["user"] = username
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request){
	utils.ExecuteTemplate(w, "register.html",nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	error := ""
	if username == "" || password == "" {
		error = "Login Invalid !"
	}

	if error != "" {
		utils.ExecuteTemplate(w, "register.html", error)	
		return
	}
	cost := bcrypt.DefaultCost
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	} 

	models.Client.Set("user:" + username, hashPass, 0)
	http.Redirect(w, r, "/login", 302)	
}