package position

import (
	"math"
	"testing"
)

func TestCoordinateTransformConsistency(t *testing.T) {
	// Test radar position (approximate Swiss radar)
	radarPos := RadarPosition{
		Latitude:  46.5197,
		Longitude: 6.6323,
		Height:    430.0,
		SIC:       33,
		SAC:       40,
	}

	// Test cases from real data
	testCases := []struct {
		name       string
		rangeNM    float64
		azimuthDeg float64
		xNM        float64
		yNM        float64
	}{
		{"Case1", 96.1640625, 31.3330078125, 50.0078125, 82.140625},
		{"Case2", 124.78125, 33.55224609375, 68.96875, 103.9921875},
		{"Case3", 98.4765625, 36.67236328125, 58.8125, 78.984375},
		{"Case4", 98.4296875, 39.17724609375, 62.1796875, 76.3046875},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert polar coordinates to WGS84
			rangeM := NauticalMilesToMeters(tc.rangeNM)
			latPolar, lonPolar := PolarToWGS84(radarPos, rangeM, tc.azimuthDeg)

			// Convert Cartesian coordinates to WGS84
			xM := NauticalMilesToMeters(tc.xNM)
			yM := NauticalMilesToMeters(tc.yNM)
			latCart, lonCart := CartesianToWGS84(radarPos, xM, yM)

			// Calculate the distance between the two WGS84 positions
			distance := GeodeticDistance(latPolar, lonPolar, latCart, lonCart)

			// Verify they are very close (within 10 meters - expected due to ASTERIX quantization)
			if distance > 10.0 {
				t.Errorf("Coordinate transformation inconsistency for %s:", tc.name)
				t.Errorf("  Polar->WGS84:     %.8f°, %.8f°", latPolar, lonPolar)
				t.Errorf("  Cartesian->WGS84: %.8f°, %.8f°", latCart, lonCart)
				t.Errorf("  Distance between: %.2f meters", distance)

				// Debug info
				t.Logf("Input data:")
				t.Logf("  Range: %.4f NM (%.0f m), Azimuth: %.6f°", tc.rangeNM, rangeM, tc.azimuthDeg)
				t.Logf("  X: %.4f NM (%.0f m), Y: %.4f NM (%.0f m)", tc.xNM, xM, tc.yNM, yM)

				// Verify Cartesian conversion
				computedRange := math.Sqrt(xM*xM + yM*yM)
				computedAzimuth := math.Atan2(xM, yM) * 180.0 / math.Pi
				if computedAzimuth < 0 {
					computedAzimuth += 360
				}
				t.Logf("Computed from Cartesian:")
				t.Logf("  Range: %.0f m, Azimuth: %.6f°", computedRange, computedAzimuth)
				t.Logf("  Range diff: %.2f m, Azimuth diff: %.6f°",
					computedRange-rangeM, computedAzimuth-tc.azimuthDeg)
			} else {
				t.Logf("%s: Transformations match within %.2f meters ✓", tc.name, distance)
			}
		})
	}
}

func TestPolarCartesianConversion(t *testing.T) {
	testCases := []struct {
		name       string
		rangeM     float64
		azimuthDeg float64
	}{
		{"North", 1000, 0},
		{"East", 1000, 90},
		{"South", 1000, 180},
		{"West", 1000, 270},
		{"NorthEast", 1414.21, 45},
		{"SouthWest", 1414.21, 225},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert polar to Cartesian
			azimuthRad := tc.azimuthDeg * math.Pi / 180.0
			x := tc.rangeM * math.Sin(azimuthRad)
			y := tc.rangeM * math.Cos(azimuthRad)

			// Convert back to polar
			computedRange := math.Sqrt(x*x + y*y)
			computedAzimuthRad := math.Atan2(x, y)
			computedAzimuthDeg := computedAzimuthRad * 180.0 / math.Pi
			if computedAzimuthDeg < 0 {
				computedAzimuthDeg += 360
			}

			// Check accuracy
			rangeDiff := math.Abs(computedRange - tc.rangeM)
			azimuthDiff := math.Abs(computedAzimuthDeg - tc.azimuthDeg)

			if rangeDiff > 0.01 || azimuthDiff > 0.01 {
				t.Errorf("Polar<->Cartesian conversion error for %s:", tc.name)
				t.Errorf("  Input:    Range=%.2f, Azimuth=%.2f°", tc.rangeM, tc.azimuthDeg)
				t.Errorf("  Computed: Range=%.2f, Azimuth=%.2f°", computedRange, computedAzimuthDeg)
				t.Errorf("  Diffs:    Range=%.4f, Azimuth=%.4f°", rangeDiff, azimuthDiff)
			}
		})
	}
}

func BenchmarkPolarToWGS84(b *testing.B) {
	radarPos := RadarPosition{
		Latitude:  46.5197,
		Longitude: 6.6323,
		Height:    430.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PolarToWGS84(radarPos, 100000.0, 45.0)
	}
}

func BenchmarkCartesianToWGS84(b *testing.B) {
	radarPos := RadarPosition{
		Latitude:  46.5197,
		Longitude: 6.6323,
		Height:    430.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CartesianToWGS84(radarPos, 70710.0, 70710.0)
	}
}
