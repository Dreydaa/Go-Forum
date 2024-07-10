package forum

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitiateDatabase() {
	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		panic(err)
	}
	db = database
}

func FillDatabase() {
	if _, err := os.Stat("./forum.db"); err == nil {
		return
	}

	CreateTables()
	CreateTablesLike()
	CreateTablesLikeComment()
	CreateTablesPost()
	CreateTablesCategories()
	CreateTablesCategories2()
	CreateTablesComment()
	FillTables()
}

func CreateTables() {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, email TEXT, session_id TEXT, password TEXT ) ") // auto-incrémentation de l'i
	if err != nil {
		fmt.Println("bug ici")
		fmt.Println(err)
		return
	}
	statement.Exec()
}

func CreateTablesPost() {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, title TEXT, categories TEXT, categories2 TEXT, content TEXT, like INTEGER, dislike INTEGER, CommentsNumb INTEGER, date TEXT NOT NULL) ")
	if err != nil {
		fmt.Println("Errreur ici (post)")
		fmt.Println(err)
		return
	}
	statement.Exec()
}

func CreateTablesLike() {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS like (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, post_id INT, like INTEGER, dislike INTEGER) ")
	if err != nil {
		fmt.Println("Errreur ici (like)")
		fmt.Println(err)
		return
	}
	statement.Exec()
}

func CreateTablesLikeComment() {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS likeComment (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, post_id INT, like INTEGER, dislike INTEGER) ")
	if err != nil {
		fmt.Println("Errreur ici (like)")
		fmt.Println(err)
		return
	}
	statement.Exec()
}

func CreateTablesComment() {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS comment (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, content TEXT, like INTEGER, dislike INTEGER, post_id INT, date TEXT NOT NULL) ")
	if err != nil {
		fmt.Println("Errreur ici (comment)")
		fmt.Println(err)
		return
	}
	statement.Exec()
}

func CreateTablesCategories() {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT)")
	if err != nil {
		fmt.Println("Errreur ici (categories)")
		fmt.Println(err)
		return
	}
	statement.Exec()
}

func CreateTablesCategories2() {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS categories2 (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT)")
	if err != nil {
		fmt.Println("Errreur ici (categories)")
		fmt.Println(err)
		return
	}
	statement.Exec()
}

func FillTables() {

	InsertCategories("VALORANT")
	InsertCategories("League_Of_Legends")
	InsertCategories("Sport") /*  */
	InsertCategories("Cinema")
	InsertCategories("Actualité")
	InsertCategories("Economie")
	InsertCategories("Basketball")
	InsertCategories("Meteo")
	InsertCategories("Paysans")
	InsertCategories("Riche")
	InsertCategories("Puants")

	InsertCategories2("VALORANT")
	InsertCategories2("League_Of_Legends")
	InsertCategories2("Sport")
	InsertCategories2("Cinema")
	InsertCategories2("Actualité")
	InsertCategories2("Economie")
	InsertCategories2("Basketball")
	InsertCategories2("Meteo")
	InsertCategories2("Paysans")
	InsertCategories2("Riche")
	InsertCategories2("Puants")
}

func insertUser(username, email, session_id, password string) error {
	_, err := db.Exec("INSERT INTO users (email, username, session_id, password) VALUES (?, ?, ?, ?)", username, email, session_id, password)
	if err != nil {
		fmt.Println("Error insert user")
		fmt.Println(err)
		return err
	}
	return nil
}

func InsertPost(username string, title string, categories string, categories2 string, content string, like int, dislike int, commentsNumb int) error {

	_, err := db.Exec("INSERT INTO posts (username, title, categories, categories2, content, like, dislike, CommentsNumb, date) VALUES (?,?,?,?,?,?,?,?,?)", username, title, categories, categories2, content, like, dislike, commentsNumb, time.Now().Format("2006-01-02 : 15:04"))
	if err != nil {
		fmt.Println("Error insert post")
		fmt.Println(err)
		return err
	}

	return nil
}

func InsertComment(username string, content string, like int, dislike int, post_id int) error {
	_, err := db.Exec("INSERT INTO comment (username, content, like, dislike, post_id, date) VALUES (?, ?, ?, ?, ?, ?)", username, content, like, dislike, post_id, time.Now().Format("2006-01-02 : 15:04"))
	if err != nil {
		fmt.Println("Error insert comment")
		fmt.Println(err)
		return err
	}
	return nil
}

