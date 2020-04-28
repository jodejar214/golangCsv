# CrowdStrike ECD team homework assignment

## Task
The task is to compute some summary statistics from a collection of CSV files accessed via URLs.

Your program should read data consisting of fname, lname, and age in CSV format from a set of given URLs. How the URLs are delivered to the program is up to you - command-line arguments or from a flat file are the simplest choices. The program should read from all the URLs provided, with a reasonable level of concurrency. Assume that there will be hundreds of files to read and network latency will be a consideration, so reading them sequentially will be too slow. The response from the URL may or may not be well-formed or even available. You could get back a well-formed file, you could get back garbage, or a 404 error, among other failure modes, so be sure that your implementation does some sanity checking.

After reading all the data, your program should print the average age, the median age, and the name of a person of the median age, as well as the clock time spent reading all the data files and the amount of concurrency utilized (e.g. if you spun up 4 threads and the total clock time was 90 seconds, report that).

Treat the input files as comprising one single logical data set that's been partitioned into multiple files, not as independent and distinct inputs. For example, if each input file contained names and ages for the residents of each US state, then you'd have 50 input files and your output would give the average and median age for the whole country. The summary stats should be reported over all the records in all the input files, not one line per input file. 

It’s likely that there will be multiple records that have the median age, in which case any of the names of people of that age are acceptable as output. On the other hand, if there are an even number of records, then the median will be the average of the middle two and might not represent an actual record. This unlikely to happen but possible, describe how your code handles this situation in the readme. It should also separately list any of the URLs that failed to be read for whatever reason.

Sample CSV input files are provided. They can be read from local disk using file:// URLs, but assume that for actual use the files will be coming over the network.

## Language and tools
You should use a common programming language that you're quite familiar with. Please show your best style and cleanest code. The cloud teams at CrowdStrike internally make heavy use of Go, Python, and Scala; implementations that use one of these languages will be easiest to evaluate. Java and C# are also very reasonable choices.

You should use a minimum of external dependencies. Facilities provided by the standard libraries of your language and deployment platform of choice are always acceptable. However, please don't ask us to install random code from github. Code that *requires* an IDE to build and run is not acceptable, although you're of course welcome to use your tools of choice to implement the homework.

## External References
Please treat this assignment as you would a work task. You may use any external references that you care to: documentation, books, stack overflow, advice from a peer, whatever. However, you must implement the assignment yourself, and you represent that the work product you submit is wholly your own output.

## Deliverables
You should provide the source code for your application, along with any other files needed to build, test, and run the program. The program is not expected to be ready for use at scale in production, but it should be clear, easy to read, easy to run, and be code that you'd be comfortable putting up for code review with a peer engineer.
You should provide a brief design document as a readme file that discusses how to build and run your program, why your program looks the way it does, and how you tested it. Please indicate if there are design directions that you considered and rejected, and why you rejected them. In addition, it should also answer the following questions:
1. What assumptions did you make in your design? Why?
2. How would you change your program if it had to process many files where each file was over 10M records?
3. How would you change your program if it had to process data from more than 20K URLs?
