package pdfparser

import (
	"math"
	"pdf-parser/utils"
	"regexp"
	"sort"
	"strings"

	"github.com/ledongthuc/pdf"
)

func getPDF(path string) []pdf.Text {
	f, r, err := pdf.Open(path)
	utils.Catch(err)
	defer f.Close()

	sections := []pdf.Text{}

	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {continue}

		var lastTextStyle pdf.Text
		texts := p.Content().Text

		sum := 0.
		len := 0.
		for _, text := range texts {
			if isSameSentence(lastTextStyle, text) {
				sum += text.W
				len++

				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				lastTextStyle.W = sum / len

				sections = addString(lastTextStyle, sections)
				lastTextStyle = text
				sum = text.W
				len = 1
			}
		}
		sections = addString(lastTextStyle, sections)
	}

	sort.SliceStable(sections, func(i, j int) bool {
		if(math.Abs(sections[i].Y - sections[j].Y) < sections[i].FontSize / 2){
			return sections[i].X < sections[j].X
		}
		return sections[i].Y > sections[j].Y
	})
	
	return sections
}

func isSameSentence(prev pdf.Text, cur pdf.Text) bool{
	if(prev.S == "") {return false}
	styleCheck := math.Abs(prev.FontSize - cur.FontSize) < 1 && prev.Font == cur.Font
	heightCheck := math.Abs(prev.Y - cur.Y) < prev.FontSize * 2.5
	return styleCheck && heightCheck
}

func addString(cur pdf.Text, sections []pdf.Text) []pdf.Text {
	cur.S = cleanString(cur.S)
	if(cur.S == ""){return sections}
	sections = append(sections, cur)
	return sections
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