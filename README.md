# go training

This is a training course for the `Go` programming language designed to onboard backend engineers to the re-engagement
team.

### Introduction
This example project follows the [golang-standards project layout](https://github.com/golang-standards/project-layout). As such, `main` files are under the `/cmd` module, code meant for sharing is under `/pkg` and, private application code will be under `/internal`.

### Getting Started

1. [Setup Go](https://golang.org/doc/tutorial/getting-started#install).
2. Install [GoLand](https://www.jetbrains.com/go/) on your machine.
3. Email US Help Desk for a GoLand license key. If necessary, list me as your manager for approval in the email.
4. Fork this repository and clone it locally on your machine. Open the repo in GoLand.
5. Make sure you're able to run `go run cmd/pmgr/main.go` on your machine.
5. Go through and complete the [Go Tutorial](https://golang.org/doc/tutorial/getting-started) and [A Tour of Go](https://tour.golang.org/welcome/1).
6. While you're going through the tutorial and the tour, make sure you stop and answer the questions in `QUESTIONS.md`.

### Final Project

After completing the tutorial, your final project will be to implement a basic CLI password manager in `Go`.

Your application should be able to:

1. display a useful help message when the user supplies the `-help` flag. (Look into the `Flag` package). 
2. take in 4 subarguments: `add`, `update`, `get`, `delete` that will add, update, get, or delete an entry in your password manager.
    - `./program add foo bar` will add an entry for account `foo` with password `bar` if one doesn't exist
    - `./program get foo` will get `foo`'s password if it exists.
    - `./program update foo newbar` will update `foo`'s password to `newbar`.
    - `./program delete foo` will delete `foo` if it exists.
3. encrypt the entries using [bcrypt](https://linuxhint.com/golang-crypto-package/)
4. handle error conditions gracefully.
5. have tests that can be run (and pass!) with `go test`.
6. build and install using `go build ./cmd/pmgr` and `go install ./cmd/pmgr`.

Upon completion of the project, a code review session will be scheduled to review the quality of your final project,
review your answers to the `QUESTIONS.md`, and provide constructive feedback to use in your HelloFresh Go coding.

### Extra material
Several presentations from the internal HF lectures given by Diego in 2020:
- [Introduction to Go, W1](https://drive.google.com/file/d/1ZhAhTBWsLB_ElgrRJiS80o6OSJtX-5mM/view)
- [Go is a little different, W2](https://drive.google.com/file/d/1_BdQldvtZ3UAe_pFxgeQ-OylWm1uw-oa/view)
- [Multiple things at the time, W3](https://drive.google.com/file/d/10sjjGZ5kIKLEhT1HN-q4xUrsEOMqLyEC/view)
- [Don't forget the tests, W4](https://drive.google.com/file/d/1uhC3V9tJS2l4Nu3V39DauRStwY8WcmMb/view)

### Contributing

If you've found an issue, please notifty me and create a PR to fix it if you wish. I'll review it and merge it if it is ready.
