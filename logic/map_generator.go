package logic

import (
	simplex "github.com/ojrac/opensimplex-go"
	"math"
	rpc "projectx-server/rpc/generated"
	"time"
)

type TerrainGenerator interface {
	GenerateTerrain(width int, height int) []*rpc.Vector3D
}

type SimplexTerrainGenerator struct {
	octaves     int
	persistence float64
}

func NewSimplexMapGenerator(octaves int, persistence float64) SimplexTerrainGenerator {
	return SimplexTerrainGenerator{
		octaves:     octaves,
		persistence: persistence,
	}
}

func (s SimplexTerrainGenerator) GenerateTerrain(width int, height int) (result []*rpc.Vector3D) {
	var generator simplex.Noise
	generator = simplex.New(time.Now().Unix())

	pixels := make([][]float64, width)
	maxNoise := 0.0
	minNoise := 0.0

	for x := 0; x < width; x++ {
		pixels[x] = make([]float64, height)

		for y := 0; y < height; y++ {
			noise := 0.0
			freq := 1.0

			for octave := 1; octave <= s.octaves; octave++ {
				nx := float64(x) / float64(width)
				ny := float64(y) / float64(height)

				noise += (1 / float64(octave)) * generator.Eval2(nx*freq, ny*freq)
				freq *= 2.0
			}

			pixels[x][y] = noise
			maxNoise = math.Max(noise, maxNoise)
			minNoise = math.Min(noise, minNoise)
		}
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			normalized := (pixels[x][y] - minNoise) / (maxNoise - minNoise)
			normalized = math.Pow(normalized, s.persistence)

			result = append(result, &rpc.Vector3D{
				X: float32(x),
				Y: float32(y),
				Z: float32(normalized),
			})
		}
	}

	return result
}
