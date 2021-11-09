package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type General struct {
	Artists  string `json:"artists"`
	Location string `json:"locations"`
	Dates    string `json:"dates"`
	Relation string `json:"relation"`
}

type All struct {
	Id             int      `json:"id"`
	Image          string   `json:"image"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	Locations      string   `json:"locations"`
	ConcertDates   string   `json:"concertDates"`
	Relations      string   `json:"relations"`
	DatesLocations map[string][]string
}

type City struct {
	Index []Index `json:"index"`
}

type Index struct {
	ID             int64               `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Coordinate interface{}

func UnmarshalCity(data []byte) (City, error) {
	var r City
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *City) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", hello)
	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal(err)
		return
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}
	response, err := http.Get("https://groupietrackers.herokuapp.com/api")
	if err != nil {
		http.Error(w, "400 Bad Request", 400)
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	general := General{}
	err1 := json.Unmarshal(responseData, &general)
	if err1 != nil {
		http.Error(w, "500 Internal Server Error", 500)
		return
	}
	all := []All{}
	response1, errr := http.Get(general.Artists)
	if errr != nil {
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	responseData1, err2 := ioutil.ReadAll(response1.Body)
	if err2 != nil {
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	err3 := json.Unmarshal(responseData1, &all)
	if err3 != nil {
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	cityRes, errC := http.Get(general.Relation)
	if errC != nil {
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	responseDataRes, errC2 := ioutil.ReadAll(cityRes.Body)
	if errC2 != nil {
		http.Error(w, "500 Internal Server Error", 500)
		return
	}
	city, err := UnmarshalCity(responseDataRes)
	for i := 0; i < len(all); i++ {
		all[i].DatesLocations = city.Index[i].DatesLocations
	}
	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, all)
}
