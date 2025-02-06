package home

import (
	"github.com/mholzen/play-go/controls"
	"github.com/mholzen/play-go/fixture"
	"github.com/mholzen/play-go/patterns"
)

type Home struct {
	FreedomPars fixture.AddressableChannelFixtures
	TomeShine   fixture.AddressableChannelFixtures
	ColorStrip  fixture.AddressableChannelFixtures
	ParCans     fixture.AddressableChannelFixtures
	Universe    fixture.AddressableChannelFixtures
}

func GetHome() Home {
	universe := fixture.NewAddressableFixtures[fixture.ChannelFixture]()
	return Home{
		FreedomPars: universe.AddFixtures(fixture.NewFreedomPar, 65, 81, 97, 113),
		TomeShine:   universe.AddFixtures(fixture.NewTomeshine, 1, 17, 33, 49),
		ColorStrip:  universe.AddFixtures(fixture.NewColorstripMini, 129),
		ParCans:     universe.AddFixtures(fixture.NewParCan, 140),
		Universe:    *universe,
	}
}

func GetRootSurface(universe fixture.Fixtures[fixture.Fixture], clock *controls.Clock) controls.Container {
	surface := controls.NewList(5)

	// dial map
	dialFixtures := fixture.NewObservableFixtures(universe.Clone())
	// dialFixtures := fixture.NewIndividualObservableFixtures(universe.Clone())

	dialMap := fixture.NewObservableDialMapForAllChannels(dialFixtures)
	dialList := controls.NewDialList2(dialMap)

	// rainbow
	rainbowFixtures := fixture.NewIndividualObservableFixtures(universe.Clone())
	rainbowControls := patterns.NewRainbowControls(clock)
	rainbowControls.Rainbow(&rainbowFixtures.AddressableFixtures)

	// mux
	mux := controls.NewMux[fixture.FixtureValues]()
	mux.Add("dials", dialFixtures)
	mux.Add("rainbow", rainbowFixtures)

	mux.SetSource("dials")

	// link mux emitter to universe fixture
	fixture.LinkObservableToFixture(mux, &universe)

	surface.SetItem(0, mux)
	surface.SetItem(1, dialList)
	surface.SetItem(2, dialMap)
	surface.SetItem(3, rainbowControls)
	// surface.SetItem(4, rainbowControls.GetContainer())

	return surface
}
