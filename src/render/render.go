package render

import (
	"fmt"
	"math/rand"
	"sort"

	"go-engine/src/load"
	"go-engine/src/pkg"
	"go-engine/src/world"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

<<<<<<< HEAD
func DrawFrameView(chunk *pkg.Chunk, chunkPos rl.Vector3) {
	verts := chunk.Vertices
  inds := chunk.Indices

	for i := 0; i < len(inds); i += 3 {
		i0 := int(inds[i]) * 3
		i1 := int(inds[i+1]) * 3
		i2 := int(inds[i+2]) * 3

		v0 := rl.NewVector3(verts[i0], verts[i0+1], verts[i0+2])
		v1 := rl.NewVector3(verts[i1], verts[i1+1], verts[i1+2])
		v2 := rl.NewVector3(verts[i2], verts[i2+1], verts[i2+2])

		x0 := rl.Vector3{
			X: v0.X + chunkPos.X,
			Y: v0.Y + chunkPos.Y,
			Z: v0.Z + chunkPos.Z,
		}

		x1 := rl.Vector3{
			X: v1.X + chunkPos.X,
			Y: v1.Y + chunkPos.Y,
			Z: v1.Z + chunkPos.Z,
		}

		x2 := rl.Vector3{
			X: v2.X + chunkPos.X,
			Y: v2.Y + chunkPos.Y,
			Z: v2.Z + chunkPos.Z,
		}

		rl.DrawLine3D(x0, x1, rl.Red)
		rl.DrawLine3D(x1, x2, rl.Red)
		rl.DrawLine3D(x2, x0, rl.Red)
	}
}
=======
var cloudColor rl.Color = world.BlockTypes["Cloud"].Color
>>>>>>> 940ec6882166068be7b17332f081b411320e56dc

func RenderVoxels(game *load.Game) {
	cam := game.Camera.Position

	rl.SetShaderValue(game.Shader, rl.GetShaderLocation(game.Shader, "viewPos"), []float32{cam.X, cam.Y, cam.Z}, rl.ShaderUniformVec3)

	// --- Global collection of transparencies ---
	var transparentItems []pkg.TransparentItem

	// --- Round 1: solids ---
	for coord, chunk := range game.ChunkCache.Active {
		// Converts chunk coordinate to actual position
		chunkPos := rl.NewVector3(
			float32(coord.X*pkg.ChunkSize),
			0,
			float32(coord.Z*pkg.ChunkSize),
		)

		if chunk.IsOutdated {
			BuildChunkMesh(game, chunk, chunkPos)
			chunk.IsOutdated = false // reset flag → do not rebuild each frame
		}

		// If the chunk has mesh, draw directly
		if chunk.Model.MeshCount > 0 && chunk.Model.Meshes != nil {
			rl.DrawModel(chunk.Model, chunkPos, 1.0, rl.White)
		}

		if GameState.DrawFrames {
			DrawFrameView(chunk, chunkPos)
		}

		// --- Passage 2: plants per chunk (without global sorting) ---
		for _, voxel := range chunk.SpecialVoxels {
			pos := rl.NewVector3(
				chunkPos.X+float32(voxel.Position.X),
				chunkPos.Y+float32(voxel.Position.Y),
				chunkPos.Z+float32(voxel.Position.Z),
			)

			if voxel.Type == "Plant" {
				// The plant model is a little off center so I just adjust it a little
				pos.Z -= 0.8
				rl.DrawModel(voxel.Model, pos, 0.4, rl.White)
				continue
			}

			transparentItems = append(transparentItems, pkg.TransparentItem{
				Position:       pos,
				Type:           voxel.Type,
				Color:          world.BlockTypes[voxel.Type].Color,
				IsSurfaceWater: voxel.Type == "Water" && voxel.IsSurface,
			})
		}
	}

	// --- Back-to-front sorting ---
	if game.Camera.Position.Y >= float32(pkg.CloudHeight) {
		sort.Slice(transparentItems, func(i, j int) bool {
			di := rl.Vector3Length(rl.Vector3Subtract(transparentItems[i].Position, cam))
			dj := rl.Vector3Length(rl.Vector3Subtract(transparentItems[j].Position, cam))
			return di > dj // furthest first
		})
	}

	// --- Round 3: transparent ---
	rl.SetBlendMode(rl.BlendAlpha)
	rl.DisableDepthMask()
	//rl.BeginShaderMode(game.Shader)
	for _, it := range transparentItems {
		switch it.Type {
		case "Water":
			p := rl.NewVector3(it.Position.X+0.5, it.Position.Y+0.5, it.Position.Z+0.5)

			// calculates light intensity for this position
			//lightIntensity := calculateLightIntensity(p, game.LightPosition)
			//litColor := applyLighting(it.Color, lightIntensity)

			rl.DrawPlane(p, rl.NewVector2(1.0, 1.0), it.Color)
		case "Cloud":
			/*
				lightIntensity := calculateLightIntensity(it.Position, game.LightPosition)
				litColor := applyLighting(it.Color, lightIntensity)
			*/
			p := rl.NewVector3(it.Position.X, float32(pkg.CloudHeight), it.Position.Z)

<<<<<<< HEAD
			if GameState.RenderClouds {
				/*CloudMesh := rl.Mesh{}

				vertices := []float32{}
				size := float32(1)
				offset := int32(len(vertices) / 3)

				vertices = append(vertices, 
				    p.X, p.Y, p.Z+size,
				    p.X+size, p.Y, p.Z+size,
				    p.X+size, p.Y+size, p.Z+size,
				    p.X, p.Y+size, p.Z+size,
				)

				indices := []uint16{}
				indices = append(indices, 
				    uint16(offset+0), uint16(offset+1), uint16(offset+2),
				    uint16(offset+0), uint16(offset+2), uint16(offset+3),
				)

				normals := []float32{}
				for i := 0; i < 4; i++ {
				    normals = append(normals, 0, 0, 1)
				}
				//texcoords := []float32{}

				CloudMesh.VertexCount = int32(len(vertices) / 3)
				CloudMesh.TriangleCount = int32(len(indices) / 3)

				CloudMesh.Vertices = &vertices[0]
				CloudMesh.Indices = &indices[0]
				CloudMesh.Normals = &normals[0]
				//CloudMesh.Texcoords = &texcoords[0]

				rl.UploadMesh(&CloudMesh, false)
				rl.LoadModelFromMesh(CloudMesh)*/

				rl.DrawCube(p, 1.0, 1.0, 1.0, it.Color)
=======
			if ShowClouds {
				rl.DrawCube(p, 1.0, 0.0, 1.0, cloudColor)
>>>>>>> 940ec6882166068be7b17332f081b411320e56dc
			}
		}
	}

	//rl.EndShaderMode()
	rl.EnableDepthMask()
	rl.SetBlendMode(rl.BlendMode(0))
}

func applyUnderwaterEffect(game *load.Game) {
	waterLevel := int(float64(pkg.WorldHeight)*pkg.WaterLevelFraction) + 1

	coord := world.ToChunkCoord(game.Camera.Position)
	chunk := game.ChunkCache.Active[coord]

	if chunk == nil {
		return
	}

	localX := int(game.Camera.Position.X) - coord.X*pkg.ChunkSize
	localY := int(game.Camera.Position.Y) - coord.Y*pkg.WorldHeight
	localZ := int(game.Camera.Position.Z) - coord.Z*pkg.ChunkSize

	if localX >= 0 && localX < pkg.ChunkSize &&
		localY >= 0 && localY < pkg.WorldHeight &&
		localZ >= 0 && localZ < pkg.ChunkSize {

			voxel := chunk.Voxels[localX][localY][localZ]
			if voxel.Type != "Water" || game.Camera.Position.Y > float32(waterLevel)-0.5 { return }
			
			// apply blue overlay
			rl.SetBlendMode(rl.BlendMode(0))
			rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.NewColor(0, 0, 255, 100))

			if shouldRain == 1 {
				rl.SetMusicVolume(load.RainSound, currentVolume*0.3)
			}
	}
}

var shouldRain = 0
var nextWeatherChange = 0
var targetFogDensity = float32(1)

func colorToVec4(c rl.Color) []float32 {
	return []float32{
		float32(c.R) / 255.0,
		float32(c.G) / 255.0,
		float32(c.B) / 255.0,
		float32(c.A) / 255.0,
	}
}

func RenderGame(game *load.Game) {
	rl.BeginDrawing()
	rl.ClearBackground(skyColor)
	//rl.ClearBackground(rl.NewColor(134, 13, 13, 255))  Red

	load.UpdateTimer()

	// Only change the climate when the scheduled time comes
	if load.ElapsedSeconds >= nextWeatherChange {
		shouldRain = rand.Intn(4)
		nextWeatherChange = load.ElapsedSeconds + 60 // shedules next change
	}

	updateSkyColor()

	updateFog(game, targetFogDensity, colorToVec4(skyColor))

	updateAmbient()

	locAmbient := rl.GetShaderLocation(game.Shader, "ambient")
	rl.SetShaderValue(game.Shader, locAmbient, []float32{load.Ambient}, rl.ShaderUniformFloat)

	rl.BeginMode3D(game.Camera)

	//	Begin drawing solid blocks and then transparent ones (avoid flickering)
	RenderVoxels(game)

	updateRainAudio(0.0)

	targetFogDensity = 1
	// clean up particles when it stops raining.
	pkg.RainDrops = nil

	if shouldRain == 1 {
		targetFogDensity = 1.5
		updateRainAudio(1.0)

		// only initializes if there are not enough particles yet.
		if len(pkg.RainDrops) < 400 {
			initRain(game, 400, 40)
		}

		updateRain(game, 40)
		drawRain()
	}

	rl.EndMode3D()

	applyUnderwaterEffect(game)

	if ShowConfigMenu {
		ShouldUpdateCamera = false
		menuWidth := int32(rl.GetScreenWidth()) / 2
		menuHeight := int32(rl.GetScreenHeight()) / 2
		menuX := (int32(rl.GetScreenWidth()) - menuWidth) / 2
		menuY := (int32(rl.GetScreenHeight()) - menuHeight) / 2

		// Visible menu area
		menuBounds := rl.NewRectangle(float32(menuX), float32(menuY), float32(menuWidth), float32(menuHeight))

		contentHeight := float32(700) // Actual height of the content, including what is not visible.
		contentBounds := rl.NewRectangle(0, 0, float32(menuWidth-20), contentHeight)

		gui.ScrollPanel(menuBounds, "Game Settings", contentBounds, &menuScroll, &menuView)

		renderConfigMenu(menuX, menuY, menuWidth)
	} else if ShowEscMenu {
		ShouldUpdateCamera = false
		menuWidth := int32(rl.GetScreenWidth()) / 2
		menuHeight := int32(rl.GetScreenHeight()) / 2
		menuX := (int32(rl.GetScreenWidth()) - menuWidth) / 2
		menuY := (int32(rl.GetScreenHeight()) - menuHeight) / 2

		menuBounds := rl.NewRectangle(float32(menuX), float32(menuY), float32(menuWidth), float32(menuHeight))

		contentHeight := float32(200)
		contentBounds := rl.NewRectangle(0, 0, float32(menuWidth-20), contentHeight)

		gui.ScrollPanel(menuBounds, "Pause Menu", contentBounds, &menuScroll, &menuView)
		renderEscMenu(menuX, menuY, menuWidth)
	} else {
		ShouldUpdateCamera = true
	}

	if GameState.ShowFPS {
		rl.DrawFPS(10, 30)
	}

	if GameState.FrameRateLimit != GameState.OldFrameRateLimit {
		rl.SetTargetFPS(int32(GameState.FrameRateLimit))
	}

	if GameState.ShowPosition {
		positionText := fmt.Sprintf("Player's position: (%.2f, %.2f, %.2f)", game.Camera.Position.X, game.Camera.Position.Y, game.Camera.Position.Z)
		rl.DrawText(positionText, 10, 5, 20, rl.DarkGreen)
	}

	rl.EndDrawing()
}
