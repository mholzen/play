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
	universe := fixture.NewFixturesGeneric[fixture.Fixture]()
	return Home{
		FreedomPars: universe.AddFixtures(fixture.NewFreedomPar, 65, 81, 97, 113),
		TomeShine:   universe.AddFixtures(fixture.NewTomeshine, 1, 17, 33, 49),
		ColorStrip:  universe.AddFixtures(fixture.NewColorstripMini, 129),
		ParCans:     universe.AddFixtures(fixture.NewParCan, 140),
		Universe:    *universe,
	}
}

func GetRootSurface(universe fixture.FixturesInterface[fixture.FixtureI], clock *controls.Clock) controls.Container {
	surface := controls.NewList(3)

	// dials
	dialFixtures := fixture.NewObservableFixtures(universe.Clone())
	channelDials := fixture.NewObservableDialMapForAllChannels(dialFixtures.GetChannels(), dialFixtures)
	surface.SetItem(0, channelDials)

	// rainbow
	rainbowFixtures := fixture.NewIndividualObservableFixtures(universe.Clone())
	patterns.Rainbow(&rainbowFixtures.FixturesGeneric, clock)

	// mux
	mux := controls.NewMux[fixture.FixtureValues]()
	mux.Add("dials", dialFixtures)
	mux.Add("rainbow", rainbowFixtures)
	surface.SetItem(2, mux)

	mux.SetSource("dials")

	// link mux emitter to universe fixture
	fixture.LinkObservableToFixture(mux, &universe)

	return surface
}
