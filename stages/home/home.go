package home

import (
	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/mholzen/play-go/patterns"
)

type Home struct {
	FreedomPars fixture.Fixtures
	TomeShine   fixture.Fixtures
	ColorStrip  fixture.Fixtures
	ParCans     fixture.Fixtures
	Universe    fixture.Fixtures
}

func GetHome() Home {
	universe := fixture.NewFixtures()
	return Home{
		FreedomPars: universe.AddFixtures(fixture.NewFreedomPar, 65, 81, 97, 113),
		TomeShine:   universe.AddFixtures(fixture.NewTomeshine, 1, 17, 33, 49),
		ColorStrip:  universe.AddFixtures(fixture.NewColorstripMini, 129),
		ParCans:     universe.AddFixtures(fixture.NewParCan, 140),
		Universe:    universe,
	}
}

func GetRootSurface(universe fixture.Fixtures, clock *controls.Clock) controls.Container {
	surface := controls.NewList(3)

	dialFixtures := fixture.NewFixturesFromFixtures(universe)
	dialFixtureEmitter := fixture.NewFixtureEmitter(dialFixtures)
	// dialFixtureEmitter := controls.NewEmitterValue(dialFixtures)

	channelDials := fixture.NewDialMapAllFixtures(dialFixtures)
	surface.SetItem(0, channelDials)

	rainbowFixtures := fixture.NewFixturesFromFixtures(universe)
	rainbowFixtureEmitter := fixture.NewFixtureEmitter(rainbowFixtures)
	patterns.Rainbow(rainbowFixtures, clock)

	mux := controls.NewMux[fixture.FixtureValues]()
	mux.Add("dials", &dialFixtureEmitter)
	mux.Add("rainbow", &rainbowFixtureEmitter)
	mux.SetSource("rainbow")
	surface.SetItem(2, &mux)

	// link mux emitter to universe fixture
	fixture.LinkEmitterToFixture(&mux, universe)

	return surface
}
