package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gocv.io/x/gocv"
	"html/template"
	"image"
	"image/color"
	"log"
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// resizeMatrixByWidth resizes a matrix to conform to a provided maxWidth
func resizeMatrixByWidth(src gocv.Mat, maxWidth int) (gocv.Mat, error) {
	// Allocate a new matrix to hold the final thumbnail
	dest := gocv.NewMat()

	// Max Width for thumbnail
	mw := float64(maxWidth)
	// Width
	w := src.Cols()
	// Height
	h := src.Rows()
	// Aspect ratio
	as := float64(h) / float64(w)
	// New Width
	nw := mw
	// New Height
	nh := nw * as
	// Resize the source matrix to
	gocv.Resize(src, &dest, image.Pt(int(nw), int(nh)), 0, 0,
		gocv.InterpolationDefault)
	// Return the buffer
	return dest, nil
}

// rectCenter gets the center point of an image.Rectangle
func rectCenter(r1 image.Rectangle) image.Point {
	return image.Pt(r1.Min.X+r1.Dx()/2, r1.Min.Y+r1.Dy()/2)
}

// distance gets the distance between the center of two rectangles
func distance(r1 image.Rectangle, r2 image.Rectangle) float64 {
	return math.Sqrt(math.Pow(float64(rectCenter(r1).X-rectCenter(r2).X),
		2) + math.Pow(float64(rectCenter(r1).Y-rectCenter(r2).Y), 2))
}

func focalLength(knownPixel float64, knownDistance float64, knownWidth float64) float64 {
	return (knownPixel * knownDistance) / knownWidth
}

func distanceTo(knownWidth float64, pixelWidth float64, focalLength float64) float64 {
	return (knownWidth * focalLength) / pixelWidth
}

// drawDetectionCorners draws the lines at the corners of a detection and a cross-hair in the middle
func drawDetectionCorners(mat *gocv.Mat, name string, rect image.Rectangle, clr color.RGBA, point image.Point) error {
	// Define a local instance of the Bounds rectangle

	// Define an inset distance variable
	const insetDistance = -10
	// Create the inset rect to make the corners appear to have depth
	insetRect := rect.Inset(insetDistance)
	// Define macro variables to hold rect dimensions
	w := insetRect.Dx()
	h := insetRect.Dy()
	// Find the center point
	center := rectCenter(insetRect)
	// Define the four corners of the detection area
	corners := [4]image.Point{
		// Bottom Right
		image.Pt(1, 1),
		// Top Right
		image.Pt(1, -1),
		// Top Left
		image.Pt(-1, 1),
		// Bottom Left
		image.Pt(-1, -1),
	}
	// Define a line length
	line := 30
	// Define the line thickness
	const lineWeight = 2
	localOffsets := [4][4]int{
		// Bottom Right
		{-line, 0, 0, -line},
		// Top Right
		{0, line, -line, 0},
		// Top Left
		{line, 0, 0, -line},
		// Bottom Left
		{0, line, line, 0},
	}
	// Draw all four corners
	for i, loc := range localOffsets {
		// Compute the point corresponding to the correct corner
		local := center.Add(image.Pt((w/2)*corners[i].X, (h/2)*corners[i].Y))
		// Draw both lines to enclose the pseudo-square
		gocv.Line(mat, local, local.Add(image.Pt(loc[0], loc[1])), clr, lineWeight)
		gocv.Line(mat, local, local.Add(image.Pt(loc[2], loc[3])), clr, lineWeight)
	}

	// Define constants for ui
	const fontScale = 1.3
	const fontThickness = 2
	// Get text metrics for dynamically rendered text
	textDimensions := gocv.GetTextSize(name, gocv.FontHersheyDuplex, fontScale, fontThickness)
	// Draw the offset text
	gocv.PutText(mat, name,
		rect.Min.Add(image.Pt(rect.Dx()-textDimensions.X, textDimensions.Y)).Add(image.Pt(-16, 16)),
		gocv.FontHersheyDuplex, fontScale,
		clr, fontThickness)
	// Draw a crosshair in the middle
	return nil
}

func Crosshair(mat *gocv.Mat, point image.Point, clr color.RGBA) {
	line := 5
	lineWeight := 2
	gocv.Line(mat, image.Pt(point.X, point.Y-line), image.Pt(point.X, point.Y+line), clr, lineWeight/2)
	gocv.Line(mat, image.Pt(point.X-line, point.Y), image.Pt(point.X+line, point.Y), clr, lineWeight/2)
}

// createThumbnail generates a thumbnail for the client
func createThumbnail(src gocv.Mat, size int) ([]byte, error) {
	// Resize the input matrix to have a max width of 360 pixels
	matrix, err := resizeMatrixByWidth(src, size)

	if err != nil {
		return nil, err
	}
	// Convert to jpg
	buf := matToBase64(matrix)
	defer matrix.Close()
	// Return the buffer
	return []byte(buf), nil
}

type Landmarks struct {
	RightEye struct {
		Xa int `json:"xa"`
		Ya int `json:"ya"`
		Xb int `json:"xb"`
		Yb int `json:"yb"`
	} `json:"rightEye"`
	LeftEye struct {
		Xa int `json:"xa"`
		Ya int `json:"ya"`
		Xb int `json:"xb"`
		Yb int `json:"yb"`
	} `json:"leftEye"`
	Nose struct {
		Xa int `json:"x"`
		Ya int `json:"y"`
	} `json:"nose"`
}

