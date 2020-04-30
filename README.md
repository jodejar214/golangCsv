# Design Document

## Code Design

### How To Build/Run:

```
go run processData.go {file with list of urls}
```

***Example:***
```
go run processData.go urls.txt
```

---

### Implementation

I seperated the entire process into different functions that I thought grouped the different functionality well. Also, I wanted the different parts to be cleanly seperated into concise and distinct functions for easy debugging and readability.

#### Reading File With List of URLs:
For this part, the function would check the command line input and get the second argument to know what file to read for the list of urls. The program ignores any arguments following the second one. Once it knows what file to get it tries to read the file. If successful, it parses the data into a list of urls to be processed. This part will panic with an error if there is no second argument given in the command line, there is an error reading the file, or the file has no data.

#### Processing CSV Data Files:
To process the CSV, a GET http request is done to retrieve the data from the URL and if it successfully gets the data, the golang csv library is used to read the csv file retrieved from the URL. If the http request returned an error or a status code related to an unsuccessful request, the URL is saved to the list of files that failed to be read. The URL was also added to that list if the data returned is not a valid csv format readable by the golang csv library or does not have the required number of fields per record in the file.

If the data retrieved from the URL is readable and the correct format, the program would then iterate over the records and organize the values for each field. The fname and lname value are concatenated to form the full name and the age is converted into an integer. All values were trimmed to remove leading and trailing spaces. Furthermore, since csv files do not type check the values of the field, I added a check to see if the age value is of type integer. If not, the record would be skipped over and the program would continue processing the rest of the records in the file. From the record values, two maps would be formed - a map for the frequency of each age in the dataset and a map storing one name for each age in the dataset. These maps are stored in a CsvResult struct to make the code more compact by reducing the number of parameters and return values that had to be given in the function headers as well as make the results more organized by putting them into one object.

The program uses goroutines to concurrently retrieve the data from each URL and process the data. Each url corresponded to one goroutine. I limited the code to only allow a maximum of 100 goroutines to run concurrently. Once the maximum amount of goroutines are running the program waits until all those goroutines are done to start up the next batch of goroutines. This is done until all the URLs have been processed. Each goroutine would return a value into either the results channel which contains the CsvResult struct that stores the two maps constructed from the csv data or the bad csvs channel which contains the URLs that failed to be read. Once all the URLs have been processed, the values returned to the results channel are aggregated into a final map for each map and the URLs returned to the bad csvs channel are aggregated into a final list.

#### Calculating Statistics:
With all the csv data organized into a map, the program can now calculate the average and median age values. For the average, it iterates through the map of age frequencies to get the sum of the ages in the dataset and divides it by the total number of records. Before dividing, I converted the sum and total which were integer values into float64 values so that the program would return a more accurate average.

To find the median, there are several cases to keep track of. First, if the total number of records is even, then the median value is between two age values, but if the total number of records is odd, then the median value simply the age value in the middle when all age values are sorted. Also, since I use a frequency map rather than a list of individual age values, the program cannot simply iterate over each age value and return the at the middle index. To find the median, I had to sort the map of age frequencies in ascending order of age and keep track of the current index that was incremented by the frequency of the ages the program iterated over in the map. 

To find the median in the case of an odd amount of records, the current index was compared to the index the median is supposed to be in. There were three possible cases that would determine if the current index needed to be incremented or the median age was found.
1. Current Index == Median Index --> Median age is the current key in map
2. Current Index > Median Index --> Median age is the previous key in map
3. Current Index < Median Index --> Increment current index

To find the median in the case of an even amount of records, a similar method was used but instead of having just one median index, we needed to keep track of two indexes - the low median index and the high median index. This resulted in having four possible cases that would determine if the current index needed to be incremented or the two values needed to calculate the median age were found.
1. Current Index == Low Median Index & Current Index < High Median Index --> Median age is between current key and next key in map
2. Current Index > Low Median Index & Current Index == High Median Index --> Median age is between previous key and current key in map
3. Current Index > Low Median Index & Current Index > High Median Index --> Median age is the previous key in map
4. Current Index < Low Median Index & Current Index < High Median Index --> Increment current index

#### Printing Results:
When the median age is found, the program then gets the name corresponding to the age in the other map and prints the two values out along with the calculated average. The list of bad URLs that failed to be read are also iterated over and printed out along with the stats. Finally, the time it took to complete the entire program is printed at the end. All other debugging information printed out during testing was removed to speed up the program. Only information printed out besides the results are errors and which URL the program is processing so that the error can be matched to a URL.

---

### Testing

For testing my implementation, I first ran tests to check functionality of each part of the process - reading the input file with urls to retrieve, retrieving the data from the urls, organizing the data, calculating the average and median - on one csv dataset. I tested several successful cases on different data sets for each part as well as error cases for each part.

#### Reading File With List of URLs:
For testing the reading of the input file with the list of urls to process, I had uploaded the data files onto Github so that I could test retrieving the raw csv files over the network. I then listed a subset of the urls in the input file to test retrieving individual data files first and then retrieved multiple data files. Additionally, I listed urls for csv files that did not exist on Github as well as random urls that did not link to csv files such as www.google.com. I also listed strings that were not real urls such as random words and an empty string. These urls or strings were added to test the error handling of the http request and reading the response data. I also tested with urls that had leading and trailing spaces to ensure that these were trimmed and the data could be retrieved.

