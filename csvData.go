package main

import (
	"encoding/csv"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	// "sync"
	"time"
)

/**
Tests:
- given test files
- test files on git and get raw data somehow
- random urls that shouldn't give you data like google.com

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
        log.Println("Error reading file: ", err.Error())
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
func getCsv(csvUrl string) (map[int]int, map[int]string, error) {
	resp, err := http.Get(csvUrl)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	log.Println(data)

	ageCount := make(map[int]int, 0)
	ageToName := make(map[int]string, 0)
	// csvRows := strings.Split(string(data), "\n")
	// for _, row := range(csvRows) {
	// 	log.Println(row)
	// }

	return ageCount, ageToName, nil
}

//compile csv data from each request into overall maps for final stat calculations
func compileData() {

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

	_, _, gerr := getCsv(urlList[0])
	if gerr != nil {
		panic(gerr)
	}

	totalTime := time.Now().Sub(start)
	log.Println("Total Runtime: ", totalTime)
}