package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sort"
	"sync"
	"time"
)

/**
Error Cases:
- file or url does not exist
- returned data is not csv format
- returned data is not in fname,lname,age format
- age is not integer
- returned data is empty

Think About:
- how to retrieve list of urls or files to read? --> file of urls
- how to store data and use to get stats? --> 2 maps = (age -> count, age -> name) / have one for each file and compile into overall 
- how to deal with median value in between two data points? --> check if total records is even or odd, if even get

- concurrency = how many files to do concurrently? --> find out how many files and divide by memory size?
- how to aggregate concurrency results from each file? 
- keep track of urls you could not get data from
**/

//struct for results channels to pass as args
type CsvResults struct {
	CountMap map[int]int
	NameMap map[int]string
	AgeList []int
}

//reads input file of urls to call to get csv files
func readInputFile() []string {
	//check for file argument
	if len(os.Args) < 2 {
		err := errors.New("Missing required argument! Please provide file containing list of urls!")
        panic(err)
    }

    //read file
    fmt.Println("Attempting to read file...")
    data, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        panic(err)
    }

    //check if data exists
    if len(data) == 0 {
    	err := errors.New("File of urls does not contain any data.")
        panic(err)
    }

    //make into list to iterate later
    urlList := strings.Split(string(data), "\n")
    fmt.Println("Need to retrieve", len(urlList), "csv files.")
    return urlList
}

//get data from urls in batches using concurrency
func getCsvInBatches(urlList []string) (chan *CsvResults, chan string) {
	//limit number of workers to 100
	maxNumOfGoRoutines := len(urlList)
	if len(urlList) > 100 {
		maxNumOfGoRoutines = 100
	}

	//objects for goroutines
	resultsChan := make(chan *CsvResults, len(urlList))
	badCsvsChan := make(chan string, len(urlList))
	var wg sync.WaitGroup

	//process csvs concurrently
	numGoroutines := 0
	numUrlsProcessed := 0
	for _, csvUrl := range(urlList) {
		wg.Add(1)
		go retrieveCsvDataFromUrl(csvUrl, resultsChan, badCsvsChan, &wg)
		numGoroutines++
		numUrlsProcessed++
		if numGoroutines == maxNumOfGoRoutines || numUrlsProcessed == len(urlList) {
			wg.Wait()
			numGoroutines = 0
		}
	}
	close(resultsChan)
	close(badCsvsChan)
	return resultsChan, badCsvsChan
}

//calls the given url and reads the data in the response
func retrieveCsvDataFromUrl(csvUrl string, resultsChan chan *CsvResults, badCsvsChan chan string, wg *sync.WaitGroup){
	defer wg.Done()
	fmt.Println("Retrieving data from: ", csvUrl)
	resp, err := http.Get(csvUrl)
	if err != nil {
		fmt.Println("Error with HTTP request:",  err.Error())
		badCsvsChan <- csvUrl
		return
	}

	defer resp.Body.Close()

	//check if http returned data succesfully
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("HTTP request was not able to retrieve data. StatusCode:", resp.StatusCode)
		badCsvsChan <- csvUrl
		return
	}

	//read csv
	reader := csv.NewReader(resp.Body)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = 3
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading csv file:", err.Error())
		badCsvsChan <- csvUrl
		return
	}

	organizeData(data, resultsChan)
	return
}

//organize data from http request into maps
func organizeData(data [][]string, resultsChan chan *CsvResults) {
	results := &CsvResults{
		CountMap: make(map[int]int, 0),
		NameMap: make(map[int]string, 0),
		AgeList: make([]int, 0),
	}

	for i, row := range(data[1:]) {
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
			fmt.Println("Row", i+1, "does not have a valid value for age.")
			continue
		}

		//update count for given age
		if val, found := results.CountMap[age]; found {
			results.CountMap[age] = val + 1
		} else {
			results.CountMap[age] = 1
		}

		//only set name associated to age on first occurence of age
		if _, found := results.NameMap[age]; !found {
			results.NameMap[age] = name
		}

		results.AgeList = append(results.AgeList, age)
	}
	resultsChan <- results
	return
}

