package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"log"
	"math"
)

const (
	screenWidth  = 1000
	screenHeight = 800
	centerX      = screenWidth / 2
	centerY      = screenHeight / 2
)

type Planet struct {
	Name   string
	Color  color.RGBA
	Radius float64
	Orbit  float64
	Speed  float64
	Angle  float64
}

type Game struct {
	planets     []Planet
	angleMoon   float64
	orbitTilt   float64
	orbitRotate float64
	scale       float64
	offsetX     float64
	offsetY     float64
}

func (g *Game) Update() error {
	// ↑↓ 仰角压缩
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.orbitTilt -= 0.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.orbitTilt += 0.01
	}
	if g.orbitTilt < 0.1 {
		g.orbitTilt = 0.1
	}
	if g.orbitTilt > 1.0 {
		g.orbitTilt = 1.0
	}

	// ←→ 左右旋转
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.orbitRotate -= 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.orbitRotate += 0.02
	}

	// +/- 缩放
	if ebiten.IsKeyPressed(ebiten.KeyEqual) || ebiten.IsKeyPressed(ebiten.KeyKPAdd) {
		g.scale *= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyMinus) || ebiten.IsKeyPressed(ebiten.KeyKPSubtract) {
		g.scale *= 0.98
	}

	// WASD 平移
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.offsetY -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.offsetY += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.offsetX -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.offsetX += 5
	}

	// R 重置视角
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.orbitTilt = 0.5
		g.orbitRotate = 0
		g.scale = 1
		g.offsetX = 0
		g.offsetY = 0
	}

	// 公转角度递增
	for i := range g.planets {
		g.planets[i].Angle += g.planets[i].Speed
	}

	// 月亮自转
	g.angleMoon += 0.05

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	// 太阳（中心）
	sx := float64(centerX) + g.offsetX
	sy := float64(centerY) + g.offsetY
	ebitenutil.DrawCircle(screen, sx, sy, 25*g.scale, color.RGBA{255, 200, 0, 255})

	var earthX, earthY float64

	for _, p := range g.planets {
		// 椭圆轨道线（简洁方式）
		for a := 0.0; a < 2*math.Pi; a += 0.05 {
			theta1 := a + g.orbitRotate
			theta2 := a + 0.05 + g.orbitRotate

			x1 := sx + g.scale*p.Orbit*math.Cos(theta1)
			y1 := sy + g.scale*p.Orbit*math.Sin(theta1)*g.orbitTilt
			x2 := sx + g.scale*p.Orbit*math.Cos(theta2)
			y2 := sy + g.scale*p.Orbit*math.Sin(theta2)*g.orbitTilt

			ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.RGBA{100, 100, 100, 120})
		}

		// 星球位置
		theta := p.Angle + g.orbitRotate
		x := sx + g.scale*p.Orbit*math.Cos(theta)
		y := sy + g.scale*p.Orbit*math.Sin(theta)*g.orbitTilt

		scaleFactor := 1.0 - 0.3*((y-sy)/(g.scale*p.Orbit))
		r := p.Radius * g.scale * scaleFactor
		ebitenutil.DrawCircle(screen, x, y, r, p.Color)

		if p.Name == "Earth" {
			earthX, earthY = x, y
		}
	}

	// 月亮绕地球
	moonX := earthX + 40*g.scale*math.Cos(g.angleMoon)
	moonY := earthY + 40*g.scale*math.Sin(g.angleMoon)*g.orbitTilt
	ebitenutil.DrawCircle(screen, moonX, moonY, 5*g.scale, color.RGBA{180, 180, 180, 255})

	// FPS 显示
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", ebiten.CurrentFPS()))

	// 操作说明
	ebitenutil.DebugPrintAt(screen,
		"↑↓：俯仰角   ←→：旋转视角   +/-：缩放\nWASD：平移视角   R：重置视角",
		10, screenHeight-50,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		planets: []Planet{
			{"Mercury", color.RGBA{200, 200, 200, 255}, 4, 60, 0.04, 0},
			{"Venus", color.RGBA{255, 165, 0, 255}, 6, 100, 0.02, 0},
			{"Earth", color.RGBA{0, 150, 255, 255}, 8, 160, 0.01, 0},
			{"Mars", color.RGBA{255, 80, 80, 255}, 7, 220, 0.007, 0},
			{"Jupiter", color.RGBA{230, 180, 100, 255}, 12, 300, 0.005, 0},
			{"Saturn", color.RGBA{210, 180, 140, 255}, 10, 370, 0.003, 0},
			{"Uranus", color.RGBA{173, 216, 230, 255}, 9, 430, 0.002, 0},
			{"Neptune", color.RGBA{100, 149, 237, 255}, 9, 490, 0.001, 0},
		},
		orbitTilt:   0.5,
		orbitRotate: 0,
		scale:       1.0,
		offsetX:     0,
		offsetY:     0,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("🌌 太阳系模拟器")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
