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

	route := mux.NewRouter()

	// path folder public (js,css,images)
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	//routing
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/add-project", formAddProject).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/myproject-detail/{name}", myProjectDetail).Methods("GET")

	fmt.Println("server running on port 5000")
	http.ListenAndServe("localhost:5000", route)

}

// menampilkan page dan data
func formAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/add-project.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

// function untuk memasukkan data
func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Project Name: " + r.PostForm.Get("inputProjectName")) //get value dari tag input name
	fmt.Println("Startdate: " + r.PostForm.Get("inputStartdate"))
	fmt.Println("Enddate: " + r.PostForm.Get("inputEnddate"))
	fmt.Println("Description: " + r.PostForm.Get("inputDescription"))

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

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

	id, _ := strconv.Atoi(mux.Vars(r)["name"])
	fmt.Println(id)

	response := map[string]interface{}{
		"Title":       "Hallo",
		"Id":          "id",
		"Description": "test description",
	}

	tmpl.Execute(w, response)
}
