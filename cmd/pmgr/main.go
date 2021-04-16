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
	getTextFromFile("userData.txt", users)
	fmt.Println("user map from file: ", users)

	handleUserInput(users)
}

func handleUserInput(u userMap) {
	switch os.Args[1] {
	case "add":
		addUserEntry(os.Args[2], os.Args[3])
		fmt.Println("user map: ", u)
	case "update":
		u.updatePassword(os.Args[2], os.Args[3], "userData.txt")
	case "get":
		fmt.Println("getting")
		u.getPasswordFromMap(os.Args[2])
	case "delete":
		fmt.Println("deleting.")
		u.deleteUserEntry(os.Args[2], "userData.txt")
	default:
		fmt.Printf("You entered an invalid option")
	}
}

func addUserEntry(username string, password string) {
	u := userMap{}
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

	saveUserData("userData.txt", u)
}

func createKeyValuePairsAsString(m map[string]string) string {
	//A Buffer is a variable-sized buffer of bytes with Read and Write methods. The zero value for Buffer is an empty buffer ready to use.
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=%s\n", key, value)
	}
	fmt.Println("Value of B: ", b)
	fmt.Println("B as Bytes: ", b.Bytes())
	return b.String()
}

func saveUserData(filename string, userInfo userMap) {
	keyValuePairsAsByteSlice := createKeyValuePairsAsString(userInfo)
	fmt.Println("Returned value: ", keyValuePairsAsByteSlice)

	file, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	keyValuePairs := fmt.Sprintf("\n%s", keyValuePairsAsByteSlice)
	if _, err := file.WriteString(keyValuePairs); err != nil {
		fmt.Println(err)
	}
}

func getTextFromFile(filename string, u userMap) {
	byteSlice, _ := os.ReadFile(filename)

	newSlice := strings.Split((string(byteSlice)), "\n")
	fmt.Println("new slice: ", newSlice)
	fmt.Println("length: ", len(newSlice))

	for _, value := range newSlice {
		if value != "" {
			result := strings.Split(value, "=")
			fmt.Println("Result: ", result)
			u[result[0]] = result[1]
		}
	}
}

func (u userMap) getPasswordFromMap(key string) {
	fmt.Println(u[key])
}

func (u userMap) updatePassword(username string, newPassword string, fileName string) {
	fmt.Println("*******  Updating ")
	fmt.Println("Same Username: ", username)
	fmt.Println("updated pwd: ", newPassword)
	newHashedPassword, error := HashPassword(newPassword)
	fmt.Println("New hashed password", newHashedPassword)

	fmt.Println("user map coming in: ", u)

	file, e := os.OpenFile(fileName, os.O_RDWR, 0755)
	fmt.Println("error? ", e)

	if error != nil {
		fmt.Println("An Error Occurred")
		os.Exit(1)
	}


	u[username] = newPassword

	fmt.Println("user map going out: ", u)

	err := file.Truncate(0)
	fmt.Println("error: ", err)
	saveUserData("userData.txt", u)
}

func (u userMap) deleteUserEntry(username string, fileName string) {
	fmt.Println("user map coming in: ", u)
	file, e := os.OpenFile(fileName, os.O_RDWR, 0755)
	fmt.Println("error? ", e)

	fmt.Println("Deleting Username: ", username)
	delete(u, username)
	fmt.Println("user map going out: ", u)

	err := file.Truncate(0)
	fmt.Println("error: ", err)
	saveUserData("userData.txt", u)

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
