package main

import (
	"encoding/csv"
	"fmt"
	"github.com/opesun/goquery"
	"github.com/tealeg/xlsx"
	"io"
	"os"
)

var Vocs []*Voc

func main1() {
	GetInfoFromChaDanCi(&Voc{Spell: "abuse"})
	Read()
	Write()
}

func GetInfoFromChaDanCi(voc *Voc) {
	var url = "http://www.chadanci.com/s/" + voc.Spell
	p, err := goquery.ParseUrl(url)
	if err != nil {
		panic(err)
	} else {
		voc.Pronunciation = p.Find(".trs").Text()
		//voc.WordForm = p.Find("#dictc_PWDECMEC dl div").Text()
		wordForm := ""
		for _, v := range p.Find("#dictc_PWDECMEC dl") {
			n := goquery.Nodes{v}
			for _, o := range n.Find("div") {
				wordForm += goquery.Nodes{o}.Text() + "|"
			}
		}
		voc.WordForm = wordForm
		sample := ""
		for _, v := range p.Find("#dictc_xglj dl div ol li.li_sent") {
			//for _, e := range v.Child {
			n := goquery.Nodes{v}
			sample += n.Find(".ee_title").Text() + "|"
			sample += n.Find("").Last().Text() + "||"
		}
		voc.Sample = sample
		fmt.Println(voc)
	}
}

//Read .csv
func Read() {
	file, err := os.Open("cet4.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comment = '#' //可以设置读入文件中的注释符
	reader.Comma = ','   //默认是逗号，也可以自己设置
	//还可以设置以下信息
	//FieldsPerRecord  int  // Number of expected fields per record
	//LazyQuotes       bool // Allow lazy quotes
	//TrailingComma    bool // Allow trailing comma
	//TrimLeadingSpace bool // Trim leading space
	//line             int
	//column           int
	fout, err := os.OpenFile("out.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer fout.Close()
	for {
		voc := new(Voc)
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		voc.Spell = record[0]
		voc.Translation = record[1]
		GetInfoFromChaDanCi(voc)
		Vocs = append(Vocs, voc)
	}

}

func Write() {
	f, err := os.Create("cet4_result.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f)

	for _, voc := range Vocs {
		w.Write([]string{voc.Spell, voc.Pronunciation, voc.Translation, voc.WordForm, voc.Sample})
	}
	w.Flush()
	fmt.Println(len(Vocs))
	/*
		for _, voc := range Vocs {
			fmt.Printf("%s %s\n", voc.Spell, voc.Translation)
		}
	*/
}

//Read .xlsx
func Read1() {
	excelFileName := "cet4.xlsx"
	xlFile, error := xlsx.OpenFile(excelFileName)
	if error != nil {
		fmt.Println("err")
	}
	for _, sheet := range xlFile.Sheets {
		fmt.Println(sheet.Name)
		for _, row := range sheet.Rows {
			voc1 := new(Voc)
			voc1.Spell = row.Cells[0].Value
			voc1.Translation = row.Cells[1].Value
			Vocs = append(Vocs, voc1)
		}
	}
}

/*
func ExampleScrape() {
	doc, err := goquery.NewDocument("http://metalsucks.net")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".reviews-wrap article .review-rhs").Each(func(i int, s *goquery.Selection) {
		band := s.Find("h3").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
}
*/