func InsertLike(username string, post_id int, like int, dislike int) {
	_, err := db.Exec("INSERT INTO like (username, post_id, like, dislike) VALUES (?, ?, ?, ?)", username, post_id, like, dislike)
	if err != nil {
		fmt.Println("Error insert like")
		fmt.Println(err)
		return
	}
}
func InsertLikeComment(username string, post_id int, like int, dislike int) {
	_, err := db.Exec("INSERT INTO likeComment (username, post_id, like, dislike) VALUES (?, ?, ?, ?)", username, post_id, like, dislike)
	if err != nil {
		fmt.Println("Error insert like")
		fmt.Println(err)
		return
	}
}

func InsertCategories(title string) error {
	_, err := db.Exec("INSERT INTO categories (title) VALUES (?)", title)
	if err != nil {
		fmt.Println("Error insert categories")
		fmt.Println(err)
		return err
	}
	return nil
}

func InsertCategories2(title string) error {
	_, err := db.Exec("INSERT INTO categories2 (title) VALUES (?)", title)
	if err != nil {
		fmt.Println("Error insert categories2")
		fmt.Println(err)
		return err
	}
	return nil
}

func GetCategories() []string {
	categories := make([]string, 0)

	res, err := db.Query("SELECT title FROM categories")
	if err != nil {
		panic(err)
	}
	for res.Next() {
		var cat string
		res.Scan(&cat)
		categories = append(categories, cat)
	}
	return categories
}

func GetCategories2() []string {
	categories2 := make([]string, 0)

	res, err := db.Query("SELECT title FROM categories2")
	if err != nil {
		panic(err)
	}
	for res.Next() {
		var cat string
		res.Scan(&cat)
		categories2 = append(categories2, cat)
	}
	return categories2
}

func AddLikeDB(post_id int, username string) {
	fmt.Print("addlikeddb is called\n")
	var like int
	var dislike int
	res := db.QueryRow("SELECT like FROM like WHERE username =? AND post_id =?", username, post_id).Scan(&like)
	db.QueryRow("SELECT dislike FROM like WHERE username=? AND post_id =?", username, post_id).Scan(&dislike)

	if dislike == 0 {
		if like == 0 {
			if res != nil {
				db.Exec("INSERT INTO like (username, post_id, like, dislike) VALUES (?, ?, ?, ?)", username, post_id, 1, 0)
				db.Exec("UPDATE posts SET like = like + 1 WHERE id =?", post_id)
			} else if res == nil {
				db.Exec("UPDATE posts SET like = like + 1 WHERE id=?", post_id)
				db.Exec("UPDATE like SET like = 1 WHERE username =? AND post_id=?", username, post_id)
			}
		} else if like == 1 {
			db.Exec("UPDATE posts SET like = like - 1 WHERE id=?", post_id)
			db.Exec("UPDATE like SET like = 0 WHERE username =? AND post_id=?", username, post_id)
		}
	}
}

func AddDislikeDB(post_id int, username string) {
	fmt.Print("addlikeddb is called\n")
	var like int
	var dislike int
	res := db.QueryRow("SELECT like FROM like WHERE username =? AND post_id =?", username, post_id).Scan(&like)
	db.QueryRow("SELECT dislike FROM like WHERE username=? AND post_id =?", username, post_id).Scan(&dislike)

	if like == 0 {
		if dislike == 0 {
			if res != nil {
				db.Exec("INSERT INTO like (username, post_id, like, dislike) VALUES (?, ?, ?, ?)", username, post_id, 1, 0)
				db.Exec("UPDATE like SET dislike = dislike + 1 WHERE id =?", post_id)
			} else if res == nil {
				db.Exec("UPDATE posts SET dislike = dislike + 1 WHERE id=?", post_id)
				db.Exec("UPDATE like SET dislike = 1 WHERE username =? AND post_id=?", username, post_id)
			}
		} else if dislike == 1 {
			db.Exec("UPDATE posts SET dislike = dislike - 1 WHERE id=?", post_id)
			db.Exec("UPDATE like SET dislike =  0 WHERE username =? AND post_id=?", username, post_id)
		}
	}
}

func GetComments(post_id string) []Comments {
	comments := make([]Comments, 0)

	db, err := db.Query("SELECT * FROM comment WHERE post_id =?", post_id)
	if err != nil {
		fmt.Println("error in getcomments function")
	}

	fmt.Println("post id:", post_id)

	for db.Next() {
		var comment Comments
		errScan := db.Scan(&comment.Id, &comment.Username, &comment.Content, &comment.Like, &comment.Dislike, &comment.Post_id, &comment.Date)
		if errScan != nil {
			panic("an error has occurred on GetComments :" + errScan.Error())
		}
		comments = append(comments, comment)

	}

	return comments
}

