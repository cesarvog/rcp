package main

import (
	"log"
	"net/http"
	"strings"
	"errors"
	"time"
	"io/ioutil"
	"os"

	"github.com/gorilla/mux"
)

const (
	ErrorMessage = "Something is wrong :("
	TmpDir = "tmp/"
	FilePrefix = "f_"
)

func main() {
	//service
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v", Copy).Methods(http.MethodPut)
	router.HandleFunc("/v", Paste).Methods(http.MethodGet)
	router.HandleFunc("/v", Erase).Methods(http.MethodDelete)
	log.Println("Service started in port " + os.Getenv("PORT"))

	go PurgeRoutine()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func filename(id string) string {
	return TmpDir + FilePrefix + id
}

func getAuth(w http.ResponseWriter, r *http.Request) (string, error) {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Bearer" {
		http.Error(w, "Auth Error visit: ", http.StatusUnauthorized) //TODO link com docs
		return "", errors.New("auth error")
	}

	return auth[1], nil
}

func Copy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	auth, err := getAuth(w, r)
	if err != nil {
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ErrorMessage, http.StatusInternalServerError)
		return
	}
	err = ioutil.WriteFile(filename(auth), content, 0644)

	if err != nil {
		log.Println(err)
		http.Error(w, ErrorMessage, http.StatusInternalServerError)
	}
}


func Paste(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	auth, err := getAuth(w, r)
	if err != nil {
		return
	}

	content, err := ioutil.ReadFile(filename(auth))
	if err != nil {
		log.Println(err)
		http.Error(w, ErrorMessage, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func Erase(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	auth, err := getAuth(w, r)
	if err != nil {
		return
	}
	
	err = os.Remove(filename(auth))

	if err != nil {
		log.Println(err)
		http.Error(w, ErrorMessage, http.StatusInternalServerError)
	}
}

func MustDeleteFile(info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}

	if time.Now().Sub(info.ModTime()) < time.Hour {
		return false
	}

	if !strings.HasPrefix(info.Name(), FilePrefix) {
		return false
	}

	return true
}

func PurgeRoutine() {
	for {
		time.Sleep(time.Hour)

		files, err := ioutil.ReadDir(TmpDir)
		if err != nil {
			continue
		}

		for _, f := range files {
			if MustDeleteFile(f) {
				DeleteFile(f.Name())
			}
		}
	}
}

func DeleteFile(f string) {
	err := os.Remove(TmpDir + f)	
	if err != nil {
		log.Println("Erro deleting file", err.Error())
	}
}


