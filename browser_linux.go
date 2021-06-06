package hackbrowserdata

import (
	"crypto/sha1"

	"github.com/godbus/dbus/v5"
	keyring "github.com/ppacher/go-dbus-keyring"
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

var rootProfile = "/home/*/"

func (b webkit) ProfilePath() string {
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
		return unsupported
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
	// what is d-bus @https://dbus.freedesktop.org/
	var secretKey []byte
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	svc, err := keyring.GetSecretService(conn)
	if err != nil {
		return nil, err
	}
	session, err := svc.OpenSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	collections, err := svc.GetAllCollections()
	if err != nil {
		return nil, err
	}
	for _, col := range collections {
		items, err := col.GetAllItems()
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			label, err := item.GetLabel()
			if err != nil {
				continue
			}
			if label == b.Storage() {
				se, err := item.GetSecret(session.Path())
				if err != nil {
					return nil, err
				}
				secretKey = se.Value
			}
		}
	}
	if secretKey == nil {
		return nil, errDbusSecretIsEmpty
	}
	var salt = []byte("saltysalt")
	// @https://source.chromium.org/chromium/chromium/src/+/master:components/os_crypt/os_crypt_linux.cc
	key := pbkdf2.Key(secretKey, salt, 1, 16, sha1.New)
	return key, nil
}

func (b gecko) MasterSecretKey() ([]byte, error) {
	return nil, nil
}
