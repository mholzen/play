# Play

A DMX lighting control system written in Go, designed for stage lighting and live performances. Play provides real-time control over DMX fixtures through multiple interfaces including web API, terminal controls, and programmatic patterns.

## Features

### Lighting Control
- **DMX Interface**: Direct hardware control via serial/USB DMX adapters
- **Multiple Fixture Support**: Freedom Par, TomeShine moving heads, LED strips, Par cans
- **Real-time Rendering**: 40ms refresh rate for smooth lighting transitions

### Beat-Synchronized Patterns
- **Master Clock**: BPM-based timing system with beats, bars, and phrases
- **Pattern Library**: Rainbow effects, color chases, movement sequences, and transitions
- **Easing Functions**: Smooth transitions with customizable curves
- **Live Tapping**: Tempo detection and synchronization

### Multiple Control Interfaces
- **Web API**: RESTful HTTP interface for remote control
- **Terminal UI**: Live keyboard controls with real-time feedback
- **Programmatic**: Direct Go API for custom lighting shows

### Advanced Effects
- **Color Management**: Full RGBAW+UV color mixing with named color presets
- **Movement Patterns**: Automated pan/tilt sequences for moving head fixtures
- **Transition System**: Smooth interpolation between lighting states
- **Pattern Multiplexing**: Layer and combine multiple effect sources

## Architecture

### Core Components

- **`controls/`**: Control abstractions (dials, clocks, triggers, containers)
- **`fixture/`**: DMX fixture models and hardware abstraction
- **`patterns/`**: Lighting effect generators and sequencers
- **`stages/`**: Venue-specific fixture configurations
- **`cmd/`**: Executable applications

### Fixture Models

Currently supports:
- **Freedom Par**: RGBAW+UV LED pars with dimming and strobe
- **TomeShine**: Moving head with pan/tilt, color mixing, and effects
- **Colorstrip Mini**: RGB LED strips with chase patterns
- **Par Can**: Traditional tungsten fixtures

## Installation

### Prerequisites

- Go 1.23 or later
- DMX hardware interface (USB or serial)
- ALSA libraries (Linux) or CoreAudio (macOS)

### Build from Source

```bash
git clone https://github.com/mholzen/play.git
cd play
go mod download
make build
```

### Docker

```bash
make build-docker
make run-docker
```

## Usage

### Web Server Mode

Start the web API server for remote control:

```bash
./bin/server
```

The server runs on port 1300 and provides:
- Static web interface at `http://localhost:1300`
- REST API at `http://localhost:1300/api/v2/root`

### Terminal Control Mode

Launch interactive terminal controls:

```bash
make live
# or directly:
(cd cmd/live && go run live.go)
```

**Terminal Controls:**
- `t` - Tap tempo
- `s` - Sync clock to tapped tempo
- `r` - Reset clock
- `l` - Toggle beat logging
- `m<bpm>` - Set specific BPM (e.g., m120)
- `+/-` - BPM adjustment
- `Ctrl+C` - Quit

### Clock-Only Mode

Run just the master clock with terminal controls:

```bash
make clock
# or:
(cd cmd/clock && go run clock.go)
```

## API Reference

### REST Endpoints

The web API provides hierarchical access to all controls:

- `GET /api/v2/root` - List all available controls
- `GET /api/v2/root/{path}` - Get control/container value
- `POST /api/v2/root/{path}` - Set control value

**Example API Calls:**

```bash
# List all controls
curl http://localhost:1300/api/v2/root/

# Get dial values
curl http://localhost:1300/api/v2/root/dials

# Set a color control
curl -X POST http://localhost:1300/api/v2/root/dials/r -d "255"

# Control pattern multiplexer
curl -X POST http://localhost:1300/api/v2/root/mux/source -d "rainbow"
```

## Configuration

### DMX Hardware Setup

The system auto-detects DMX interfaces. Supported devices:
- ENTTEC DMX USB Pro
- Generic FTDI-based DMX adapters
- Any serial DMX interface

## License

Open source project. Check repository for license details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

---

For questions or support, see the project repository or create an issue.
