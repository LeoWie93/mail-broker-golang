package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"leowie93/go-mail-broker/internal/dialer"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type MailPost struct {
	Action string
	To     string
	Body   string
}

type Options struct {
	Template string
}

var optionsMap = map[string]*Options{
	"v1": &Options{Template: "v1.html"},
	"v2": &Options{Template: "v2.html"},
}

func exchangeAction(action string) (options *Options, err error) {
	if options, ok := optionsMap[action]; ok {
		return options, nil
	}

	return nil, fmt.Errorf("Given action is not valid: %s", action)
}

func handlePostMail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//unmarshall post
	decoder := json.NewDecoder(r.Body)
	var mailPost MailPost
	err := decoder.Decode(&mailPost)
	if err != nil {
		http.Error(w, "Invalid params", http.StatusUnprocessableEntity)
		return
	}

	//exchange action for options
	options, err := exchangeAction(mailPost.Action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	//TODO parse folders and subfiles
	tmpl := template.Must(template.ParseFiles(
		"./templates/components/base.html",
		"./templates/components/redHeader.html",
		"./templates/components/blueHeader.html",
		"./templates/"+options.Template))

	//TODO when include the .Body into a definition. The html tags are parsed out. Why?
	var buf bytes.Buffer
	tmpl.Execute(&buf, mailPost)

	//for debuggin
	s := buf.String()
	fmt.Println("the body: ")
	fmt.Println(s)

	// create Mail
	m := gomail.NewMessage()
	from := "noreply@h2g.ch"

	m.SetHeader("From", from)
	m.SetHeader("To", mailPost.To)
	m.SetHeader("Subject", "Todo define subject by action")

	m.SetBody("text/html", s)

	// send mail
	if err := dialer.Dialer.DialAndSend(m); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

// TODO is a init function usefull in this case?
func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		panic(err)
	}

	dialer.InitDialer()

	http.HandleFunc("/email-broker", handlePostMail)

	fmt.Println("Server starting on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
