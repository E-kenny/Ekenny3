package main
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//Book type
type Book struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Author string `json:"author"`
	Publish_at string `json:"publish_at"`
}

var db *sql.DB

func loadBooks(w http.ResponseWriter, r *http.Request){
	var allbooks = []Book{}
	rows, err := db.Query("SELECT * FROM bless")
	if err != nil {
		panic(err)
	}else{
		fmt.Fprintln(w, "successfully selected")
	}
	defer rows.Close()
    defer db.Close()


	for rows.Next() {
		book:=Book{}
		err := rows.Scan(&book.ID,&book.Name,&book.Author,&book.Publish_at)
		if err != nil {
			panic(err)
		}
		allbooks = append(allbooks,book)

	}

	w.Header().Set("Content-Type","application/json")
	
	bk, err := json.Marshal(allbooks)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	
	w.Write(bk)
	
}

func createBook(w http.ResponseWriter, r *http.Request){
	insert , err := db.Prepare("INSERT INTO mypostgres ($1,$2,$3,$4)") 
	fmt.Fprintln(w, "successfully created")
	w.Header().Set("Content-Type","application/json")
	if err != nil {
		panic(err)
	}else{
		fmt.Fprintln(w, "successfully created")

	}
	  
	update , err := insert.Exec(1,"Ekene","Alchemist",time.Now().String())
	 if err != nil {
	 	panic(err)
	 }else {
	 	fmt.Fprintln(w,"d% created successfully",update)
 }
	defer insert.Close()
	defer db.Close()
	
}


func updateBook(w http.ResponseWriter, r *http.Request){
	queries := mux.Vars(r) 
		fmt.Fprintf(w, "Category: %v", queries["id"])
	    

	UpdateStatement := `
	   UPDATE postgres
	   SET Name = &1
	   WHERE ID = &2	
   `
	   UpdateResult, UpdateResultErr := db.Exec(UpdateStatement,"Mathew",1)
	   if UpdateResultErr != nil {
		   panic(UpdateResultErr)

	   }
	   UpdateRecord, UpdateRecordErr := UpdateResult.RowsAffected()
	   if UpdateRecordErr != nil {
		   panic(UpdateRecordErr)
	   }else{
		   fmt.Fprintln(w," Updated successfully",UpdateRecord )
	   }
		
	   defer db.Close()

			   
}


func deleteBook(w http.ResponseWriter, r *http.Request){	
		queries := mux.Vars(r) 
		fmt.Fprintf(w, "Category: %v\n", queries["id"])
	    
	
	DeleteStatement := `
		DELETE postgres
		WHERE ID = &1	
	`
		DeleteResult, DeleteResultErr := db.Exec(DeleteStatement,1)
		if DeleteResultErr != nil {
			panic(DeleteResultErr)

		}
		DeleteRecord, DeleteRecordErr := DeleteResult.RowsAffected()
		if DeleteRecordErr != nil {
			panic(DeleteRecordErr)
		}else{fmt.Printf("%d deleted successfully",DeleteRecord)}
		
		defer db.Close()
			
	
	}

	func logger(next http.HandlerFunc)http.HandlerFunc  {
		return func (w http.ResponseWriter,r *http.Request)  {
			fmt.Printf("This request was received from %v",r.URL)
			next(w,r)
		}
	}

func main()  {
	db, err := sql.Open("postgres", "user=postgres password=Ekenny2468 host=127.0.0.1 port=5432 dbname=postgres sslmode=disable")
if err != nil {
	panic(err)
}else{
	fmt.Println("The connection to the DB was successfully initialized ")
}


err = db.Ping()
  if err != nil {
    panic(err)
  }

  fmt.Println("Successfully connected!")

DBCreate :=` CREATE  TABLE  mypostgres
 (
	  
	id INT,
	name TEXT,
	author TEXT,
	publish_at TEXT UNIQUE NOT NULL
		  
 )
 WITH
(
	OIDS=FALSE
)
TABLESPACE pg_default;
ALTER TABLE	mypostgres
OWNER TO postgres;
`
 _, err = db.Exec(DBCreate)
if err != nil {
	panic(err)
}else{
	fmt.Println("The table was successfully created")
}
fmt.Println("The table was successfully created")

router :=mux.NewRouter()
api := router.PathPrefix("/api/v1").Subrouter()
api.HandleFunc("/Books",logger(loadBooks)).Methods(http.MethodGet)
api.HandleFunc("/Books",logger(createBook)).Methods(http.MethodPost)
api.HandleFunc("/Books/id/{id}",logger(updateBook)).Methods(http.MethodPatch)
api.HandleFunc("/Books/id/{id}",logger(deleteBook)).Methods(http.MethodDelete)
fmt.Println("server started successfully")
 log.Fatalln(http.ListenAndServe(":8080",router))
}

