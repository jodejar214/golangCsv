# Design Document

## Code Design

### How To Build/Run:
> go run csvData.go {file with list of urls}

***Example:***
go run csvData.go urls.txt

### Implementation

I seperated the entire process into different functions that I thought grouped the different functionality well. Also, I wanted the different parts to be cleanly seperated into concise and distinct functions for easy debugging and readability.

#### Reading File With List of URLs:
For this part, the function would check the command line input and get the second argument to know what file to read for the list of urls. The program ignores any arguments following the second one. Once it knows what file to get it tries to read the file. If successful, it parses the data into a list of urls to be processed. This part will panic with an error if there is no second argument given in the command line, there is an error reading the file, or the file has no data.

#### Processing CSV Data Files:
To process the CSV, a GET http request is done to retrieve the data from the URL and if it successfully gets the data, the golang csv library is used to read the csv file retrieved from the URL. If the http request returned an error or a status code related to an unsuccessful request, the URL is saved to the list of files that failed to be read. The URL was also added to that list if the data returned is not a valid csv format readable by the golang csv library or does not have the required number of fields per record in the file.

If the data retrieved from the URL is readable and the correct format, the program would then iterate over the records and organize the values for each field. The fname and lname value are concatenated to form the full name and the age is converted into an integer. All values were trimmed to remove leading and trailing spaces. Furthermore, since csv files do not type check the values of the field, I added a check to see if the age value is of type integer. If not, the record would be skipped over and the program would continue processing the rest of the records in the file. From the record values, two maps would be formed - a map for the frequency of each age in the dataset and a map storing one name for each age in the dataset. These maps are stored in a CsvResult struct to reduce the number of parameters and return values that had to be given in the function headers as well as make the results more organized into one object.

The program uses goroutines to concurrently retrieve the data from each URL and process the data. Each url corresponded to one goroutine. I limited the code to only allow a maximum of 100 goroutines to run concurrently. Once the maximum amount of goroutines are running the program waits until all those goroutines are done to start up the next batch of goroutines. This is done until all the URLs have been processed. Each goroutine would return a value into either the results channel which contains the CsvResult struct that stores the two maps constructed from the csv data or the bad csvs channel which contains the URLs that failed to be read. Once all the URLs have been processed, the values returned to the results channel are aggregated into a final map for each map and the URLs returned to the bad csvs channel are aggregated into a final list.

#### Calculating Statistics:
With all the csv data organized into a map, the program can now calculate the average and median age values. For the average, it iterates through the map of age frequencies to get the sum of the ages in the dataset and divides it by the total number of records. Before dividing, I converted the sum and total which were integer values into float64 values so that the program would return a more accurate average.

To find the median, there are several cases to keep track of. First, if the total number of records is even, then the median value is between two age values, but if the total number of records is odd, then the median value simply the age value in the middle when all age values are sorted. Also, since I use a frequency map rather than a list of individual age values, the program cannot simply iterate over each age value and return the at the middle index. To find the median, I had to sort the map of age frequencies in ascending order of age and keep track of the current index that was incremented by the frequency of the ages the program iterated over in the map. 

To find the median in the case of an odd amount of records, the current index was compared to the index the median is supposed to be in. There were three possible cases that would determine if the current index needed to be incremented or the median age was found.
1. Current Index == Median Index --> Median age is the current key in map
2. Current Index > Median Index --> Median age is the previous key in map
3. Current Index < Median Index --> Increment current index

To find the median in the case of an even amount of records, a similar method was used but instead of having just one median index, we needed to keep track of two indexes - the low median index and the high median index. This resulted in having four possible cases that would determine if the current index needed to be incremented or the two values needed to calculate the median age were founde.
1. Current Index == Low Median Index & Current Index < High Median Index --> Median age is between current key and next key in map
2. Current Index > Low Median Index & Current Index == High Median Index --> Median age is between previous key and current key in map
3. Current Index > Low Median Index & Current Index > High Median Index --> Median age is the previous key in map
4. Current Index < Low Median Index & Current Index < High Median Index --> Increment current index

#### Printing Results:
When the median age is found, the program then gets the name corresponding to the age in the other map and prints the two values out along with the calculated average. The list of bad URLs that failed to be read are also iterated over and printed out along with the stats. Finally, the time it took to complete the entire program is printed at the end.


### Testing
For testing my implementation, I first ran tests to check functionality of each part of the process - reading the input file with urls to retrieve, retrieving the data from the urls, organizing the data, calculating the average and median - on one csv dataset depending on the case I was testing. I tested several successful cases on different data sets for each part as well as error cases for each part.

#### Reading File With List of URLs:
For testing the reading of the input file with the list of urls to process, I had uploaded the data files onto Github so that I could test retrieving the raw csv files over the network. I then listed a subset of the urls in the input file to test retrieving individual data files as well as multiple data files. Additionally, I listed urls for csv files that did not exist on Github as well as random urls that did not link to csv files such as "www.google.com". I also listed strings that were not real urls such as random words and an empty string. These urls or strings were added to test the error handling of the http request and reading the response data.

#### Processing CSV data files:
For csv data files, I used the given data files as well created my own test files. 

#### Calculating Average and Median:
For some of the custom test files, I used smaller data sets to make it easier to verify the correctness of the calculated average and median values. I used datasets with both an even and odd amount of records so that I could test the two different cases of finding a median value. Additionally, for the datasets with an even amount of records I tested the cases where the median age is between the same age values and where the median age is between different age values. Also, since I used a frequency map for ages to calculate the median rather than using a list of the individual ages recorded in the dataset, I had to test the different points at which the index could how much the median index could be overlapped when incrementing by frequencies. For the datasets with an odd amount of records, I had to verify the correctness of the median in two cases - when the frequency increments the index we are iterating to be equal to the median index and when the frequency increments the index to be more than the median index. For the datasets with an even amount of records, I had to verify the correctness of the median in three cases since 

I compared the median values calculated from the frequency map with the median values found by using a list with all the individual age values listed to verify the correctness.

#### Sequential Processing:
Once I was sure that each part was working how I wanted it to, I added the functionality to handle data from multiple datasets sequentially. I did end-to-end tests for the whole process for both success and error cases.


#### Concurrent Processing:
Once I was sure that my implementation was able to correctly handle multiple datasets sequentially, I added concurrency. I first tested only needing to process one url and then moved onto processing multiple urls that consisted of csvs that 

Finally, with concurrency implemented, I wanted to ensure that the number of goroutines I used could process a large amount of urls quickly and correctly so I tested out several different max number of goroutines that could run concurrently. The dataset I did this load test with was file3.csv which has 10000 lines of code and had it listed hundreds of times in the input url file which so that it had to be read 1600 times.

Tests:
- csv data with trailing and leading spaces

### Design Considerations
Aggregating in parts after each batch of goroutines finished or as goroutines finish
Using sync map to store aggregated results
Use list of ages along with age frequency map but wanted to save space
Number of goroutines to use
How to handle some errors like incorrect type for csv data
Dealing with incorrect types for csv fields

## Questions
1. What assumptions did you make in your design? Why?

   Some assumptions I made in my design are that the 

2. How would you change your program if it had to process many files where each file was over 10M records?


3. How would you change your program if it had to process data from more than 20K URLs?

   I don't think I would need to change 
