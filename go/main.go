package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type vec3d struct {
	x, y, z float64
}

type triangle struct {
	p       []vec3d
	shading uint8
}

type mat4x4 struct {
	m [][]float64
}

type mesh struct {
	tris []triangle
}

func (m mesh) Len() int {
	return len(m.tris)
}

func (m mesh) Swap(i, j int) {
	m.tris[i], m.tris[j] = m.tris[j], m.tris[i]
}

func (m mesh) Less(i, j int) bool {
	z1 := (m.tris[i].p[0].z + m.tris[i].p[1].z + m.tris[i].p[2].z) / 3
	z2 := (m.tris[j].p[0].z + m.tris[j].p[1].z + m.tris[j].p[2].z) / 3
	return z1 > z2
}

func multiplymatrixvector(i, o *vec3d, m *mat4x4) {
	o.x = i.x*m.m[0][0] + i.y*m.m[1][0] + i.z*m.m[2][0] + m.m[3][0]
	o.y = i.x*m.m[0][1] + i.y*m.m[1][1] + i.z*m.m[2][1] + m.m[3][1]
	o.z = i.x*m.m[0][2] + i.y*m.m[1][2] + i.z*m.m[2][2] + m.m[3][2]
	w := i.x*m.m[0][3] + i.y*m.m[1][3] + i.z*m.m[2][3] + m.m[3][3]
	if w != 0 {
		o.x = o.x / w
		o.y = o.y / w
		o.z = o.z / w
	}
}

func main() {

	// fmt.Println("\n\n\nNEW")
	angle := 79.0
	doTheThing(angle)
	for {
		time.Sleep(time.Millisecond * 200)
		angle += 0.1
		doTheThing(angle)
	}
}

