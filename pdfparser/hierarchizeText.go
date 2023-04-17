package pdfparser

import (
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type TextNode struct {
	Value    string
	Font     string
	Children []TextNode
	Position Rect
}

func hierarchizeText(texts []Section, titleFonts []Font) TextNode {
	titleFonts = sortFonts(titleFonts)

	topNode := TextNode{
		Children: recursiveClassify(texts, titleFonts),
	}

	return topNode
}

func sortFonts(titleFonts []Font) []Font {
	var indices []int

	baseFonts := make([]string, len(titleFonts))
	for i, v := range titleFonts {
		baseFonts[i] = strings.Split(v.name, "-")[0]
	}

	// purging of type same base fonts
	m := make(map[string][]Font)

	for i, v := range titleFonts {
		baseFont := strings.Split(v.name, "-")[0] + "-" + strconv.Itoa(int(v.size))

		if _, ok := m[baseFont]; ok {
			indices = append(indices, i)
			m[baseFont] = append(m[baseFont], v)
		} else {
			m[baseFont] = []Font{v}
		}
	}

	for k, v := range m {
		if len(v) <= 1 {
			delete(m, k)
		}
	}

	for i := len(indices) - 1; i > -1; i-- {
		titleFonts = append(titleFonts[:indices[i]], titleFonts[indices[i]+1:]...)
	}

	sort.SliceStable(titleFonts, func(i, j int) bool {
		if math.Round(titleFonts[i].size) == math.Round(titleFonts[j].size) {
			return titleFonts[i].width > titleFonts[j].width
		}
		return titleFonts[i].size > titleFonts[j].size
	})

	// creating an order map
	order := []string{"Bold_CAPS", "Bold", "BoldItalic_CAPS", "BoldItalic", "BoldOblique_CAPS", "BoldOblique", "Base_CAPS", "Base", "Italic_CAPS", "Italic", "Oblique_CAPS", "Oblique"}
	orderMap := make(map[string]int)

	for i, v := range order {
		orderMap[v] = i
	}

	for k, v := range m {
		// sorting according to font style
		sort.SliceStable(v, func(i, j int) bool {
			return orderMap[getStyle(v[i])] < orderMap[getStyle(v[j])]
		})

		idx := findIndex(titleFonts, k)
		if idx+1 < len(titleFonts) {
			v = append(v, titleFonts[idx+1:]...)
		}

		titleFonts = append(titleFonts[:idx], v...)
	}

	return titleFonts
}

func getStyle(font Font) string {
	fontAttr := strings.Split(font.name, "-")
	if len(fontAttr) <= 1 {
		return "Base"
	}

	if fontAttr[1] == "CAPS" {
		return "Base-CAPS"
	}

	re := regexp.MustCompile(`\W*((?i)bold|bolditalic|boldoblique|italic|oblique(?-i))\W*`)
	match := re.FindStringSubmatch(fontAttr[1])

	var style string
	if match != nil {
		style = match[1]
	} else {
		style = "Base"
	}

	if len(fontAttr) > 2 && fontAttr[len(fontAttr)-1] == "CAPS" {
		style += "-CAPS"
	}

	return style
}

func findIndex(titlefonts []Font, k string) int {
	for i, v := range titlefonts {
		if strings.Split(v.name, "-")[0]+"-"+strconv.Itoa(int(v.size)) == k {
			return i
		}
	}

	return len(titlefonts)
}

func recursiveClassify(texts []Section, titleFonts []Font) []TextNode {
	if titleFonts == nil {
		return transformToNodes(texts)
	}

	var nodes []TextNode

	for iTitle, vTitle := range titleFonts {
		var tempTitleFonts []Font
		if iTitle == len(titleFonts)-1 { // if its the last title font
			tempTitleFonts = nil
		} else {
			tempTitleFonts = titleFonts[iTitle+1:]
		}

		lastTitle := 0
		cond := false

		for i, v := range texts {
			if vTitle.name == v.text.Font && vTitle.size == v.text.FontSize {
				if i != 0 {
					if nodes == nil && !(vTitle.name == texts[lastTitle].text.Font && vTitle.size == texts[lastTitle].text.FontSize) {
						nodes = append(nodes, TextNode{
							Children: recursiveClassify(texts[lastTitle:i], tempTitleFonts),
						})
					} else {
						nodes = append(nodes, TextNode{
							Value:    texts[lastTitle].text.S,
							Font:     vTitle.name + "-" + strconv.Itoa(int(vTitle.size)),
							Position: v.position,
							Children: recursiveClassify(texts[lastTitle+1:i], tempTitleFonts),
						})
					}
				}

				cond = true
				lastTitle = i
			}
		}

		if cond {
			nodes = append(nodes, TextNode{
				Value:    texts[lastTitle].text.S,
				Font:     vTitle.name + "-" + strconv.Itoa(int(vTitle.size)),
				Position: texts[lastTitle].position,
				Children: recursiveClassify(texts[lastTitle+1:], tempTitleFonts),
			})

			return nodes
		}
	}

	return transformToNodes(texts)
}

func transformToNodes(texts []Section) []TextNode {
	var nodes = make([]TextNode, len(texts))

	for i, v := range texts {
		nodes[i] = TextNode{
			Value:    v.text.S,
			Font:     v.text.Font + "-" + strconv.Itoa(int(v.text.FontSize)),
			Position: v.position,
			Children: nil,
		}
	}

	return nodes
}
