package generation

import (
	simplex "github.com/ojrac/opensimplex-go"
	log "github.com/sirupsen/logrus"
	"math"
)

type TerrainGenerator interface {
	GenerateTerrain(width int, height int, offsetX, offsetY float64) []float32
	SetSeed(seed int64)
}

type SimplexTerrainGenerator struct {
	config    TerrainGeneratorConfig
	generator simplex.Noise
}

type TerrainGeneratorConfig struct {
	Octaves     int
	Persistence float64
	ScaleFactor float64
	Normalize   bool
}

func NewSimplexTerrainGenerator(config TerrainGeneratorConfig, seed int64) *SimplexTerrainGenerator {
	log.WithField("module", "terrain_generator").
		WithField("config", config).
		Info("Simplex terrain generator initialized")

	return &SimplexTerrainGenerator{
		config:    config,
		generator: simplex.New(seed),
	}
}

func (s *SimplexTerrainGenerator) SetSeed(seed int64) {
	s.generator = simplex.New(seed)
}

func (s SimplexTerrainGenerator) GenerateTerrain(width, height int, offsetX, offsetY float64) (result []float32) {
	pixels := make([][]float64, width)
	maxNoise := 0.0
	minNoise := 0.0

	for x := 0; x < width; x++ {
		pixels[x] = make([]float64, height)

		x := x
		go func() {
			for y := 0; y < height; y++ {
				noise := 0.0
				freq := 2.0

				for octave := 0; octave < s.config.Octaves; octave++ {
					// Freq is always growing
					freq = math.Pow(2, float64(octave))
					amplitude := math.Pow(s.config.Persistence, float64(octave))

					nx := (float64(x) + offsetX) / float64(width)
					ny := (float64(y) + offsetY) / float64(height)

					// Multiply by amplitude and map noise value from [-1;1] to [0;1]
					noise += amplitude * (s.generator.Eval2(nx*freq, ny*freq))
				}

				pixels[x][y] = noise
			}
		}()
	}

	//Calculate min/max noise
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			noise := pixels[x][y]
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
			//normalized = math.Pow(normalized, s.config.Persistence)

			result = append(result, float32(normalized*s.config.ScaleFactor))
		}
	}

	return result
}
