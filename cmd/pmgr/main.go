package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type userMap map[string][]byte
type Account struct {
	Username string `json:"name"`
	Password []byte `json:"password"`
}

type Vault struct {
	Accounts []Account `json:"account"`
}

func main() {
	filename := "test.json"
	hashedPassphrase := createHash("p@S$w0rd")
	runCommandLineProgram(filename, hashedPassphrase)
}

func runCommandLineProgram(fileName string, hashedPassphrase string) {
	addHelperFlagText()

	switch os.Args[1] {
	case "add":
		validateCommandLineArguments(4)
		addUserEntryToFile(os.Args[2], os.Args[3], fileName, hashedPassphrase)
	case "update":
		validateCommandLineArguments(4)
		updatePassword(os.Args[2], os.Args[3], fileName, hashedPassphrase)
	case "get":
		validateCommandLineArguments(3)
		getPassword(os.Args[2], hashedPassphrase)
	case "delete":
		deleteUserEntry(os.Args[2], fileName)
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

func addUserEntryToFile(username string, password string, filename string, hashedPassphrase string) {
	// create new vault with existing data
	vaultWithExistingData := Vault{}
	dataFile, _ := ioutil.ReadFile(filename)
	_ = json.Unmarshal([]byte(dataFile), &vaultWithExistingData)

	for _, v := range vaultWithExistingData.Accounts {
		if v.Username == username {
			fmt.Println("Oops! Looks like this account already exists")
			os.Exit(1)
		}
	}
	encryptedPasswordAsByteSlice := encrypt([]byte(password), hashedPassphrase)

	// create new account entry
	newAccount := Account{Username: username, Password: encryptedPasswordAsByteSlice}
	vaultWithExistingData.Accounts = append(vaultWithExistingData.Accounts, newAccount)

	// update json file with new data
	updateJsonFile(vaultWithExistingData, filename)
}

func updatePassword(username string, newPassword string, filename string, hashedPassphrase string) {
	vaultWithExistingData := Vault{}
	dataFile, _ := ioutil.ReadFile(filename)
	_ = json.Unmarshal([]byte(dataFile), &vaultWithExistingData)

	newVault := Vault{}
	found := false
	encryptedPasswordAsByteSlice := encrypt([]byte(newPassword), hashedPassphrase)
	for _, v := range vaultWithExistingData.Accounts {

		if v.Username == username {
			found = true
			newVault.Accounts = append(newVault.Accounts, Account{Username: v.Username, Password: encryptedPasswordAsByteSlice})
		} else {
			newVault.Accounts = append(newVault.Accounts, Account{Username: v.Username, Password: v.Password})
		}
	}

	printErrorMessageIfNecessary(found)

	// update json file with new data
	updateJsonFile(newVault, filename)
	fmt.Println("successfully updated")
}

func getPassword(username string, hashedPassphrase string) {
	vaultWithExistingData := Vault{}
	dataFile, _ := ioutil.ReadFile("test.json")
	_ = json.Unmarshal([]byte(dataFile), &vaultWithExistingData)

	found := false
	for _, v := range vaultWithExistingData.Accounts {
		if v.Username == username {
			fmt.Println("password: ", string(decrypt(v.Password, hashedPassphrase)))
			found = true
		}
	}

	printErrorMessageIfNecessary(found)
	fmt.Println("successfully retrieved password")
}

func deleteUserEntry(username string, filename string) string {
	vaultWithExistingData := Vault{}
	dataFile, _ := ioutil.ReadFile(filename)
	_ = json.Unmarshal([]byte(dataFile), &vaultWithExistingData)

	newVault := Vault{}
	found := false
	for _, v := range vaultWithExistingData.Accounts {
		if v.Username == username {
			found = true
		} else {
			newVault.Accounts = append(newVault.Accounts, Account{Username: v.Username, Password: v.Password})
		}
	}

	printErrorMessageIfNecessary(found)
	updateJsonFile(newVault, filename)
	successMessage := "successfully deleted"
	fmt.Println(successMessage)
	return successMessage
}

// **** Helper functions
func printErrorMessageIfNecessary(found bool){
	if !found {
		fmt.Println("Oops! Looks like that username doesn't exist")
		os.Exit(1)
	}
}
func updateJsonFile(newVault Vault, filename string){
	// update json file with new data
	out, error := json.MarshalIndent(newVault, "", " ")

	if error != nil {
		fmt.Println(error)
	}

	error = ioutil.WriteFile(filename, out, 0644)
	if error != nil {
		fmt.Println(error)
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