package main

import (
	"math"
	"sort"
	"strconv"
)

type TextNode struct{
	Value string
	Font string
	Position []float64
	Children []TextNode
}

func hierarchizeText(texts []TextType, titleFonts []Font) TextNode{

	titleFonts = sortFonts(texts, titleFonts)

	topNode := TextNode{
		Children: recursiveClassify(texts, titleFonts),
	}
	return topNode
}

func sortFonts(texts []TextType, titleFonts []Font) []Font{

	var extras []Font
	var indeces []int
	// purging of type type-Bold
	for i, v := range titleFonts{
		if v.name[len(v.name) - 4:] == "Bold"{
			for j := 0; j < len(titleFonts); j++{
				if titleFonts[j].name  + "-Bold" == v.name{
					indeces = append(indeces, i)
					extras = append(extras, v)
				}
			}
		}
	}

	for i := len(indeces) - 1; i > -1; i--{
		titleFonts = append(titleFonts[:indeces[i]], titleFonts[indeces[i]+1:]...)
	}

	cpy := make([]Font, len(titleFonts))
	copy(cpy, titleFonts)
	for i, v := range cpy{
		cpy[i].size = findWidth(v, texts)
	}

	sort.SliceStable(titleFonts, func(i, j int) bool {
		if math.Round(titleFonts[i].size) == math.Round(titleFonts[j].size){
			if math.Round(cpy[i].size) != math.Round(cpy[i].size){
				if titleFonts[i].name == titleFonts[j].name + "-Bold"{
					return true
				}else if titleFonts[j].name == titleFonts[i].name + "-Bold"{
					return false
				}
			}
			return cpy[i].size > cpy[j].size
		}
		return titleFonts[i].size > titleFonts[j].size
	})

	for _, v := range extras{
		originalFont := v.name[:len(v.name) - 5]
		idx := sort.Search(len(titleFonts), func(i int) bool { return titleFonts[i].name == originalFont })
		titleFonts = append(titleFonts[:idx + 1], titleFonts[idx:]...)
		titleFonts[idx] = v
	}

	return titleFonts
}

func findWidth(font Font, texts []TextType) float64{
	maxWidth := -1.
	for _, v := range texts{
		if font.name == v.text.Font && font.size == v.text.FontSize && v.text.W > maxWidth{
			return v.text.W
		}
	}
	return maxWidth
}

func recursiveClassify(texts []TextType, titleFonts []Font) []TextNode{
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
			if vTitle.name == v.text.Font && vTitle.size == v.text.FontSize{
				if i != 0 {
					nodes = append(nodes, TextNode{
						Value: texts[lastTitle].text.S,
						Font: vTitle.name + "-" + strconv.Itoa(int(vTitle.size)),
						Position: []float64{v.text.X, v.text.Y},
						Children: recursiveClassify(texts[lastTitle + 1:i], tempTitleFonts),
				})}
				cond = true
				lastTitle = i
			}
		}
		if cond {
			nodes = append(nodes, TextNode{
				Value: texts[lastTitle].text.S,
				Font: vTitle.name + "-" + strconv.Itoa(int(vTitle.size)),
				Position: []float64{texts[lastTitle].text.X, texts[lastTitle].text.Y},
				Children: recursiveClassify(texts[lastTitle + 1:], tempTitleFonts),
			})
			return nodes
		}
	}

	return transformToNodes(texts)
}

func transformToNodes(texts []TextType) []TextNode{
	var nodes = make([]TextNode, len(texts))

	for i, v := range texts{
		nodes[i] = TextNode{
			Value: v.text.S,
			Font: v.text.Font + "-" + strconv.Itoa(int(v.text.FontSize)),
			Position: []float64{v.text.X, v.text.Y},
			Children: nil,
		}
	}

	return nodes
}