type Prediction struct {
	Name      string    `json:"name"`
	Top       int       `json:"top"`
	Right     int       `json:"right"`
	Bottom    int       `json:"bottom"`
	Left      int       `json:"left"`
	Distance  float64   `json:"distance"`
	Landmarks Landmarks `json:"landmarks"`
}

// func FindFocalLength() float64 {
//
// }

// matToBase64 converts a gocv matrix into a base64 encoded string
func matToBase64(src gocv.Mat) string {
	// Encode the matrix into a jpg
	encoded, err := gocv.IMEncodeWithParams(".jpeg", src, []int{gocv.IMWriteJpegOptimize})
	if err != nil {
		return ""
	}
	// Close the image when the function exits
	defer encoded.Close()
	// Allocate a response buffer
	return string(encoded.GetBytes())
}

func main() {

	u := url.URL{Scheme: "ws", Host: "10.0.1.2:8765", Path: "/recognize"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	var predicts []Prediction

	done := make(chan struct{})
	c.SetCloseHandler(func(code int, text string) error {
		time.Sleep(time.Second * 5)
		fmt.Println("Reconnecting...")
		c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		return nil
	})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			var predictions []Prediction
			err = json.Unmarshal(message, &predictions)
			if err != nil {
				log.Println("read:", err)
				break
			}
			predicts = predictions
		}
	}()

	bufChan := make(chan []byte, 8)
	var frame bytes.Buffer

	m := sync.RWMutex{}

	go func() {
		for {
			select {
			case out := <-bufChan:
				m.Lock()
				frame.Reset()
				frame.Write(out)
				m.Unlock()
			}
		}
	}()

	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
		data := ""
		for {
			m.RLock()
			data = "--frame\r\n  Content-Type: image/jpeg\r\n\r\n" + frame.String() + "\r\n\r\n"
			m.RUnlock()
			time.Sleep(40 * time.Millisecond)
			_, err := w.Write([]byte(data))
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("./index.html")
		err := t.Execute(w, "index")
		if err != nil {
			return
		}
	})

	mat := gocv.NewMat()

	capture, err := gocv.OpenVideoCapture(0)
	if err != nil {
		fmt.Printf("Open device on '/dev/video%d' failed. (%s)\n", 0, err.Error())
		return
	}
	if ok := capture.Read(&mat); !ok {
		fmt.Printf("Failed to read from /dev/video%d\n", 0)
	}

	fmt.Printf("Cap settings: Pixel Format: %s, Gain: %.2f\n", capture.CodecString(),
		capture.Get(gocv.VideoCaptureGain))

	go func() {
		fmt.Println("Server running on :6688")
		err := http.ListenAndServe("10.0.1.2:6688", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()
	cycle := 0
	const scale = 2
	//
	colorWhite := color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}

	for {
		if !capture.IsOpened() {
			break
		}
		if ok := capture.Read(&mat); !ok {
			fmt.Printf("Failed to read from /dev/video%d\n", 0)
		}

		if mat.Empty() {
			fmt.Printf("Empty read from /dev/video%d\n", 0)
			continue
		}

		cycle++

		ref := gocv.NewMat()
		mat.CopyTo(&ref)
		for _, prediction := range predicts {

			rect := image.Rect(prediction.Left*scale, prediction.Top*scale, prediction.Right*scale, prediction.Bottom*scale)
			rightEye := image.Rect(prediction.Landmarks.RightEye.Xa*scale, prediction.Landmarks.RightEye.Ya*scale,
				prediction.Landmarks.RightEye.Xb*scale, prediction.Landmarks.RightEye.Yb*scale)

			leftEye := image.Rect(prediction.Landmarks.LeftEye.Xa*scale, prediction.Landmarks.LeftEye.Ya*scale,
				prediction.Landmarks.LeftEye.Xb*scale, prediction.Landmarks.LeftEye.Yb*scale)

			Crosshair(&ref, rectCenter(rightEye), colorWhite)
			Crosshair(&ref, rectCenter(leftEye), colorWhite)
			noseTop := rightEye.Union(leftEye)
			noseBottom := image.Pt(prediction.Landmarks.Nose.Xa*scale, prediction.Landmarks.Nose.Ya*scale)
			Crosshair(&ref, rectCenter(noseTop), colorWhite)

			// fmt.Println(image.Pt(prediction.Landmarks.Nose.Xa, prediction.Landmarks.Nose.Ya))
			err = drawDetectionCorners(&ref, prediction.Name, rect, color.RGBA{
				R: uint8((prediction.Distance) * 255),
				G: uint8((1 - prediction.Distance) * 255),
				B: 0,
				A: 255,
			}, noseBottom)
			if err != nil {
				return
			}
			Crosshair(&ref, noseBottom, color.RGBA{
				R: 255,
				G: 255,
				B: 255,
				A: 255,
			})

		}

		var refBuf []byte
		refBuf, err = createThumbnail(mat, mat.Cols()/scale)
		if err != nil {
			fmt.Println(err)
			break
		}

		if cycle > 4 {
			cycle = 0
			err = c.WriteMessage(websocket.BinaryMessage, refBuf)
			if err != nil {
				log.Println("write:", err)
			}
		}

		var fancy []byte
		fancy, err = createThumbnail(ref, ref.Cols())
		if err != nil {
			fmt.Println(err)
			break
		}

		tr := time.NewTimer(time.Millisecond * 20)
		select {
		case bufChan <- fancy:
			tr.Stop()
		case <-tr.C:
			break
		}
	}
}
