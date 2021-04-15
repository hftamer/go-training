package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

type userMap map[string]string

func main() {
	users := userMap{}
	fmt.Println("user map: ", users)
	getTextFromFile("userData.txt", users)
	fmt.Println("user map: ", users)

	handleUserInput(users)

}

func handleUserInput(u userMap) {
	switch os.Args[1] {
	case "add":
		u.addUserEntry(os.Args[2], os.Args[3])
		fmt.Println("user map: ", u)
	case "update":
		u.updatePassword("user1", "4567")
	case "get":
		fmt.Println("getting")
		getPasswordFromMap(os.Args[2], u)
	case "delete":
		fmt.Println("deleting.")
		u.deleteUserEntry(os.Args[2])
	default:
		fmt.Printf("You entered an invalid option")
	}
}

func (u userMap) addUserEntry(username string, password string) {
	fmt.Println("*******  Adding ")
	fmt.Println("Username: ", username)
	fmt.Println("pwd: ", password)
	hashedPassword, error := HashPassword(password)
	fmt.Println("hashed pwd: (leaving out for now) ", hashedPassword)
	if error != nil {
		fmt.Println("An Error Occurred")
		os.Exit(1)
	}
	u[username] = password

	for key, value := range u {
		fmt.Println("key", key, "value", value)
	}

	keyValuePairsAsByteSlice := createKeyValuePairsAsString(u)
	fmt.Println("Returned value: ", keyValuePairsAsByteSlice)
	saveTextToFile("userData.txt", keyValuePairsAsByteSlice)
}

func createKeyValuePairsAsString(m map[string]string) string {
	//A Buffer is a variable-sized buffer of bytes with Read and Write methods. The zero value for Buffer is an empty buffer ready to use.
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	fmt.Println("Value of B: ", b)
	fmt.Println("B as Bytes: ", b.Bytes())
	return b.String()
}

func saveTextToFile(filename string, data string) {
	file, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	keyValuePairs := fmt.Sprintf("%s \n", data)
	if _, err := file.WriteString(keyValuePairs); err != nil {
		fmt.Println(err)
	}
}

func getTextFromFile(filename string, u userMap) {
	byteSlice, _ := os.ReadFile(filename)

	newSlice := strings.Split((string(byteSlice)),"\n")

	for _, value := range newSlice{
		result := strings.Split(value,"=")
		u[result[0]]= result[1]
	}
}

func getPasswordFromMap(key string, u userMap) {
	fmt.Println(u[key])
}


func (u userMap) updatePassword(username string, newPassword string) {
	fmt.Println("*******  Updating ")
	fmt.Println("Same Username: ", username)
	fmt.Println("updated pwd: ", newPassword)
	newHashedPassword, error := HashPassword(newPassword)
	if error != nil {
		fmt.Println("An Error Occurred")
		os.Exit(1)
	}
	u[username] = newHashedPassword
}

func (u userMap) deleteUserEntry(username string) {
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
