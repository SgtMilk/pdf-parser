# pdf-parser
A small server to parse a pdf and return it's content. The parser will return hierarchized data in JSON format with the positioning and style in the original pdf.

## Run the server
The server only has one route: `POST /parse-pdf`. It takes a pdf file and returns the JSON-formatted data. To run the server, you can either run the Docker container or if you have go installed, you can run `go run *.go` in your command line.

## I only want to parse a pdf pwease 😩
The `ParsePDF(filename string)` is exported for your convinience! It will export the JSON data in string format.

