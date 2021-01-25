package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var configPath string

func init() {
	configPath, _ = os.UserConfigDir()
	configPath += "/lyr/"
}

func ReadToken() (string, error) {
	data, err := ioutil.ReadFile(configPath + "key")
	if err != nil {
		return "", errors.New(fmt.Sprintf(
			"No config file found\nCreate a file called key at %s with your genius api token as its content\n",
			configPath))
	}
	return string(data), nil
}

func SetToken(token string) error {

	keyPath := configPath + "key"

	err := os.MkdirAll(configPath, 0755)
	if err != nil {
		return errors.New("Unable to create config directory at " + "\n" + configPath)
	}

	file, err := os.Create(keyPath)
	if err != nil {
		return errors.New("Unable to create config file at " + keyPath + "\n" + err.Error())
	}

	_, err = file.Write([]byte(token))
	if err != nil {
		return errors.New("Unable to write to config file at " + keyPath + "\n" + err.Error())
	}

	fmt.Println("Successfully wrote token to " + keyPath)
	return nil
}
