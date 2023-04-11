# pdf-parser
A small server to parse a pdf and return it's content. The parser will return hierarchized data in JSON format with the positioning and style in the original pdf.

## Run the server
The server only has one route: `POST /parse-pdf`. It takes a pdf file and returns the JSON-formatted data. To run the server, you can either run the Docker container or if you have go installed, you can run `go run *.go` in your command line.

### A little docker tutorial
You can run the docker container by first building it, then running it
```sh
docker build -t "pdf-parser:Dockerfile" .
docker run -p 8080:8080 pdf-parser:Dockerfile
```

## I only want to parse a pdf pwease 😩
The `ParsePDF(filename string)` is exported for your convinience! It will export the JSON data in []byte format.

## Testing with Postman
If you would like to test the server, there are two routes:
- `GET /ping` is the testing route, to see if the base server works. It will return `pong` in json format.
- POST /pdf-parser is the route for extracting data from a pdf file. You can attach a file to the request by going to the `Body` tab, and selecting the form-data format. Then, in the key, select File type and type file. In the value, select your file.