func doTheThing(fTheta float64) {
	var vCamera vec3d

	img := options{
		bgcolor:  color.RGBA{255, 255, 255, 255},
		fgcolor:  color.RGBA{0, 0, 0, 255},
		width:    200,
		height:   200,
		offset:   0,
		filename: "image.png",
	}
	img.init()
	img.clear()

	// meshCube := readObj("ship.obj")
	meshCube := readObj("tea.obj")
	// meshCube := mesh{
	// 	// south
	// 	[]triangle{
	// 		{[]vec3d{{0.0, 0.0, 0.0}, {0.0, 1.0, 0.0}, {1.0, 1.0, 0.0}}},
	// 		{[]vec3d{{0.0, 0.0, 0.0}, {1.0, 1.0, 0.0}, {1.0, 0.0, 0.0}}},
	// 		// east
	// 		{[]vec3d{{1.0, 0.0, 0.0}, {1.0, 1.0, 0.0}, {1.0, 1.0, 1.0}}},
	// 		{[]vec3d{{1.0, 0.0, 0.0}, {1.0, 1.0, 1.0}, {1.0, 0.0, 1.0}}},
	// 		// north
	// 		{[]vec3d{{1.0, 0.0, 1.0}, {1.0, 1.0, 1.0}, {0.0, 1.0, 1.0}}},
	// 		{[]vec3d{{1.0, 0.0, 1.0}, {0.0, 1.0, 1.0}, {0.0, 0.0, 1.0}}},
	// 		// west
	// 		{[]vec3d{{0, 0, 1}, {0, 1, 1}, {0, 1, 0}}},
	// 		{[]vec3d{{0, 0, 1}, {0, 1, 0}, {0, 0, 0}}},
	// 		// top
	// 		{[]vec3d{{0, 1, 0}, {0, 1, 1}, {1, 1, 1}}},
	// 		{[]vec3d{{0, 1, 0}, {1, 1, 1}, {1, 1, 0}}},
	// 		// bottom
	// 		{[]vec3d{{1, 0, 1}, {0, 0, 1}, {0, 0, 0}}},
	// 		{[]vec3d{{1, 0, 1}, {0, 0, 0}, {1, 0, 0}}},
	// 	}}

	// projection matrix
	fNear := 0.1
	fFar := 1000.0
	fFov := 90.0
	fAspectRatio := 100.0 / 100.0 // height/width
	fFovRad := 1.0 / math.Tan(fFov*0.5/180.0*3.14159)

	matProj := mat4x4{
		[][]float64{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		}}
	matProj.m[0][0] = fAspectRatio * fFovRad
	matProj.m[1][1] = fFovRad
	matProj.m[2][2] = fFar / (fFar - fNear)
	matProj.m[3][2] = (-fFar * fNear) / (fFar - fNear)
	matProj.m[2][3] = 1.0
	matProj.m[3][3] = 0.0

	// fTheta := 3.8

	matRotZ := mat4x4{
		[][]float64{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		}}

	matRotZ.m[0][0] = math.Cos(fTheta)
	matRotZ.m[0][1] = math.Sin(fTheta)
	matRotZ.m[1][0] = -math.Sin(fTheta)
	matRotZ.m[1][1] = math.Cos(fTheta)
	matRotZ.m[2][2] = 1
	matRotZ.m[3][3] = 1

	matRotX := mat4x4{
		[][]float64{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		}}

	matRotX.m[0][0] = 1
	matRotX.m[1][1] = math.Cos(fTheta / 2)
	matRotX.m[1][2] = math.Sin(fTheta / 2)
	matRotX.m[2][1] = -math.Sin(fTheta / 2)
	matRotX.m[2][2] = math.Cos(fTheta / 2)
	matRotX.m[3][3] = 1

	trisreadytoraster := mesh{}

	for _, tri := range meshCube.tris {

		trirotatedz := triangle{
			[]vec3d{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			}, 0}

		multiplymatrixvector(&tri.p[0], &trirotatedz.p[0], &matRotZ)
		multiplymatrixvector(&tri.p[1], &trirotatedz.p[1], &matRotZ)
		multiplymatrixvector(&tri.p[2], &trirotatedz.p[2], &matRotZ)

		trirotatedxz := triangle{
			[]vec3d{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			}, 0}

		multiplymatrixvector(&trirotatedz.p[0], &trirotatedxz.p[0], &matRotX)
		multiplymatrixvector(&trirotatedz.p[1], &trirotatedxz.p[1], &matRotX)
		multiplymatrixvector(&trirotatedz.p[2], &trirotatedxz.p[2], &matRotX)

		tritranslated := triangle{
			[]vec3d{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			}, 0}

		tritranslated = trirotatedxz
		distance := 30.0
		tritranslated.p[0].z = trirotatedxz.p[0].z + distance
		tritranslated.p[1].z = trirotatedxz.p[1].z + distance
		tritranslated.p[2].z = trirotatedxz.p[2].z + distance

		var normal, line1, line2 vec3d
		line1.x = tritranslated.p[1].x - tritranslated.p[0].x
		line1.y = tritranslated.p[1].y - tritranslated.p[0].y
		line1.z = tritranslated.p[1].z - tritranslated.p[0].z

		line2.x = tritranslated.p[2].x - tritranslated.p[0].x
		line2.y = tritranslated.p[2].y - tritranslated.p[0].y
		line2.z = tritranslated.p[2].z - tritranslated.p[0].z

		normal.x = line1.y*line2.z - line1.z*line2.y
		normal.y = line1.z*line2.x - line1.x*line2.z
		normal.z = line1.x*line2.y - line1.y*line2.x

		var l float64
		l = math.Sqrt(normal.x*normal.x + normal.y*normal.y + normal.z*normal.z)
		normal.x /= l
		normal.y /= l
		normal.z /= l

		// if normal.z < 0 {

		cameraAlign := normal.x*(tritranslated.p[0].x-vCamera.x) + normal.y*(tritranslated.p[0].y-vCamera.y) + normal.z*(tritranslated.p[0].z-vCamera.z)
		if cameraAlign < 0 {

			lightdirection := vec3d{0.0, 0.0, -1.0}
			l = math.Sqrt(lightdirection.x*lightdirection.x + lightdirection.y*lightdirection.y + lightdirection.z*lightdirection.z)
			lightdirection.x /= l
			lightdirection.y /= l
			lightdirection.z /= l

			lightdp := normal.x*lightdirection.x + normal.y*lightdirection.y + normal.z + lightdirection.z
			// fmt.Println(lightdp)

			triprojected := triangle{
				[]vec3d{
					{0, 0, 0},
					{0, 0, 0},
					{0, 0, 0},
				}, uint8(-lightdp * 255)}

			multiplymatrixvector(&tritranslated.p[0], &triprojected.p[0], &matProj)
			multiplymatrixvector(&tritranslated.p[1], &triprojected.p[1], &matProj)
			multiplymatrixvector(&tritranslated.p[2], &triprojected.p[2], &matProj)

			triprojected.p[0].x += 1.0
			triprojected.p[1].x += 1.0
			triprojected.p[2].x += 1.0
			triprojected.p[0].y += 1.0
			triprojected.p[1].y += 1.0
			triprojected.p[2].y += 1.0

			triprojected.p[0].x *= 0.5 * float64(img.width)  // 100 // width
			triprojected.p[0].y *= 0.5 * float64(img.height) // 100 // height
			triprojected.p[1].x *= 0.5 * float64(img.width)  // 100 // width
			triprojected.p[1].y *= 0.5 * float64(img.height) // 100 // height
			triprojected.p[2].x *= 0.5 * float64(img.width)  //  100 // width
			triprojected.p[2].y *= 0.5 * float64(img.height) // 100 // height

			trisreadytoraster.tris = append(trisreadytoraster.tris, triprojected)
		}

		sort.Sort(trisreadytoraster)

		for _, tri := range trisreadytoraster.tris {
			p0x := int(math.Round(tri.p[0].x))
			p0y := int(math.Round(tri.p[0].y))
			p1x := int(math.Round(tri.p[1].x))
			p1y := int(math.Round(tri.p[1].y))
			p2x := int(math.Round(tri.p[2].x))
			p2y := int(math.Round(tri.p[2].y))
			img2 := img
			img2.fgcolor = color.RGBA{tri.shading, 0, 0, 255}
			tofill := []point{
				{p0x, p0y, 0},
				{p1x, p1y, 0},
				{p2x, p2y, 0},
			}
			filltri(tofill, img2)
		}

		// line3(point{p0x, p0y, 0}, point{p1x, p1y, 0}, img)
		// line3(point{p1x, p1y, 0}, point{p2x, p2y, 0}, img)
		// line3(point{p2x, p2y, 0}, point{p0x, p0y, 0}, img)

		// drawline(lnpts(point{p0x, p0y, 0}, point{p1x, p1y, 0}), img2)
		// drawline(lnpts(point{p1x, p1y, 0}, point{p2x, p2y, 0}), img2)
		// drawline(lnpts(point{p2x, p2y, 0}, point{p0x, p0y, 0}), img2)
		// }
	}

	output, _ := os.Create(img.filename)
	defer output.Close()
	fmt.Println("Saving...")
	png.Encode(output, img.image)

}

