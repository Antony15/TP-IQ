package main
 
import (
    "log"
    "regexp"
    "strings"
    "net/http"
    "io/ioutil"
    "text/template"	
    "gorilla/mux"
    "github.com/grokify/html-strip-tags-go"
)

var tmpl = template.Must(template.ParseGlob("template/*"))

type WordCount struct {
	Name string
	Count int
}
func wordCount(str string) map[string]int {
    wordList := strings.Fields(str)
    counts := make(map[string]int)
    for _, word := range wordList {
        _, ok := counts[word]
        if ok {
            counts[word] += 1
        } else {
            counts[word] = 1
        }
    }
    return counts
}

func Index(w http.ResponseWriter, r *http.Request){
	var wordArr []WordCount
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	if r.Method == "POST"{
		//~ url := "https://reuters.com"	 
		url := r.FormValue("url")	 
		req, _ := http.NewRequest("GET", url, nil)	 
		res, _ := http.DefaultClient.Do(req)	 
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)	
		strLine := string(body) 
		stripped := strip.StripTags(strLine)  
		for index,element := range wordCount(stripped){
			index = reg.ReplaceAllString(index, "")
			wordArr = append(wordArr,WordCount{Name:index,Count:element})
		}
	}   
    tmpl.ExecuteTemplate(w, "Index", wordArr)
} 
func main() {
    log.Println("Server started on: http://localhost:8080")
	r := mux.NewRouter().StrictSlash(true)
    r.HandleFunc("/", Index)
  	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}   
}
