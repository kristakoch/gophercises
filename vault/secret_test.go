package main

import (
	"testing"
)

func TestNewFileVault(t *testing.T) {
	t.Run("invalid key", func(t *testing.T) {
		key, fp := "92725", "secrets.data"
		_, err := NewFileVault(key, fp)
		if err == nil {
			t.Errorf("should have received error creating new file vault with key '%s'", key)
		}
	})

	t.Run("invalid filepath", func(t *testing.T) {
		key, fp := "6D696E6E69657468656D6F6F63686572", ""
		_, err := NewFileVault(key, fp)
		if err == nil {
			t.Error("should have received error creating new file vault with empty filepath")
		}
	})

	t.Run("happy path", func(t *testing.T) {
		key, fp := "6D696E6E69657468656D6F6F63686572", "secrets.data"
		_, err := NewFileVault(key, fp)
		if err != nil {
			t.Errorf("failed to create new file vault with key '%s' and filepath '%s'", key, fp)
		}
	})
}

func TestGetSet(t *testing.T) {
	key, fp := "6D696E6E69657468656D6F6F63686572", "secrets.data"
	fv, err := NewFileVault(key, fp)
	if err != nil {
		t.Errorf("failed to create new file vault with key '%s' and filepath '%s'", key, fp)
	}

	secrets := []struct {
		key string
		val string
	}{
		{"evernote", "123456"},
		{"garage code", "8366"},
		{"bank pin", "9994"},
	}

	// Test store.
	for _, s := range secrets {
		if err := fv.Set(s.key, s.val); err != nil {
			t.Errorf("failed to set new secret with key '%s' and val '%s', err: %s", s.key, s.val, err)
		}
	}

	// Test lookup.
	searchKey := "garage code"
	want := "8366"

	got, err := fv.Get(searchKey)
	if err != nil {
		t.Errorf("failed to get secret with key %s, err: %s", searchKey, err)
	}

	if got != want {
		t.Errorf("got %s want %s when looking up value for key %s", got, want, searchKey)
	}
}

func TestListAll(t *testing.T) {
	key, fp := "6D696E6E69657468656D6F6F63686572", "secrets.data"
	fv, err := NewFileVault(key, fp)
	if err != nil {
		t.Errorf("failed to create new file vault with key '%s' and filepath '%s'", key, fp)
	}

	if err := fv.ListAll(); err != nil {
		t.Errorf("failed to list all in file vault, err: '%s'", err)
	}
}

func TestDelete(t *testing.T) {
	key, fp := "6D696E6E69657468656D6F6F63686572", "secrets.data"
	fv, err := NewFileVault(key, fp)
	if err != nil {
		t.Errorf("failed to create new file vault with key '%s' and filepath '%s'", key, fp)
	}

	deleteKey := "bank pin"
	if err := fv.Delete(deleteKey); err != nil {
		t.Errorf("failed while attempting to delete entry with key '%s', err: '%s'", deleteKey, err)
	}

	if _, err := fv.Get(deleteKey); err == nil {
		t.Errorf("should have been error getting deleted entry with key '%s'", deleteKey)
	}
}