func AddComents(username string, content string, like int, dislike int, post_id int) error {
	_, err := db.Exec("INSERT INTO comment (username, content, like, dislike, post_id, date) VALUES (?,?,?,?,?,?)", username, content, like, dislike, post_id, time.Now().Format("2006-01-02 : 15:04"))
	if err != nil {
		fmt.Println("Error insert comment")
		fmt.Println(err)
		return err
	}
	return nil
}

func AddLikeCommentDB(post_id int, username string) {
	fmt.Print("addlikeddb is called\n")
	var like int
	var dislike int
	res := db.QueryRow("SELECT like FROM likeComment WHERE username =? AND post_id =?", username, post_id).Scan(&like)
	db.QueryRow("SELECT dislike FROM likeComment WHERE username=? AND post_id =?", username, post_id).Scan(&dislike)
	fmt.Println("likecomment :", like)
	fmt.Println("dislikecomment:", dislike)

	if dislike == 0 {
		if like == 0 {
			if res != nil {
				db.Exec("INSERT INTO likeComment (username, post_id, like, dislike) VALUES (?, ?, ?, ?)", username, post_id, 1, 0)
				db.Exec("UPDATE comment SET like = like + 1 WHERE id =?", post_id)
			} else if res == nil {
				db.Exec("UPDATE comment SET like = like + 1 WHERE id=?", post_id)
				db.Exec("UPDATE likeComment SET like = 1 WHERE username =? AND post_id=?", username, post_id)
			}
		} else if like == 1 {
			db.Exec("UPDATE comment SET like = like - 1 WHERE id=?", post_id)
			db.Exec("UPDATE likeComment SET like = 0 WHERE username =? AND post_id=?", username, post_id)
		}
	}
}

func AddDislikeCommentDB(post_id int, username string) {
	fmt.Print("addlikeddb is called\n")
	var like int
	var dislike int
	res := db.QueryRow("SELECT like FROM likeComment WHERE username =? AND post_id =?", username, post_id).Scan(&like)
	db.QueryRow("SELECT dislike FROM likeComment WHERE username=? AND post_id =?", username, post_id).Scan(&dislike)
	fmt.Println("likecomment :", like)
	fmt.Println("dislikecomment:", dislike)

	if like == 0 {
		if dislike == 0 {
			if res != nil {
				db.Exec("INSERT INTO likeComment (username, post_id, like, dislike) VALUES (?, ?, ?, ?)", username, post_id, 1, 0)
				db.Exec("UPDATE likeComment SET dislike = dislike + 1 WHERE id =?", post_id)
			} else if res == nil {
				db.Exec("UPDATE comment SET dislike = dislike + 1 WHERE id=?", post_id)
				db.Exec("UPDATE likeComment SET dislike = 1 WHERE username =? AND post_id=?", username, post_id)
			}
		} else if dislike == 1 {
			db.Exec("UPDATE comment SET dislike = dislike - 1 WHERE id=?", post_id)
			db.Exec("UPDATE likeComment SET dislike =  0 WHERE username =? AND post_id=?", username, post_id)
		}
	}
}

/* func GetPostIDLiked(username string) []string {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil
	}
	defer db.Close()

	res, err := db.Query("SELECT ID FROM like WHERE username = ?", username)
	if err != nil {
		return nil
	}

	likedposts := []string{}

	for res.Next() {
		var likePost string
		errScan := res.Scan(&likePost)
		if errScan != nil {
			panic("an error has occurred on GetUsernameByPost :" + errScan.Error())
		}
		likedposts = append(likedposts, likePost)
	}
	for _, idPostLiked := range likedposts {
		fmt.Println("likedPost:", idPostLiked)
	}
	return likedposts
}

func GetPostFromIdLike(idposts []string) []Post {
	db, err := sql.Open("sqlite3", "./forum.db")
	posts := []Post{}
	if err != nil {
		return nil
	}
	defer db.Close()

	for _, idpost := range idposts {
		res, err := db.Query("SELECT Id, username, title, categories, categories2, content, like, dislike, CommentsNumb, date FROM posts WHERE id = ?", idpost)
		if err != nil {
			return nil
		}

		var post Post
		if res.Next() {
			errScan := res.Scan(&post.Id, &post.Username, &post.Title, &post.Categories, &post.Categories2, &post.Content, &post.Like, &post.Dislike, &post.CommentsNumb, &post.Date)
			if errScan != nil {
				panic("an error has occurred on GetPostFromIdLike :" + errScan.Error())
			}
			return []Post{post}
		}
	}
	return posts
} */
