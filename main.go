package main

import (
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	PrevMousePosition Vector
	diaglo            float32 = 1900
	X, Y              float32 = 300, 200
	IsQ               bool
	count             int = 0
	ponit                 = [8][4]float32{{-50, -50, 0, 1}, {-50, 50, 0, 1}, {50, 50, 0, 1}, {50, -50, 0, 1}, {-50, -50, 100, 1}, {50, -50, 100, 1}, {-50, 50, 100, 1}, {50, 50, 100, 1}}
	ponitCopy             = [8][4]float32{{-50, -50, 0, 1}, {-50, 50, 0, 1}, {50, 50, 0, 1}, {50, -50, 0, 1}, {-50, -50, 100, 1}, {50, -50, 100, 1}, {-50, 50, 100, 1}, {50, 50, 100, 1}}
)

type Vector [2]float64

func (v *Vector) Sub(s Vector) Vector {
	return Vector{v[0] - s[0], v[1] - s[1]}
}

const (
	screenWidth  = 640
	screenHeight = 480
)

func Mux(matrix [4][4]float32, vect [4]float32) [4]float32 {
	return [4]float32{
		matrix[0][0]*vect[0] + matrix[0][1]*vect[1] + matrix[0][2]*vect[2] + matrix[0][3],
		matrix[1][0]*vect[0] + matrix[1][1]*vect[1] + matrix[1][2]*vect[2] + matrix[1][3],
		matrix[2][0]*vect[0] + matrix[2][1]*vect[1] + matrix[2][2]*vect[2] + matrix[2][3],
		matrix[3][0]*vect[0] + matrix[3][1]*vect[1] + matrix[3][2]*vect[2] + matrix[3][3]}
}

func Mux2(matrix [4][4]float32, vect [4]float32) [4]float32 {
	return [4]float32{
		matrix[0][0]*vect[0] + matrix[1][0]*vect[1] + matrix[2][0]*vect[2] + matrix[3][0]*vect[3],
		matrix[0][1]*vect[0] + matrix[1][1]*vect[1] + matrix[2][1]*vect[2] + matrix[3][1]*vect[3],
		matrix[0][2]*vect[0] + matrix[1][2]*vect[1] + matrix[2][2]*vect[2] + matrix[3][2]*vect[3],
		matrix[0][3]*vect[0] + matrix[1][3]*vect[1] + matrix[2][3]*vect[2] + matrix[3][3]*vect[3]}
}

func Get2DXY(v [4]float32) (float32, float32) {
	x := v[0] / v[3]
	y := v[1] / v[3]
	return x, y
}
func Mult(matrix [4][4]float32, other [4][4]float32) [4][4]float32 {
	var newMat [4][4]float32
	newMat[0][0] = matrix[0][0]*other[0][0] + matrix[0][1]*other[1][0] + matrix[0][2]*other[2][0] + matrix[0][3]*other[3][0]
	newMat[1][0] = matrix[1][0]*other[0][0] + matrix[1][1]*other[1][0] + matrix[1][2]*other[2][0] + matrix[1][3]*other[3][0]
	newMat[2][0] = matrix[2][0]*other[0][0] + matrix[2][1]*other[1][0] + matrix[2][2]*other[2][0] + matrix[2][3]*other[3][0]
	newMat[3][0] = matrix[3][0]*other[0][0] + matrix[3][1]*other[1][0] + matrix[3][2]*other[2][0] + matrix[3][3]*other[3][0]

	newMat[0][1] = matrix[0][0]*other[0][1] + matrix[0][1]*other[1][1] + matrix[0][2]*other[2][1] + matrix[0][3]*other[3][1]
	newMat[1][1] = matrix[1][0]*other[0][1] + matrix[1][1]*other[1][1] + matrix[1][2]*other[2][1] + matrix[1][3]*other[3][1]
	newMat[2][1] = matrix[2][0]*other[0][1] + matrix[2][1]*other[1][1] + matrix[2][2]*other[2][1] + matrix[2][3]*other[3][1]
	newMat[3][1] = matrix[3][0]*other[0][1] + matrix[3][1]*other[1][1] + matrix[3][2]*other[2][1] + matrix[3][3]*other[3][1]

	newMat[0][2] = matrix[0][0]*other[0][2] + matrix[0][1]*other[1][2] + matrix[0][2]*other[2][2] + matrix[0][3]*other[3][2]
	newMat[1][2] = matrix[1][0]*other[0][2] + matrix[1][1]*other[1][2] + matrix[1][2]*other[2][2] + matrix[1][3]*other[3][2]
	newMat[2][2] = matrix[2][0]*other[0][2] + matrix[2][1]*other[1][2] + matrix[2][2]*other[2][2] + matrix[2][3]*other[3][2]
	newMat[3][2] = matrix[3][0]*other[0][2] + matrix[3][1]*other[1][2] + matrix[3][2]*other[2][2] + matrix[3][3]*other[3][2]

	newMat[0][3] = matrix[0][0]*other[0][3] + matrix[0][1]*other[1][3] + matrix[0][2]*other[2][3] + matrix[0][3]*other[3][3]
	newMat[1][3] = matrix[1][0]*other[0][3] + matrix[1][1]*other[1][3] + matrix[1][2]*other[2][3] + matrix[1][3]*other[3][3]
	newMat[2][3] = matrix[2][0]*other[0][3] + matrix[2][1]*other[1][3] + matrix[2][2]*other[2][3] + matrix[2][3]*other[3][3]
	newMat[3][3] = matrix[3][0]*other[0][3] + matrix[3][1]*other[1][3] + matrix[3][2]*other[2][3] + matrix[3][3]*other[3][3]

	return newMat

}
func drawEbitenText(screen *ebiten.Image) {
	c := color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xff}
	x, y := Get2DXY(ponit[0])
	x1, y1 := Get2DXY(ponit[1])
	x2, y2 := Get2DXY(ponit[2])
	x3, y3 := Get2DXY(ponit[3])
	x4, y4 := Get2DXY(ponit[4])
	x5, y5 := Get2DXY(ponit[5])
	x6, y6 := Get2DXY(ponit[6])
	x7, y7 := Get2DXY(ponit[7])
	ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x1), float64(y1), c)
	ebitenutil.DrawLine(screen, float64(x1), float64(y1), float64(x2), float64(y2), c)
	ebitenutil.DrawLine(screen, float64(x2), float64(y2), float64(x3), float64(y3), c)
	ebitenutil.DrawLine(screen, float64(x3), float64(y3), float64(x), float64(y), c)
	ebitenutil.DrawLine(screen, float64(x4), float64(y4), float64(x), float64(y), c)
	ebitenutil.DrawLine(screen, float64(x3), float64(y3), float64(x5), float64(y5), c)
	ebitenutil.DrawLine(screen, float64(x4), float64(y4), float64(x5), float64(y5), c)
	ebitenutil.DrawLine(screen, float64(x1), float64(y1), float64(x6), float64(y6), c)
	ebitenutil.DrawLine(screen, float64(x2), float64(y2), float64(x7), float64(y7), c)
	ebitenutil.DrawLine(screen, float64(x6), float64(y6), float64(x7), float64(y7), c)
	ebitenutil.DrawLine(screen, float64(x4), float64(y4), float64(x6), float64(y6), c)
	ebitenutil.DrawLine(screen, float64(x5), float64(y5), float64(x7), float64(y7), c)
}

