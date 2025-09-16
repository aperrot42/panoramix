# RecTools

Radar data processing tool for aviation surveillance systems. Decodes .if container files containing ASTERIX messages.

## Quick Start

**Build:**
```bash
go build .
```

**Run:**
```bash
./rectools -filename data.if
./rectools -filename data.if -json
./rectools -filename data.if -limit 100
```

## What it does

- Reads .if radar recording files
- Decodes ASTERIX CAT 034 (system status) and CAT 048 (radar targets)
- Extracts aircraft data from Mode S BDS registers
- Outputs text or JSON format

## Documentation

For detailed information:
- ASTERIX specifications: [EUROCONTROL](https://www.eurocontrol.int/asterix)
- Mode S documentation: see ED-73F - MOPS for Secondary Surveillance Radar Mode S Transponders
- Configuration: `radar_config.sample.yaml`