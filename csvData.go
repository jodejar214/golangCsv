package main

import (
	"encoding/csv"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	// "sync"
	"time"
)

/**
Tests:
- given test files
- test files on git and get raw data somehow
- random urls that shouldn't give you data like google.com
- csv data with trailing and leading spaces

Error Cases:
- file or url does not exist
- returned data is not csv format
- returned data is not in fname,lname,age format
- returned data is empty


Error Handling:
- log errors & ignore

Think About:
- how to retrieve list of urls or files to read? --> file of urls
- concurrency = how many files to do concurrently? --> find out how many files and divide by memory size?
- how to store data and use to get stats? --> 2 maps = (age -> count, age -> name) / have one for each file and compile into overall 
- how to aggregate concurrency results from each file? 
- how to deal with median value in between two data points? --> check if total records is even or odd, if even get
- keep track of urls you could not get data from
**/

//reads input file of urls to call to get csv files
func readInputFile() ([]string, error){
	//check for file argument
	if len(os.Args) < 2 {
		err := errors.New("Missing required argument! Please provide file containing list of urls!")
        log.Println(err.Error())
        return nil, err
    }

    //read file
    log.Println("Attempting to read file...")
    data, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        log.Println("Error reading file:", err.Error())
        return nil, err
    }

    //check if data exists
    if len(data) == 0 {
    	err := errors.New("File of urls does not contain any data.")
        log.Println(err.Error())
        return nil, err
    }

    //make into list to iterate later
    urlList := strings.Split(string(data), "\n")
    log.Println("Successfully retrieved list of urls. Need to retrieve", len(urlList), "csv files.")
    return urlList, nil
}

//get data from urls in batches using concurrency
func getCsvInBatches(urlList []string) {

}

//calls the given url and reads the data in the response
func retrieveCsvDataFromUrl(csvUrl string) (map[string]int, map[string]string, error) {
	log.Println("Retrieving data from: ", csvUrl)
	resp, err := http.Get(csvUrl)
	if err != nil {
		log.Println("Error with http request:",  err.Error())
		return nil, nil, err
	}

	defer resp.Body.Close()

	//check if http returned data succesfully
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Println("Http request returned:", resp.StatusCode)
		return nil, nil, err
	}

	//read csv
	reader := csv.NewReader(resp.Body)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	ageCount, ageToName := organizeData(data)
	return ageCount, ageToName, nil
}

//organize data from http request into maps
func organizeData(data [][]string) (map[string]int, map[string]string) {
	ageCount := make(map[string]int, 0)
	ageToName := make(map[string]string, 0)
	for i, row := range(data) {
		//check if csv data is valid
		if len(row) != 3 {
			log.Println("Row " + string(i) + " is not valid format for data.")
			continue
		}
		if _, serr := strconv.Atoi(row[2]); serr != nil {
			log.Println("Row " + string(i) + " does not have a valid value for age.")
			continue
		}

		//parse data
		fname := strings.Trim(row[0], " ")
		lname := strings.Trim(row[1], " ")
		name := fname + " " + lname
		age := strings.Trim(row[2], " ")

		//update count for given age
		if val, found := ageCount[age]; found {
			ageCount[age] = val + 1
		} else {
			ageCount[age] = 1
		}

		//only set name associated to age on first occurence of age
		if _, found := ageToName[age]; !found {
			ageToName[age] = name
		}
	}

	return ageCount, ageToName
}

//aggregate csv data from each request into overall maps for final stat calculations
func aggregateData() {

}

//calculates average from data compiled from all csvs
func getAverageAge() {

}

//finds median age and name from data compiled from all csvs
func getMedianAgeAndName() {

}

func main() {
	start := time.Now()

	urlList, rerr := readInputFile()
	if rerr != nil {
		panic(rerr)
	}

	countMap, nameMap, gerr := retrieveCsvDataFromUrl(urlList[0])
	if gerr != nil {
		panic(gerr)
	}

	log.Println(countMap)
	log.Println(nameMap)

	totalTime := time.Now().Sub(start)
	log.Println("Total Runtime: ", totalTime)
}