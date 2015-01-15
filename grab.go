package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func GetInfo(voc *Voc) {
	url := "http://dict.youdao.com/search?q=" + voc.Spell + "&keyfrom=dict.index"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}
	getPronounce(voc, doc)
	getTranslation(voc, doc)
	getForms(voc, doc)
	getSamples(voc, doc)
	fmt.Println(voc)
}

func getPronounce(voc *Voc, doc *goquery.Document) {
	doc.Find(".pronounce").Each(func(i int, s *goquery.Selection) {
		voc.Pronunciation += strings.TrimSpace(s.Nodes[0].FirstChild.Data) + ":"
		voc.Pronunciation += strings.TrimSpace(s.Find(".phonetic").Text()) + ";"
	})
	voc.Pronunciation = strings.TrimSuffix(voc.Pronunciation, ";")
}

func getTranslation(voc *Voc, doc *goquery.Document) {
	liList := doc.Find("ul").Nodes[1]
	for i := liList.FirstChild.NextSibling; i != nil; i = i.NextSibling {
		if i.FirstChild == nil || len(strings.TrimSpace(i.Data)) == 0 || len(strings.TrimSpace(i.FirstChild.Data)) == 0 {
			continue
		}
		voc.Translation += strings.TrimSpace(i.FirstChild.Data) + "|"
	}

	voc.Translation = strings.TrimSuffix(voc.Translation, "|")
}

func getForms(voc *Voc, doc *goquery.Document) {
	form := ""
	nodes := doc.Find("p.additional").Nodes
	if nodes != nil && len(nodes) > 0 {
		dom := nodes[0]
		if dom.FirstChild != nil {
			form = dom.FirstChild.Data
			if form != "span" {
				re, _ := regexp.Compile("\\s{1,}")
				form = re.ReplaceAllString(form, " ")
				form = strings.Trim(form, "[")
				form = strings.Trim(form, "]")
				voc.Form += form
			}
		}
	}
}

func getSamples(voc *Voc, doc *goquery.Document) {
	sample := ""
	doc.Find("#bilingual ul li p").Each(func(i int, s *goquery.Selection) {
		if !s.HasClass("example-via") {
			if len(strings.TrimSpace(s.Text())) > 0 {
				sample += s.Text()
			}
		}
	})
	re, _ := regexp.Compile("\\s{2,}")
	sample = re.ReplaceAllString(sample, "|")
	sample = strings.TrimSuffix(sample, "|")
	voc.Sample = sample
}
