package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type userMap map[string]string

func main() {
	//fmt.Println(pmgr.GetQuote())
	//fmt.Println("Argumenets", os.Args)

	users := userMap{}
	fmt.Println("user map: ", users)

	users.handleUserInput()

	fmt.Println("user map: ", users)

	//myPwd := "shubham"
	//providedHash, _ := HashPassword(myPwd)
	//fmt.Println("Password :", myPwd)
	//fmt.Println("Hash :", providedHash)
	//
	//isMatch := CheckPasswordHash(myPwd, providedHash)
	//fmt.Println("Matched ?:", isMatch)

}

func (u userMap)handleUserInput(){
	switch os.Args[1] {
	case "add":
		fmt.Println("Adding")
		u.addUserEntry(os.Args[2], os.Args[3])
		fmt.Println("user map: ", u)
		u.updatePassword(os.Args[2], "4567")
		fmt.Println("user map: ", u)
		u.deleteUserEntry(os.Args[2])
	case "update":
		u.updatePassword("jesse@hellofresh.com", "4567")
	case "get":
		fmt.Println("getting")
	case "delete":
		fmt.Println("deleting.")
	default:
		fmt.Printf("You entered an invalid option")
	}
}

func (u userMap) addUserEntry(username string, password string){
	fmt.Println("Username: ", username)
	fmt.Println("pwd: ", password)
	hashedPassword, error := HashPassword(password)
	if error != nil {
		fmt.Println("An Error Occurred")
		os.Exit(1)
	}
	u[username] = hashedPassword
}

func (u userMap) updatePassword(username string, newPassword string){
	fmt.Println("Same Username: ", username)
	fmt.Println("updated pwd: ", newPassword)
	newHashedPassword, error := HashPassword(newPassword)
	if error != nil {
		fmt.Println("An Error Occurred")
		os.Exit(1)
	}
	u[username] = newHashedPassword
}

func (u userMap) deleteUserEntry(username string){
	fmt.Println("Deleting Username: ", username)
	delete(u, username)
}


func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
