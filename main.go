package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"leowie93/go-mail-broker/internal/options"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

var smtpDialer *gomail.Dialer

type MailPost struct {
	Action string
	To     string
	Body   template.HTML
}

func handlePostMail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// unmarshall post body
	var mailPost MailPost
	err := json.NewDecoder(r.Body).Decode(&mailPost)
	if err != nil {
		http.Error(w, "Invalid params", http.StatusUnprocessableEntity)
		return
	}

	options, err := options.ExchangeAction(mailPost.Action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	tmpl := template.Must(template.ParseGlob("./templates/components/*"))
	tmpl, _ = tmpl.ParseFiles("./templates/" + options.Template)

	var buf bytes.Buffer
	tmpl.Execute(&buf, mailPost)

	//TODO debug code
	bodyString := buf.String()
	fmt.Println("the body: ")
	fmt.Println(bodyString)

	m := gomail.NewMessage()
	from := "noreply@codingforest.ch"

	m.SetHeader("From", from)
	m.SetHeader("To", mailPost.To)
	m.SetHeader("Subject", "Todo define subject by action")

	m.SetBody("text/html", bodyString)

	if err := smtpDialer.DialAndSend(m); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func main() {

	//TODO implement my own .env parser
	if err := godotenv.Load(".env.local"); err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		panic(err)
	}
	smtpDialer = gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))

	mux := http.NewServeMux()

	postMailHandler := http.HandlerFunc(handlePostMail)
	mux.Handle("/email-broker", postMailHandler)

	fmt.Println("Server starting on port :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
