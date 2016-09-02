package controllers_test

import (
	"github.com/vsukhin/numbers/controllers"
	"github.com/vsukhin/numbers/logger"
	"testing"
)

// TestLogger is a test logger implementation
type TestLogger struct {
}

// Println prints data to logger
func (logger *TestLogger) Println(args ...interface{}) {
}

// Fatalf prints and stops the execution
func (logger *TestLogger) Fatalf(query string, args ...interface{}) {
}

func TestNumbers(t *testing.T) {
	logger.Log = new(TestLogger)
}

func TestObjectContorllerImplementationParseQuery_MissingURLs(t *testing.T) {
	objectController := controllers.NewObjectControllerImplementation()

	query := make(map[string][]string)
	urls, err := objectController.ParseQuery(query)
	if err == nil {
		t.Error("ObjectContorllerImplementation ParseQuery should return error for missing URLs")
	}
	if urls != nil {
		t.Error("ObjectContorllerImplementation ParseQuery should return empty URLs")
	}
}

func TestObjectContorllerImplementationParseQuery_WrongURLs(t *testing.T) {
	objectController := controllers.NewObjectControllerImplementation()

	query := map[string][]string{"u": {"http:/localhost", "localhost"}}
	urls, err := objectController.ParseQuery(query)
	if err == nil {
		t.Error("ObjectContorllerImplementation ParseQuery should return error for wrong URLs")
	}
	if urls != nil {
		t.Error("ObjectContorllerImplementation ParseQuery should return empty URLs")
	}
}

func TestObjectContorllerImplementationParseQuery_OkURLs(t *testing.T) {
	objectController := controllers.NewObjectControllerImplementation()

	query := map[string][]string{"u": {"http://localhost", "http://www.example.com"}}
	urls, err := objectController.ParseQuery(query)
	if err != nil {
		t.Error("ObjectContorllerImplementation ParseQuery should not return error for URLs")
	}
	if urls == nil {
		t.Error("ObjectContorllerImplementation ParseQuery should not return empty URLs")
	}
}

func TestObjectContorllerImplementationParseData_WrongJSON(t *testing.T) {
	objectController := controllers.NewObjectControllerImplementation()

	body := "{, , ,}"
	object, err := objectController.ParseData([]byte(body))
	if err == nil {
		t.Error("ObjectContorllerImplementation ParseData should return error for object")
	}
	if object != nil {
		t.Error("ObjectContorllerImplementation ParseData should return empty object")
	}
}

func TestObjectContorllerImplementationParseData_OkJSON(t *testing.T) {
	objectController := controllers.NewObjectControllerImplementation()

	body := "{ \"numbers\": [1, 2, 3, 4, 5]}"
	object, err := objectController.ParseData([]byte(body))
	if err != nil {
		t.Error("ObjectContorllerImplementation ParseData should not return error for object")
	}
	if object == nil {
		t.Error("ObjectContorllerImplementation ParseData should not return empty object")
	}
}
