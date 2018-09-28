package main

import _ "image/gif"
import _ "image/jpeg"
import _ "image/png"
import cf "github.com/lucasb-eyer/go-colorful"

import (
	"flag"
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"sort"
)

var sortFlag bool
var grid int
var allowDups bool

type palette struct {
	m map[cf.Color]int
	c []cf.Color
}

func newPalette() *palette {
	var p palette
	p.m = make(map[cf.Color]int)
	return &p
}

func (p *palette) add(c cf.Color) {
	if !allowDups {
		if _, ok := p.m[c]; ok {
			return
		}
	}
	p.m[c] = 1
	p.c = append(p.c, c)
}

func (p *palette) addImage(img image.Image) {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y += grid {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x += grid {
			p.add(cf.MakeColor(img.At(x, y)))
		}
	}
}

func (p *palette) Less(i, j int) bool {
	hi, si, li := p.c[i].Hsl()
	hj, sj, lj := p.c[j].Hsl()

	if true {
		//magi := hi + si + li
		//magj := hj + sj + lj
		a := p.c[i]
		b := p.c[j]
		magi := math.Sqrt(a.R*a.R + a.G*a.G + a.B*a.B)
		magj := math.Sqrt(b.R*b.R + b.G*b.G + b.B*b.B)
		return magi < magj
	}

	if hi < hj {
		return true
	} else if hj < hi {
		return false
	}

	if si < sj {
		return true
	} else if sj < si {
		return false
	}

	if li < lj {
		return true
	} else {
		return false
	}
}

func (p *palette) Len() int {
	return len(p.c)
}

func (p *palette) Swap(i, j int) {
	p.c[i], p.c[j] = p.c[j], p.c[i]
}

func main() {
	flag.BoolVar(&sortFlag, "sort", true, "sort the palette")
	flag.IntVar(&grid, "grid", 1, "sample grid spacing")
	flag.BoolVar(&allowDups, "allowdups", false, "allow duplicate colors")

	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("usage: img2pal images...")
	}

	pal := newPalette()

	for i := range flag.Args() {
		f, _ := os.Open(flag.Arg(i))
		img, _, err := image.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		pal.addImage(img)
	}

	if sortFlag {
		sort.Sort(pal)
	}

	fmt.Println("; Paint.NET Palette File")
	for i := range pal.c {
		r, g, b, _ := pal.c[i].RGBA()
		fmt.Printf("FF%02X%02X%02X\n", r/256, g/256, b/256)
	}
}
