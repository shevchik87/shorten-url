package main

import (
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
	"crypto/md5"
	"encoding/hex"
	b64 "encoding/base64"
	"net/url"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/bson"
	"log"
	"fmt"
	"context"
	"os"
	"time"
)

type InputLongUrl struct {
	Url string
}

type ResponceUrl struct {
	LongUrl string
	ShortUrl string

}

type Url struct {
	OriginUrl string
	Hash string
	Created string
}



func main() {

	var baseTinyUrl = os.Getenv("BASE_TINY_URL")
	client, err := mongo.Connect(context.TODO(), os.Getenv("MONGODB_URL"))

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("shorten").Collection("urls")

	fmt.Println("Connected to MongoDB!")


	r := mux.NewRouter()
	tmpl := template.Must(template.ParseFiles("index.html"))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		data := InputLongUrl{
			Url: r.FormValue("long-url"),
		}

		shortHash, err := ShortenUrl(data)
		dataResponse := ResponceUrl{
			LongUrl: data.Url,
		}
		if err != nil {
			dataResponse.ShortUrl = shortHash

		} else {
			dataResponse.ShortUrl = baseTinyUrl+shortHash
			filter := bson.D{{"hash", shortHash}}
			var result Url

			err = collection.FindOne(context.TODO(), filter).Decode(&result)
			if err != nil {
				log.Println(err)
			}
			var currentTime = time.Now()
			if (Url{}) == result{
				dataToInsert := Url{dataResponse.LongUrl, shortHash, currentTime.Format("2006-01-02 15:04:05")}
				_, err := collection.InsertOne(context.TODO(), dataToInsert)
				if err != nil {
					log.Fatal(err)
				}

			}

		}

		tmpl.Execute(w, dataResponse)
	})
	r.HandleFunc("/{hash}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hash := vars["hash"]

		filter := bson.D{{"hash", hash}}
		var result Url

		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			w.WriteHeader(404)
			return
		}
		http.Redirect(w, r, result.OriginUrl, 301)

	})

	http.ListenAndServe(":80", r)
}

func ShortenUrl(inputUrl InputLongUrl) (string, error)  {

	_, err := url.ParseRequestURI(inputUrl.Url)
	if err != nil {
		return "No valid url", err
	}
	hash := GetMD5Hash(inputUrl.Url)
	baseHash := b64.StdEncoding.EncodeToString([]byte(hash))
	return baseHash[:6], nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
