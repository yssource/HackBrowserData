package main

import (
	"fmt"

	hbd "github.com/moond4rk/hack-browser-data"
)

func main() {
	FirefoxPassword()
}

func FirefoxPassword() {
	b, err := hbd.NewBrowser(hbd.Firefox)
	if err != nil {
		panic(err)
	}
	password, err := b.GetBrowsingData(hbd.Password)
	if err != nil {
		panic(err)
	}
	fmt.Println(password.(*hbd.GeckoPassword))
	bookmark, err := b.GetBrowsingData(hbd.Bookmark)
	if err != nil {
		panic(err)
	}
	fmt.Println(bookmark.(*hbd.GeckoBookmark))
	history, err := b.GetBrowsingData(hbd.History)
	if err != nil {
		panic(err)
	}
	fmt.Println(history.(*hbd.GeckoHistory))
	cookie, err := b.GetBrowsingData(hbd.Cookie)
	if err != nil {
		panic(err)
	}
	fmt.Println(cookie.(*hbd.GeckoCookie))
	download, err := b.GetBrowsingData(hbd.Download)
	if err != nil {
		panic(err)
	}
	fmt.Println(download.(*hbd.GeckoDownload))
}

func ChromePassword() {
	b, err := hbd.NewBrowser(hbd.Chrome)
	if err != nil {
		panic(err)
	}
	password, err := b.GetBrowsingData(hbd.Password)
	if err != nil {
		panic(err)
	}
	fmt.Println(password.(*hbd.WebkitPassword))
	bookmark, err := b.GetBrowsingData(hbd.Bookmark)
	if err != nil {
		panic(err)
	}
	fmt.Println(bookmark.(*hbd.WebkitBookmark))
	history, err := b.GetBrowsingData(hbd.History)
	if err != nil {
		panic(err)
	}
	fmt.Println(history.(*hbd.WebkitHistory))

	creditCard, err := b.GetBrowsingData(hbd.CreditCard)
	if err != nil {
		panic(err)
	}
	fmt.Println(creditCard.(*hbd.WebkitCreditCard))

	download, err := b.GetBrowsingData(hbd.Download)
	if err != nil {
		panic(err)
	}
	fmt.Println(download.(*hbd.WebkitDownload))

	cookie, err := b.GetBrowsingData(hbd.Cookie)
	if err != nil {
		panic(err)
	}
	var _ = cookie.(*hbd.WebkitCookie)
	fmt.Println(cookie.(*hbd.WebkitCookie))
}