#### Processing CSV data files:
For csv data files, I used the given data files as well created my own test files. For my custom data files, I made both readable and unreadable csv files so that I could test both successful and error cases. For the readable files, I made a few that were small datasets to make verifying caluclations easier and to allow easier manipulation of the records to get the cases I wanted to test. Additionally, I created test datafiles for unusual cases such as when the csv file has the incorrect number of fields, has no data, has values with leading and trailing spaces, only has the field names in the data, has empty lines before and after showing the data, has the wrong fields, and has records where the value for the age field was not an integer to ensure that the program handled those correctly. 

#### Calculating Average and Median:
For some of the custom test files, I used smaller data sets to make it easier to verify the correctness of the calculated average and median values. I used datasets with both an even and odd amount of records so that I could test the two different cases of finding a median value. Additionally, for the datasets with an even amount of records I tested the cases where the median age is between the same age values and where the median age is between different age values. Also, since I used a frequency map for ages to calculate the median rather than using a list of the individual ages recorded in the dataset, I had to test the various cases that were used to determine the median age of the dataset when iterating through the map. I compared the median values calculated from the frequency map with the median values found by using a list with all the individual age values listed to verify the correctness. For the name printed out with the median age, I tested both the cases when the median age was in the map and when there was no name associated to the median age.

#### Sequential and Concurrent Processing:
Once I was sure that each part was working how I wanted it to for individual URLs, I added the functionality to handle data from multiple URLs. These multiple URLs consisted of a mix of csvs that were readable and unreadable to see if could handle the mix. First, I tested this with the program processing the data in the URLs sequentially. I performed these as initial end-to-end tests for both success and error cases to verify the correctness of the process and that the program could run completely. 

Once I was sure that my implementation was able to correctly handle multiple datasets sequentially, I added concurrency. I first tested only needing to process one url and then moved onto processing multiple URLs like I did when testing sequentially. Finally, with concurrency implemented and verified with a small amount of URLs, I wanted to ensure that the number of goroutines I chose to use as a limit could process a large amount of URLs quickly and correctly so I tested out several different max number of goroutines that could run concurrently. The dataset I did this load test with was file3.csv which has 10000 lines of code and had listed it a total of 1600 times so that it would have to be processed hundreds of times. After verifying several limits such as 50, 100, 150, 200, and 400 goroutines, I found that 100 goroutines concurrently resulted in the best time on average. Thus, I set that as the limit for concurrency in my program.

---

### Design Considerations

While implementing this program, I considered several options for various parts of the process. One item I considered was the number of goroutines to allow running concurrently. I chose a high number at first since it seems likely that the more goroutines the program has happening concurrently, the faster the program will finish. However, while testing various limits, the higher numbers I tested ended up being slower than the lower ones. From the tests, I ended up on settling with 100 goroutines as the limit. Another design possibility was rather than ignoring the age value that were the incorrect type in the readable csv files and continuing to process the remaining records in the file, I would consider the csv file an unreadable file and stop processing the data in it. However, most of the records could be fine and the incorrect type would not be found until that record is processed so I felt that since most of the data is fine, the file should still be considered readable. One other design consideration is keeping a list of the individual age along with the maps in the results. I already had it for testing and using it to find the median would be simpler than using frequency maps, but I thought the list would be extraneous and a waste of space so I did not keep it in the final implementation. One last consideration I had was to add more concurrency to the program such as aggregating the final results when each batch finishes rather than when all the URLs processed or adding it to the final result objects once the goroutine is finished. This would reduce the space required by the channels to hold the processed results especially with a large number of URLs to process. However, I felt that this added unnecessary complexity to the program and would not provide much benefit in increasing the runtime. It is easier to iterate over the results in the channel once all results are there and would not need to worry about using locks or sync maps to ensure that shared objects are not modified by multiple goroutines at the same time.

---

## Questions
1. What assumptions did you make in your design? Why?

  	One assumption I made is that the list of URLs will come in a file that has each URL on its own line. This is the most organized method I can think of to list the URLs especially if there are hundreds of URLs to process. Having the URLs listed with commas on one line would make the file difficult to read and organize when considering what URLs to list. Another assumption I made is that the URL can return any possible status code and any status code that is not 200-299 is not considered a success case. When handling the http response, I used this assumption to determine which status codes are for errors but there are too many and I know that the 200s are usually the successful requests while ones in the 400s and 500s are errors.

2. How would you change your program if it had to process many files where each file was over 10M records?

	If the files have that many records, I would change my program to better organize the records. Instead of iterating over each record one by one, I would use a divide and conquer approach to divide the 10M records into batches of thousands of records and process those batches concurrently. Then, once the results of all the batches are created, they would be aggregated as the results for the entire file. This would improve the speed at which a file is processed and decrease the runtime.

3. How would you change your program if it had to process data from more than 20K URLs?

	If there are that many URLs to process data from, I don't think I would need to change my program. The concurrency I have when processing the URLs in batches of 100 was shown to perform better than larger batches when doing my tests so I don't think changing the amount of work being done concurrently will improve the processing of the data.
