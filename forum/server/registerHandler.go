package forum

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gofrs/uuid"
)

type Message struct {
	Message string
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("register handler")
	switch r.Method {
	case http.MethodGet:
		// GET : renvoyer la page resgister
		tmpl, err := template.ParseFiles("./ui/templates/register.html")
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

		session_id, err := uuid.NewV4()
		if err != nil {
			fmt.Println("uuid problem")
			fmt.Println(err)
			return
		}

		db.Exec("UPDATE users SET session_id=? WHERE username=?", session_id, username)

		if UserMayExist(email, username) {
			// l'utilisateur existe déjà => message d'erreur
			data := Message{"user already exist"}
			tmpl, err := template.ParseFiles("./ui/templates/register.html")
			if err != nil {
				fmt.Println("register handler, template méthode get")
				fmt.Println(err)
				return
			}
			tmpl.Execute(w, data)
		} else {
			// l'utilisateur n'existe pas : on le rajoute à la db
			session_idSTR := session_id.String()
			insertUser(email, username, session_idSTR, password)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
func UserMayExist(email, username string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? OR username = ?", email, username).Scan(&count)
	if err != nil {
		fmt.Println("error at usermayexist function")
		fmt.Println(err)
	}
	return count > 0
}
