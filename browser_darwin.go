package hackbrowserdata

import (
	"bytes"
	"crypto/sha1"
	"os/exec"

	"github.com/pkg/errors"

	"golang.org/x/crypto/pbkdf2"
)

// Storage is key stored in macOS keychain access
func (b webkit) Storage() string {
	switch b {
	case Chrome, ChromeBeta:
		return "Chrome"
	case Chromium:
		return "Chromium"
	case Edge:
		return "Microsoft Edge"
	case Brave:
		return "Brave"
	case Opera, OperaGX:
		return "Opera"
	case Vivaldi:
		return "Vivaldi"
	default:
		return unsupported
	}
}

func (b gecko) Storage() string {
	return unsupported
}

const rootProfile = "/Users/*/Library/Application Support/"

// ProfilePath is the path where the browser stores data
func (b webkit) ProfilePath() string {
	switch b {
	case Chrome:
		return rootProfile + "Google/Chrome/*/"
	case Chromium:
		return rootProfile + "Chromium/*/"
	case ChromeBeta:
		return rootProfile + "Google/Chrome Beta/*/"
	case Edge:
		return rootProfile + "Microsoft Edge/*/"
	case Brave:
		return rootProfile + "BraveSoftware/Brave-Browser/*/"
	case Opera:
		return rootProfile + "com.operasoftware.Opera/"
	case OperaGX:
		return rootProfile + "com.operasoftware.OperaGX/"
	case Vivaldi:
		return rootProfile + "Vivaldi/*/"
	default:
		return unsupported
	}
}

func (b gecko) ProfilePath() string {
	switch b {
	case Firefox:
		return rootProfile + "Firefox/Profiles/*.default*/"
	case FirefoxBeta:
		return rootProfile + "Firefox/Profiles/*.default-beta*/"
	case FirefoxDev:
		return rootProfile + "Firefox/Profiles/*.dev-edition-default*/"
	case FirefoxNightly:
		return rootProfile + "Firefox/Profiles/*.default-nightly*/"
	case FirefoxESR:
		return rootProfile + "Firefox/Profiles/*.default-esr*/"
	default:
		return unsupported
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
	// $ security find-generic-password -wa 'Chrome'
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
		return nil, ErrWrongSecurityCommand
	}
}

func (b gecko) MasterSecretKey() ([]byte, error) {
	return nil, nil
}
