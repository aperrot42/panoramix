# Panoramix

ASTERIX radar message decoder library and CLI for aviation surveillance data.

## Quick Start

```bash
go build -o panoramix .
```

```bash
# Decode raw ASTERIX binary data
panoramix --filename data.ast
panoramix --filename data.ast --json
panoramix --filename data.ast --json --limit 100

# Add computed WGS84 position (requires radar_config.yaml)
panoramix --filename data.ast --json --position --radar-config radar_config.yaml

# Add FSPEC field list
panoramix --filename data.ast --json --fspec-fields
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--filename` | *(required)* | Input file containing raw ASTERIX data |
| `--json` | `false` | Output JSON instead of text |
| `--limit` | `0` | Max messages to decode (0 = unlimited) |
| `--position` | `false` | Compute WGS84 position for CAT 048 |
| `--fspec-fields` | `false` | List available FSPEC fields per message |
| `--radar-config` | `radar_config.yaml` | Radar station positions config file |

## Supported ASTERIX Categories

- **CAT 034** — Monoradar Service Messages (system status)
- **CAT 048** — Monoradar Target Reports (radar plots/tracks)

## BDS Register Decoding

Mode S BDS registers from I048/250 are decoded into typed structs:

| BDS | Name |
|-----|------|
| 1,0 | Data Link Capability Report |
| 1,7 | Common Usage GICB Capability Report |
| 2,0 | Aircraft Identification |
| 3,0 | ACAS Active Resolution Advisory |
| 4,0 | Selected Vertical Intention |
| 4,4 | Meteorological Routine Air Report |
| 5,0 | Track and Turn Report |
| 6,0 | Heading and Speed Report |

## Position Computation

The `--position` flag computes WGS84 coordinates from CAT 048 data:

- **Polar** (I048/040): range + azimuth → Vincenty direct formula
- **Cartesian** (I048/042): X/Y offset → WGS84
- **Altitude**: 3D radar height (I048/110) preferred, then flight level (I048/090)

Requires a radar config file with station positions:

```yaml
radars:
  "40/33":  # SAC/SIC
    name: "DOLS"
    position:
      latitude_deg: 46.42567939
      longitude_deg: 6.09995219
      altitude_m: 1734.13
```

## Package Structure

```
pkg/
  asterix/     # ASTERIX wire format decoder (CAT 034, CAT 048)
  modes/bds/   # Mode S BDS register decoder (typed structs)
  transform/
    position/  # Polar/Cartesian → WGS84 coordinate transform
    fspec/     # FSPEC field enumeration
```

## References

- [EUROCONTROL ASTERIX specifications](https://www.eurocontrol.int/asterix)
- ED-73F — MOPS for Secondary Surveillance Radar Mode S Transponders