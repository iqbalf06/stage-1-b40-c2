package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/add-project", formAddProject).Methods("GET")
	route.HandleFunc("/add-project", middleware.UploadFile(addProject)).Methods("POST")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")

	route.HandleFunc("/update-project/{id}", updateProject).Methods("GET")
	route.HandleFunc("/submit-update/{id}", middleware.UpdateFile(submitUpdate)).Methods("POST")

	route.HandleFunc("/register", registerForm).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")
	route.HandleFunc("/login", loginForm).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")
	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("server running on port 5000")
	http.ListenAndServe("localhost:5000", route)

}

type SessionData struct {
	IsLogin   bool
	UserName  string
	FlashData string
}

var Data = SessionData{}

type Project struct {
	ID             int
	ProjectName    string
	StartDate      time.Time
	EndDate        time.Time
	Description    string
	Nodejs         string
	React          string
	Java           string
	Python         string
	Duration       string
	IsLogin        bool
	Author         string
	Image          string
	Date_StartDate string
	Date_EndDate   string
	Techno         []string
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	//menyimpan session ke dalam cookie
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	//kondisi untuk Login
	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {

			flashes = append(flashes, f1.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	data, _ := connection.Conn.Query(context.Background(),
		"SELECT tb_projects.id, tb_projects.name, description, duration, image, technologies, tb_user.name as author FROM tb_projects LEFT JOIN tb_user ON tb_projects.author_id = tb_user.id ORDER BY id DESC")

	var result []Project
	for data.Next() {
		var each = Project{}

		var err = data.Scan(&each.ID, &each.ProjectName, &each.Description, &each.Duration, &each.Image, &each.Techno, &each.Author)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	resData := map[string]interface{}{
		"DataSession": Data,
		"Projects":    result,
	}

	w.WriteHeader(http.StatusOK)

	tmpl.Execute(w, resData)
}

func formAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/add-project.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}
	//menyimpan session ke dalam cookie
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	//kondisi untuk Login
	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	data := map[string]interface{}{
		"DataSession": Data,
	}

	tmpl.Execute(w, data)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	var projectName = r.PostForm.Get("inputProjectName")
	var startDate = r.PostForm.Get("inputStartdate")
	var endDate = r.PostForm.Get("inputEnddate")
	var description = r.PostForm.Get("inputDescription")

	var techno []string
	techno = r.Form["technologies"]

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	author := session.Values["ID"].(int)
	fmt.Println((author))

	layout := ("2006-01-02")
	startDateParse, _ := time.Parse(layout, startDate)
	endDateParse, _ := time.Parse(layout, endDate)

	hours := endDateParse.Sub(startDateParse).Hours()
	days := hours / 24
	weeks := math.Round(days / 7)
	months := math.Round(days / 30)
	years := math.Round(days / 365)

	var duration string

	if years > 0 {
		duration = strconv.FormatFloat(years, 'f', 0, 64) + "year"
	} else if months > 0 {
		duration = strconv.FormatFloat(months, 'f', 0, 64) + " Month"
	} else if weeks > 0 {
		duration = strconv.FormatFloat(weeks, 'f', 0, 64) + " Week"
	} else if days > 0 {
		duration = strconv.FormatFloat(days, 'f', 0, 64) + " Day"
	} else if hours > 0 {
		duration = strconv.FormatFloat(hours, 'f', 0, 64) + " Hour"
	} else {
		duration = "0 Days"
	}

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, start_date, end_date, description, author_id, duration, image, technologies) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", projectName, startDate, endDate, description, author, duration, image, techno)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	data := map[string]interface{}{
		"DataSession": Data,
	}

	tmpl.Execute(w, data)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/project-detail.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var ProjectDetail = Project{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, description, start_date, end_date, duration, image, technologies FROM tb_projects WHERE id=$1", id).Scan(&ProjectDetail.ID, &ProjectDetail.ProjectName, &ProjectDetail.Description, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Image, &ProjectDetail.Techno)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	ProjectDetail.Date_StartDate = ProjectDetail.StartDate.Format("2 Jan 2006")
	ProjectDetail.Date_EndDate = ProjectDetail.EndDate.Format("2 Jan 2006")

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	data := map[string]interface{}{
		"Project":     ProjectDetail,
		"DataSession": Data,
	}

	tmpl.Execute(w, data)

}

func deleteProject(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func updateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/update-project.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	var updateProject = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, description, start_date, end_date, duration FROM tb_projects WHERE id = $1", index).Scan(&updateProject.ID, &updateProject.ProjectName, &updateProject.Description, &updateProject.StartDate, &updateProject.EndDate, &updateProject.Duration)

	updateProject.Date_StartDate = updateProject.StartDate.Format("2006-01-02")
	updateProject.Date_EndDate = updateProject.EndDate.Format("2006-01-02")

	data := map[string]interface{}{
		"DataSession":   Data,
		"UpdateProject": updateProject,
	}

	tmpl.Execute(w, data)
}

func submitUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var projectName = r.PostForm.Get("inputProjectName")
	var startDate = r.PostForm.Get("inputStartdate")
	var endDate = r.PostForm.Get("inputEnddate")
	var description = r.PostForm.Get("inputDescription")
	var techno []string
	techno = r.Form["technologies"]

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	author := session.Values["ID"].(int)
	fmt.Println((author))

	layout := ("2006-01-02")
	startDateParse, _ := time.Parse(layout, startDate)
	endDateParse, _ := time.Parse(layout, endDate)

	hours := endDateParse.Sub(startDateParse).Hours()
	days := hours / 24
	weeks := math.Round(days / 7)
	months := math.Round(days / 30)
	years := math.Round(days / 365)

	var duration string

	if years > 0 {
		duration = strconv.FormatFloat(years, 'f', 0, 64) + "year"
	} else if months > 0 {
		duration = strconv.FormatFloat(months, 'f', 0, 64) + " Month"
	} else if weeks > 0 {
		duration = strconv.FormatFloat(weeks, 'f', 0, 64) + " Week"
	} else if days > 0 {
		duration = strconv.FormatFloat(days, 'f', 0, 64) + " Day"
	} else if hours > 0 {
		duration = strconv.FormatFloat(hours, 'f', 0, 64) + " Hour"
	} else {
		duration = "0 Days"
	}

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name = $1, description = $2, start_date = $3, end_date = $4, duration = $5, image = $6, technologies= $7 WHERE id = $8", projectName, description, startDateParse, endDateParse, duration, image, techno, id)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func registerForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/register.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func loginForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/login.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		fmt.Println("Email belum terdaftar")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : Email belum terdaftar" + err.Error()))
		return
	}

	fmt.Println(user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		fmt.Println("Password salah")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : Password salah" + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	session.Values["Name"] = user.Name
	session.Values["Email"] = user.Email
	session.Values["ID"] = user.ID
	session.Values["IsLogin"] = true
	session.Options.MaxAge = 10800 // 3 hours

	session.AddFlash("Login Successful", "message")

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