func aggregateResults(resultsChan chan *CsvResults, badCsvsChan chan string) (map[int]int, map[int]string, []int, []string) {
	//aggregate result objects
	finalCountMap := make(map[int]int, 0)
	finalNameMap := make(map[int]string, 0)
	finalAgeList := make([]int, 0)
	badCsvs := make([]string, 0)

	//aggregate data results
	for result := range(resultsChan) {
		for k, v := range(result.CountMap) {
			if val, found := finalCountMap[k]; found {
				finalCountMap[k] = val + v
			} else {
				finalCountMap[k] = v
			}
		}

		for k, v := range(result.NameMap) {
			if _, found := finalNameMap[k]; !found {
				finalNameMap[k] = v
			}
		}

		finalAgeList = append(finalAgeList, result.AgeList...)
	}

	//aggregate bad csvs
	for badCsv := range(badCsvsChan) {
		badCsvs = append(badCsvs, badCsv)
	}

	return finalCountMap, finalNameMap, finalAgeList, badCsvs
}

//calculate average for aggregated dataset
func getAverageAge(countMap map[int]int) {
	sum := 0
	totalCount := 0
	for k, v := range(countMap) {
		sum += k * v
		totalCount += v
	}
	avg := float64(sum) / float64(totalCount)
	fmt.Println("The average age is:", avg)
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

	//check if median age is value in dataset
	if totalCount % 2 == 1 {
		//median is the existing middle value of dataset
		curInd := 0
		medInd := totalCount / 2
		medAge := 0
		fmt.Println("Total Count:",totalCount)
		fmt.Println("Median Index:",medInd)
		for i, k := range(keyList) {
			if curInd == medInd {
				fmt.Println("Option1")
				fmt.Println("Current Index:",curInd)
				medAge = k
				break
			} else if curInd > medInd {
				fmt.Println("Option2")
				fmt.Println("Current Index:",curInd)
				medAge = keyList[i-1]
				break
			} else {
				val := countMap[k]
				curInd += val
			}
		}
		fmt.Println(ageList[medInd])
		fmt.Println("The median age is:", medAge)
		fmt.Println("A name corresponding to the median age is:", nameMap[medAge])
	} else {
		//median is in between two values of dataset and may not exist in dataset
		curInd := 0
		medIndLow := (totalCount / 2) - 1
		medIndHigh := medIndLow + 1
		medAgeLow := 0
		medAgeHigh := 0

		fmt.Println("Total Count:",totalCount)
		fmt.Println("Median Index Low:",medIndLow)
		fmt.Println("Median Index High:",medIndHigh)
		for i, k := range(keyList) {
			if curInd == medIndLow && curInd < medIndHigh {
				fmt.Println("Option1")
				fmt.Println("Current Index:",curInd)
				medAgeLow = k
				medAgeHigh = keyList[i+1]
				break
			} else if curInd > medIndLow && curInd == medIndHigh {
				fmt.Println("Option2")
				fmt.Println("Current Index:",curInd)
				medAgeLow = keyList[i-1]
				medAgeHigh = k
				break
			} else if curInd > medIndLow && curInd > medIndHigh {
				fmt.Println("Option3")
				fmt.Println("Current Index:",curInd)
				medAgeLow = keyList[i-1]
				medAgeHigh = keyList[i-1]
				break
			} else {
				val := countMap[k]
				curInd += val
			}
		}
		fmt.Println(ageList[medIndLow])
		fmt.Println("The median low is:", medAgeLow)
		fmt.Println(ageList[medIndHigh])
		fmt.Println("The median high is:", medAgeHigh)
		fmt.Println((ageList[medIndLow] + ageList[medIndHigh])/2.0)
		medAge := (medAgeLow + medAgeHigh) / 2.0
		fmt.Println("The median age is:", medAge)
		if  val, found := nameMap[medAge]; found {
			fmt.Println("A name corresponding to the median age is:", val)
		} else {
			fmt.Println("A name cannot be determined since the median age is not an age in the dataset.")
		}
	}
}

func processCsvData() {
	//read input file to get csv file to process
	urlList := readInputFile()

	//process the csv files
	resultsChan, badCsvsChan := getCsvInBatches(urlList)
	finalCountMap, finalNameMap, finalAgeList, badCsvs := aggregateResults(resultsChan, badCsvsChan)

	//Print out bad csvs
	if len(badCsvs) > 0 {
		fmt.Println("\n\nThese CSVs returned errors when attempting to retrieve data:")
		for _, c := range(badCsvs) {
			fmt.Println(c)
		}
	}

	//calculate stats and print out results
	if len(finalCountMap) == 0 {
		fmt.Println("No data was retrieved from the collection of csv files. Cannot caluclate average and median.")
	} else {
		fmt.Println("\n\n----------------Results-----------------")
		getAverageAge(finalCountMap)
		getMedianAgeAndName(finalCountMap, finalNameMap, finalAgeList)
	}
}

func main() {
	start := time.Now()
	processCsvData()
	totalTime := time.Now().Sub(start)
	fmt.Println("Total Runtime: ", totalTime)
}