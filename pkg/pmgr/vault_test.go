package pmgr

import (
	"path"
	"testing"
)

func getTestVault(t *testing.T) Vault {
	vaultPath := path.Join("./testdata", "test-vault.json")

	vault, err := LoadVault(vaultPath)
	if err != nil {
		t.Fatal("error loading test vault: ", err)
	}

	return vault
}

func TestAddAccountWithNoCollision(t *testing.T) {
	vault := getTestVault(t)

	err := vault.AddAccount("foo", "bar")
	if err != nil {
		t.Error("adding to the vault failed: ", err)
	}
}

func TestAddAccountWithCollision(t *testing.T) {
	vault := getTestVault(t)

	err := vault.AddAccount("foo", "bar")
	if err != nil {
		t.Error("adding to the vault failed: ", err)
	}

	err = vault.AddAccount("foo", "baz")
	if err == nil {
		t.Error("adding to a vault with a conflict should raise an error")
	}
}

func TestGetAccountThatExists(t *testing.T) {
	vault := getTestVault(t)

	err := vault.AddAccount("foo", "bar")
	if err != nil {
		t.Error("adding to the vault failed: ", err)
	}

	pwd, err := vault.GetAccountPassword("foo")
	if err != nil {
		t.Error("getting a existing password should not fail: ", err)
	}

	if pwd != "bar" {
		t.Error("password did not match")
	}
}

func TestGetAccountThatDoesNotExist(t *testing.T) {
	vault := getTestVault(t)

	_, err := vault.GetAccountPassword("foo")
	if err == nil {
		t.Error("getting a non-existent password should fail")
	}
}

func TestUpdateAccountThatExists(t *testing.T) {
	vault := getTestVault(t)

	err := vault.AddAccount("foo", "bar")
	if err != nil {
		t.Error("adding to the vault failed: ", err)
	}

	err = vault.UpdateAccount("foo", "baz")
	if err != nil {
		t.Error("updating an existing account failed: ", err)
	}

	pwd, err := vault.GetAccountPassword("foo")
	if err != nil {
		t.Error("getting a existing password should not fail: ", err)
	}

	if pwd != "baz" {
		t.Error("password did not match")
	}
}

func TestUpdateAccountThatDoesNotExist(t *testing.T) {
	vault := getTestVault(t)

	err := vault.UpdateAccount("foo", "baz")
	if err == nil {
		t.Error("updating non-existent account password should fail")
	}
}

func TestDeleteAccountThatExists(t *testing.T) {
	vault := getTestVault(t)

	err := vault.AddAccount("foo", "bar")
	if err != nil {
		t.Error("adding to the vault failed: ", err)
	}

	err = vault.DeleteAccount("foo")
	if err != nil {
		t.Error("deleting an existing account should not fail with: ", err)
	}

	if vault.IsExistingAccount("foo") {
		t.Error("delete did not delete the account")
	}
}

func TestDeleteAccountThatDoesNotExist(t *testing.T) {
	vault := getTestVault(t)

	if vault.IsExistingAccount("foo") {
		t.Fatal("test vault should be empty")
	}

	err := vault.DeleteAccount("foo")
	if err != nil {
		t.Error("deleting non-existent account should not be an error")
	}

}