package main

import (
	"go-engine/src/load"
	"go-engine/src/render"
	"go-engine/src/world"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func cleanupTempFiles() {
	for _, f := range load.TempPlantFiles {
		os.Remove(f) // Ignore errors, as they may no longer exist.
	}
}

/*func clearCustomMesh(model rl.Model) {
	model.Meshes.Vertices = nil
	model.Meshes.Normals = nil
	model.Meshes.Texcoords = nil
}

func GenMeshCustom() rl.Mesh {
	mesh := rl.Mesh{
		TriangleCount: 1,
		VertexCount:   3,
	}

	var vertices, normals, texcoords []float32

	vertices = addCoord(vertices, 0, 0, 0)
	vertices = addCoord(vertices, 1, 0, 2)
	vertices = addCoord(vertices, 2, 0, 0)
	mesh.Vertices = unsafe.SliceData(vertices)

	normals = addCoord(normals, 0, 1, 0)
	normals = addCoord(normals, 0, 1, 0)
	normals = addCoord(normals, 0, 1, 0)
	mesh.Normals = unsafe.SliceData(normals)

	texcoords = addCoord(texcoords, 0, 0)
	texcoords = addCoord(texcoords, 0.5, 1)
	texcoords = addCoord(texcoords, 1, 0)
	mesh.Texcoords = unsafe.SliceData(texcoords)

	rl.UploadMesh(&mesh, false)

	return mesh
}

func addCoord(slice []float32, values ...float32) []float32 {
	for _, value := range values {
		slice = append(slice, value)
	}
	return slice
}*/

func main() {
	rl.InitAudioDevice()

	game := load.InitGame()

	// Main game loop
	for !rl.WindowShouldClose() {
		// Toggle menu
		if rl.IsKeyPressed(rl.KeyTab) {
			render.ShowEscMenu = false
			render.ShowConfigMenu = !render.ShowConfigMenu

			if render.ShowConfigMenu {
				rl.EnableCursor()
			}
		}

		if rl.IsKeyPressed(rl.KeyEscape) {
			render.ShowConfigMenu = false
			render.ShowEscMenu = !render.ShowEscMenu

			if render.ShowEscMenu {
				rl.EnableCursor()
			}
		}

	
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) && render.ShouldUpdateCamera {
			rl.DisableCursor()
		}

		// Update the camera only when it is not in the menu.
		if render.ShouldUpdateCamera {
			rl.UpdateCamera(&game.Camera, game.CameraMode)
			load.MainPlayer.Camera.Position.Y += 10.0
		}

		// Manage chunks based on player's position
		world.ManageChunks(game.Worley, game.BiomeSelector, game.Camera.Position, game.ChunkCache, game.Perlin1, game.Perlin2, game.Perlin3)

		//  Draw
		render.RenderGame(&game)
	}

	rl.UnloadShader(game.Shader)					// unload shaders
	rl.UnloadMusicStream(load.RainSound)	// Unload source sound data
	rl.CloseAudioDevice()									// Close audio device

	// After the loop ends:
	defer rl.CloseWindow()

	cleanupTempFiles()
}
