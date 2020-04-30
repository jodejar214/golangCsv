# Design Document

You should provide the source code for your application, along with any other files needed to build, test, and run the program. The program is not expected to be ready for use at scale in production, but it should be clear, easy to read, easy to run, and be code that you'd be comfortable putting up for code review with a peer engineer.
You should provide a brief design document as a readme file that discusses how to build and run your program, why your program looks the way it does, and how you tested it. Please indicate if there are design directions that you considered and rejected, and why you rejected them.

## Implementation Design

### How To Build/Run
> go run csvData.go {file with list of urls}

***Example:***
go run csvData.go urls.txt

### Functionality



### Testing
For testing my implementation, I first ran tests to check functionality of each part of the process - reading the input file with urls to retrieve, retrieving the data from the urls, organizing the data, calculating the average and median - on one csv dataset depending on the case I was testing. I tested several successful cases on different data sets for each part as well as error cases for each part.

#### Reading File With List of URLs
For testing the reading of the input file with the list of urls to process, I had uploaded the data files onto Github so that I could test retrieving the raw csv files over the network. I then listed a subset of the urls in the input file to test retrieving individual data files as well as multiple data files. Additionally, I listed urls for csv files that did not exist on Github as well as random urls that did not link to csv files such as "www.google.com". I also listed strings that were not real urls such as random words and an empty string. These urls or strings were added to test the error handling of the http request and reading the response data.

#### Processing CSV data files
For test data, I used the given data files as well created my own data files. For some of the custom data files, I used smaller data sets to make it easier to verify the correctness of the calculated average and median values. 

#### Calculating Average and Median
I used datasets with both an even and odd amount of records so that I could test the two different cases of finding a median value. Additionally, for the datasets with an even amount of records I tested the cases where the median age is between the same age values and where the median age is between different age values. Also, since I used a frequency map for ages to calculate the median rather than using a list of the individual ages recorded in the dataset, I had to test the different points at which the index could how much the median index could be overlapped when incrementing by frequencies. For the datasets with an odd amount of records, I had to verify the correctness of the median in two cases - when the frequency increments the index we are iterating to be equal to the median index and when the frequency increments the index to be more than the median index. For the datasets with an even amount of records, I had to verify the correctness of the median in three cases since 

I compared the median values calculated from the frequency map with the median values found by using a list with all the individual age values listed to verify the correctness.

#### Sequential Processing
Once I was sure that each part was working how I wanted it to, I added the functionality to handle data from multiple datasets sequentially. I did end-to-end tests for the whole process for both success and error cases.


#### Concurrent Processing
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

## Questions
1. What assumptions did you make in your design? Why?

Some assumptions I made in my design are that the 

2. How would you change your program if it had to process many files where each file was over 10M records?


3. How would you change your program if it had to process data from more than 20K URLs?

I don't think I would need to change 
