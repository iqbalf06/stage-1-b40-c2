package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter() //variabel route dari mux router

	// route path folder public (js,css,images)
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	//routing, menjalankan function home
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/add-project", formAddProject).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/myproject-detail/{index}", myProjectDetail).Methods("GET")
	route.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET")

	fmt.Println("server running on port 5000")
	http.ListenAndServe("localhost:5000", route)

}

// menampilkan page dan data
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil { //untuk handle error. agar tau errornya
		w.Write([]byte("message : " + err.Error())) //byte untuk menampilkan string jika ada error
		return
	}

	//mengirimkan data string ke dalam interface
	response := map[string]interface{}{
		"Projects": dataProject,
	}

	tmpl.Execute(w, response) //menampilkan response dari views
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

// type data untuk variabel object
type Project struct {
	ProjectName string
	Description string
}

// deklarasi variabel global = array/slice
var dataProject = []Project{}

// function untuk memasukkan dan menangkap data
func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Project Name: " + r.PostForm.Get("inputProjectName")) //get value data dari tag input name
	// fmt.Println("Startdate: " + r.PostForm.Get("inputStartdate"))
	// fmt.Println("Enddate: " + r.PostForm.Get("inputEnddate"))
	// fmt.Println("Description: " + r.PostForm.Get("inputDescription"))

	//variabel object untuk menampung data dari tag input.
	var projectName = r.PostForm.Get("inputProjectName")
	// var startdate = r.PostForm.Get("innputStartdate")
	// var enddate = r.PostForm.Get("innputEnddate")
	var description = r.PostForm.Get("inputDescription")

	//pemanggilan type struct dan variabel global dan object diatas, sama seperti object di JS
	var newProject = Project{ //type struct dari Project
		ProjectName: projectName,
		Description: description,
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

func myProjectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/myproject-detail.html")

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
			}
		}
	}

	// //mengirimkan data string ke dalam interface
	// response := map[string]interface{}{
	// 	"ProjectName": "Dumbways Mobile Apps 2022",
	// 	"Index":       index,
	// }

	data := map[string]interface{}{
		"Project": ProjectDetail,
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
