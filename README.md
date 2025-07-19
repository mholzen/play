# Play

Play is a DMX lighting control system written in Go, designed for live stage
performances or for ambient lights in a home setup.

The goal of this project is to "play" lighting as one does an instrument,
through both traditional controls (such buttons or dials), but also by writing
and executing code during a live performance.

It provides a terminal application responding to keyboard events, as well as a
server with a REST API, designed to work seamlessly with a web front end
available at [play-editor](https://github.com/mholzen/play-editor).

It also serves as an opportunity for me to continue to refine my software
engineering skills.


## Getting Started

### Prerequisites

- Go 1.23 or later
- a DMX hardware interface (USB or serial)
- one or more DMX compatible light fixtures

### DMX controllers

Play uses [akualab/dmx](github.com/akualab/dmx) which currently supports:
- Enttec DMX USB Pro
- DMX King UltraDMX Micro

### Run from source

```bash
brew install go
go install github.com/githubnemo/CompileDaemon@latest
make run
```

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

## Usage

### Configuration

The DMX universe is currently configured via code in `stages/home.go`

### Web Server

Start the web API server for remote control:

```bash
make build
./bin/server
```

The server runs on port 1300 and provides the backend for [play-editor](https://github.com/mholzen/play-editor)

### Terminal App

Launch interactive terminal controls:

```bash
make live
```

Press `?` for a list of keyboard controls:

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

This can be used to tap beats and detect tempo.