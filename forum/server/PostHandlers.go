package forum

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

type Post struct {
	Id           int
	Username     string
	Title        string
	Categories   string
	Categories2  string
	Content      string
	Like         int
	Dislike      int
	CommentsNumb int
	Date         string
}

type Posts struct {
	PostAll []Post
}

type Comments struct {
	Id       int
	Username string
	Content  string
	Like     int
	Dislike  int
	Post_id  int
	Date     string
}

func HandleMainPage(w http.ResponseWriter, r *http.Request) {
	username, _, _ := GetCookies(r)
	/* 	fmt.Println(username) */

	/* 	most_liked := r.FormValue("most_Liked")
	   	fmt.Println(most_liked) */
	category := r.FormValue("")
	/* likeFilter := r.FormValue("like-filter")
	datefilter := r.FormValue("created-post-filter")  */

	/* Mycreatedpost := r.FormValue("my-created-post") */

	/* Mylikedpost := r.FormValue("my-liked-post") */

	/* posts := GetPostFromUser(username, GetAllPosts()) */

	/* if category == "" {
		category = "All"
	} */
	posts, err := GetPosts(category)
	if err != nil {
		fmt.Println("error getting posts:", err)
		return
	}
	/* 	categories := GetCategories()
	   	categories2 := GetCategories() */

	// username, _, _ := GetCookies(r)

	/* if likeFilter == "most_liked" {
		posts = GetFilterLike("DESC")
	}

	if datefilter == "most_recent" {
		posts = GetDatePost("DESC")
	} */

	tmpl, err := template.ParseFiles("./ui/templates/index.html")
	if err != nil {
		fmt.Println("error parsing template:", err)
		return
	}

	datas := struct {
		Posts       []Post
		Categories  []string
		Categories2 []string
	}{
		Posts: posts,
		/* Categories:  categories,
		Categories2: categories2, */
	}

	if username == "connected" {
		datas.Posts = GetPostCreated(username, posts)
	}

	/* if Mylikedpost == "my-liked-post" {
		postsId := GetPostsId(username)
		fmt.Println("postId:", postsId)
		datas.Posts = GetLikedPost(postsId, posts)
	} */

	if err = tmpl.Execute(w, datas); err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		return
	}

}

func GetAllPosts() []Post {
	posts, err := GetPosts("All")
	if err != nil {
		fmt.Println("error getting posts:", err)
		return nil
	}
	return posts
}

func GetPostFromUser(username string, posts []Post) []Post {
	var postsUser []Post
	for _, p := range posts {
		if p.Username == username {
			postsUser = append(postsUser, p)
		}
	}
	return postsUser
}

func HandleCreatePost(w http.ResponseWriter, r *http.Request) {

	CheckSession(w, r)

	tmpl, err := template.ParseFiles("./ui/templates/create-post.html")
	if err != nil {
		fmt.Println("error parsing template:", err)
		return
	}

	categories := GetCategories()
	categories2 := GetCategories2()
	data := struct {
		Categories  []string
		Categories2 []string
	}{
		Categories:  categories,
		Categories2: categories2,
	}

	tmpl.Execute(w, data)
}
func GetCookies(r *http.Request) (string, string, string) {
	cookies := r.Cookies()
	var username, token string
	error := "probleme GETcookie"
	// Parcourir la liste des cookies
	for _, cookie := range cookies {
		// Vérifier le nom du cookie
		if cookie.Name == "Username" {
			// Traitement pour le cookie "User"
			username = cookie.Value
		} else if cookie.Name == "Token" {
			// Traitement pour le cookie "Token"
			token = cookie.Value
		}
		// Ajoutez d'autres conditions si vous avez plus de cookies à vérifier
	}
	if username == "" {
		return "prob", "prob", error
	}
	if token == "" {
		return "prob", "prob", error
	}
	return username, token, error
}

func HandlerAddPost(w http.ResponseWriter, r *http.Request) {

	title := r.FormValue("Title")
	content := r.FormValue("text")
	categories := r.FormValue("Categories")
	categories2 := r.FormValue("Categories2")

	username, _, _ := GetCookies(r)
	// fmt.Println(title, content, categories)
	/* fmt.Println("title :", title)
	fmt.Println("content :", content)
	fmt.Println("categories:", categories) */
	InsertPost(username, title, categories, categories2, content, 0, 0, 0)

	http.Redirect(w, r, "/index", http.StatusSeeOther)
}

func GetPosts(category string) ([]Post, error) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT Id, username, title, categories, categories2, content, like, dislike, CommentsNumb, date FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Id, &post.Username, &post.Title, &post.Categories, &post.Categories2, &post.Content, &post.Like, &post.Dislike, &post.CommentsNumb, &post.Date); err != nil {
			return nil, err
		}

		res, err := db.Query("SELECT COUNT(*) AS content FROM comment WHERE post_id = ?", post.Id)
		if err != nil {
			return nil, err
		}
		defer res.Close()

		if res.Next() {
			var commentCount int
			if err := res.Scan(&commentCount); err != nil {
				return nil, err
			}
			post.CommentsNumb = commentCount
		}

		if post.Categories == category || category == "All" || post.Categories2 == category {
			posts = append(posts, post)
		}

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

/* func AddLike(w http.ResponseWriter, r *http.Request) {

	CheckSession(w, r)

	username, _, _ := GetCookies(r)

	idPostStr := r.FormValue("idpost")
	idPost, _ := strconv.Atoi(idPostStr)

	AddLikeDB(idPost, username)
	http.Redirect(w, r, "/index", http.StatusSeeOther)
}

func AddDislike(w http.ResponseWriter, r *http.Request) {

	CheckSession(w, r)

	username, _, _ := GetCookies(r)

	idPostStr := r.FormValue("idpost2")
	idPost, _ := strconv.Atoi(idPostStr)

	AddDislikeDB(idPost, username)
	http.Redirect(w, r, "/index", http.StatusSeeOther)
} */

