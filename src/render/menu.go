package render

import (
	"fmt"
	"math"

	"go-engine/src/load"
	"go-engine/src/pkg"
	"go-engine/src/world"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/constraints"
)

var ShowConfigMenu = false
var ShowEscMenu = false
var ShouldUpdateCamera = true

type ConfigState struct {
	ShowFPS bool
	ShowPosition bool
	RenderClouds bool
	FrameRateLimit int
	OldFrameRateLimit int
	CloudHeight int
	DrawFrames bool
}

var menuScroll rl.Vector2
var menuView rl.Rectangle

var GameState = ConfigState{
	ShowFPS: true,
	ShowPosition: true,
	RenderClouds: true,
	FrameRateLimit: 120,
	OldFrameRateLimit: 120,
	CloudHeight: 100,
	DrawFrames: false,
}

func renderEscMenu(menuX, menuY, width int32) {
	rl.BeginScissorMode(
		int32(menuView.X),
		int32(menuView.Y),
		int32(menuView.Width),
		int32(menuView.Height),
	)

	if gui.Button(rl.NewRectangle(float32(menuX + 20), float32(menuY + 50), float32(width - 40), 40.0), "Exit Game") {
		rl.WindowShouldClose()
	}

	gui.Button(rl.NewRectangle(float32(menuX + 20), float32(menuY + 100), float32(width - 40), 40.0), "Game Settings")

	//gui.Button(rl.NewRectangle(float32(menuX + 20), float32(menuY + 150), float32(width - 40), 40.0), "Credits and Attribution")

	rl.EndScissorMode()
}

func renderConfigMenu(menuX, menuY, width int32) {
	rl.BeginScissorMode(
		int32(menuView.X),
		int32(menuView.Y),
		int32(menuView.Width),
		int32(menuView.Height),
	)

	offsetY := int32(menuScroll.Y)

	newButton(menuX+20, menuY+40+offsetY, float32(width-40), 40.0, &GameState.ShowPosition, "Show Player Position")

	newButton(menuX+20, menuY+90+offsetY, float32(width-40), 40.0, &GameState.ShowFPS, "Show FPS")

	newButton(menuX+20, menuY+140+offsetY, float32(width-40), 40.0, &GameState.DrawFrames, "Draw Framed View")

	//Y = Y + 60 + 30
	newGuiSlider(menuX+20, menuY+190+offsetY, float32(width-40), 40.0,
		&GameState.FrameRateLimit, 30, 120,
		fmt.Sprintf("FPS Limit: %d", GameState.FrameRateLimit),
	)

	newButton(menuX+20, menuY+270+offsetY, float32(width-40), 40.0, &GameState.RenderClouds, "Clouds")

	newGuiSlider(menuX+20, menuY+320+offsetY, float32(width-40), 40.0,
		&pkg.CloudHeight, 30, 120,
		fmt.Sprintf("Cloud Height: %d", pkg.CloudHeight),
	)

	newGuiSlider(menuX+20, menuY+400+offsetY, float32(width-40), 40.0,
		&pkg.ChunkDistance, 1, 10,
		fmt.Sprintf("View Distance: %d", pkg.ChunkDistance),
	)

	newButton(menuX+20, menuY+480+offsetY, float32(width-40), 40.0, &world.StopChunkLoading, "Stop Chunk Loading")

	newGuiSlider(menuX+20, menuY+530+offsetY, float32(width-40), 40.0,
		&load.FogCoefficient, 0.0, 0.1,
		fmt.Sprintf("Fog Density: %.3f", load.FogCoefficient),
	)

	newGuiSlider(menuX+20, menuY+620+offsetY, float32(width-40), 40.0,
		&baseVolume, 0.0, 1.0,
		fmt.Sprintf("Sound FX Volume: %.3f", baseVolume),
	)

	rl.EndScissorMode()
}

func newButton(menuX, menuY int32, buttonWidth, buttonHeight float32, isOn *bool, text string) {
	// raygui button
	if gui.Button(rl.NewRectangle(float32(menuX), float32(menuY), buttonWidth, buttonHeight),
		func() string {
			if *isOn {
				return fmt.Sprintf("%s: ON", text)
			}
			return fmt.Sprintf("%s: OFF", text)
		}()) {
		*isOn = !*isOn
	}
}

func newGuiSlider[T constraints.Integer | constraints.Float](menuX, menuY int32, barWidth, barHeight float32, value *T, minVal, maxVal float32, text string) {
	floatVal := float32(*value)

	rl.DrawText(text, menuX, menuY, 20, rl.DarkGray)

	floatVal = gui.Slider(rl.NewRectangle(float32(menuX), float32(menuY+30), barWidth, barHeight),
		"", "",
		floatVal, minVal, maxVal,
	)

	switch any(value).(type) {
	case *int:
		*value = T(int(math.Floor(float64(floatVal))))
	case *float32:
		*value = T(floatVal)
	}
}
