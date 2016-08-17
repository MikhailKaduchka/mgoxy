package main

//Simple proxy application

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func main() {

	http.HandleFunc("/index.html", indexHandler)
	http.HandleFunc("/go/", urlHandler)

	http.ListenAndServe(":8080", nil)
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