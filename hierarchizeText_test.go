package pdfparser

import (
	"reflect"
	"strconv"
	"testing"
)

func TestSortFonts(t *testing.T) {
	fonts := []font{
		{
			name:  "a",
			size:  12,
			width: 5,
		},
		{
			name:  "b",
			size:  11,
			width: 5,
		},
		{
			name:  "a-Bold",
			size:  12,
			width: 5,
		},
		{
			name:  "a-Italic",
			size:  12,
			width: 5,
		},
		{
			name:  "a",
			size:  11,
			width: 4,
		},
		{
			name:  "a-BoldItalic",
			size:  12,
			width: 5,
		},
		{
			name:  "b",
			size:  13,
			width: 5,
		},
	}

	baseline := []font{
		{
			name:  "b",
			size:  13,
			width: 5,
		},
		{
			name:  "a-Bold",
			size:  12,
			width: 5,
		},
		{
			name:  "a-BoldItalic",
			size:  12,
			width: 5,
		},
		{
			name:  "a",
			size:  12,
			width: 5,
		},
		{
			name:  "a-Italic",
			size:  12,
			width: 5,
		},
		{
			name:  "b",
			size:  11,
			width: 5,
		},
		{
			name:  "a",
			size:  11,
			width: 4,
		},
	}
	sortedFonts := sortFonts(fonts)

	if !reflect.DeepEqual(baseline, sortedFonts) {
		for _, v := range sortedFonts {
			t.Log(v.name, strconv.Itoa(int(v.size)))
		}

		t.Fatal("Wrong order for sorting")
	}
}
