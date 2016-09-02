package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/vsukhin/numbers/logger"
	"github.com/vsukhin/numbers/models"
)

const (
	timeout        = time.Millisecond * 500
	queryParameter = "u"
)

// ObjectControllerInterface is interface of object managing controller
type ObjectControllerInterface interface {
	CleverGet(w http.ResponseWriter, r *http.Request)
	StupidGet(w http.ResponseWriter, r *http.Request)
	ParseQuery(query map[string][]string) ([]string, error)
	ParseData(data []byte) (models.ByNumbers, error)
}

// ObjectControllerImplementation is implementation of object managing controller
type ObjectControllerImplementation struct {
}

var _ ObjectControllerInterface = &ObjectControllerImplementation{}

// NewObjectControllerImplementation is a constructor for ObjectControllerImplementation
func NewObjectControllerImplementation() (objectController ObjectControllerInterface) {
	return &ObjectControllerImplementation{}
}

// CleverGet returns a object ref. by id
func (objectcontroller *ObjectControllerImplementation) CleverGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		logger.Log.Println("Method is not supported: " + r.Method)
		http.Error(w, "not supported method", http.StatusNotFound)
		return
	}

	urls, err := objectcontroller.ParseQuery(r.URL.Query())
	if err != nil {
		http.Error(w, "not valid query parameters", http.StatusBadRequest)
		return
	}

	type item struct {
		result models.ByNumbers
		err    error
	}

	var result models.ByNumbers
	var lock sync.RWMutex

	data := make(chan item, len(urls))
	done := make(chan bool)

	go func(data chan item, done chan bool) {
		existing := make(map[int]bool)

		for {
			select {
			case <-done:
				return
			case it := <-data:
				lock.RLock()
				raw := make(models.ByNumbers, len(result))
				copy(raw, result)
				lock.RUnlock()

				for _, number := range it.result {
					if !existing[number] {
						raw = append(raw, number)
						existing[number] = true
					}
				}

				sort.Sort(raw)

				lock.Lock()
				result = raw
				lock.Unlock()
			}
		}
	}(data, done)

	ch := make(chan item, len(urls))
	for i := 0; i < len(urls); i++ {
		go func(url string) {
			client := http.Client{
				Timeout: timeout,
			}

			response, err := client.Get(url)
			if err != nil {
				ch <- item{err: err}
				return
			}
			data, err := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			if err != nil {
				ch <- item{err: err}
				return
			}

			result, err := objectcontroller.ParseData(data)
			if err != nil {
				ch <- item{err: err}
				return
			}

			ch <- item{result: result}
		}(urls[i])
	}

	var errors []error
	var count int
	var expired bool

	processtime := time.Millisecond * 15
	timeoutChannel := time.After(timeout - processtime)

	for {
		stop := false
		select {
		case <-timeoutChannel:
			done <- true
			close(data)
			expired = true
			stop = true
		case it := <-ch:
			count++
			if len(urls) == count {
				stop = true
			}

			if it.err != nil {
				errors = append(errors, it.err)
			} else {
				data <- it
			}
		}
		if stop {
			break
		}
	}

	if len(errors) != 0 {
		logger.Log.Println(objectcontroller.getCombinedError(errors))
	}

	if !expired {
		select {
		case <-timeoutChannel:
			done <- true
			close(data)
		}
	}

	lock.RLock()
	response := models.Object{Numbers: result}
	lock.RUnlock()
	err = objectcontroller.renderJSON(w, http.StatusOK, response)
	if err != nil {
		http.Error(w, "can't return JSON", http.StatusInternalServerError)
		return
	}
}

// StupidGet returns a object ref. by id
func (objectcontroller *ObjectControllerImplementation) StupidGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		logger.Log.Println("Method is not supported: " + r.Method)
		http.Error(w, "not supported method", http.StatusNotFound)
		return
	}

	urls, err := objectcontroller.ParseQuery(r.URL.Query())
	if err != nil {
		http.Error(w, "not valid query parameters", http.StatusBadRequest)
		return
	}

	type item struct {
		result models.ByNumbers
		err    error
	}

	ch := make(chan item, len(urls))
	for i := 0; i < len(urls); i++ {
		go func(url string) {
			client := http.Client{
				Timeout: timeout,
			}

			response, err := client.Get(url)
			if err != nil {
				ch <- item{err: err}
				return
			}
			data, err := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			if err != nil {
				ch <- item{err: err}
				return
			}

			result, err := objectcontroller.ParseData(data)
			if err != nil {
				ch <- item{err: err}
				return
			}

			ch <- item{result: result}
		}(urls[i])
	}

	var raw models.ByNumbers
	var errors []error
	var count int

	processtime := time.Millisecond * 25
	timeoutChannel := time.After(timeout - processtime)

	for {
		stop := false
		select {
		case <-timeoutChannel:
			stop = true
		case it := <-ch:
			count++
			if len(urls) == count {
				stop = true
			}

			if it.err != nil {
				errors = append(errors, it.err)
			} else {
				raw = append(raw, it.result...)
			}
		}
		if stop {
			break
		}
	}

	if len(errors) != 0 {
		logger.Log.Println(objectcontroller.getCombinedError(errors))
	}

	var result models.ByNumbers
	existing := make(map[int]bool)
	for _, number := range raw {
		if !existing[number] {
			result = append(result, number)
			existing[number] = true
		}
	}

	sort.Sort(result)

	response := models.Object{Numbers: result}
	err = objectcontroller.renderJSON(w, http.StatusOK, response)
	if err != nil {
		http.Error(w, "can't return JSON", http.StatusInternalServerError)
		return
	}
}

// ParseQuery parse query and validates parameters
func (objectcontroller *ObjectControllerImplementation) ParseQuery(query map[string][]string) ([]string, error) {
	if _, ok := query[queryParameter]; !ok {
		logger.Log.Println("URLs are not found")
		return nil, errors.New("Not found URLs")
	}

	urls := query[queryParameter]

	for _, rawurl := range urls {
		_, err := url.ParseRequestURI(rawurl)
		if err != nil {
			logger.Log.Println(err)
			return nil, err
		}
	}

	return urls, nil
}

// ParseData parse data
func (objectcontroller *ObjectControllerImplementation) ParseData(data []byte) (models.ByNumbers, error) {
	var result models.Object

	err := json.Unmarshal(data, &result)
	if err != nil {
		logger.Log.Println(err)
		return nil, err
	}

	return result.Numbers, nil
}

func (objectcontroller *ObjectControllerImplementation) getCombinedError(errs []error) error {
	var result []string

	for _, err := range errs {
		result = append(result, err.Error())
	}

	return errors.New(strings.Join(result, ";"))
}

// renderJSON renders object in JSON
func (objectcontroller *ObjectControllerImplementation) renderJSON(w http.ResponseWriter, code int, object interface{}) error {
	js, err := json.Marshal(object)
	if err != nil {
		logger.Log.Println(err)
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(js)

	return nil
}
