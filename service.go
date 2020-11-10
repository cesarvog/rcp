package main

import (
	"log"
	"net/http"
	"strings"
	"errors"
	"time"
	"io"
	"io/ioutil"
	"os"

	"github.com/gorilla/mux"
	"context"
	"github.com/go-redis/redis/v8"
)

const (
	ErrorMessage = "Something is wrong :("
)

var rclient *redis.Client
var ctx = context.Background()


func main() {
	//redis
	client, err := rclientNew()
	rclient = client;
	if err != nil {
		panic(err)
	}
	log.Println("Redis started")

	//service
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v", Copy).Methods(http.MethodPut)
	router.HandleFunc("/v", Paste).Methods(http.MethodGet)
	router.HandleFunc("/v", Erase).Methods(http.MethodDelete)
	log.Println("Service started")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func rclientNew() (*redis.Client, error) {
	url := os.Getenv("REDIS_URL")
	opt, _ := redis.ParseURL(url)
	rdb := redis.NewClient(opt)

	err := pingRedis(rdb)

	return rdb, err
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
	auth, err := getAuth(w, r)
	if err != nil {
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ErrorMessage, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = rclient.Set(ctx, auth, content, time.Hour).Err()
	if err != nil {
		http.Error(w, ErrorMessage, http.StatusInternalServerError)
	}
}


func Paste(w http.ResponseWriter, r *http.Request) {
	auth, err := getAuth(w, r)
	if err != nil {
		return
	}

	val, err := rclient.Get(ctx, auth).Result()
	if err != nil {
		http.Error(w, ErrorMessage, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, val)
}

func Erase(w http.ResponseWriter, r *http.Request) {
	auth, err := getAuth(w, r)
	if err != nil {
		return
	}

	err = rclient.Set(ctx, auth, "", time.Second).Err()
	if err != nil {
		http.Error(w, ErrorMessage, http.StatusInternalServerError)
	}
}

func pingRedis(client *redis.Client) error {
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}
