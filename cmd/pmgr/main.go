package main

import (
	"fmt"
	"os"
)

func main() {
	//fmt.Println(pmgr.GetQuote())
	fmt.Println("Argumenets", os.Args)

	handleUserInput()

}

func handleUserInput(){
	switch os := os.Args[1]; os {
	case "add":
		fmt.Println("Adding")
	case "update":
		fmt.Println("updating.")
	case "get":
		fmt.Println("getting")
	case "delete":
		fmt.Println("deleting.")
	default:
		fmt.Printf("You entered an invalid option")
	}
}

func addPassword(){

}
