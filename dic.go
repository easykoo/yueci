package main

import (
	"encoding/csv"
	"fmt"
	//"github.com/tealeg/xlsx"
	"io"
	"os"
)

var Vocs []*Voc

func main() {
	Read()
	Write()
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
		//voc.Translation = record[1]
		GetInfo(voc)
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
		w.Write([]string{voc.Spell, voc.Pronunciation, voc.Translation, voc.Form, voc.Sample})
	}
	w.Flush()
	fmt.Println(len(Vocs))
	/*
		for _, voc := range Vocs {
			fmt.Printf("%s %s\n", voc.Spell, voc.Translation)
		}
	*/
}

/*
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

*/
