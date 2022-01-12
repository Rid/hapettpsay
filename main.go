package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	hapesay "github.com/Rid/hapesay/v2"
)

type Hapesay struct {
	// Type of hape
	typ string
	// What the hape should say
	say string
}

type htmlVars struct {
	Hapes map[int]map[string]string
}

func hapeList() []string {
	hapes, err := hapesay.Hapes()
	if err != nil {
		return hapesay.HapesInBinary()
	}
	list := make([]string, 0)
	for _, hape := range hapes {
		list = append(list, hape.HapeFiles...)
	}
	return list
}

func serveTemplate(w http.ResponseWriter, req *http.Request) {
	html := &htmlVars{Hapes: make(map[int]map[string]string, 0)}
	lp := filepath.Join("templates", "layout.html")

	hapes := hapeList()

	hape := &Hapesay{typ: "mobile", say: "Hello"}

	route := strings.Split(req.URL.Path, "/")

	if len(route) > 2 {
		hape.typ = route[1]
		hape.say = route[2]

		say, err := hapesay.Say(
			hape.say,
			hapesay.Type(hape.typ),
			hapesay.BallonWidth(40),
		)

		if err != nil {
			say, _ = hapesay.Say(
				"Error 404: Hape not found",
				hapesay.Type("mobile"),
				hapesay.BallonWidth(15),
			)
			w.WriteHeader(404)
		}

		html.Hapes[0] = map[string]string{hape.typ: say}
	} else {
		if route[1] == "" {
			hape.say = "You can make me say anything, just type it in the url"
		} else {
			hape.say = route[1]
		}

		say, _ := hapesay.Say(
			hape.say,
			hapesay.Type("mobile"),
			hapesay.BallonWidth(15),
		)
		html.Hapes[0] = map[string]string{"mobile": say}

		counter := 1
		for _, hapeFile := range hapes {
			if hapeFile == "mobile" {
				continue
			}
			say, _ := hapesay.Say(
				hape.say,
				hapesay.Type(hapeFile),
				hapesay.BallonWidth(15),
			)

			html.Hapes[counter] = map[string]string{hapeFile: say}
			counter++
		}
	}

	// res := strings.Split(say, "\n")
	// fmt.Fprintf(w, "%v", res[2])
	tmpl, _ := template.ParseFiles(lp)
	tmpl.ExecuteTemplate(w, "layout", html)
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := r.RemoteAddr
		fwdAddress := r.Header.Get("X-Forwarded-For")
		if fwdAddress != "" {
			// Got X-Forwarded-For
			ipAddress = fwdAddress // If it's a single IP, then awesome!

			// If we got an array... grab the first IP
			ips := strings.Split(fwdAddress, ", ")
			if len(ips) > 1 {
				ipAddress = ips[0]
			}
		}
		fmt.Printf("%s %s %s\n", ipAddress, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveTemplate)

	log.Println("Listening on :8000...")
	err := http.ListenAndServe(":8000", Log(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}
