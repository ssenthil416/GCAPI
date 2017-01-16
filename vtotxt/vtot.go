package main

import (
	"fmt"
	"os"
	//    "io"
	"log"
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"time"
	//    "html/template"

	"sample/demo/caption"
)

const (
	DEFAULT    = "http://www.google.com"
	FACEBOOK   = "http://www.facebook.com"
	LINKEDIN   = "http://www.linkedin.com"
	YAHOO      = "http://www.yahoo.com"
	APPLE      = "http://www.apple.com"
)

type vtotxt struct {
	Stotxt string `json:"stotxt,omitempty"`
	Lpage  string `json:"lpage,omitempty"`
}

var data vtotxt

//Record from microphone to wav file
func rec() {
	os.Remove("test.wav")
	cmd := exec.Command("rec", "--encoding", "signed-integer", "--bits", "16", "--channels", "1", "--rate", "16000", "test.wav")
	st := time.Now()
	c := make(chan bool, 1)
	go func() {
		cmd.Run()
		<-c
		fmt.Println("Stopping Recording!!!!!!!!!!!")
	}()
	for {
		d := time.Since(st)
		if d.Seconds() > 3.0 {
			c <- true
			break
		}
	}
}

//Get URL based on content
func getURL(in string) {
	data.Lpage = DEFAULT

	if strings.Contains(in, "yahoo") {
		data.Lpage = YAHOO
		return
	}

	if strings.Contains(in, "linked in") {
		data.Lpage = LINKEDIN
		return
	}

	if strings.Contains(in, "apple") {
		data.Lpage = APPLE
		return
	}

	if strings.Contains(in, "facebook") {
		data.Lpage = FACEBOOK
		return
	}
}

//Laod page with new URL
func lpHandler(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, data.Lpage, http.StatusSeeOther)
}

//Index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	jdata, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return
	}

	//fmt.Println(string(jdata))
	w.Write(jdata)
}

//Loadpage webserver
func loadPage() {
	http.HandleFunc("/loadPage", lpHandler)
	if err := http.ListenAndServe(":8880", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//Index webserver 
func index() {
	http.HandleFunc("/index", indexHandler)
	if err := http.ListenAndServe(":9090", nil); if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//Main
func main() {
        //Initialise
	data = vtotxt{}
     
        // Start webserver
	go index()
	go loadPage()

        //Loop to record and covert to txt and load page
	for {
		//fmt.Println("Recording!!!!!!")
		rec()
		//fmt.Println("Speech to Text process")
		data.Stotxt = caption.Get("./test.wav")
		//fmt.Println("After Processing =", data.Stotxt)
		getURL(data.Stotxt)
 
		//cmd := exec.Command("curl","-X","POST","http://localhost:9090/index")
		//cmd.Run()

                //Sleep for 5 seconds
		st := time.Now()
		for {
			d := time.Since(st)
			if d.Seconds() > 5.0 {
				break
			}
		}
	}
}
