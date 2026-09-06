package main

import (
	"errors"
	"os"

	"github.com/leonelquinteros/gotext"

	"github.com/zalando/go-keyring"

	"github.com/Jguer/yay/v13/pkg/text"
)

type ErrKeyring struct {
	inner   error
	service string
	key     string
}

func (e *ErrKeyring) Error() string {
	return gotext.Get("Error with keyring of service %s with key %s. err: %s", e.service, e.key, e.inner.Error())
}

func handleKeyring(logger *text.Logger,
) error {
	aurUsername := os.Getenv("AUR_USERNAME")
	aurPassword := os.Getenv("AUR_PASSWORD")

	// surface error when variables are not set
	if aurUsername == "" || aurPassword == "" {
		return errors.New(
			gotext.Get("Please set AUR_USERNAME and AUR_PASSWORD environment variables for voting"),
		)
	}

	err := keyring.Set("yay", "user", aurUsername)
	if err != nil {
		return &ErrKeyring{inner: err, service: "yay", key: "user"}
	}

	err = keyring.Set("yay", "password", aurPassword)
	if err != nil {
		return &ErrKeyring{inner: err, service: "yay", key: "password"}
	}

	// check if creds are readable
	secret, err := keyring.Get("yay", "user")
	if err != nil && secret == aurUsername {
		return &ErrKeyring{inner: err, service: "yay", key: "user"}
	}
	secret, err = keyring.Get("yay", "password")
	if err != nil && secret == aurPassword {
		return &ErrKeyring{inner: err, service: "yay", key: "password"}
	}

	logger.Println(gotext.Get("Saved credentials successfully to keyring"))
	return nil
}
