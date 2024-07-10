package main

import (
	"fmt"
	"net/http"

	f "forum/server"
)

func main() {
	f.InitiateDatabase()
	f.FillDatabase()

	static := http.FileServer(http.Dir("ui"))
	http.Handle("/static/", http.StripPrefix("/static/", static))
	http.HandleFunc("/static/css/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		static.ServeHTTP(w, r)
	})

	http.HandleFunc("/register", f.RegisterHandler)
	http.HandleFunc("/login", f.LoginHandler)
	http.HandleFunc("/index", f.HandleMainPage)
	http.HandleFunc("/create-post", f.HandleCreatePost)
	http.HandleFunc("/submit_post", f.HandlerAddPost)
	/* http.HandleFunc("/AddLike", f.AddLike)
	http.HandleFunc("/AddDislike", f.AddDislike) */
	http.HandleFunc("/ShowComments", f.ShowComments)
	http.HandleFunc("/CreateComments", f.CreateComments)
	http.HandleFunc("/submit-comment", f.Submit_Comment)
	/* http.HandleFunc("/AddLikeComment", f.AddLikeComment)
	http.HandleFunc("/AddDislikeComment", f.AddDislikeComment) */
	/* http.HandleFunc("/", f.QCQHandler)
	http.HandleFunc("/500", f.CCCHandler) */
	http.HandleFunc("/logout", f.LogoutHandler)

	fmt.Println("starting server at port 8000 : http://localhost:8000/index")
	http.ListenAndServe(":8000", nil)
}
