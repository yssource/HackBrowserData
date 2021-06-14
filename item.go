package hackbrowserdata

type item int

const (
	Password item = iota + 1
	Bookmark
	History
	Download
	Cookie
	CreditCard
)

func getAllItemer() (items []Itemer) {
	for i := 0; i <= int(CreditCard); i++ {
		if item(i).Name() != unsupportedItem {
			items = append(items, item(i))
		}
	}
	return items
}

type Itemer interface {
	Name() string
	FileName(browser Client) string
	Data(browser Client) BrowsingData
}

const unsupportedItem = "unsupported item"

func (i item) Name() string {
	switch i {
	case Password:
		return "Password"
	case Bookmark:
		return "Bookmark"
	case History:
		return "History"
	case Cookie:
		return "Cookie"
	case Download:
		return "Download"
	case CreditCard:
		return "CreditCard"
	default:
		return unsupportedItem
	}
}

// FileName return filename by browser type
func (i item) FileName(client Client) string {
	const (
		chromePasswordFile = "Login Data"
		chromeHistoryFile  = "History"
		chromeDownloadFile = "History"
		chromeCookieFile   = "Cookies"
		chromeBookmarkFile = "Bookmarks"
		chromeCreditFile   = "Web Data"
		firefoxCookieFile  = "cookies.sqlite"
		firefoxKey4File    = "key4.db"
		firefoxLoginFile   = "logins.json"
		firefoxDataFile    = "places.sqlite"
	)
	switch client.(type) {
	case webkit:
		switch i {
		case Password:
			return chromePasswordFile
		case Bookmark:
			return chromeBookmarkFile
		case Cookie:
			return chromeCookieFile
		case History:
			return chromeHistoryFile
		case Download:
			return chromeDownloadFile
		case CreditCard:
			return chromeCreditFile
		}
	case gecko:
		switch i {
		case Password:
			return firefoxLoginFile + "|" + firefoxKey4File
		case Bookmark:
			return firefoxDataFile
		case Cookie:
			return firefoxCookieFile
		case History:
			return firefoxDataFile
		case Download:
			return firefoxDataFile
		default:
			return unsupportedItem
		}
	}
	return unsupported
}

func (i item) Data(client Client) BrowsingData {
	switch client.(type) {
	case webkit:
		switch i {
		case Password:
			return &WebkitPassword{}
		case Cookie:
			return &WebkitCookie{}
		case Bookmark:
			return &WebkitBookmark{}
		case History:
			return &WebkitHistory{}
		case CreditCard:
			return &WebkitCreditCard{}
		case Download:
			return &WebkitDownload{}
		}
	case gecko:
		switch i {
		case Password:
			return &GeckoPassword{}
		case Cookie:
			return &GeckoCookie{}
		case Bookmark:
			return &GeckoBookmark{}
		case History:
			return &GeckoHistory{}
		case Download:
			return &GeckoDownload{}
		}
	}
	return nil
}
