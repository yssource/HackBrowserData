package main

import (
	"fmt"

	hbd "github.com/moond4rk/hack-browser-data"
)

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
