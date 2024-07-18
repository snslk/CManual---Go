package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os/exec"
	"net/http"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234s"
	dbname   = "pucsd"
)

var db *sql.DB
var tmpl = template.Must(template.ParseFiles("template/login.html"))
var filetmpl = template.Must(template.ParseFiles("template/file.html"))
var signuptmpl = template.Must(template.ParseFiles("template/signup.html"))

//**********************************************************************************

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}


	fmt.Println("Successfully connected to the database!")

	r := mux.NewRouter()
	r.HandleFunc("/", loginHandler).Methods("GET", "POST")
	r.HandleFunc("/file",fileHandler).Methods("GET","POST")
	r.HandleFunc("/signup",signupHandler).Methods("GET","POST")

	log.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
//***************************************************************

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("name")
		password := r.FormValue("rollno")

		if checkUser(username, password)==true {
			http.Redirect(w,r,"/file?username="+username,http.StatusSeeOther)
		} else {
			fmt.Fprintf(w, "User not found or incorrect password %s.",username)
		}
	} else {
		tmpl.Execute(w, nil)
	}
}

//******************************************************************

func checkUser(username, password string) bool {
	var dbUsername, dbPassword string
	query := "select userid,name from users where userid=$1 and name=$2 ;"
	err:= db.QueryRow(query,password,username).Scan(&dbPassword,&dbUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			return false 
		}
		log.Fatalf("Error querying database: %v", err)
	}

	return dbPassword == password
}

//*******************************************************************
func fileHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost{
		name := r.FormValue("text")
		
		cmd:= exec.Command("man","-a",name)

		output,err := cmd.Output()

		if err != nil{
		fmt.Fprintf(w,"<h1>No Manpage for the argument\n</h1> <a href=/file?username=%s> search other </a>",http.StatusSeeOther)


			fmt.Println("Error:",err)
			return
		}

		fmt.Fprintf(w,"%s",string(output))
	}else{	
		filetmpl.Execute(w,nil)
	}
}
//******************************************************************
func signupHandler(w http.ResponseWriter ,r *http.Request){
	
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		userid := r.FormValue("rollno")
		email := r.FormValue("email")
		mobile := r.FormValue("number")

	if insertUser( name,userid,email,mobile )==true {
		fmt.Fprintf(w,"<a href=/>login</a>")
		return
	}else{
		fmt.Fprintf(w,"Error signin\n")
	}
}	else {
		signuptmpl.Execute(w, nil)
	}
}

//************************************************************************
func insertUser(name,userid,email,mobile string) bool {
	query := "insert into users values($1,$2,$3,$4);"
	_,err:= db.Query(query,userid,name,email,mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			return false 
		}
		log.Fatalf("Error querying database: %v", err)
	}

	return true
}

//************************************************************************