type point struct {
	x, y, z int
}

func line3(p1, p2 point, o options) {
	horizΔ := p2.x - p1.x
	vertiΔ := p2.y - p1.y
	if horizΔ < 0 {
		horizΔ = -horizΔ
	}
	if vertiΔ < 0 {
		vertiΔ = -vertiΔ
	}
	xdir, ydir := 1, 1
	if p1.x > p2.x {
		xdir = -1
	}
	if p1.y > p2.y {
		ydir = -1
	}
	wrong := horizΔ - vertiΔ
	// fmt.Println("horiz:", horizΔ, "verti:", vertiΔ)
	// fmt.Println()
	// counter := 0
	for {
		// counter++
		// if counter > 500 {
		// 	break
		// }
		o.image.Set(int(p1.x), int(p1.y), o.fgcolor)
		if p1.x == p2.x && p1.y == p2.y {
			break
		}
		wcopy := wrong * 2
		// fmt.Println()
		// fmt.Println("wrong:", wrong, wcopy < horizΔ, wcopy > -vertiΔ)
		// fmt.Println(p1.x, p1.y, "-", p2.x, p2.y)
		// SUPER SMOOSH
		if wcopy < horizΔ { // !!
			p1.y += ydir
			wrong += horizΔ
		}
		if wcopy > -vertiΔ {
			p1.x += xdir
			wrong -= vertiΔ
		}
	}
	// fmt.Println("done")
}

type options struct {
	fgcolor, bgcolor      color.RGBA
	width, height, offset int
	image                 *image.RGBA
	filename              string
}

func (o *options) init() {
	o.image = image.NewRGBA(image.Rect(0, 0, o.width, o.height))
}

func (o options) clear() {
	for x := 0; x <= o.width; x++ {
		for y := 0; y <= o.height; y++ {
			o.image.Set(x, y, o.bgcolor)
		}
	}
}

func lnpts(p1, p2 point) []point {
	output := []point{}
	horizΔ := p2.x - p1.x
	vertiΔ := p2.y - p1.y
	if horizΔ < 0 {
		horizΔ = -horizΔ
	}
	if vertiΔ < 0 {
		vertiΔ = -vertiΔ
	}
	xdir, ydir := 1, 1
	if p1.x > p2.x {
		xdir = -1
	}
	if p1.y > p2.y {
		ydir = -1
	}
	wrong := horizΔ - vertiΔ
	for {
		output = append(output, point{p1.x, p1.y, 0})
		if p1.x == p2.x && p1.y == p2.y {
			break
		}
		wcopy := wrong * 2
		if wcopy < horizΔ {
			p1.y += ydir
			wrong += horizΔ
		}
		if wcopy > -vertiΔ {
			p1.x += xdir
			wrong -= vertiΔ
		}
	}
	return output
}

func drawpixels(pts []point, o options) {
	for i := 0; i < len(pts)-1; i++ {
		o.image.Set(pts[i].x, pts[i].y, o.fgcolor)
	}
}

