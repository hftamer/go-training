package pmgr

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Vault struct {
	AccountsByName map[string]string `json:"vault"`
}

func GetVaultPath() string {
	return "pmgr-vault.json"
}

func LoadVault(path string) (Vault, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		// ReadFile will only error if the file doesn't exist
		// In that case, return a new vault with no error
		return Vault{
			AccountsByName: make(map[string]string),
		}, nil
	}

	vault := Vault{
		AccountsByName: make(map[string]string),
	}
	if err = json.Unmarshal(file, &vault); err != nil {
		return Vault{}, err
	}

	return vault, nil
}

func (vault Vault) Save(path string) error {
	data, err := json.MarshalIndent(vault, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
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

func (vault Vault) GetAccount(name string) (string, error) {
	if acc, found := vault.AccountsByName[name]; found {
		return acc, nil
	}

	return "", errors.New("account " + name + " doesn't exist")
}