package hackbrowserdata

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/moond4rk/hack-browser-data/internal/utils"
)

type Browser struct {
	clients    []Client
	items      []Itemer
	outputType outputType
	outputDir  string
	compress   bool
	fileAppend bool
	silence    bool
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
	if len(b.clients) <= 0 {
		return nil, errors.New("clients must be set")
	}
	if len(b.items) <= 0 {
		return nil, errors.New("items must be set")
	}
	return b, nil
}

func (b *Browser) Run() error {
	output := NewOutPutter(b.outputType)

	err := createDir(b.outputDir)
	if err != nil {
		return err
	}
	for _, client := range b.clients {
		for _, itemer := range b.items {
			fmt.Println(client.Name(), itemer.Name())
			data, err := client.GetBrowsingData(itemer)
			// Handle error while browser not installed
			if err != nil {
				if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "not support") {
					fmt.Println(err)
					continue
				} else {
					return err
				}
			}
			filename := strings.ToLower(b.outputDir + "/" + client.Name() + "_" + itemer.Name() + b.outputType.String())
			f, err := output.CreateFile(filename)
			err = output.Write(data, f)
			if err != nil {
				return err
			}
		}
	}
	if b.compress {
		if err = utils.Compress(b.outputDir); err != nil {
			return err
		}
	}
	return nil
}

// WithClient gets a single browser
func WithClient(cname string) Option {
	return func(b *Browser) error {
		names := BrowserMap()
		if client, ok := names[cname]; ok {
			b.clients = append(b.clients, client)
		} else {
			return fmt.Errorf("not support this browser client: %s", client.Name())
		}
		return nil
	}
}

func WithOutputJson() Option {
	return func(b *Browser) error {
		b.outputType = OutputJson
		return nil
	}
}

func WithOutputCSV() Option {
	return func(b *Browser) error {
		b.outputType = OutputCSV
		return nil
	}
}

func WithOutputDir(dir string) Option {
	return func(b *Browser) error {
		b.outputDir = dir
		return nil
	}
}

func WithFileAppend(append bool) Option {
	return func(b *Browser) error {
		b.fileAppend = append
		return nil
	}

}

// WithAllClients is a option to add all browser clients
func WithAllClients() Option {
	return func(b *Browser) error {
		b.clients = getAllClients()
		return nil
	}
}

func WithItemer(itemname string) Option {
	return func(b *Browser) error {
		names := ItemerMap()
		if itemer, ok := names[itemname]; ok {
			b.items = append(b.items, itemer)
		} else {
			return fmt.Errorf("not support this itemer: %s", itemer.Name())
		}
		return nil
	}
}

func WithAllItemers() Option {
	return func(b *Browser) error {
		b.items = getAllItemer()
		return nil
	}
}

func WithCompress(compress bool) Option {
	return func(b *Browser) error {
		b.compress = compress
		return nil
	}
}

// ListAllBrowserName get all browsers supported by this OS
func ListAllBrowserName() (names []string) {
	clients := getAllClients()
	for _, v := range clients {
		names = append(names, strings.ToLower(v.Name()))
	}
	return names
}

func ListAllItemerName() (names []string) {
	itmes := getAllItemer()
	for _, v := range itmes {
		names = append(names, strings.ToLower(v.Name()))
	}
	return names
}

func BrowserMap() map[string]Client {
	var names = make(map[string]Client)
	clients := getAllClients()
	for i, name := range ListAllBrowserName() {
		names[name] = clients[i]
	}
	return names
}

func ItemerMap() map[string]Itemer {
	var names = make(map[string]Itemer)
	items := getAllItemer()
	for i, name := range ListAllItemerName() {
		names[name] = items[i]
	}
	return names
}

// getAllClients using the profile path is exist to get all browser are supported by this OS
func getAllClients() []Client {
	var clients []Client
	for i := 0; i <= int(Vivaldi); i++ {
		if webkit(i).ProfilePath() != unsupported {
			clients = append(clients, webkit(i))
		}
	}
	for i := 0; i <= int(FirefoxESR); i++ {
		if gecko(i).ProfilePath() != unsupported {
			clients = append(clients, gecko(i))
		}
	}
	return clients
}

func createDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.Mkdir(dir, 0700)
	}
	return nil
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
		return nil, errors.Wrapf(err, "browser: %s maybe not exist", b.Name())
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
	if itemer.FileName(b) == unsupportedItem {
		return nil, fmt.Errorf("browser %s not support export %s", b.Name(), itemer.Name())
	}
	if itemer == Password {
		paths := strings.Split(itemer.FileName(b), "|")
		if err := utils.CopyFileToLocal(b.ProfilePath(), paths[0]); err != nil {
			return nil, errors.Wrapf(err, "browser: %s maybe not exist", b.Name())
		}
		if err := utils.CopyFileToLocal(b.ProfilePath(), paths[1]); err != nil {
			return nil, errors.Wrapf(err, "browser: %s maybe not exist", b.Name())
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
			return nil, errors.Wrapf(err, "browser: %s maybe not exist", b.Name())
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
	ErrWrongSecurityCommand = errors.New("macOS wrong security command")
	errDbusSecretIsEmpty    = errors.New("dbus secret key is empty")
)
