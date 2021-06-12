package hackbrowserdata

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/moond4rk/hack-browser-data/internal/utils"
)

type Browser struct {
	clientList  []Client
	itemerList  []Itemer
	name        string
	storage     string
	profilePath string
	keyPath     string
	outputType  outputType
	outputDir   string
	filename    string
}

type Option func(*Browser) error

func NewBrowser(options ...Option) (*Browser, error) {
	b := &Browser{}
	for _, option := range options {
		err := option(b)
		if err != nil {
			return nil, err
		}
	}
	if len(b.clientList) <= 0 {
		return nil, errors.New("clientlist must be set")
	}
	if len(b.itemerList) <= 0 {
		return nil, errors.New("itemerlist must be set")
	}
	return b, nil
}

func (b *Browser) Run() error {
	outputter := NewOutPutter(b.outputType)
	for _, client := range b.clientList {
		for _, itemer := range b.itemerList {
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

// WithClient gets a single browser
func WithClient(cname string) Option {
	return func(b *Browser) error {
		names := NameBrowserMap()
		if client, ok := names[cname]; ok {
			b.clientList = append(b.clientList, client)
		} else{
			return fmt.Errorf("not support this browser %s", client.Name())
		}
		return nil
	}
}

// WithAllClients is a option to add all browser clients
func WithAllClients() Option {
	return func(b *Browser) error {
		b.clientList = getAllClients()
		return nil
	}
}

// ListAllBrowserName get all browsers supported by this OS
func ListAllBrowserName() (names []string) {
	list := getAllClients()
	for _, v := range list {
		names = append(names, strings.ToLower(v.Name()))
	}
	return names
}

func NameBrowserMap() map[string]Client {
	var names = make(map[string]Client)
	clients := getAllClients()
	for i, name := range ListAllBrowserName() {
		names[name] = clients[i]
	}
	return names
}

// getAllClients using the profile path is exist to get all browser are supported by this OS
func getAllClients() []Client {
	var clientList []Client
	for i := 0; i <= int(Vivaldi); i++ {
		if webkit(i).ProfilePath() != unsupported {
			clientList = append(clientList, webkit(i))
		}
	}
	for i := 0; i <= int(FirefoxESR); i++ {
		if gecko(i).ProfilePath() != unsupported {
			clientList = append(clientList, gecko(i))
		}
	}
	return clientList
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
		return "ChromeBeta"
	case Edge:
		return "Edge"
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

type Client interface {
	Name() string

	Storage() string

	ProfilePath() string

	KeyFilePath() string

	MasterSecretKey() ([]byte, error)

	GetBrowsingData(item Itemer) (BrowsingData, error)
}

func (b gecko) Name() string {
	switch b {
	case Firefox:
		return "Firefox"
	case FirefoxBeta:
		return "Firefox-Beta"
	case FirefoxDev:
		return "Firefox-Dev"
	case FirefoxNightly:
		return "Firefox-Nightly"
	case FirefoxESR:
		return "Firefox-ESR"
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
	err = utils.CopyFileToLocal(b.ProfilePath(), itemer.FileName(Chrome))
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
		if err := utils.CopyFileToLocal(b.ProfilePath(), paths[0]); err != nil {
			fmt.Println(err)
		}
		if err := utils.CopyFileToLocal(b.ProfilePath(), paths[1]); err != nil {
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
		if err := utils.CopyFileToLocal(b.ProfilePath(), itemer.FileName(b)); err != nil {
			fmt.Println(err)
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
	errDbusSecretIsEmpty    = errors.New("dbus secret key is empty")
)
