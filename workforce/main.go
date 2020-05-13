package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type Role int

const (
	//Occupational Group, Non-specific
	Group Role = iota + 1
	//Generic Occupational Role
	Generic
	//Specific Occupational Role, very specific
	Specific
)

func (r Role) String() string {
	return fmt.Sprintf(" %v Role ", [...]string{"Group, Generic, Specific"}[r-1])
}

type Skill struct {
	Name      string `json: "name"`
	Shortcode string `json: "shortcode"`
	Role      Role   `role: "role"`
}

type Check int

const (
	API Check = iota + 1
	Photo
	Video
	Location
	Reviews
	Legal
)

func (c Check) String() string {
	return [...]string{"API, Photo, Video, Location, Reviews"}[c]
}

type Industry struct {
	Name      string          `json:"name"`
	Shortcode string          `json: "shortcode"`
	Tasks     map[Skill]Check `json:"tasks"`
}

func SkillHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//Lookup vars["shortcode"]
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "applicaton/json")

	//Tasks
	//if empty, return all available skill shortcodes
	//if param["shortcode"] exists, return data of the skill
	//encoded as JSON response := json.NewEncoder(w).Encode(skill)
	fmt.Fprintf(w, "%s", json.NewEncoder(w).Encode(availableSkill))
}

func IndustryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Industry Shortcode: %v\n", vars["shortcode"])
}

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Check Shortcode: %v\n", vars["shortcode"])
}

var availableSkill []Skill
var availableIndustries []Industry
var availableChecks []Check

func isDir(filename string) bool {
	_, err := ioutil.ReadFile(filename)
	if err != nil {
		return true
	}
	return false
}

func buildSkills() []Skill {
	files := []string{}
	skills := []Skill{}

	e := filepath.Walk(
		"skills",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			//fmt.Println(path, info.Size())
			//Exclude dir
			if !isDir(path) {
				files = append(files, path)
			}
			return nil
		},
	)
	if e != nil {
		log.Println(e)
	}

	for _, f := range files {
		fmt.Printf("Importing %s \n", f)
		//Deserialize YAML into Skill{}

		yamlFile, err := ioutil.ReadFile(f)

		if err != nil {
			log.Panic(err)
		}

		var currentSkill Skill

		e := yaml.Unmarshal(yamlFile, &currentSkill)

		if e != nil {
			log.Panic(e)
		}

		skills = append(skills, currentSkill)
	}

	return skills
}

func buildIndustries() []Industry {
	return []Industry{}
}

func buildChecks() []Check {
	return []Check{}
}

func main() {
	//Initialize

	availableSkill = buildSkills()
	// availableIndustries = buildIndustries()
	// availableChecks = buildChecks()

	//Setup Router
	r := mux.NewRouter()
	r.HandleFunc("/skills/{shortcode}", SkillHandler).Methods("GET")
	r.HandleFunc("/industries/{shortcode}", IndustryHandler).Methods("GET")
	r.HandleFunc("/checks/{shortcode}", CheckHandler).Methods("GET")
	http.Handle("/", r)

	//Start Server
	fmt.Printf("Starting Workforce Server ...\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
