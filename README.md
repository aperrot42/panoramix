# RecTools

Radar data processing tool for aviation surveillance systems. Decodes .if container files containing ASTERIX messages.

## Quick Start

**Build:**
```bash
go build .
```

**Run:**
```bash
./rectools --filename data.if
./rectools --filename data.if --json
./rectools --filename data.if --limit 100
./rectools --filename data.if --limit 100 --position
```

## What it does

- Reads .if radar recording files
- Decodes ASTERIX CAT 034 (system status) and CAT 048 (radar targets)
- Extracts aircraft data from Mode S BDS registers
- Outputs text or JSON format

### Additional filters 

#### Position

Position filter (--position) adds a computed wgs84 position extracted from CAT048 data using the best available positionning data :

```json
    "computed": {
        "position": {
            "wgs84_position": {
                "latitude_deg": 47.78745995705229,
                "longitude_deg": 7.3358434913021755,
                "position_source": "POLAR",
                "altitude_ft": 10000,
                "altitude_source": "FLIGHT_LEVEL"
            },
            "timestamp": "2025-05-15T12:00:00.000Z",
            "radar_sic": 33,
            "radar_sac": 40
        }
    }
```



## Documentation

For detailed information:
- ASTERIX specifications: [EUROCONTROL](https://www.eurocontrol.int/asterix)
- Mode S documentation: see ED-73F - MOPS for Secondary Surveillance Radar Mode S Transponders
- Configuration: `radar_config.sample.yaml`