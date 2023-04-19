package pdfparser

func catch(e error) {
	if e != nil {
		panic(e)
	}
}
