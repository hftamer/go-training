package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type userMap map[string][]byte
type Account struct {
	Username string `json:"name"`
	Password string `json:"password"`
}

// vault is a struct that has a slice of type account
type Vault struct {
	Accounts []Account `json:"account"`
}

func main() {
	userData := userMap{}
	mainVault := Vault{}

	newAccount := Account{Username: "hi", Password: "bye"}
	fmt.Println("new account", newAccount)
	fmt.Println("vault: ", mainVault)


	mainVault.Accounts = append(mainVault.Accounts, newAccount)
	fmt.Println("vault: ", mainVault)
	fmt.Printf("vault: %+v", mainVault)
	fmt.Printf("*** vault: %p\n", mainVault.Accounts)


	mainVault.addAccount()
	mainVault.addAccount()

	fmt.Println("vault: ", mainVault)
	fmt.Printf("AGAINNNNNvault: %+v", mainVault)
	filename := "userData.txt"
	hashedPassphrase := createHash("p@S$w0rd")
	populateUserMapWithDataFromFile(filename, userData)
	runCommandLineProgram(userData, filename, hashedPassphrase, &mainVault, mainVault)

	fmt.Println("vault: ", mainVault)
	fmt.Printf("last one: %+v", mainVault)

	out, error := json.MarshalIndent(mainVault, "", " ")
	fmt.Println("output:", string(out))

	fmt.Println("errors:", error)

	//err := ioutil.WriteFile("test.json", file, 0644)
	//fmt.Println("errors:", err)

}

func (pointerToVault *Vault) addAccount() {
	newAccount2 := Account{Username: "ciao", Password: "arrivederci"}
	pointerToVault.Accounts = append(pointerToVault.Accounts, newAccount2)
	fmt.Println("in change: ", newAccount2)
	fmt.Printf("in change: %+v", newAccount2)
	fmt.Printf("*** vault: %p\n", newAccount2)
}


func runCommandLineProgram(userData userMap, fileName string, hashedPassphrase string, pointerToVault *Vault, mainVault Vault) {
	addHelperFlagText()

	switch os.Args[1] {
	case "add":
		validateCommandLineArguments(4)
		userData.addUserEntryToFile(os.Args[2], os.Args[3], fileName, hashedPassphrase, pointerToVault, mainVault)

	case "update":
		validateCommandLineArguments(4)
		userData.updatePassword(os.Args[2], os.Args[3], fileName, hashedPassphrase)
	case "get":
		validateCommandLineArguments(3)
		userData.getPasswordFromMap(os.Args[2], hashedPassphrase)
	case "delete":
		userData.deleteUserEntry(os.Args[2], fileName)
	default:
		fmt.Printf("You entered an invalid option")
	}
}

// *** main functions
func addHelperFlagText(){
	boolArgPtr := flag.Bool("help", false, "Give instructions on how to use the program")
	flag.Parse()

	if *boolArgPtr {
		fmt.Println("Welcome to password manager!")
		fmt.Println("To create a new username and password use: $ go run main.go add {username} {password}")
		fmt.Println("To retrieve your password use: $ go run main.go get {username}")
		fmt.Println("To update your password use: $ go run main.go update {username} {newPassword}")
		fmt.Println("To delete your username and password use: $ go run main.go delete {username}")
		os.Exit(1)
	}
}

func validateCommandLineArguments(expectedLength int) {
	if len(os.Args) != expectedLength {
		errorMessage := fmt.Sprintf("Wrong number of arguments: expected %v, got %v", expectedLength, len(os.Args))
		fmt.Println(errorMessage)
		os.Exit(1)
	}
}

func (userData userMap) addUserEntryToFile(username string, password string, filename string, hashedPassphrase string, pointerToMainVault *Vault, mainVault Vault) {

	// create new account entry
	newAccount := Account{Username: username, Password: password}
	fmt.Printf("new account: %+v", newAccount)

	// add new account entry to vault
	pointerToMainVault.Accounts = append(pointerToMainVault.Accounts, newAccount)
	pointerToMainVault.Accounts = append((*pointerToMainVault).Accounts, newAccount)
	fmt.Printf("*** main vault: %+v", pointerToMainVault)
	fmt.Println("*****")
	fmt.Printf("***** main vault: %+v", mainVault)






	if len(userData[username]) != 0 {
		log.Fatal("Oops! Looks like that username already exists")
		os.Exit(1)
	}

	newUserData := userMap{}


	fmt.Println("hashed passphrase: ", hashedPassphrase)

	passwordAsByteSlice := []byte(password)
	fmt.Println("password as byte slice: ", passwordAsByteSlice)

	encryptedByteSlice :=  encrypt(passwordAsByteSlice, hashedPassphrase)

	fmt.Println("encrypted Byte Slice: ", string(encryptedByteSlice))

	newUserData[username] = encryptedByteSlice
	saveUserData(filename, newUserData)
	fmt.Println("successfully added")
}

func (userData userMap) updatePassword(username string, newPassword string, fileName string, hashedPassphrase string) {
	userData.checkForExistingUser(username)

	fmt.Println("updated pwd: ", newPassword)


	file, _ := os.OpenFile(fileName, os.O_RDWR, 0755)

	passwordAsByteSlice := []byte(newPassword)
	fmt.Println("password as byte slice: ", passwordAsByteSlice)

	encryptedByteSlice :=  encrypt(passwordAsByteSlice, hashedPassphrase)
	userData[username] = encryptedByteSlice

	err := file.Truncate(0)
	handleFileTruncatingError(err)
	saveUserData(fileName, userData)
	fmt.Println("successfully updated")
}

func (userData userMap) getPasswordFromMap(username string, hashedPassphrase string) string {
	userData.checkForExistingUser(username)
	password := userData[username]
	decryptedPassword := string(decrypt(password, hashedPassphrase))

	fmt.Println(decryptedPassword)
	fmt.Println("successfully retreived password")
	return decryptedPassword
}

func (userData userMap) deleteUserEntry(username string, fileName string) string {
	userData.checkForExistingUser(username)
	clearTextFile(fileName)
	delete(userData, username)
	saveUserData(fileName, userData)
	successMessage := "successfully deleted"
	fmt.Println(successMessage)
	return successMessage
}

// **** Helper functions
func covertUserDataMapToString(userData userMap) string {
	userDataAsByteSlice := new(bytes.Buffer)
	for username, password := range userData {
		fmt.Fprintf(userDataAsByteSlice, "%s=%s,", username, password)
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
	userDataAsStringSlice := strings.Split((string(userDataAsByteSlice)), ",")
	fmt.Println("userData as slice: ", len(userDataAsStringSlice))

	if len(userDataAsStringSlice) != 0 {
		for _, value := range userDataAsStringSlice {
			if value != "" {
				result := strings.Split(value, "=")
				u[result[0]] = []byte(result[1])
			}
		}
	}
}

func (userData userMap) checkForExistingUser(username string) {
	value := userData[username]
	if len(value) == 0 {
		fmt.Println("Oops! Looks like that username doesn't exist")
		os.Exit(1)
	}
}

func clearTextFile(fileName string) {
	file, _ := os.OpenFile(fileName, os.O_RDWR, 0755)
	err := file.Truncate(0)
	handleFileTruncatingError(err)
}

func handleFileTruncatingError(e error) {
	if e != nil {
		fmt.Println("An Error Occurred while clearing the file")
		os.Exit(1)
	}
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}