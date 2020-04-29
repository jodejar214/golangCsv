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
	"sort"
	"time"
)

/**
Tests:
- given test files
- test files on git and get raw data somehow
- random urls that shouldn't give you data like google.com
- csv data with trailing and leading spaces
- even number of data points
- odd number of data points with median between same age value
- odd number of data points with median between diff age value

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
func getCsvInBatches(urlList []string) (map[int]int, map[int]string, []int) {
	finalCountMap := make(map[int]int, 0)
	finalNameMap := make(map[int]string, 0)
	finalAgeList := make([]int, 0)

	for _, csvUrl := range(urlList) {
		countMap, nameMap, ageList, rerr := retrieveCsvDataFromUrl(csvUrl)
		if rerr != nil {
			log.Println("Error retrieving data:", rerr.Error())
			continue
		}
		
		for k, v := range(countMap) {
			if val, found := finalCountMap[k]; found {
				finalCountMap[k] = val + v
			} else {
				finalCountMap[k] = v
			}
		}

		for k, v := range(nameMap) {
			if _, found := finalNameMap[k]; !found {
				finalNameMap[k] = v
			}
		}

		finalAgeList = append(finalAgeList, ageList...)
	}
	return finalCountMap, finalNameMap, finalAgeList
}

//calls the given url and reads the data in the response
func retrieveCsvDataFromUrl(csvUrl string) (map[int]int, map[int]string, []int, error) {
	log.Println("Retrieving data from: ", csvUrl)
	resp, err := http.Get(csvUrl)
	if err != nil {
		log.Println("Error with HTTP request:",  err.Error())
		return nil, nil, nil, err
	}

	defer resp.Body.Close()

	//check if http returned data succesfully
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Println("HTTP request returned:", resp.StatusCode)
		err := errors.New("HTTP request did not return csv data successfully.")
		return nil, nil, nil, err
	}

	//read csv
	reader := csv.NewReader(resp.Body)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, nil, nil, err
	}

	ageCount, ageToName, ageList, oerr := organizeData(data)
	return ageCount, ageToName, ageList, oerr
}

//organize data from http request into maps
func organizeData(data [][]string) (map[int]int, map[int]string, []int, error) {
	ageCount := make(map[int]int, 0)
	ageToName := make(map[int]string, 0)
	ageList := make([]int, 0)

	if len(data[0]) != 3 {
		err := errors.New("Csv file data is not formatted correctly.")
		log.Println("Error organizing data from csv:", err.Error())
		return nil, nil, nil, err
	}

	for i, row := range(data[1:]) {
		if len(row) != 3 {
			log.Println("Row", i+1, "is not valid format for data.")
			continue
		}

		//parse data
		fname := strings.Trim(row[0], " ")
		lname := strings.Trim(row[1], " ")
		name := fname + " " + lname
		ageStr := strings.Trim(row[2], " ")
		age := 0

		//check if age is valid integer
		if val, serr := strconv.Atoi(ageStr); serr == nil {
			age = val
		} else {
			log.Println("Row", i+1, "does not have a valid value for age.")
			continue
		}

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

		ageList = append(ageList, age)
	}

	return ageCount, ageToName, ageList, nil
}

//calulate stats to be outputted
func calculateStats(countMap map[int]int, nameMap map[int]string, ageList []int) {
	log.Println("----------------Results-----------------")
	getAverageAge(countMap)
	getMedianAgeAndName(countMap, nameMap, ageList)
}

//calculate average for aggregated dataset
func getAverageAge(countMap map[int]int) {
	sum := 0
	totalCount := 0
	for k, v := range(countMap) {
		sum += k * v
		totalCount += v
	}
	avg := sum / totalCount
	log.Println("The average age is:", avg)
}

//find median age and name for aggregated dataset
func getMedianAgeAndName(countMap map[int]int, nameMap map[int]string, ageList []int) {
	//find median age index
	totalCount := 0
	keyList := make([]int,0)
	for k, v := range(countMap) {
		totalCount += v
		keyList = append(keyList, k)
	}

	//sort keys to get median value
	sort.Ints(keyList)
	sort.Ints(ageList)
	log.Println(keyList)
	log.Println(ageList)
	log.Println(countMap)
	//check if median age is value in dataset
	if totalCount % 2 == 1 {
		//median is the existing middle value of dataset
		curInd := 0
		medInd := totalCount / 2
		medAge := 0
		log.Println("Total Count:",totalCount)
		log.Println("Median Index:",medInd)
		for _, k := range(keyList) {
			if curInd >= medInd {
				medAge = k
				break
			} else {
				val := countMap[k]
				curInd += val
			}
		}
		log.Println(ageList[medInd])
		log.Println("The median age is:", medAge)
		log.Println("A name corresponding to the median age is:", nameMap[medAge])
	} else {
		//median is in between two values of dataset and may not exist in dataset
		curInd := 0
		medIndLow := (totalCount / 2) - 1
		medAgeLow := 0
		medAgeHigh := 0

		log.Println("Total Count:",totalCount)
		log.Println("Median Index Low:",medIndLow)
		for i, k := range(keyList) {
			log.Println("Current Index:",curInd)
			if curInd == medIndLow {
				medAgeLow = k
				medAgeHigh = keyList[i+1]
				break
			} else if curInd > medIndLow {
				medAgeLow = keyList[i-1]
				medAgeHigh = k
				break
			} else {
				val := countMap[k]
				curInd += val
			}
		}
		log.Println(ageList[medIndLow])
		log.Println("The median low is:", medAgeLow)
		log.Println(ageList[medIndLow+1])
		log.Println("The median high is:", medAgeHigh)
		log.Println((ageList[medIndLow] + ageList[medIndLow+1])/2.0)
		medAge := (medAgeLow + medAgeHigh) / 2.0
		log.Println("The median age is:", medAge)
		if  val, found := nameMap[medAge]; found {
			log.Println("A name corresponding to the median age is:", val)
		} else {
			log.Println("A name cannot be determined since the median age is not an age in the dataset.")
		}
	}
}

func processData() {
	urlList, rerr := readInputFile()
	if rerr != nil {
		panic(rerr)
	}

	finalCountMap, finalNameMap, finalAgeList := getCsvInBatches(urlList)
	calculateStats(finalCountMap, finalNameMap, finalAgeList)
}

func main() {
	start := time.Now()
	processData()
	totalTime := time.Now().Sub(start)
	log.Println("Total Runtime: ", totalTime)
}