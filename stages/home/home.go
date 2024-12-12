package home

import (
	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/mholzen/play-go/patterns"
)

type Home struct {
	FreedomPars fixture.Fixtures2
	TomeShine   fixture.Fixtures2
	ColorStrip  fixture.Fixtures2
	ParCans     fixture.Fixtures2
	Universe    fixture.Fixtures2
}

func GetHome() Home {
	universe := fixture.NewFixturesGeneric[fixture.Fixture]()
	return Home{
		// FreedomPars: universe.AddFixtures(fixture.NewFreedomPar, 65, 81, 97, 113),
		TomeShine:  universe.AddFixtures(fixture.NewTomeshine, 1, 17, 33, 49),
		ColorStrip: universe.AddFixtures(fixture.NewColorstripMini, 129),
		ParCans:    universe.AddFixtures(fixture.NewParCan, 140),
		Universe:   *universe,
	}
}

func GetRootSurface(universe fixture.FixturesInterface[fixture.FixtureI], clock *controls.Clock) controls.Container {
	surface := controls.NewList(3)

	dialFixtures := fixture.NewObservableFixtures2(universe)

	channelDials := fixture.NewObservableDialMapForAllChannels(dialFixtures.GetChannels(), dialFixtures)
	surface.SetItem(0, channelDials)

	mux := controls.NewMux[fixture.FixtureValues]()
	mux.Add("dials", dialFixtures)

	rainbowFixtures := fixture.NewIndividualObservableFixtures2(universe)
	patterns.Rainbow(&rainbowFixtures.FixturesGeneric, clock)
	mux.Add("rainbow", rainbowFixtures)

	mux.SetSource("dials")
	surface.SetItem(2, mux)

	// link mux emitter to universe fixture
	// fixture.LinkEmitterToFixture(&mux, universe)
	fixture.LinkObservableToFixture(mux, &universe)

	return surface
}
