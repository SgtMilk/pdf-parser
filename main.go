package main

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/ledongthuc/pdf"
)

func main(){
	data , err := getPDF("./source/resume-en.pdf")
	catch(err)

	for _, s := range data {
		fmt.Println(s.S)
	}
}

func catch(e error){
	if(e != nil) {panic(e)}
}

func isSameSentence(prev pdf.Text, cur pdf.Text) bool{
	if(prev.S == "") {return false}
	styleCheck := math.Abs(prev.FontSize - cur.FontSize) < 1
	heightCheck := math.Abs(prev.Y - cur.Y) < prev.FontSize
	return styleCheck && heightCheck
}

func cleanString(s string) string{
	re := regexp.MustCompile(`[^[:print:]À-ÿ]`)
	s = re.ReplaceAllLiteralString(s, "")
	re = regexp.MustCompile(`/s+`)
	s = re.ReplaceAllLiteralString(s, " ")
	s = strings.ReplaceAll(s, "\n", "")
	if(s != "" && s[0] == ' ') {return s[1:]}
	return s
}

func addString(cur pdf.Text, sections []pdf.Text) []pdf.Text {
	cur.S = cleanString(cur.S)
	if(cur.S == ""){return sections}
	sections = append(sections, cur)
	return sections
}

func getPDF(path string) ([]pdf.Text, error) {
	f, r, err := pdf.Open(path)
	// remember close file
	if err != nil {
		return nil, err
	}

	defer f.Close()

	totalPage := r.NumPage()

	sections := []pdf.Text{}

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		var lastTextStyle pdf.Text
		texts := p.Content().Text
		for _, text := range texts {
			if isSameSentence(lastTextStyle, text) {
				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				sections = addString(lastTextStyle, sections)
				lastTextStyle = text
			}
		}
		sections = addString(lastTextStyle, sections)
	}

	sort.SliceStable(sections, func(i, j int) bool {
		if(math.Abs(sections[i].Y - sections[j].Y) < sections[i].FontSize / 2){
			return sections[i].X > sections[j].X
		}
		return sections[i].Y > sections[j].Y
	})
	
	return sections, nil
}
