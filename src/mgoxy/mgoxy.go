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
  http.HandleFunc("/index.html", indexHandler)
  http.HandleFunc("/go/", urlHandler)

  port := os.Getenv("PORT")

  http.ListenAndServe(":" + port, nil)

  fmt.Println("Proxy started on port: ", port)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, templateStr)
}

func urlHandler(w http.ResponseWriter, r *http.Request) {

  url := r.URL.Path[len("/go/"):]

  resp, err := http.Get("http://" + url)
  if err != nil {
    log.Println("Couldn't reach url. Error: ", err)
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
        log.Println("Reading error: ", err)
        //Stop loading data
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
}

const BUF_SIZE = 512

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
