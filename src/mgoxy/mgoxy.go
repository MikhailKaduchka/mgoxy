package main

//Simple proxy application

import (
	"fmt"
	"net/http"
	"os"
	"bufio"
	"log"
)

func main() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/go/", urlHandler)

	port := os.Getenv("PORT")

	//If port not set then use default port 8080
	if len(port) == 0 {
		port = "8080"
	}

	http.ListenAndServe(":" + port, nil)

	fmt.Println("Proxy started on port: ", port)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, templateStr)
}

func urlHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path[len("/go/"):]

	params := r.URL.RawQuery

	askPageUrl := "http://" + url

	if len(params) > 0 {
		askPageUrl += "?" + params
	}

	resp, err := http.Get(askPageUrl)
	if err != nil {
		panic("Couldn't reach url \" " + askPageUrl + "\":" + err.Error())
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	readPage := true

	for readPage {

		bytes := make([]byte, BUF_SIZE, BUF_SIZE)

		for i := 0; i < BUF_SIZE; i++ {
			line, err := reader.ReadByte()
			if err == nil {
				bytes[i] = line
			} else {
				readPage = false
				break
			}
		}

		w.Write(bytes)

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		} else {
			log.Println("Damn, no flush");
		}
	}

	log.Println("Got page \"" + askPageUrl + "\" from " + r.RemoteAddr)
}

const BUF_SIZE = 512

const templateStr = `
<html>
<head>
<title>Mini Go Proxy</title>
<script type="text/javascript">
    function redirect () {
        location.href = '/go/' + document.getElementById("url").value;
    };
</script>
</head>
<body>
    <div id="content" style="width:400px; margin:0 auto;">
<input id="url" type="text" size="50" required="">
<button id="launch" onclick="redirect();">go</button>
  </div>
</body>
</html>
`
