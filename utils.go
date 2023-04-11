package main

import "fmt"

func printCyan(text string){
	colorReset := "\033[0m"
    colorCyan := "\033[36m"

	fmt.Println(string(colorCyan) + text + string(colorReset))
}