package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Student struct {
	gorm.Model
	Name   string `json:"name"`
	Rollno string `json:"rollno"`
	City   string `json:"city"`
}

func allStudents(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var students []Student
	db.Find(&students)
	fmt.Println("{}", students)

	json.NewEncoder(w).Encode(students)
}

func singleStudent(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	rollno := vars["rollno"]

	var student Student
	db.Where("rollno = ?", rollno).Find(&student)
	json.NewEncoder(w).Encode(student)
}

func newStudent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New Student Endpoint Hit")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	rollno := vars["rollno"]
	city := vars["city"]

	fmt.Println(name)
	fmt.Println(rollno)
	fmt.Println(city)

	db.Create(&Student{Name: name, Rollno: rollno, City: city})
	fmt.Fprintf(w, "New Student Successfully Created")
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	rollno := vars["rollno"]

	var student Student
	db.Where("rollno = ?", rollno).Find(&student)
	db.Delete(&student)

	fmt.Fprintf(w, "Successfully Deleted Student")
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	rollno := vars["rollno"]
	city := vars["city"]

	var student Student
	db.Where("rollno = ?", rollno).Find(&student)

	student.Name = name
	student.City = city

	db.Save(&student)
	fmt.Fprintf(w, "Successfully Updated Student")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/students", allStudents).Methods("GET")
	myRouter.HandleFunc("/student/{rollno}", singleStudent).Methods("GET")
	myRouter.HandleFunc("/student/{rollno}", deleteStudent).Methods("DELETE")
	myRouter.HandleFunc("/student/{name}/{rollno}/{city}", updateStudent).Methods("PUT")
	myRouter.HandleFunc("/student/{name}/{rollno}/{city}", newStudent).Methods("POST")
	log.Fatal(http.ListenAndServe(":3002", myRouter))
}

func initialMigration() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Student{})
}

func main() {
	fmt.Println("Go ORM")
	initialMigration()

	handleRequests()
}
