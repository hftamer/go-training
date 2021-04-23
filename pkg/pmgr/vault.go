package pmgr

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

// NOTE: this is a __horrible__ security practice
// Please do not emulate this in production!
var cryptoPassword = "mySuperL33tpassword"

type Vault struct {
	AccountsByName map[string]string `json:"vault"`
	path           string
}

func GetVaultPath() string {
	return "pmgr-vault.json"
}

func LoadVault(path string) (Vault, error) {
	vault := Vault {
		AccountsByName: make(map[string]string),
		path: path,
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		// ReadFile will only error if the file doesn't exist
		// The file won't exist if this is the first time we are creating the vault
		// In that case, return a new vault with no error
		return vault, nil
	}

	if err = json.Unmarshal(file, &vault); err != nil {
		return Vault{}, err
	}

	return vault, nil
}

func (vault Vault) Save() error {
	data, err := json.MarshalIndent(vault, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(vault.path, data, 0644)
}

func (vault Vault) IsExistingAccount(name string) bool {
	_, found := vault.AccountsByName[name]

	return found
}

func (vault *Vault) AddAccount(name string, pwd string) error {
	if vault.IsExistingAccount(name) {
		return errors.New(name + " account already exists")
	}

	return vault.addPassword(name, pwd)
}

func (vault Vault) GetAccountPassword(name string) (string, error) {
	encryptedPwd, found := vault.AccountsByName[name]
	if !found {
		return "", errors.New("account " + name + " doesn't exist")
	}

	decoded, err := hex.DecodeString(encryptedPwd)
	if err != nil {
		return "", err
	}

	return decrypt(decoded, cryptoPassword)
}

func (vault *Vault) UpdateAccount(name string, pwd string) error {
	if !vault.IsExistingAccount(name) {
		return errors.New("account " + name + " not found")
	}

	return vault.addPassword(name, pwd)
}

func (vault *Vault) DeleteAccount(name string) error {
	delete(vault.AccountsByName, name)

	return nil
}

func (vault *Vault) addPassword(name string, pwd string) error {
	encryptedPwd, err := encrypt([]byte(pwd), cryptoPassword)
	if err != nil {
		return err
	}

	vault.AccountsByName[name] = encryptedPwd

	return nil
}

func encrypt(data []byte, passphrase string) (string, error) {
	hash, err := hash(passphrase)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(hash)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return hex.EncodeToString(ciphertext), nil
}

func decrypt(data []byte, passphrase string) (string, error) {
	hash, err := hash(passphrase)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(hash)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("nonce size exceeds data length")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func hash(key string) ([]byte, error) {
	hasher := md5.New()
	if _, err := hasher.Write([]byte(key)); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}