package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter() //variabel route dari mux router

	connection.DatabaseConnect()

	// route path folder public (js,css,images)
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	//routing, menjalankan function home
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/add-project", formAddProject).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET") //index url params
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")

	route.HandleFunc("/update-project/{id}", updateProject).Methods("GET")
	route.HandleFunc("/submit-update/{id}", submitUpdate).Methods("POST")

	fmt.Println("server running on port 5000")
	http.ListenAndServe("localhost:5000", route)

}

// type data untuk variabel object
type Project struct {
	ID          int
	ProjectName string
	StartDate   string
	EndDate     string
	Description string
	Nodejs      string
	React       string
	Java        string
	Python      string
	Duration    string
}

// menampilkan page dan data
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil { //untuk handle error. agar tau errornya
		w.Write([]byte("message : " + err.Error())) //byte untuk menampilkan string jika ada error
		return
	}

	// memanggil variabel conn dari package connection
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, description FROM tb_projects ORDER BY id DESC")
	// fmt.Println(data)

	var result []Project //variabel result dan slice dan struct yang punya type data
	for data.Next() {    //looping data
		var each = Project{} //each untuk menampung data2 dari query.

		var err = data.Scan(&each.ID, &each.ProjectName, &each.Description)
		if err != nil {
			fmt.Println(err.Error())
			return //jika error akan berhenti disini
		}

		result = append(result, each) //menampung data result. sperti push pda JS
	}

	resData := map[string]interface{}{
		"Projects": result, //properi yang berisi result query diatas dikirim menggunakan map string
	}

	w.WriteHeader(http.StatusOK)

	tmpl.Execute(w, resData) //menampilkan response dari views
}

func formAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/add-project.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

// function untuk memasukkan dan menangkap data
func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	//variabel object untuk menampung data dari tag input.
	var projectName = r.PostForm.Get("inputProjectName")
	var startDate = r.PostForm.Get("inputStartdate")
	var endDate = r.PostForm.Get("inputEnddate")
	var description = r.PostForm.Get("inputDescription")

	// layout := ("2006-01-02")
	// startDateParse, _ := time.Parse(layout, startDate)
	// endDateParse, _ := time.Parse(layout, endDate)

	// hours := endDateParse.Sub(startDateParse).Hours()
	// days := hours / 24
	// weeks := math.Round(days / 7)
	// months := math.Round(days / 30)
	// years := math.Round(days / 365)

	// var duration string

	// if years > 0 {
	// 	duration = strconv.FormatFloat(years, 'f', 0, 64) + "year"
	// } else if months > 0 {
	// 	duration = strconv.FormatFloat(months, 'f', 0, 64) + " Month"
	// } else if weeks > 0 {
	// 	duration = strconv.FormatFloat(weeks, 'f', 0, 64) + " Week"
	// } else if days > 0 {
	// 	duration = strconv.FormatFloat(days, 'f', 0, 64) + " Day"
	// } else if hours > 0 {
	// 	duration = strconv.FormatFloat(hours, 'f', 0, 64) + " Hour"
	// } else {
	// 	duration = "0 Days"
	// }

	//mengurutkan value dari postform dan tag input
	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, start_date, end_date, description) VALUES ($1, $2, $3, $4)", projectName, startDate, endDate, description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently) //untuk meredirect kehalaman home.

}

// menampilkan page dan data
func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/project-detail.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var ProjectDetail = Project{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"]) //mengconvert string to integer dan menangkap dari params(id) dari url

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, description FROM tb_projects WHERE id=$1", id).Scan(&ProjectDetail.ID, &ProjectDetail.ProjectName, &ProjectDetail.Description)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	data := map[string]interface{}{ //variabel data
		"Project": ProjectDetail, //properti dan isinya
	}

	// fmt.Println(data)

	tmpl.Execute(w, data)

}

func deleteProject(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently) //untuk meredirect kehalaman home.

}

func updateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/update-project.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func submitUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently) //untuk meredirect kehalaman home.

}
