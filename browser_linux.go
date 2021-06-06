package hackbrowserdata

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"os/exec"

	"golang.org/x/crypto/pbkdf2"
)

// Storage is key stored Linux D-BUS Session
func (b webkit) Storage() string {
	switch b {
	case Chrome, ChromeBeta, Vivaldi:
		return "Chrome Safe Storage"
	case Chromium, Edge, Opera:
		return "Chromium Safe Storage"
	case Brave:
		return "Brave Safe Storage"
	default:
		return unsupported
	}
}

func (b gecko) Storage() string {
	return unsupported
}

func (b webkit) ProfilePath() string {
	const rootProfile = "/home/*/"
	switch b {
	case Chrome:
		return rootProfile + ".config/google-chrome/*/"
	case Chromium:
		return rootProfile + ".config/chromium/*/"
	case ChromeBeta:
		return rootProfile + ".config/google-chrome-beta/*/"
	case Edge:
		return rootProfile + ".config/microsoft-edge*/*/"
	case Brave:
		return rootProfile + ".config/BraveSoftware/Brave-Browser/*/"
	case Opera:
		return rootProfile + ".config/opera/"
	case Vivaldi:
		return rootProfile + ".config/vivaldi/*/"

	default:
		return "Unknown Browser"
	}
}

func (b gecko) ProfilePath() string {
	switch b {
	case Firefox:
		return rootProfile + ".mozilla/firefox/*.default*/"
	case FirefoxBeta:
		return rootProfile + ".mozilla/firefox/*.default-beta*/"
	case FirefoxDev:
		return rootProfile + ".mozilla/firefox/*.dev-edition-default*/"
	case FirefoxNightly:
		return rootProfile + ".mozilla/firefox/*.default-nightly*/"
	case FirefoxESR:
		return rootProfile + ".mozilla/firefox/*.default-esr*/"
	default:
		return "Unknown Browser"
	}
}

func (b webkit) KeyFilePath() string {
	return unsupported
}

func (b gecko) KeyFilePath() string {
	return unsupported
}

func (b webkit) MasterSecretKey() ([]byte, error) {
	var (
		cmd            *exec.Cmd
		stdout, stderr bytes.Buffer
	)

	// ➜ security find-generic-password -wa 'Chrome'
	if b.Storage() != unsupported {
		cmd = exec.Command("security", "find-generic-password", "-wa", b.Storage())
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
		if stderr.Len() > 0 {
			err = errors.New(stderr.String())
			return nil, err
		}
		temp := stdout.Bytes()
		chromeSecret := temp[:len(temp)-1]
		if chromeSecret == nil {
			return nil, ErrWrongSecurityCommand
		}
		var chromeSalt = []byte("saltysalt")
		// @https://source.chromium.org/chromium/chromium/src/+/master:components/os_crypt/os_crypt_mac.mm;l=157
		key := pbkdf2.Key(chromeSecret, chromeSalt, 1003, 16, sha1.New)
		return key, nil
	} else {
		return nil, errors.New(unsupported)
	}
}

func (b gecko) MasterSecretKey() ([]byte, error) {
	return nil, nil
}
