package specform

import (
	"fmt"
	"math"
)

// CosineSimilarity returns the cosine similarity between two vectors.
// It assumes both vectors are of equal length and non-zero.
func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		panic(fmt.Sprintf("CosineSimilarity: vector length mismatch (%d != %d)", len(a), len(b)))
	}

	dot := 0.0
	normA := 0.0
	normB := 0.0

	for i := 0; i < len(a); i++ {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0.0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
