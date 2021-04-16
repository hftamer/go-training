package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
)

type userMap map[string]string

func main() {
	userData := userMap{}
	filename := "userData.txt"
	populateUserMapWithDataFromFile(filename, userData)
	runCommandLineProgram(userData, filename)
}

func runCommandLineProgram(userData userMap, fileName string) {
	switch os.Args[1] {
	case "add":
		validateCommandLineArguments(4)
		userData.addUserEntryToFile(os.Args[2], os.Args[3], fileName)
	case "update":
		validateCommandLineArguments(4)
		userData.updatePassword(os.Args[2], os.Args[3], fileName)
	case "get":
		validateCommandLineArguments(3)
		userData.getPasswordFromMap(os.Args[2])
	case "delete":
		userData.deleteUserEntry(os.Args[2], fileName)
	default:
		fmt.Printf("You entered an invalid option")
	}
}

func validateCommandLineArguments(expectedLength int) {
	if len(os.Args) != expectedLength {
		errorMessage := fmt.Sprintf("Wrong number of arguments: expected %v, got %v", expectedLength, len(os.Args))
		fmt.Println(errorMessage)
		os.Exit(1)
	}
}

func (userData userMap) addUserEntryToFile(username string, password string, filename string) {
	if userData[username] != "" {
		log.Fatal("Oops! Looks like that username already exists")
		os.Exit(1)
	}

	newUserData := userMap{}
	hashedPassword,_ := HashPassword(password)
	fmt.Println("hashed pwd: (leaving out for now) ", hashedPassword)
	newUserData[username] = password
	saveUserData(filename, newUserData)
	fmt.Println("successfully added")
}

func covertUserDataMapToString(userData userMap) string {
	userDataAsByteSlice := new(bytes.Buffer)
	for username, password := range userData {
		fmt.Fprintf(userDataAsByteSlice, "%s=%s\n", username, password)
	}
	return userDataAsByteSlice.String()
}

func saveUserData(filename string, userInfo userMap) {
	userDataAsOneString := covertUserDataMapToString(userInfo)

	file, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	newUserData := fmt.Sprintf("\n%s", userDataAsOneString)
	if _, err := file.WriteString(newUserData); err != nil {
		fmt.Println(err)
	}
}

func populateUserMapWithDataFromFile(filename string, u userMap) {
	userDataAsByteSlice, _ := os.ReadFile(filename)
	userDataAsStringSlice := strings.Split((string(userDataAsByteSlice)), "\n")

	for _, value := range userDataAsStringSlice {
		if value != "" {
			result := strings.Split(value, "=")
			u[result[0]] = result[1]
		}
	}
}

func (userData userMap) getPasswordFromMap(username string) string {
	userData.checkForExistingUser(username)
	fmt.Println(userData[username])
	fmt.Println("successfully retried password")
	return userData[username]
}

func (userData userMap) checkForExistingUser(username string) {
	value := userData[username]
	if value == "" {
		fmt.Println("Oops! Looks like that username doesn't exist")
		os.Exit(1)
	}
}

func (userData userMap) updatePassword(username string, newPassword string, fileName string) {
	userData.checkForExistingUser(username)

	fmt.Println("updated pwd: ", newPassword)
	newHashedPassword,_ := HashPassword(newPassword)
	fmt.Println("New hashed password", newHashedPassword)

	file,_ := os.OpenFile(fileName, os.O_RDWR, 0755)
	userData[username] = newPassword

	err := file.Truncate(0)
	handleFileTruncatingError(err)
	saveUserData(fileName, userData)
	fmt.Println("successfully updated")
}

func (userData userMap) deleteUserEntry(username string, fileName string) string {
	userData.checkForExistingUser(username)
	clearTextFile(fileName)
	delete(userData, username)
	saveUserData(fileName, userData)
	successMessage := "successfully deleted"
	return successMessage
}
func clearTextFile(fileName string){
	file,_ := os.OpenFile(fileName, os.O_RDWR, 0755)
	err := file.Truncate(0)
	handleFileTruncatingError(err)
}

func handleFileTruncatingError(e error){
	if e != nil {
		fmt.Println("An Error Occurred while clearing the file")
		os.Exit(1)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
