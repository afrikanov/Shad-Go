// +build !solution

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

var dataPath string

type AthleteStruct struct {
	Athlete string `json:"athlete"`
	Age     int    `json:"age"`
	Country string `json:"country"`
	Year    int    `json:"year"`
	Date    string `json:"date"`
	Sport   string `json:"sport"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

type AthleteInfo struct {
	Athlete       string               `json:"athlete"`
	Country       string               `json:"country"`
	Medals        MedalsAtYear         `json:"medals"`
	MedalsByYears map[int]MedalsAtYear `json:"medals_by_year"`
}

type MedalsAtYear struct {
	Gold   int `json:"gold"`
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
	Total  int `json:"total"`
}

var athletes []AthleteStruct

func parseJSON(w *http.ResponseWriter) {
	content, err := ioutil.ReadFile(dataPath)
	if err != nil {
		http.Error(*w, "Error", 400)
	}
	_ = json.Unmarshal(content, &athletes)
}

func searchByName(w *http.ResponseWriter, searchName string, sport string) AthleteInfo {
	var variants []AthleteStruct
	currentCountry := ""
	targetName := searchName
	i := 0
	for _, value := range athletes {
		if value.Athlete == targetName {
			if sport == "" || value.Sport == sport {
				if i == 0 {
					currentCountry = value.Country
					variants = append(variants, value)
				} else {
					if value.Country == currentCountry {
						variants = append(variants, value)
					}
				}
				i++
			}
		}
	}
	if len(variants) == 0 {
		http.Error(*w, "athlete not found", 404)
		(*w).WriteHeader(404)
		return AthleteInfo{}
	}
	var currentInfo AthleteInfo
	currentInfo.Athlete = targetName
	currentInfo.Country = currentCountry
	var medalsByYearMap = make(map[int]MedalsAtYear)
	for _, value := range variants {
		currentInfo.Medals.Gold += value.Gold
		currentInfo.Medals.Silver += value.Silver
		currentInfo.Medals.Bronze += value.Bronze
		currentInfo.Medals.Total += value.Total
		medalsByYearMap[value.Year] = MedalsAtYear{
			Gold:   medalsByYearMap[value.Year].Gold + value.Gold,
			Silver: medalsByYearMap[value.Year].Silver + value.Silver,
			Bronze: medalsByYearMap[value.Year].Bronze + value.Bronze,
			Total:  medalsByYearMap[value.Year].Total + value.Total,
		}
	}
	currentInfo.MedalsByYears = medalsByYearMap
	return currentInfo
}

func topInSports(w *http.ResponseWriter, currentSport string, limit int) []AthleteInfo {
	usedNames := make(map[string]bool)
	currentAthletes := make([]AthleteInfo, 0)
	for _, value := range athletes {
		if value.Sport == currentSport && !usedNames[value.Athlete] {
			currentAthletes = append(currentAthletes, searchByName(w, value.Athlete, value.Sport))
			usedNames[value.Athlete] = true
		}
	}
	if len(currentAthletes) == 0 {
		http.Error(*w, "sport not found", 404)
		(*w).WriteHeader(404)
		return make([]AthleteInfo, 0)
	}
	sort.Slice(currentAthletes, func(i, j int) bool {
		if currentAthletes[i].Medals.Gold == currentAthletes[j].Medals.Gold {
			if currentAthletes[i].Medals.Silver == currentAthletes[j].Medals.Silver {
				if currentAthletes[i].Medals.Bronze == currentAthletes[j].Medals.Bronze {
					return currentAthletes[i].Athlete < currentAthletes[j].Athlete
				}
				return currentAthletes[i].Medals.Bronze > currentAthletes[j].Medals.Bronze
			}
			return currentAthletes[i].Medals.Silver > currentAthletes[j].Medals.Silver
		}
		return currentAthletes[i].Medals.Gold > currentAthletes[j].Medals.Gold
	})
	if len(currentAthletes) < limit {
		return currentAthletes
	}
	return currentAthletes[:limit]
}

type TopCountriesInYear struct {
	Country string `json:"country"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

func topCountriesInYear(w *http.ResponseWriter, year int, limit int) []TopCountriesInYear {
	countriesMap := make(map[string]TopCountriesInYear)
	for _, value := range athletes {
		if value.Year == year {
			countriesMap[value.Country] = TopCountriesInYear{
				Country: value.Country,
				Gold:    countriesMap[value.Country].Gold + value.Gold,
				Silver:  countriesMap[value.Country].Silver + value.Silver,
				Bronze:  countriesMap[value.Country].Bronze + value.Bronze,
				Total:   countriesMap[value.Country].Total + value.Total,
			}
		}
	}
	if len(countriesMap) == 0 {
		http.Error(*w, "year not found", 404)
		(*w).WriteHeader(404)
		return make([]TopCountriesInYear, 0)
	}
	sortedCountries := make([]TopCountriesInYear, 0)
	countryUsed := make(map[string]bool)
	for _, value := range athletes {
		if !countryUsed[value.Country] {
			sortedCountries = append(sortedCountries, countriesMap[value.Country])
			countryUsed[value.Country] = true
		}
	}
	sort.Slice(sortedCountries, func(i, j int) bool {
		if sortedCountries[i].Gold == sortedCountries[j].Gold {
			if sortedCountries[i].Silver == sortedCountries[j].Silver {
				if sortedCountries[i].Bronze == sortedCountries[j].Bronze {
					return sortedCountries[i].Country < sortedCountries[j].Country
				}
				return sortedCountries[i].Bronze > sortedCountries[j].Bronze
			}
			return sortedCountries[i].Silver > sortedCountries[j].Silver
		}
		return sortedCountries[i].Gold > sortedCountries[j].Gold
	})
	sortedNonNull := make([]TopCountriesInYear, 0)
	for _, value := range sortedCountries {
		if value.Country != "" {
			sortedNonNull = append(sortedNonNull, value)
		}
	}
	if len(sortedNonNull) < limit {
		return sortedNonNull
	}
	return sortedNonNull[:limit]
}

func handler(w http.ResponseWriter, r *http.Request) {
	urlPath := "http://" + r.Host + r.URL.String()
	u, err := url.Parse(urlPath)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	parseJSON(&w)
	if len(q["name"]) > 0 {
		currentInfo := searchByName(&w, q["name"][0], "")
		if currentInfo.Athlete != "" {
			currentInfoJSON, _ := json.Marshal(currentInfo)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, _ = fmt.Fprintln(w, string(currentInfoJSON))
		}
	} else if len(q["sport"]) > 0 {
		limit := 3
		if len(q["limit"]) > 0 {
			limit, err = strconv.Atoi(q["limit"][0])
			if err != nil {
				http.Error(w, "invalid limit", 400)
				w.WriteHeader(400)
			}
		}
		currentTop := topInSports(&w, q["sport"][0], limit)
		if len(currentTop) > 0 {
			currentTopJSON, _ := json.Marshal(currentTop)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, _ = fmt.Fprintln(w, string(currentTopJSON))
		}
	} else if len(q["year"]) > 0 {
		limit := 3
		if len(q["limit"]) > 0 {
			limit, err = strconv.Atoi(q["limit"][0])
			if err != nil {
				http.Error(w, "invalid limit", 400)
				w.WriteHeader(400)
			}
		}
		year, _ := strconv.Atoi(q["year"][0])
		currentTop := topCountriesInYear(&w, year, limit)
		if len(currentTop) > 0 {
			currentTopJSON, _ := json.Marshal(currentTop)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, _ = fmt.Fprintln(w, string(currentTopJSON))
		}
	}
}

func main() {
	portPtr := flag.Int("port", 8000, "port string")
	dataPtr := flag.String("data", "", "data string")
	flag.Parse()
	portNumber := *portPtr
	dataPath = *dataPtr
	http.HandleFunc("/", handler)
	localAddress := "localhost:" + strconv.Itoa(portNumber)
	log.Fatal(http.ListenAndServe(localAddress, nil))
}
