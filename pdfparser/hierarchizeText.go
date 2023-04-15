package pdfparser

import (
	"sort"
	"strconv"
	"strings"

	"github.com/ledongthuc/pdf"
)

type TextNode struct{
	Value string
	Font string
	Position []float64
	Children []TextNode
}

func hierarchizeText(texts []pdf.Text, titleFonts []Font) TextNode{

	titleFonts = sortFonts(titleFonts)

	topNode := TextNode{
		Children: recursiveClassify(texts, titleFonts),
	}
	return topNode
}

func sortFonts(titleFonts []Font) []Font{
	var indeces []int

	baseFonts := make([]string, len(titleFonts))
	for i, v := range titleFonts{
		baseFonts[i] = strings.Split(v.name, "-")[0]
	}

	// purging of type same base fonts
	m := make(map[string][]Font)
	for i, v := range titleFonts{
		baseFont := strings.Split(v.name, "-")[0] + "-" + strconv.Itoa(int(v.size))

		if _, ok := m[baseFont]; ok {
			indeces = append(indeces, i)
			m[baseFont] = append(m[baseFont], v)
		}else{
			m[baseFont] = []Font{v};
		}
	}

	for k, v := range m{
		if len(v) < 2 {delete(m, k)}
	}

	for i := len(indeces) - 1; i > -1; i--{
		titleFonts = append(titleFonts[:indeces[i]], titleFonts[indeces[i]+1:]...)
	}

	indeces = nil

	sort.SliceStable(titleFonts, func(i, j int) bool {
		// if math.Round(titleFonts[i].size) == math.Round(titleFonts[j].size){
		// 	return titleFonts[i].width > titleFonts[j].width
		// }
		return titleFonts[i].size > titleFonts[j].size
	})

	// creating an order map
	order := []string{"Bold", "BoldItalic", "BoldOblique", "Base", "Italic", "Oblique"}
	orderMap := make(map[string]int)
	for i, v := range order{orderMap[v] = i}

	for k, v := range m{
		// sorting according to font style
		sort.SliceStable(v, func(i, j int) bool {
			return orderMap[getStyle(v[i])] < orderMap[getStyle(v[j])]
		})

		idx := findIndex(titleFonts, k)
		if idx + 1 < len(titleFonts){v = append(v, titleFonts[idx + 1:]...)}
		titleFonts = append(titleFonts[:idx], v...)
	}

	return titleFonts
}

func getStyle(font Font)string{
	fontAttr := strings.Split(font.name, "-")
	if len(fontAttr) < 2 {return "Base"}
	return fontAttr[1]
}

func findIndex(titlefonts []Font, k string) int{
	for i, v := range titlefonts{
		if strings.Split(v.name, "-")[0] + "-" + strconv.Itoa(int(v.size)) == k {
			return i
		}
	}
	return len(titlefonts)
}

func recursiveClassify(texts []pdf.Text, titleFonts []Font) []TextNode{
	if(titleFonts == nil){return transformToNodes(texts)}

	var nodes []TextNode = nil

	for iTitle, vTitle := range titleFonts{
		isLastTitle := iTitle == len(titleFonts) - 1
		var tempTitleFonts []Font
		if isLastTitle{
			tempTitleFonts = nil
		}else{tempTitleFonts = titleFonts[iTitle + 1:]}

		lastTitle := 0
		cond := false
		for i, v := range texts{
			if vTitle.name == v.Font && vTitle.size == v.FontSize{
				if i != 0 {
					nodes = append(nodes, TextNode{
						Value: texts[lastTitle].S,
						Font: vTitle.name + "-" + strconv.Itoa(int(vTitle.size)),
						Position: []float64{v.X, v.Y},
						Children: recursiveClassify(texts[lastTitle + 1:i], tempTitleFonts),
				})}
				cond = true
				lastTitle = i
			}
		}
		if cond {
			nodes = append(nodes, TextNode{
				Value: texts[lastTitle].S,
				Font: vTitle.name + "-" + strconv.Itoa(int(vTitle.size)),
				Position: []float64{texts[lastTitle].X, texts[lastTitle].Y},
				Children: recursiveClassify(texts[lastTitle + 1:], tempTitleFonts),
			})
			return nodes
		}
	}

	return transformToNodes(texts)
}

func transformToNodes(texts []pdf.Text) []TextNode{
	var nodes = make([]TextNode, len(texts))

	for i, v := range texts{
		nodes[i] = TextNode{
			Value: v.S,
			Font: v.Font + "-" + strconv.Itoa(int(v.FontSize)),
			Position: []float64{v.X, v.Y},
			Children: nil,
		}
	}

	return nodes
}