type Game struct {
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		IsQ = !IsQ
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		diaglo += 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		diaglo -= 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		X -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		Y -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		Y += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		X += 5
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	drawEbitenText(screen)
	mx, my := ebiten.CursorPosition()
	newV := Vector{float64(mx), float64(my)}
	diff := newV.Sub(PrevMousePosition)
	if IsQ {
		count += int(diff[0])
	} else {
		count += int(diff[1])
	}

	PrevMousePosition = newV
	//count++
	tihuan := [4][4]float32{}
	angle := float64(count%360) * math.Pi / 180
	//平移
	tihuan3 := [4][4]float32{}
	tihuan3[0][0] = 1
	tihuan3[1][1] = 1
	tihuan3[2][2] = 1
	tihuan3[3][0] = X
	tihuan3[3][1] = Y
	tihuan3[3][3] = 1

	//Y轴旋转
	if IsQ {
		tihuan[0][0] = float32(math.Cos(angle))
		tihuan[0][2] = float32(math.Sin(angle))
		tihuan[1][1] = 1
		tihuan[2][0] = -float32(math.Sin(angle))
		tihuan[2][2] = float32(math.Cos(angle))
		tihuan[3][3] = 1
	} else {
		//X轴旋转
		tihuan[0][0] = 1
		tihuan[1][1] = float32(math.Cos(angle))
		tihuan[1][2] = -float32(math.Sin(angle))
		tihuan[2][1] = float32(math.Sin(angle))
		tihuan[2][2] = float32(math.Cos(angle))
		tihuan[3][3] = 1
	}

	tihuan = Mult(tihuan, tihuan3)

	//移动z轴
	tihuan1 := [4][4]float32{}
	tihuan1[0][0] = 1
	tihuan1[1][1] = 1
	tihuan1[2][2] = 1
	//tihuan1[3][0] = -40
	tihuan1[3][2] = diaglo
	tihuan1[3][3] = 1
	//旋转
	// dd := 0.1
	// tihuan11 := [4][4]float32{}
	// tihuan11[0][0] = 1
	// tihuan11[1][1] = float32(math.Cos(dd))
	// tihuan11[1][2] = -float32(math.Sin(dd))
	// tihuan11[2][1] = float32(math.Sin(dd))
	// tihuan11[2][2] = float32(math.Cos(dd))
	// tihuan11[3][3] = 1
	// tihuan1 = Mult(tihuan1, tihuan11)

	tihuan2 := [4][4]float32{}
	tihuan2[0][0] = 1
	tihuan2[1][1] = 1
	tihuan2[2][2] = 1
	tihuan2[2][3] = 1.0 / diaglo

	tihuan = Mult(tihuan, tihuan1)
	tihuan = Mult(tihuan, tihuan2)

	for i := range ponit {
		ponit[i] = Mux2(tihuan, ponitCopy[i])

	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{}
	os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("2D to 3D")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
