package main
 
import (
    "log"
    "strconv"
    "strings"
    "net/http"
    "text/template"	
    "gorilla/mux"
)
type ColumnFinder struct{
	ColNum int
	Column []string
}
var tmpl = template.Must(template.ParseGlob("template/*"))

// Search array by value & get its index
func array_search(element string, data []string) (int) {
   for k, v := range data {
       if element == v {
           return k
       }
   }
   return -1    //not found.
}

func Index(w http.ResponseWriter, r *http.Request){
	var CF ColumnFinder
	if r.Method == "POST"{
		x				:= 0
		Alb 			:= make([]string, 26)
		values 			:= make([]string, 18278)
        start 			:= strings.ToUpper(r.FormValue("start"))
        rows 			:= r.FormValue("rows")	
        row, _ 			:= strconv.Atoi(rows)	
        cols 			:= r.FormValue("columns") 
        col, _ 			:= strconv.Atoi(cols) 
        CF.ColNum		 = col
        total_values 	:= col*row
        log.Println(total_values)
        for i,j := 'A',0; i <= 'Z'; i++ {
			Alb[j] = string(i)			
			j++
		}
		for _,v := range Alb {
			values[x] = string(v)
			x++
		}
		for _,v := range Alb {
			for _,vv := range Alb {
				values[x] = string(v+vv)
				x++
			}
		}					
		for _,v := range Alb {
			for _,vv := range Alb {
				for _,vvv := range Alb {
					values[x] = string(v+vv+vvv)
					x++
				}
			}
		}
		start_index := array_search(start,values)
		if start_index >=0 {
		  for i:=start_index;i<=(start_index + (total_values - 1 )); i++ {	
			  CF.Column = append(CF.Column,values[i])		  
		  }			
		}
		log.Println(CF)	
		tmpl.AddFunc("inc", func(num int, step int) int {
			return num % step
		})	
		tmpl.ExecuteTemplate(w, "Index", CF)			   		
	} else {  
		tmpl.ExecuteTemplate(w, "Index", nil)
	}
} 

// main Function
func main() {
    log.Println("Server started on: http://localhost:8080")
	r := mux.NewRouter().StrictSlash(true)
    r.HandleFunc("/", Index)
  	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}   
}
