package main

//Simple proxy application

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
)

func main() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/index.html", indexHandler)
	http.HandleFunc("/go/", urlHandler)

	port := os.Getenv("PORT")

	http.ListenAndServe(":" + port, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, templateStr)
}

func urlHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path[len("/go/"):]

	resp, err := http.Get("http://" + url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	//TODO: Streaming mp3 support
	body, err := ioutil.ReadAll(resp.Body)

	w.Write(body)
}

const templateStr = `
<html>
<head>
<title>PROXYGO</title>
<script type="text/javascript">
    function redirect () {
        location.href = '/go/' + document.getElementById("url").value;
    };
</script>
</head>
<body>
<input id="url" type="text"/>
<button id="launch" onclick="redirect();">go</button>

</body>
</html>
`