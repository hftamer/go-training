package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hftamer/go-training/pkg/encrypt"
)

func New() (Vault, error) {
	encodingKey := os.Getenv("PMGR_ENCODING_KEY")

	if encodingKey == "" {
		return Vault{}, errors.New("no encoding key environment variable set")
	}

	v := Vault{
		encodingKey: encodingKey,
		data:        make(map[string]string),
	}

	err := v.loadData()
	if err != nil {
		return v, fmt.Errorf("error loading file data: %v", err)
	}

	return v, nil
}

type Vault struct {
	encodingKey string
	data        map[string]string
}

func (v *Vault) loadData() error {
	// try to open the secrets file
	f, err := os.OpenFile("secrets", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	// copy file contents
	var sb strings.Builder
	_, err = io.Copy(&sb, f)
	if err != nil {
		return err
	}

	// decrypt file content
	decryptedData, err := encrypt.Decrypt(v.encodingKey, sb.String())
	if err != nil {
		return fmt.Errorf("file decryption failed with key: %q, on loading, %s", v.encodingKey, err)
	}

	// decode file contents into json
	r := strings.NewReader(decryptedData)
	decryptedJSON := json.NewDecoder(r)
	decryptedJSON.Decode(&v.data)
	// how to deal with the initial empty file?
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (v *Vault) SaveData() {
	var sb strings.Builder
	enc := json.NewEncoder(&sb)
	err := enc.Encode(v.data)
	if err != nil {
		fmt.Println(fmt.Errorf("can't encode file data into json: %s", err))
	}

	encryptedJSON, err := encrypt.Encrypt(v.encodingKey, sb.String())
	if err != nil {
		fmt.Println(fmt.Errorf("can't encrypt new data: %s", err))
	}

	err = ioutil.WriteFile("secrets", []byte(encryptedJSON), 0755)
	if err != nil {
		fmt.Println(fmt.Errorf("can't write to file: %s", err))
	}
}

func (v *Vault) Get(key string) (string, error) {
	val, ok := v.data[key]
	if !ok {
		return "", errors.New("no value for that key")
	}

	return val, nil
}

func (v *Vault) Add(key, value string) error {
	_, ok := v.data[key]
	if ok {
		return errors.New("this key already has a value")
	}
	v.data[key] = value
	return nil
}

func (v *Vault) Update(key, value string) error {
	_, ok := v.data[key]
	if !ok {
		return errors.New("this key doesn't exist")
	}
	v.data[key] = value
	return nil
}

func (v *Vault) Delete(key string) error {
	_, ok := v.data[key]
	if !ok {
		return errors.New("this key doesn't exist")
	}
	delete(v.data, key)
	return nil
}

func (v *Vault) PrintAll() {
	fmt.Println(v.data)
}