package pmgr

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

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

	// TODO: encrypt the password first
	vault.AccountsByName[name] = pwd

	return nil
}

func (vault Vault) GetAccountPassword(name string) (string, error) {
	if acc, found := vault.AccountsByName[name]; found {
		return acc, nil
	}

	return "", errors.New("account " + name + " doesn't exist")
}

func (vault *Vault) UpdateAccount(name string, password string) error {
	if _, found := vault.AccountsByName[name]; found {
		vault.AccountsByName[name] = password
		return nil
	} else {
		return errors.New("account " + name + " not found")
	}
}

func (vault Vault) DeleteAccount(name string) error {
	delete(vault.AccountsByName, name)

	return nil
}