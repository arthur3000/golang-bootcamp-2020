package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Todo - Main structure, csv file will have this format.
type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Records - Global array to simulate a database
var Records []Todo

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/getdata", getAPIData)
	router.HandleFunc("/readcsv", printAllData)
	log.Fatal(http.ListenAndServe(":8001", router))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Wizeline challenge</h1><p>Golang Bootcamp! Let's go!</p>")
}

func getAPIData(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://jsonplaceholder.typicode.com/todos")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(data, &Records)
	if err != nil {
		fmt.Println(err)
	}
	csvFile, err := os.Create("./data.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)

	for _, todo := range Records {
		var row []string
		row = append(row, strconv.Itoa(todo.UserID))
		row = append(row, strconv.Itoa(todo.ID))
		row = append(row, todo.Title)
		row = append(row, strconv.FormatBool(todo.Completed))
		writer.Write(row)
	}
	writer.Flush()
	fmt.Fprintf(w, "Data CSV Updated!")
}

func printAllData(w http.ResponseWriter, r *http.Request) {
	data, err := readCsv("./data.csv")
	if err != nil {
		panic(err)
	}
	var todo Todo
	var todos []Todo
	for _, record := range data {
		todo.UserID, _ = strconv.Atoi(record[0])
		todo.ID, _ = strconv.Atoi(record[1])
		todo.Title = record[2]
		todo.Completed, _ = strconv.ParseBool(record[3])
		todos = append(todos, todo)
	}

	jsonData, err := json.MarshalIndent(todos, "", "    ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Fprintf(w, string(jsonData))
}

func readCsv(filename string) ([][]string, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return [][]string{}, err
	}
	return lines, nil
}

func main() {
	fmt.Println("Golang Challange - 2nd Deliverable")
	handleRequests()
}
