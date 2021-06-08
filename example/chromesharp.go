package main

import (
	"fmt"

	hbd "github.com/moond4rk/hack-browser-data"
)

func main() {
	ChromePassword()
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
	outputter := hbd.NewOutPutter(hbd.OutputCSV)
	file, err := outputter.CreateFile(hbd.Password.Name(), true)
	if err != nil {
		panic(err)
	}
	err = outputter.Write(password, file)
	if err != nil {
		panic(err)
	}
	fmt.Println("password", password.(*hbd.GeckoPassword))
	bookmark, err := b.GetBrowsingData(hbd.Bookmark)
	if err != nil {
		panic(err)
	}
	fmt.Println("bookmark", bookmark.(*hbd.GeckoBookmark))
	history, err := b.GetBrowsingData(hbd.History)
	if err != nil {
		panic(err)
	}
	fmt.Printf("history\n %#v\n", history.(*hbd.GeckoHistory))
	cookie, err := b.GetBrowsingData(hbd.Cookie)
	if err != nil {
		panic(err)
	}
	fmt.Println("cookie", cookie.(*hbd.GeckoCookie))
	download, err := b.GetBrowsingData(hbd.Download)
	if err != nil {
		panic(err)
	}
	fmt.Println("download", download.(*hbd.GeckoDownload))
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
	outputter := hbd.NewOutPutter(hbd.OutputCSV)
	file, err := outputter.CreateFile(hbd.Password.Name()+".csv", true)
	if err != nil {
		panic(err)
	}
	err = outputter.Write(password, file)
	if err != nil {
		panic(err)
	}
	bookmark, err := b.GetBrowsingData(hbd.Bookmark)
	if err != nil {
		panic(err)
	}
	fils, err := outputter.CreateFile(hbd.Bookmark.Name()+".csv", true)
	if err != nil {
		panic(err)
	}
	err = outputter.Write(bookmark, fils)
	if err != nil {
		panic(err)
	}
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
