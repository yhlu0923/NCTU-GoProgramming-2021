package main
import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"strconv"
	"os"
	"strings"
)
//func lissajous(out io.Writer) 
//!-main
// Packages not needed by version in book.
import (
	"log"
	"net/http"
	"time"
)

//!+main

//var palette = []color.Color{color.White, color.RGBA{0x00,0xff,0x00,0xff}}
var palette = []color.Color{color.White,color.RGBA{0x00,0xff,0x00,0xff}, color.RGBA{0xff,0x00,0x00,0xff},color.RGBA{0x00,0x00,0xff,0xff}}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
	//blackIndex = 2
)
//func lissajous(out io.Writer)

func main() {
	//!-main
	// The sequence of images is deterministic unless we seed
	// the pseudo-random number generator using the current time.
	// Thanks to Randall McPherson for pointing out the omission.
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Args) > 1 && os.Args[1] == "web" {
			//!+http
			handler := func(w http.ResponseWriter, r *http.Request) {
				r.ParseForm()
				qcycles,err := strconv.Atoi(strings.Join(r.Form["cycles"],""))
				//fmt.Fprintf(os.Stdout,"Form: %v\n", r.Form)
				//fmt.Fprintf(os.Stdout,"Form[cycles]: %s\n",r.Form["cycles"])
				//fmt.Fprintf(os.Stdout,"%d\n",qcycles)
				if err != nil { qcycles = 5.0 }
				lissajous(w, float64(qcycles))
			}
			http.HandleFunc("/", handler)
			//!-http
			log.Fatal(http.ListenAndServe("localhost:8500", nil))
			return
	}
	//!+main
	//lissajous(os.Stdout)
	outfile, err := os.Create("out4.gif")
	if err != nil {
		log.Fatal(err)
	}
	lissajous(outfile, 5.0)
}
func lissajous(out io.Writer, qcycles float64) {
	const (
			//cycles  = 5     // number of complete x oscillator revolutions
			res     = 0.001 // angular resolution
			size    = 100   // image canvas covers [-size..+size]
			nframes = 64    // number of animation frames
			delay   = 8     // delay between frames in 10ms units
	)
	cycles := qcycles
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
			rect := image.Rect(0, 0, 2*size+1, 2*size+1)
			img := image.NewPaletted(rect, palette)
			for t := 0.0; t < cycles*2*math.Pi; t += res {
					x := math.Sin(t)
					y := math.Sin(t*freq + phase)
					//img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
					//        blackIndex)
				img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(t) % 4)
			}
			phase += 0.1
			anim.Delay = append(anim.Delay, delay)
			anim.Image = append(anim.Image, img)
		}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}