// optimize with:
// - https://www.avrfreaks.net/sites/default/files/triangles.c
// - http://www.sunshine2k.de/coding/java/TriangleRasterization/TriangleRasterization.html
func filltri(pts []point, o options) {
	// number of lines
	ylo, yhi := pts[0].y, pts[0].y
	for i := 0; i < len(pts); i++ {
		if pts[i].y < ylo {
			ylo = pts[i].y
		}
		if pts[i].y > yhi {
			yhi = pts[i].y
		}
	}
	lines := yhi - ylo

	// img3 := o
	// img3.fgcolor = color.RGBA{0, 255, 0, 255}
	// xlot, xhit := pts[0].x, pts[0].x
	// for i := 0; i < len(pts); i++ {
	// 	if pts[i].x < xlot {
	// 		xlot = pts[i].x
	// 	}
	// 	if pts[i].x > xhit {
	// 		xhit = pts[i].x
	// 	}
	// }
	// line3(point{xlot - 1, ylo - 1, 0}, point{xlot - 1, yhi + 1, 0}, img3)
	// line3(point{xlot - 1, yhi + 1, 0}, point{xhit + 1, yhi + 1, 0}, img3)
	// line3(point{xhit + 1, ylo - 1, 0}, point{xhit + 1, yhi + 1, 0}, img3)
	// line3(point{xlot - 1, ylo - 1, 0}, point{xhit + 1, ylo - 1, 0}, img3)
	// line3(point{xlot, ylo, 0}, point{xhit, yhi, 0}, img3)
	// line3(point{xlot, ylo, 0}, point{xhit, yhi, 0}, img3)
	// line3(point{xlot, ylo, 0}, point{xhit, yhi, 0}, img3)

	// fmt.Println(ylo, yhi)
	// get edges
	var pixels []point
	pixels = append(pixels, lnpts(pts[0], pts[1])...)
	pixels = append(pixels, lnpts(pts[1], pts[2])...)
	pixels = append(pixels, lnpts(pts[2], pts[0])...)
	// go through lines
	// fmt.Println(lines)
	// for i := 0; i < 5; i++ {
	for i := 0; i < lines; i++ {
		// fmt.Println("yhi-ylo:", yhi-ylo, "\ti:", i)
		yrow := ylo + i
		// fmt.Println("row:", yrow)
		ypts := []point{}
		for j := 0; j < len(pixels)-1; j++ {
			// find pixels in current line
			if pixels[j].y == yrow {
				// fmt.Println("\tin row:", pixels[j])
				ypts = append(ypts, pixels[j])
			}
		}
		// find extremeties
		xhi, xlo := ypts[0].x, ypts[0].x
		// fmt.Println("\tstarti xhi, xlo:", xhi, xlo)
		// fmt.Println("\t\t", ypts)
		for k := 0; k < len(ypts); k++ {
			// fmt.Println("\t\t", ypts[k])
			if ypts[k].x < xlo {
				xlo = ypts[k].x
			}
			if ypts[k].x > xhi {
				xhi = ypts[k].x
			}
		}
		// fmt.Println("\tsorted xhi, xlo:", xhi, xlo)
		drawpixels(lnpts(point{xlo, yrow, 0}, point{xhi, yrow, 0}), o)
	}
	// o.image.Set(int(pixels[i].x), int(pixels[i].y), o.fgcolor)
}

func readObj(filename string) mesh {
	vertices := []vec3d{}
	triangles := []triangle{}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		switch str[0] {
		case "v":
			var x, y, z float64
			x, err = strconv.ParseFloat(str[1], 8)
			if err != nil {
				log.Fatal(err)
			}
			y, err = strconv.ParseFloat(str[2], 8)
			if err != nil {
				log.Fatal(err)
			}
			z, err = strconv.ParseFloat(str[3], 8)
			if err != nil {
				log.Fatal(err)
			}
			vertices = append(vertices, vec3d{x, y, z})
		case "f":
			var v1, v2, v3 int
			v1, err = strconv.Atoi(str[1])
			if err != nil {
				log.Fatal(err)
			}
			v2, err = strconv.Atoi(str[2])
			if err != nil {
				log.Fatal(err)
			}
			v3, err = strconv.Atoi(str[3])
			if err != nil {
				log.Fatal(err)
			}
			triangles = append(triangles, triangle{[]vec3d{vertices[v1-1], vertices[v2-1], vertices[v3-1]}, 0})
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(vertices)
	// fmt.Println(vertices[0], vertices[1], vertices[2])
	// fmt.Println(vertices[3], vertices[4], vertices[5])
	// fmt.Println(triangles[0])
	return mesh{triangles}
}
