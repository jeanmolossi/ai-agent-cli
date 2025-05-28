package similarity

import "math"

func CosSimilarity(a, b [][]float32) (float64, error) {
	if len(a) != len(b) {
		return 0, ErrIncompatibleVector
	}

	var dot, magA, magB float64

	for i := range len(a) {
		rowA := a[i]
		rowB := b[i]

		if len(rowA) != len(rowB) {
			return 0, ErrIncompatibleRow
		}

		for j := range len(rowA) {
			x := float64(rowA[j])
			y := float64(rowB[j])

			dot += x * y
			magA += x * x
			magB += y * y
		}
	}

	if magA == 0 || magB == 0 {
		return 0, nil
	}

	sim := dot / (math.Sqrt(magA) * math.Sqrt(magB))

	return sim, nil
}
