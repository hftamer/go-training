Please find time to answer these questions while going through the [Go
tutorial](https://golang.org/doc/tutorial/getting-started).

1. What is the purpose of `go.mod` and `go.sum`? Can you think of any analogous files in other programming environments?
   Go.mod is a file that defines a modules "module path". This path is used for the root directory and dependency requirments. 
   Every dependency requirement is written as a module path with a specific semantic version. 
   Go.sum  is a file containing the expected Cryptographic checksums of the content of specific modules. 
   Every time a dependency is used, it's checksum is added to go.sum. 
   This seems like the eqivalent of `yarn` or `npm` in JavaScript and having a `yarn.lock` file. 
   
2. When running `go mod init` on a module you wish to publish, what should the module path be?
   `go mod init` creates a new `go.mod` file in the current directory. It accepts an optional arguement (the module path).
   The module path should be the path that identifies the module and acts as a prefix for package import paths within the module.
   
3. What is one reason why you'd need to run `go mod edit` command?
   `go mod edit`provides a command-line interface for editing and formatting go.mod files.
   It takes several flags, and can be used for example to: 
   - set the go version
   - add a requirement
   - format the `go.mod` file
   - print the `go.mod` file in JSON format 
    
4. What is the "blank identifier" in Go used for?
   Blank Identifier is used to define an unused variable. Every variable defined in `go` must be used, and if it isn't an error is thrown.
   To get around this, you can use a blank identifier so you wouldn't have to use a variable unneccesariy. 
   Example:
   ```
   for _, card := range cards {
    fmt.println(card)
   }
   ```
   
    In this case, the first value would give the `index`, but if the index isn't actually needed, you can omit it using a blank identifier. 
   
5. How does go know which files are test files to execute with `go test`?
   `go test` will recompile each package along with any files with names matching the file pattern `*_test.go`. 
   Any file ending with the suffix `*_test.go` gets compiled a as a separate package and get linked with the main test binary.
