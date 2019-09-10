package main

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"regexp"
)

const uploadDir = "uploads"
const uploadPage = `
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

var replaceChars = regexp.MustCompile("[^a-zA-Z0-9._-]")

func getUploadPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
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
			log.Printf("NextPart error: %v\n", err)
			return
		}

		_, params, err := mime.ParseMediaType(p.Header.Get("Content-Disposition"))
		filename, ok := params["filename"]
		if !ok {
			log.Println("missing filename in multipart")
			return
		}

		filename = replaceChars.ReplaceAllString(filename, "-")
		filename = uploadDir + "/" + filename

		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			log.Printf("open: %v\n", err)
			return
		}

		nwritten, err := io.Copy(f, p)
		if err != nil {
			log.Printf("%v: %v\n", filename, err)
			f.Close()
			return
		} else {
			log.Printf("Uploaded %v: %v bytes\n", filename, nwritten)
		}

		f.Close()
	}

	w.Header().Set("Cache-Control", "no-cache")
	io.WriteString(w, "ok\n")
}

func main() {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatal(err)
	}

	listenAt := ":8080"
	http.HandleFunc("/", getUploadPage)
	http.HandleFunc("/fup", postUploads)
	log.Printf("Start listening at %v\n", listenAt)
	log.Fatal(http.ListenAndServe(listenAt, nil))
}
