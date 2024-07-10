package forum

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/gofrs/uuid"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login handler")
	switch r.Method {
	case http.MethodGet:
		// GET : renvoyer la page login
		tmpl, err := template.ParseFiles("./ui/templates/login.html")
		if err != nil {
			fmt.Println("register handler, template méthode get")
			fmt.Println(err)
			return
		}
		tmpl.Execute(w, "data goes here")
	case http.MethodPost:
		// POST : vérifier username et email
		err := r.ParseForm()
		if err != nil {
			fmt.Println("register post problem")
			fmt.Println(err)
			return
		}
		// récupếration des valeurs email et username
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		if UserCanLog(email, username, password) {
			sessionID, err := uuid.NewV4()
			if err != nil {
				fmt.Println("uuid problem")
				fmt.Println(err)
				return
			}

			db.Exec("UPDATE users SET session_id=? WHERE username=?", sessionID, username)
			if err != nil {
				fmt.Println("register post problem")
				fmt.Println(err)
				return
			}
			setCookieHandler(w, r, username, sessionID.String())

			fmt.Println("session ID:" + sessionID.String())
			http.Redirect(w, r, "/index", http.StatusSeeOther)
			return
		} else {
			// l'utilisateur n'existe pas on renvoieLogoutHandler sur register
			data := Message{"this user isnt registed"}
			tmpl, err := template.ParseFiles("./ui/templates/register.html")
			if err != nil {
				fmt.Println("register handler, template méthode get")
				fmt.Println(err)
				return
			}
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			tmpl.Execute(w, data)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func setCookieHandler(rw http.ResponseWriter, r *http.Request, username string, sessionId string) {
	fmt.Println("setCookieHandler")
	uc := &http.Cookie{
		Name:  "Username",
		Value: username,
		Path:  "/",
	}
	http.SetCookie(rw, uc)
	c := &http.Cookie{
		Name:  "Token",
		Value: sessionId,
		Path:  "/",
	}
	http.SetCookie(rw, c)
}
func UserCanLog(email, username, password string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND username = ? AND password = ?", email, username, password).Scan(&count)
	if err != nil {
		fmt.Println("error at usermayexist function")
		fmt.Println(err)
		return false
	}
	return count > 0
}

//	func CheckSession(w http.ResponseWriter, r *http.Request) {
//		cookie, errCookie := r.Cookie("Username")
//		if errCookie != nil {
//			fmt.Println("cookie", errCookie)
//			return
//		}
//		var exist bool
//		db.QueryRow("SELECT 1 EXISTS(SELECT setSessionID FROM users WHERE setSessionID=?)", cookie.Value).Scan(&exist)
//		fmt.Println("Is User Connected ?= ", exist)
//		if !exist {
//			http.Redirect(w, r, "/index", http.StatusSeeOther)
//		}
//	}
func CheckSession(w http.ResponseWriter, r *http.Request) {
	cookie, errCookie := r.Cookie("Token") // Change "Username" to "Token" to get the session ID cookie
	if errCookie != nil {
		fmt.Println("cookie", errCookie)
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login if session ID cookie is not found
		return
	}
	var exist int
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session_id=?)", cookie.Value).Scan(&exist)
	if err != nil {
		fmt.Println("CheckSession query error:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login in case of a database query error
		return
	}
	if exist != 1 {
		fmt.Println("User not connected")
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login if user is not found
		return
	}
	fmt.Println("User Connected")
	// Proceed with handling the authenticated user
}
func setSession(u, uid string) {
	stmt, err := db.Prepare("UPDATE users SET session_id=? WHERE username=?")
	if err != nil {
		fmt.Println(err)
		return
	}
	// TODO Check erreur
	// _, errExec := stmt.Exec(uid, u)
	stmt.Exec(uid, u)
}
func clearSessionCookies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test enter clearsession")
	userCookie := http.Cookie{
		Name:     "Username",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, &userCookie)
	tokenCookie := http.Cookie{
		Name:     "Token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, &tokenCookie)
}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("testtttt")
	clearSessionCookies(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check session cookie validity or user login status here
		// For example, you can check if the session cookie exists and is valid
		_, err := r.Cookie("Username")
		if err != nil {
			// Redirect to home.html if session cookie doesn't exist
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
	})
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "ui/templates/home.html")
}
