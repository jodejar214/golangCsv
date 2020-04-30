# Design Document

## Implementation Design

### Functionality

### Testing
For testing my implementation, I first ran tests to check functionality of each part of the process - reading the input file with urls to retrieve, retrieving the data from the urls, organizing the data, calculating the average and median - on one csv dataset depending on the case I was testing. I tested several successful cases on different data sets for each part as well as error cases for each part.

For testing the reading of the input file with the list of urls to process, I had uploaded the data files onto Github so that I could test retrieving the raw csv files over the network. I then listed a subset of the urls in the input file to test retrieving individual data files as well as multiple data files. Additionally, I listed urls for csv files that did not exist on Github as well as random urls that did not link to csv files such as "www.google.com". I also listed strings that were not real urls such as random words and an empty string. These were added to test the error handling of the http retrieval and reading the response data as csvs.

For test data, I used the given data files as well created my own data files. For some of the custom data files, I used smaller data sets to make it easier to verify the correctness of the calculated average and median values. I used datasets with both an even and odd amount of records so that I could test the two different cases of finding a median value. Additionally, for the even Also, in my final implementation I chose to use a map of age frequencies rather than listing out all the ages found through the multiple csv files to save space even though finding the median from there would be simpler but I had implemented both and compared the results for both

Once I was sure that each part was working how I wanted it to, I added the functionality to handle data from multiple datasets sequentially. I did end-to-end tests for the whole process for both success and error cases.

Once I was sure that my implementation was able to correctly handle multiple datasets sequentially, I added concurrency. I first tested only needing to process one url and then moved onto processing multiple urls that consisted of csvs that 

Finally, with concurrency implemented, I wanted to ensure that the number of goroutines I used could process a large amount of urls quickly and correctly so I tested out several different max number of goroutines that could run concurrently. The dataset I did this load test with was file3.csv which has 10000 lines of code and had it listed hundreds of times in the input url file which so that it had to be read 1600 times.

Tests:
- csv data with trailing and leading spaces
- odd number of data points
- even number of data points with median between same age value
- even number of data points with median between diff age value

### How To Build/Run

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


You should provide the source code for your application, along with any other files needed to build, test, and run the program. The program is not expected to be ready for use at scale in production, but it should be clear, easy to read, easy to run, and be code that you'd be comfortable putting up for code review with a peer engineer.
You should provide a brief design document as a readme file that discusses how to build and run your program, why your program looks the way it does, and how you tested it. Please indicate if there are design directions that you considered and rejected, and why you rejected them.
