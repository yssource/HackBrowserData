package hackbrowserdata

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Browser struct {
	ClientList  []BrowserClient
	ItemerList  []Itemer
	name        string
	storage     string
	profilePath string
	keyPath     string
	outputType  outputType
	filename    string
}

type BrowserClient interface {
	Name() string

	Storage() string

	ProfilePath() string

	KeyFilePath() string

	MasterSecretKey() ([]byte, error)

	GetBrowsingData(item Itemer) (BrowsingData, error)
}

type Option func(*Browser) error

func NewBrowser(options ...Option) (*Browser, error) {
	browser := &Browser{}
	for _, option := range options {
		err := option(browser)
		if err != nil {
			return nil, err
		}
	}
	return browser, nil
}

func (b *Browser) Run() error {
	outputter := NewOutPutter(b.outputType)
	for _, client := range b.ClientList {
		for _, itemer := range b.ItemerList {
			data, err := client.GetBrowsingData(itemer)
			// Handle error
			if err != nil {
				return err
			}
			filename := client.Name() + "_" + itemer.Name() + b.outputType.String()
			f, err := outputter.CreateFile(filename, true)
			err = outputter.Write(data, f)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewBrowserList() []BrowserClient {
	var browserList []BrowserClient
	for i := 0; i <= int(Vivaldi); i++ {
		if webkit(i).ProfilePath() != unsupported {
			browserList = append(browserList, webkit(i))
		}
	}
	for i := 0; i <= int(FirefoxESR); i++ {
		if gecko(i).ProfilePath() != unsupported {
			browserList = append(browserList, gecko(i))
		}
	}
	return browserList
}

type (
	// webkit is a browser engine developed by Apple and primarily used in Safari and Chrome browser
	webkit int
	// gecko is a browser engine developed by Mozilla. It is used in the Firefox browser
	gecko int
)

const (
	Chrome webkit = iota + 1
	ChromeBeta
	Chromium
	Edge
	Speed360
	QQ
	Brave
	Opera
	OperaGX
	Vivaldi
)

const (
	Firefox gecko = iota + 1
	FirefoxBeta
	FirefoxDev
	FirefoxNightly
	FirefoxESR
)

const unsupported = "unsupported browser"

func (b webkit) Name() string {
	switch b {
	case Chrome:
		return "Chrome"
	case Chromium:
		return "Chromium"
	case ChromeBeta:
		return "Chrome Beta"
	case Edge:
		return "Microsoft Edge"
	case Speed360:
		return "360speed"
	case QQ:
		return "qq"
	case Brave:
		return "Brave"
	case Opera:
		return "Opera"
	case OperaGX:
		return "OperaGX"
	case Vivaldi:
		return "Vivaldi"
	default:
		return unsupported
	}
}

func (b gecko) Name() string {
	switch b {
	case Firefox:
		return "Firefox"
	case FirefoxBeta:
		return "Firefox Beta"
	case FirefoxDev:
		return "Firefox Dev"
	case FirefoxNightly:
		return "Firefox Nightly"
	case FirefoxESR:
		return "Firefox ESR"
	default:
		return unsupported
	}
}

// GetBrowsingData 拼接 itemer 的文件名称，
func (b webkit) GetBrowsingData(itemer Itemer) (BrowsingData, error) {
	var (
		masterKey []byte
		err       error
	)
	// TODO: store MasterSecretKey, not call function each time
	if itemer == Password || itemer == Cookie || itemer == CreditCard {
		masterKey, err = b.MasterSecretKey()
		if err != nil {
			return nil, err
		}
	}
	err = copyFileToLocal(b.ProfilePath(), itemer.FileName(Chrome))
	if err != nil {
		return nil, err
	}
	data := itemer.Data(b)
	err = data.parse(itemer, masterKey)
	if err != nil {
		return nil, err
	}
	err = os.Remove(itemer.FileName(b))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (b gecko) GetBrowsingData(itemer Itemer) (BrowsingData, error) {
	var (
		masterKey []byte
		err       error
	)
	// firefox password need key4.db
	if itemer == Password {
		paths := strings.Split(itemer.FileName(b), "|")
		if err := copyFileToLocal(b.ProfilePath(), paths[0]); err != nil {
			fmt.Println(err)
		}
		if err := copyFileToLocal(b.ProfilePath(), paths[1]); err != nil {
			fmt.Println(err)
		}
		data := itemer.Data(b)
		err = data.parse(itemer, masterKey)
		err = os.Remove(paths[0])
		if err != nil {
			return nil, err
		}
		err = os.Remove(paths[1])
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		if err := copyFileToLocal(b.ProfilePath(), itemer.FileName(b)); err != nil {

		}
		data := itemer.Data(b)
		err = data.parse(itemer, masterKey)
		if err != nil {
			return nil, err
		}
		err = os.Remove(itemer.FileName(b))
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

var (
	ErrWrongSecurityCommand = errors.New("wrong security command")
)
var (
	errItemNotSupported    = errors.New(`item not supported, default is "all", choose from history|download|password|bookmark|cookie`)
	errBrowserNotSupported = errors.New("browser not supported")
	errChromeSecretIsEmpty = errors.New("chrome secret is empty")
	errDbusSecretIsEmpty   = errors.New("dbus secret key is empty")
)

// copyToLocal copy file to local path
func copyFileToLocal(profilePath, filename string) error {
	p, err := filepath.Glob(filepath.Join(profilePath, filename))
	if err != nil {
		return err
	}
	// TODO, handle error if file not exist
	if len(p) <= 0 {
		return fmt.Errorf("find %s failed", filename)
	} else {
		src := p[0]
		locals, _ := filepath.Glob("*")
		for _, v := range locals {
			if v == filename {
				err := os.Remove(filename)
				if err != nil {
					return err
				}
			}
		}
		sourceFile, err := ioutil.ReadFile(src)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filename, sourceFile, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}
