package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"testing"
)


func clearTestConfig(t *testing.T) {
	wd,_ := os.Getwd()
	testDir := path.Join(wd, "testdata", ".config")
	err := os.RemoveAll(testDir)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Clear dir: %q", testDir)
	}
}

func TestCreateConfig(t *testing.T) {
	defer clearTestConfig(t)
	wd,_ := os.Getwd()
	testDir := path.Join(wd, "testdata", ".config", "sticky")
	_, err := OpenConfig(testDir)
	if !errors.As(err, &CredentialsNotFound{}){
		t.Fatal(err)
	}
	info, err := os.Stat(testDir)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Fatalf("Expected %q to be a directory", testDir)
	}
}

func TestConfigReadCredentials(t *testing.T) {
	wd,_ := os.Getwd()
	testDir := path.Join(wd, "testdata", "config")
	c, err := OpenConfig(testDir)
	if err != nil {
		t.Fatal(err)
	}
	expectedCredentials, err := ioutil.ReadFile(
		path.Join(testDir, "credentials.json"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if string(c.GoogleCredentials()) != string(expectedCredentials) {
		t.Fatalf(
			"Credentials do not match\n got: %10q\n wanted: %q",
			string(c.GoogleCredentials()),
			string(expectedCredentials),
		)
	}
}
