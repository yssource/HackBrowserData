package hackbrowserdata

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/tidwall/gjson"

	"github.com/moond4rk/hack-browser-data/internal/decrypt"
	"github.com/moond4rk/hack-browser-data/internal/utils"
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
		return rootProfile + "/AppData/Roaming/Mozilla/Firefox/Profiles/*.default*/"
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
	if b.KeyFilePath() == unsupported {
		return nil, fmt.Errorf("%s %s", b.Name(), unsupported)
	}
	if _, err := os.Stat(b.KeyFilePath()); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s secret key path is empty", b.Name())
	}
	keyFile, err := utils.ReadFile(b.KeyFilePath())
	if err != nil {
		return nil, err
	}
	encryptedKey := gjson.Get(keyFile, "os_crypt.encrypted_key")
	if encryptedKey.Exists() {
		pureKey, err := base64.StdEncoding.DecodeString(encryptedKey.String())
		if err != nil {
			return nil, err
		}
		masterKey, err := decrypt.DPAPI(pureKey[5:])
		if err != nil {
			return nil, err
		}
		return masterKey, nil
	}
	return nil, nil
}

func (b gecko) MasterSecretKey() ([]byte, error) {
	return nil, nil
}
