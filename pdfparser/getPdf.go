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
				if text.W != 0{
					sum += text.W
					len++
				}

				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				lastTextStyle.W = sum / len

				sections = addString(lastTextStyle, sections, totalPage - pageIndex)
				lastTextStyle = text
				if text.W != 0{
					sum = text.W
					len = 1
				}
			}
		}
		sections = addString(lastTextStyle, sections, pageIndex)
	}

	sort.SliceStable(sections, func(i, j int) bool {
		if(math.Abs(sections[i].Y - sections[j].Y) < sections[i].FontSize / 2){
			return sections[i].X < sections[j].X
		}
		return sections[i].Y > sections[j].Y
	})

	for i := 1 ; i < len(sections) ; i++{
		if isSameParagraph(sections[i], sections[i - 1]){
			temp := sections[i]
			sections[i] = sections[i - 1]
			sections[i - 1] = temp
			i++
		}
	}

	return sections
}

func isSameParagraph(cur pdf.Text, prev pdf.Text) bool{
	sizeCond := cur.FontSize - prev.FontSize > 1
	positionalCond := math.Abs(prev.Y - cur.Y) < prev.FontSize * 3
	lengthCond := len(prev.S) < 40
	return positionalCond && lengthCond && sizeCond
}

func isSameSentence(prev pdf.Text, cur pdf.Text) bool{
	if(prev.S == "") {return false}
	styleCheck := math.Abs(prev.FontSize - cur.FontSize) < 1 && prev.Font == cur.Font
	heightCheck := math.Abs(prev.Y - cur.Y) < prev.FontSize * 8
	return styleCheck && heightCheck
}

func addString(cur pdf.Text, sections []pdf.Text, page int) []pdf.Text {
	pageHeight := 800
	cur.S = cleanString(cur.S)
	if(cur.S == ""){return sections}
	cur.Y += float64(page * pageHeight)
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