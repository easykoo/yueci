package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func main() {
	GetInfo(&Voc{Spell: "Carry"})
}

func GetInfo(voc *Voc) {
	url := "http://dict.youdao.com/search?q=" + voc.Spell + "&keyfrom=dict.index"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}
	getPronounce(voc, doc)
	getTranslation(voc, doc)
	fmt.Println(voc)
}

func getPronounce(voc *Voc, doc *goquery.Document) {
	doc.Find(".pronounce").Each(func(i int, s *goquery.Selection) {
		voc.Pronunciation += strings.TrimSpace(s.Nodes[0].FirstChild.Data) + ":"
		voc.Pronunciation += strings.TrimSpace(s.Find(".phonetic").Text()) + ";"
	})
}

func getTranslation(voc *Voc, doc *goquery.Document) {
	//form:= doc.Find("p.additional").Nodes
	liList := doc.Find("ul").Nodes[1]
	for i := liList.FirstChild.NextSibling; i != nil; i = i.NextSibling {
		fmt.Println(i.Data)
		//fmt.Println(i.FirstChild.Data)
	}
	//voc.Translation += ul
}