/* func GetFilterLike(direction string) []Post {
	var posts []Post
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil
	}
	defer db.Close()
	// var query string = fmt.Sprintf("SELECT * From posts ORDER BY like %s", direction)
	// result, err := db.Query(query)
	result, err := db.Query("SELECT * From posts ORDER BY like DESC")
	if err != nil {
		panic("An error has occured:" + err.Error())
	}

	for result.Next() {
		var post Post
		errScan := result.Scan(&post.Id, &post.Username, &post.Title, &post.Categories, &post.Categories2, &post.Content, &post.Like, &post.Dislike, &post.CommentsNumb, &post.Date)
		if errScan != nil {
			panic("An error has occurred on the error scan :" + errScan.Error())
		}
		posts = append(posts, post)

	}
	return posts
} */

/* func GetDatePost(direction_time string) []Post {
	var posts []Post
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil
	}
	defer db.Close()

	result, err := db.Query(fmt.Sprintf("SELECT * FROM posts ORDER BY date %s", direction_time))
	if err != nil {
		panic("error in gatdatepost ligne:" + err.Error())
	}

	for result.Next() {
		var post Post
		errScan := result.Scan(&post.Id, &post.Username, &post.Title, &post.Categories, &post.Categories2, &post.Content, &post.Like, &post.Dislike, &post.CommentsNumb, &post.Date)
		if errScan != nil {
			panic("An error has occurred on the error scan in getDatePost :" + errScan.Error())
		}
		posts = append(posts, post)
	}
	return posts
} */

func ShowComments(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./ui/templates/show-comments.html")
	if err != nil {
		fmt.Println("error parsing template:", err)
		return
	}

	var idPost string = r.FormValue("idPost")

	if idPost == "" {
		idPost = r.URL.Query().Get("ID_POST")
	}

	comments := GetComments(idPost)

	tmpl.Execute(w, comments)
}

func CreateComments(w http.ResponseWriter, r *http.Request) {
	CheckSession(w, r)
	tmpl, err := template.ParseFiles("./ui/templates/create-comments.html")
	if err != nil {
		return
	}
	idPost := r.FormValue("id-post")
	data := struct {
		Id string
	}{
		Id: idPost,
	}

	tmpl.Execute(w, data)
}

func Submit_Comment(w http.ResponseWriter, r *http.Request) {
	username, _, _ := GetCookies(r)

	// comment := r.FormValue("comments")
	comment := r.FormValue("comment")
	idPost := r.FormValue("id-post")

	idPostInt, _ := strconv.Atoi(idPost)

	AddComents(username, comment, 0, 0, idPostInt)
	http.Redirect(w, r, "/ShowComments?ID_POST="+idPost, http.StatusSeeOther)
}

/* func CCCHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpl := template.Must(template.ParseFiles("./ui/templates/500.html"))
	tmpl.Execute(w, nil)
}

func QCQHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	tmpl := template.Must(template.ParseFiles("./ui/templates/404.html"))
	tmpl.Execute(w, nil)
} */

/* func AddLikeComment(w http.ResponseWriter, r *http.Request) {
	CheckSession(w, r)
	username, _, _ := GetCookies(r)

	idCommentStr := r.FormValue("idComment")
	idComment, _ := strconv.Atoi(idCommentStr)

	idPost := r.FormValue("idPost")

	AddLikeCommentDB(idComment, username)
	http.Redirect(w, r, "/ShowComments?ID_POST="+idPost, http.StatusSeeOther)
}

func AddDislikeComment(w http.ResponseWriter, r *http.Request) {
	CheckSession(w, r)
	username, _, _ := GetCookies(r)

	idCommentStr := r.FormValue("idComment")
	idComment, _ := strconv.Atoi(idCommentStr)

	idPost := r.FormValue("idPost")

	AddDislikeCommentDB(idComment, username)
	http.Redirect(w, r, "/ShowComments?ID_POST="+idPost+"", http.StatusSeeOther)
} */

func GetUsernameById(session_id string) string {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return ""
	}
	defer db.Close()

	var username string
	err = db.QueryRow("SELECT (*) FROM users WHERE id = ?", session_id).Scan(&username)
	if err != nil {
		return ""
	}
	return username

}

func GetPostCreated(username string, AllPosts []Post) []Post {
	var posts []Post
	for _, p := range AllPosts {
		if p.Username == username {
			posts = append(posts, p)
		}
	}
	fmt.Println("username :", username)
	fmt.Println("posts :", posts)
	return posts

}

func GetPostsId(username string) []string {
	idposts := make([]string, 0)
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		panic("An error has occurred on the file local:" + err.Error())
	}
	defer db.Close()

	res, err := db.Query("SELECT post_id FROM like WHERE username=? AND like=?", username, 1)
	if err != nil {
		panic(err)
	}

	for res.Next() {
		var idpost int
		errScan := res.Scan(&idpost)
		if errScan != nil {
			panic("an error has occurred on GetUsernameByPost :" + errScan.Error())
		}
		idpoststr := strconv.Itoa(idpost) //convert integer to string for slice
		idposts = append(idposts, idpoststr)
	}
	return idposts
}

/*
func GetLikedPost(idposts []string, allPost []Post) []Post {
	posts := make([]Post, 0)
	for _, post := range allPost {
		for _, idpost := range idposts {
			if idpost == strconv.Itoa(post.Id) {
				fmt.Println("poste name!", post.Title)
				posts = append(posts, post)
			}
		}
	}
	return posts
}
*/
