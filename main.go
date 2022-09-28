package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"
	"time"

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
	route.HandleFunc("/project-detail/{index}", projectDetail).Methods("GET") //index url params
	route.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET")

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

// deklarasi variabel global = array/slice
var dataProject = []Project{}

// menampilkan page dan data
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil { //untuk handle error. agar tau errornya
		w.Write([]byte("message : " + err.Error())) //byte untuk menampilkan string jika ada error
		return
	}

	// memanggil variabel conn dari package connection
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, description FROM tb_projects")
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

	fmt.Println(result)

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
	nodejs := r.PostForm.Get("nodejs")
	react := r.PostForm.Get("react")
	java := r.PostForm.Get("java")
	python := r.PostForm.Get("python")

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

	//pemanggilan type struct dan variabel global dan object diatas, sama seperti object di JS
	var newProject = Project{ //type struct dari Project
		ProjectName: projectName,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		Nodejs:      nodejs,
		React:       react,
		Java:        java,
		Python:      python,
		Duration:    duration,
	}

	// untuk push / append data
	dataProject = append(dataProject, newProject) //penampung dan isi data

	fmt.Println(dataProject)

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

	index, _ := strconv.Atoi(mux.Vars(r)["index"]) //mengconvert string to integer dan menangkap dari params(id) dari url

	// perulangan/looping penampung data index dan data project
	for i, data := range dataProject {
		if i == index { //kondisi index looping = index url params
			ProjectDetail = Project{
				ProjectName: data.ProjectName,
				Description: data.Description,
				StartDate:   data.StartDate,
				EndDate:     data.EndDate,
				Duration:    data.Duration,
				Nodejs:      data.Nodejs,
				React:       data.React,
				Java:        data.Java,
				Python:      data.Python,
			}
		}
	}

	data := map[string]interface{}{ //variabel data
		"Project": ProjectDetail, //properti dan isinya
	}

	// fmt.Println(data)

	tmpl.Execute(w, data)

}

func deleteProject(w http.ResponseWriter, r *http.Request) {

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	dataProject = append(dataProject[:index], dataProject[index+1:]...)
	fmt.Println(dataProject)

	http.Redirect(w, r, "/", http.StatusMovedPermanently) //untuk meredirect kehalaman home.

}
