package hackbrowserdata

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"os"
	"os/exec"

	"golang.org/x/crypto/pbkdf2"
)

// Storage use for implement BrowserClient interface
func (b webkit) Storage() string {
	return ""
}

// Storage use for implement BrowserClient interface
func (b gecko) Storage() string {
	return ""
}

var rootProfile = os.Getenv("USERPROFILE")

func (b webkit) ProfilePath() string {
	switch b {
	case Chrome:
		return rootProfile + "/AppData/Local/Google/Chrome/User Data/*/"
	case Chromium:
		return rootProfile + "/AppData/Local/Chromium/User Data/*/"
	case ChromeBeta:
		return rootProfile + "/AppData/Local/Google/Chrome Beta/User Data/*/"
	case Edge:
		return rootProfile + "/AppData/Local/Microsoft/Edge/User Data/*/"
	case Speed360:
		return rootProfile + "/AppData/Local/360chrome/Chrome/User Data/*/"
	case QQ:
		return rootProfile + "/AppData/Local/Tencent/QQBrowser/User Data/*/"
	case Brave:
		return rootProfile + "/AppData/Local/BraveSoftware/Brave-Browser/User Data/*/"
	case Opera:
		return rootProfile + "/AppData/Roaming/Opera Software/Opera Stable/"
	case OperaGX:
		return rootProfile + "/AppData/Roaming/Opera Software/Opera GX Stable/"
	case Vivaldi:
		return rootProfile + "/AppData/Local/Vivaldi/User Data/Default/"
	default:
		return unsupported
	}
}

func (b gecko) ProfilePath() string {
	switch b {
	case Firefox:
		return rootProfile + "/AppData/Roaming/Mozilla/Firefox/Profiles/*.default-release*/"
	case FirefoxBeta:
		return rootProfile + "/AppData/Roaming/Mozilla/Firefox/Profiles/*.default-beta*/"
	case FirefoxDev:
		return rootProfile + "/AppData/Roaming/Mozilla/Firefox/Profiles/*.dev-edition-default*/"
	case FirefoxNightly:
		return rootProfile + "/AppData/Roaming/Mozilla/Firefox/Profiles/*.default-nightly*/"
	case FirefoxESR:
		return rootProfile + "/AppData/Roaming/Mozilla/Firefox/Profiles/*.default-esr*/"
	default:
		return unsupported
	}
}

func (b webkit) KeyFilePath() string {
	switch b {
	case Chrome:
		return rootProfile + "/AppData/Local/Google/Chrome/User Data/Local State"
	case Chromium:
		return rootProfile + "/AppData/Local/Chromium/User Data/Local State"
	case ChromeBeta:
		return rootProfile + "/AppData/Local/Google/Chrome Beta/User Data/Local State"
	case Edge:
		return rootProfile + "/AppData/Local/Microsoft/Edge/User Data/Local State"
	case Brave:
		return rootProfile + "/AppData/Local/BraveSoftware/Brave-Browser/User Data/Local State"
	case Opera:
		return rootProfile + "/AppData/Roaming/Opera Software/Opera Stable/Local State"
	case OperaGX:
		return rootProfile + "/AppData/Roaming/Opera Software/Opera GX Stable/Local State"
	case Vivaldi:
		return rootProfile + "/AppData/Local/Vivaldi/Local State"
	default:
		return "Unknown Browser"
	}
}

func (b gecko) KeyFilePath() string {
	return ""
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
