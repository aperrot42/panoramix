package position

import (
	"testing"
	"time"
)

func TestPositionExtractorIntegration(t *testing.T) {
	extractor := NewPositionExtractor()

	extractor.AddRadarPosition(RadarPosition{
		Latitude:  46.5197,
		Longitude: 6.6323,
		Height:    430.0,
		SIC:       33,
		SAC:       40,
	})

	fl := float64(100)
	plot := RawPlot{
		SIC:   33,
		SAC:   40,
		Polar: &PolarCoord{RhoNM: 96.1640625, ThetaDeg: 31.3330078125},
		FL:    &fl,
	}

	obs, err := extractor.Extract(plot, time.Now())
	if err != nil {
		t.Fatalf("Failed to extract position: %v", err)
	}
	if obs == nil {
		t.Fatal("Expected position, got nil")
	}

	if obs.RadarSIC != 33 || obs.RadarSAC != 40 {
		t.Errorf("Expected SIC=33 SAC=40, got SIC=%d SAC=%d", obs.RadarSIC, obs.RadarSAC)
	}

	if obs.WGS84Position.PostionSource != "POLAR" {
		t.Errorf("Expected source 'POLAR', got '%s'", obs.WGS84Position.PostionSource)
	}

	if obs.WGS84Position.Latitude_deg < 45 || obs.WGS84Position.Latitude_deg > 50 {
		t.Errorf("Expected latitude between 45-50, got %.6f", obs.WGS84Position.Latitude_deg)
	}

	if obs.WGS84Position.Longitude_deg < 5 || obs.WGS84Position.Longitude_deg > 10 {
		t.Errorf("Expected longitude between 5-10, got %.6f", obs.WGS84Position.Longitude_deg)
	}

	// FL 100 * 100 = 10000 ft
	if obs.WGS84Position.AltitudeFt != 10000 {
		t.Errorf("Expected altitude 10000 ft, got %.1f ft", obs.WGS84Position.AltitudeFt)
	}

	t.Logf("WGS84: %.6f, %.6f (%.0f ft)",
		obs.WGS84Position.Latitude_deg,
		obs.WGS84Position.Longitude_deg,
		obs.WGS84Position.AltitudeFt)
}

func TestRadarRegistry(t *testing.T) {
	registry := NewRadarRegistry()

	pos1 := RadarPosition{Latitude: 46.5, Longitude: 6.6, Height: 400, SIC: 33, SAC: 40}
	pos2 := RadarPosition{Latitude: 47.4, Longitude: 8.5, Height: 500, SIC: 34, SAC: 41}

	registry.AddRadarPosition(pos1)
	registry.AddRadarPosition(pos2)

	retrieved, found := registry.GetRadarPosition(33, 40)
	if !found {
		t.Error("Should find radar position for SIC=33 SAC=40")
	}
	if retrieved.Latitude != 46.5 {
		t.Errorf("Expected latitude 46.5, got %.1f", retrieved.Latitude)
	}

	_, found = registry.GetRadarPosition(99, 99)
	if found {
		t.Error("Should not find radar position for SIC=99 SAC=99")
	}
}
