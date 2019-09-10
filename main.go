package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var uploadPage = `
<html>
<head>
<title>Upload files</title>
</head>
<body>
<form method="post" action="fup" enctype="multipart/form-data">
  <input type="file" name="f" multiple>
  <input type="submit" value="Upload">
</form>
</body>
</html>
`

func getUploadPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(uploadPage))
}

func postUploads(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		log.Println(err)
		return
	}
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatal(err)
		}

		datas, err := ioutil.ReadAll(p)
		if err != nil {
			log.Fatal(err)
		}

		// TODO: save to file

		// Content-Disposition: [form-data; name="f"; filename="svg.svg"]
		for k, v := range p.Header {
			log.Printf("%v: %v", k, v)
		}

		log.Printf("Part %q\n", datas)
	}
}

func main() {
	listenAt := "localhost:8080"
	http.HandleFunc("/", getUploadPage)
	http.HandleFunc("/fup", postUploads)
	log.Printf("Start listening at %v\n", listenAt)
	log.Fatal(http.ListenAndServe(listenAt, nil))
}
