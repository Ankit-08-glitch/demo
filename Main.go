package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type Employee struct {
	EmpId      int    `json:"empId"`
	Name       string `json:"name"`
	Department string `json:"department"`
}

type Employees struct {
	Employees []Employee `json:"emp"`
}

//Connect to database
func Connect() *gorm.DB {
	v1 := viper.New()
	v1.SetConfigFile("config.yaml")
	e := v1.ReadInConfig()
	PError(e)
	dsn := "host=" + v1.GetString("database.hostname") + " user=" + v1.GetString("database.user") + " password=" + v1.GetString("database.password") + " dbname=" + v1.GetString("database.dbname") + " port=" + v1.GetString("database.port") + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	PError(err)
	return db
}

//GetData get details of all employees
func GetData(w http.ResponseWriter, r *http.Request) {
	PError(json.NewEncoder(w).Encode(GetAllEmployee(Connect())))
}

// AddData add details of user
func AddData(w http.ResponseWriter, r *http.Request) {
	var emp Employees
	db := Connect()
	byteValue, _ := ioutil.ReadFile("emp.json")
	IError(json.Unmarshal(byteValue, &emp))
	db.Create(&emp.Employees)
}

//EditData edit details of user
func EditData(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	keys, ok := r.URL.Query()["empId"]
	if !ok || len(keys[0]) < 1 {
		log.Warn("Url Param 'key' is missing")
		return
	}
	db.Model(&Employee{}).Where("emp_id = ?", keys[0]).Update("name", "Aditya")
}

// DeleteData delete a users details
func DeleteData(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	keys, ok := r.URL.Query()["empId"]
	if !ok || len(keys[0]) < 1 {
		log.Warn("Url Param 'key' is missing")
		return
	}
	db.Delete(&Employee{}, keys[0])
}

//main function
func main() {
	v1 := viper.New()
	v1.SetConfigFile("config.yaml")
	e := v1.ReadInConfig()
	PError(e)
	http.HandleFunc("/get", GetData)
	http.HandleFunc("/add", AddData)
	http.HandleFunc("/delete", DeleteData)
	http.HandleFunc("/edit", EditData)
	http.HandleFunc("/upload", uploadFile)
	log.Fatal(http.ListenAndServe(v1.GetString("server.port"), nil))
}

// IError function to handle low level errors
func IError(err error) {
	if err != nil {
		log.Info(err)
	}
}

// PError function to handle high level errors
func PError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// CError function to handle average level errors
func CError(err error) {
	if err != nil {
		log.Error(err)
	}
}

// GetAllEmployee function to get details of all employees
func GetAllEmployee(db *gorm.DB) []Employee {
	var records []Employee
	db.Find(&records)
	return records
}

func CombinedFunctionality() []Employee {
	req, err := http.NewRequest("GET", "/add", nil)
	CError(err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddData)
	handler.ServeHTTP(rr, req)
	req, err = http.NewRequest("GET", "/edit?empId=2", nil)
	CError(err)
	handler = http.HandlerFunc(EditData)
	handler.ServeHTTP(rr, req)
	req, err = http.NewRequest("GET", "/edit", nil)
	CError(err)
	handler = http.HandlerFunc(EditData)
	handler.ServeHTTP(rr, req)
	req, err = http.NewRequest("GET", "/delete?empId=6", nil)
	CError(err)
	handler = http.HandlerFunc(DeleteData)
	handler.ServeHTTP(rr, req)
	req, err = http.NewRequest("GET", "/delete", nil)
	CError(err)
	handler = http.HandlerFunc(DeleteData)
	handler.ServeHTTP(rr, req)
	db := Connect()
	return GetAllEmployee(db)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	er := r.ParseMultipartForm(10 << 20)
	CError(er)
	file, handler, err := r.FormFile("file")
	CError(err)
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	data := csv.NewReader(file)
	record, err := data.Read()
	fmt.Println(record)
	for {
		record, err = data.Read()
		fmt.Println(record)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		for value := range record {
			fmt.Printf("%s\n", record[value])
		}
	}
}
