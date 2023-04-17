package utils

import "fmt"

func PrintCyan(text string) {
	colorReset := "\033[0m"
	colorCyan := "\033[36m"

	fmt.Println(colorCyan + text + colorReset)
}

func Catch(e error) {
	if e != nil {
		panic(e)
	}
}
