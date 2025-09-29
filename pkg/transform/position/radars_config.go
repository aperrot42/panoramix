package position

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// RadarConfig represents the radar configuration file structure
type RadarConfig struct {
	Radars map[string]RadarConfigEntry `yaml:"radars"`
}

// RadarConfigEntry represents a single radar configuration
type RadarConfigEntry struct {
	Name     string              `yaml:"name"`
	Position RadarPositionConfig `yaml:"position"`
}

// RadarPositionConfig represents radar position in config file
type RadarPositionConfig struct {
	LatitudeDeg  float64 `yaml:"latitude_deg"`
	LongitudeDeg float64 `yaml:"longitude_deg"`
	AltitudeM    float64 `yaml:"altitude_m"`
}

// LoadRadarConfig loads radar configuration from YAML file
func LoadRadarConfig(configPath string) (*RadarRegistry, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config RadarConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	registry := NewRadarRegistry()

	for sacSicKey, radarConfig := range config.Radars {
		// Parse SAC/SIC from key format "SAC/SIC"
		var sac, sic uint8
		_, err := fmt.Sscanf(sacSicKey, "%d/%d", &sac, &sic)
		if err != nil {
			return nil, fmt.Errorf("invalid SAC/SIC format in key '%s': %w", sacSicKey, err)
		}

		position := RadarPosition{
			Latitude:  radarConfig.Position.LatitudeDeg,
			Longitude: radarConfig.Position.LongitudeDeg,
			Height:    radarConfig.Position.AltitudeM,
			SIC:       sic,
			SAC:       sac,
		}

		registry.AddRadarPosition(position)
	}

	return registry, nil
}
