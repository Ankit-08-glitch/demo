package main

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetAllEmployee(t *testing.T) {
	db := Connect()
	if len(GetAllEmployee(db)) != 5 {
		t.Fail()
	}
}

func TestAddData(t *testing.T) {
	req, err := http.NewRequest("GET", "/add", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddData)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Fail()
	}
}

func TestDeleteData(t *testing.T) {
	req, err := http.NewRequest("GET", "/delete?empId=6", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteData)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Fail()
	}
}

func TestDeleteData2(t *testing.T) {
	req, err := http.NewRequest("GET", "/delete", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteData)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		fmt.Println(status)
		t.Fail()
	}
}

func TestGetData(t *testing.T) {
	req, err := http.NewRequest("GET", "/get", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetData)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Fail()
	}
}

func TestEditData(t *testing.T) {
	req, err := http.NewRequest("GET", "/edit?empId=8", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(EditData)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Fail()
	}
}

func TestEditData2(t *testing.T) {
	req, err := http.NewRequest("GET", "/edit", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(EditData)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Fail()
	}
}

func TestMain(m *testing.M) {
	v1 := viper.New()
	v1.SetConfigFile("logConfig.yaml")
	e := v1.ReadInConfig()
	PError(e)
	log.SetOutput(&lumberjack.Logger{
		Filename:   v1.GetString("filename"),
		MaxSize:    v1.GetInt("MaxSize"), // megabytes
		MaxBackups: v1.GetInt("MaxBackups"),
		MaxAge:     v1.GetInt("MaxAge"), //days
	})
	level, err := log.ParseLevel(v1.GetString("LogLevel"))
	PError(err)
	log.SetLevel(level)
	log.Debug("This is a test log entry")
	os.Exit(m.Run())
}

func TestCombinedFunctionality(t *testing.T) {
	records := CombinedFunctionality()
	if len(records) != 5 {
		t.Fail()
	}
}
