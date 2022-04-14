package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Pixel struct {
	X int   `json:"x"`
	Y int   `json:"y"`
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

var currentCanvasPixels = []Pixel{}

func generateImage(pixel_optional *Pixel) *string {
	col := color.RGBA{0, 100, 0, 255}
	background := image.NewRGBA(image.Rect(0, 0, 500, 500))
	draw.Draw(background, background.Bounds(), &image.Uniform{C: col}, image.Point{}, draw.Src)

	for _, curr_pixel := range currentCanvasPixels {
		drawPixel(curr_pixel.X, curr_pixel.Y, color.RGBA{curr_pixel.R, curr_pixel.G, curr_pixel.B, 255}, background)
	}

	f, err := os.Create("/tmp/image.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, background)
	buf := new(bytes.Buffer)
	png.Encode(buf, background)
	send := buf.Bytes()
	baseString := base64.StdEncoding.EncodeToString(send)

	return &baseString
}

func drawPixel(x int, y int, color color.RGBA, background *image.RGBA) {
	Pixel := image.Rect(x, y, x+1, y+1)
	draw.Draw(background, Pixel, &image.Uniform{C: color}, image.Point{}, draw.Src)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var newPixel Pixel
	if request.HTTPMethod == "POST" {
		body := request.Body
		json.Unmarshal([]byte(body), &newPixel)
		currentCanvasPixels = append(currentCanvasPixels, newPixel)
		return &events.APIGatewayProxyResponse{
			StatusCode:        200,
			Headers:           map[string]string{"Content-Type": "text/plain"},
			MultiValueHeaders: http.Header{"Set-Cookie": {"Ding", "Ping"}},
			Body:              *generateImage(&newPixel),
			IsBase64Encoded:   false,
		}, nil
	} else if request.HTTPMethod == "GET" {
		fmt.Println(currentCanvasPixels)
		return &events.APIGatewayProxyResponse{
			StatusCode:        200,
			Headers:           map[string]string{"Content-Type": "text/plain"},
			MultiValueHeaders: http.Header{"Set-Cookie": {"Ding", "Ping"}},
			Body:              *generateImage(&newPixel),
			IsBase64Encoded:   false,
		}, nil
	}

	return &events.APIGatewayProxyResponse{StatusCode: 500}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call
	lambda.Start(handler)
}
