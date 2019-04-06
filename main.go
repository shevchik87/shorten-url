package main

import (
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	b64 "encoding/base64"
)

type InputLongUrl struct {
	Url string
}

type ResponceUrl struct {
	LongUrl string
	ShortUrl string

}

const BASE_URL  = "http://127.0.0.1/"

func main() {
	r := mux.NewRouter()
	tmpl := template.Must(template.ParseFiles("shevchik87/shorten-url/index.html"))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}
		data := InputLongUrl{
			Url: r.FormValue("long-url"),
		}
		dataResponse := ResponceUrl{
			LongUrl: data.Url,
			ShortUrl: BASE_URL+ShortenUrl(data),

		}
		fmt.Println(dataResponse)
		tmpl.Execute(w, dataResponse)
	})

	http.ListenAndServe(":80", r)
}

func ShortenUrl(inputUrl InputLongUrl) string {

	hash := GetMD5Hash(inputUrl.Url)
	baseHash := b64.StdEncoding.EncodeToString([]byte(hash))
	return baseHash[:6]
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
