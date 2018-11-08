package routes

import (
	"../middleware"
	"../models"
	"../sessions"
	"../utils"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter()

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static",http.FileServer(http.Dir("./static"))))

	router.HandleFunc("/", middleware.AuthRequired(indexGetHandler)).Methods("GET")
	router.HandleFunc("/", middleware.AuthRequired(indexPostHandler)).Methods("POST")

	router.HandleFunc("/login", loginGetHandler).Methods("GET")
	router.HandleFunc("/login", loginPostHandler).Methods("POST")

	router.HandleFunc("/register", registerGetHandler).Methods("GET")
	router.HandleFunc("/register", registerPostHandler).Methods("POST")

	return router
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	comments, err := models.GetAllComments()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	utils.ExecuteTemplate(w, "index.html", comments)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment")

	if comment == "" {
		http.Redirect(w, r, "/",302)
		return
	}

	err := models.AddComment(comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
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

	err := models.AuthUser(username, password)

	if err != nil {
		utils.ExecuteTemplate(w, "login.html", err)
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
	err := models.CreateUser(username, password)

	if err != nil {
		utils.ExecuteTemplate(w, "register.html", err)
	}

	http.Redirect(w, r, "/login", 302)	
}