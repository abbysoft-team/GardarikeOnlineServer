package generation

import (
	"testing"
	"time"
)

func benchmarkGenerateTerrain(b *testing.B, octaves int, size int) {
	testGenerator := NewSimplexTerrainGenerator(TerrainGeneratorConfig{
		Octaves:     octaves,
		Persistence: 1,
		ScaleFactor: 1,
		Normalize:   true,
	}, time.Now().UnixNano())

	for i := 0; i < b.N; i++ {
		testGenerator.GenerateTerrain(size, size, float64(size)*float64(i), float64(size)*float64(i))
	}
}

func BenchmarkSimplexTerrainGenerator_GenerateTerrain10(b *testing.B) {
	benchmarkGenerateTerrain(b, 7, 10)
}

func BenchmarkSimplexTerrainGenerator_GenerateTerrain100(b *testing.B) {
	benchmarkGenerateTerrain(b, 7, 100)
}

func BenchmarkSimplexTerrainGenerator_GenerateTerrain1000(b *testing.B) {
	benchmarkGenerateTerrain(b, 7, 1000)
}

func BenchmarkSimplexTerrainGenerator_GenerateTerrain5000(b *testing.B) {
	benchmarkGenerateTerrain(b, 7, 5000)
}
