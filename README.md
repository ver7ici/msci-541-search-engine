# MSCI 541: Homework 5
### by Christopher McCarthy

## Instructions
1. Install [the Go programming language](https://golang.org/)
2. Clone this repository using one of the methods in the "Code" dropdown above.
3. Create an index from a source file: run  
    `./indexEngine/indexEngine.exe path/to/latimes.gz path/to/index`  
    from the root of the repository
4. Start search engine: run  
    `./retrieve/retrieve.exe path/to/index`
    from the root of the repository
5. Enter a query
6. Enter one of the following options:  
    `Q:     exit the program`  
	`N:     enter a new query`  
	`R:     show results list`  
    `1-10:  show raw document at the corresponding rank`
      
## References
* Turpin, A., Tsegay, Y., Hawking, D., & Williams, H. E. (2007). Fast Generation of Result Snippets in Web Search. SIGIR07: The 30th Annual International SIGIR Conference (pp. 127–134). Amsterdam: Association for Computing Machinery.
* [lazieburd and Nidhin David's answer for break labels](https://stackoverflow.com/a/54602693)
* [Jack's answer for using WaitGroups](https://stackoverflow.com/a/42218240)
* [Creating Unique Slices in Go by kylewbanks](https://kylewbanks.com/blog/creating-unique-slices-in-go)
* [icza's answer for parsing space delimited files](https://stackoverflow.com/a/59972879)
* [Go by Example: Switch](https://gobyexample.com/switch)
* [Go `fmt` package docs](https://pkg.go.dev/fmt)
* [Go `regexp` package docs](https://pkg.go.dev/regexp)
* [Go `strings` package docs](https://pkg.go.dev/strings)
* [Go `sort` package docs](https://pkg.go.dev/sort)
* [Go `strconv` package docs](https://pkg.go.dev/strconv)
* [Go `csv` package docs](https://pkg.go.dev/encoding/csv#Reader.ReadAll)
* [Go `math` package docs](https://pkg.go.dev/math)
* [Go module tutorial](https://go.dev/doc/tutorial/create-module)

