package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	b64 "encoding/base64"

	qrcode "github.com/skip2/go-qrcode"
)

var data string = b64.StdEncoding.EncodeToString([]byte("wangcl"))
var chans = make(map[string](chan string))

func getLoginQRCode(w http.ResponseWriter, req *http.Request) {
	var png []byte
	scanedURL := fmt.Sprintf("http://%s/%s?uuid=%s", req.Host, "scaned", data)
	png, err := qrcode.Encode(scanedURL, qrcode.Medium, 256)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(png)))
	w.Write(png)
}

func scaned(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "scaned")
}

func longPull(w http.ResponseWriter, req *http.Request) {
	uuid := req.FormValue("uuid")
	c1 := make(chan string, 1)
	chans[uuid] = c1
	log.Println(uuid)
	select {
	case res := <-c1:
		fmt.Println(res)
		http.Redirect(w, req, "http://www.google.com/"+res, 301)
	case <-time.After(10 * time.Second):
		fmt.Println("timeout 10")
	}
	delete(chans, uuid)
}

func push(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pushed")
	uuid := req.FormValue("uuid")
	if val, ok := chans[uuid]; ok {
		select {
		case val <- "ok":
			fmt.Println(val)
		case <-time.After(10 * time.Second):
			fmt.Println("timeout 10")
		}
	}
}

func main() {
	http.HandleFunc("/getLoginQRCode", getLoginQRCode)
	http.HandleFunc("/scaned", scaned)
	http.HandleFunc("/longPull", longPull)
	http.HandleFunc("/push", push)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
