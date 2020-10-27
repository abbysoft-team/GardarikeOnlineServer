package logic

import (
	simplex "github.com/ojrac/opensimplex-go"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

type TerrainGenerator interface {
	GenerateTerrain(width int, height int) []float32
}

type SimplexTerrainGenerator struct {
	config TerrainGeneratorConfig
}

type TerrainGeneratorConfig struct {
	Octaves     int
	Persistence float64
	ScaleFactor float64
	Normalize   bool
}

func NewSimplexTerrainGenerator(config TerrainGeneratorConfig) SimplexTerrainGenerator {
	log.WithField("module", "terrain_generator").
		WithField("config", config).
		Info("Simplex terrain generator initialized")

	return SimplexTerrainGenerator{
		config: config,
	}
}

func (s SimplexTerrainGenerator) GenerateTerrain(width int, height int) (result []float32) {
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

			for octave := 1; octave <= s.config.Octaves; octave++ {
				nx := float64(x) / float64(width)
				ny := float64(y) / float64(height)

				noise += (1 / float64(octave)) * generator.Eval2(nx*freq, ny*freq)
				freq = math.Pow(2, float64(octave))
			}

			pixels[x][y] = noise
			maxNoise = math.Max(noise, maxNoise)
			minNoise = math.Min(noise, minNoise)
		}
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var normalized float64
			if s.config.Normalize {
				normalized = (pixels[x][y] - minNoise) / (maxNoise - minNoise)
			} else {
				normalized = pixels[x][y]
			}
			normalized = math.Pow(normalized, s.config.Persistence)

			result = append(result, float32(normalized*s.config.ScaleFactor))
		}
	}

	return result
}
