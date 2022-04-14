build:
	mkdir -p ./netlify/functions
	GOOS=linux GOARCH=amd64 go build -o ./netlify/functions/hello ./netlify/functions/hello.go