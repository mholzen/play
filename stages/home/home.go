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
	surface := controls.NewList(4)

	// dial map
	dialFixtures := fixture.NewObservableFixtures(universe.Clone())
	dialMap := controls.ChannelsToDialMap(controls.DefaultChannelList, controls.NewDialObservableNumeric)
	dialList := controls.NewDialListFromDialMap(dialMap)

	fixture.ConnectObservablesToFixtures(dialList.GetObservables(), dialFixtures)

	// TODO: have any changes to a dial of the list apply the entirety of the channel values to dialFixtures
	// observableDialMap := controls.NewObservableFromDialMap(dialMap)
	// fixture.ConnectObservableChannelValuesToFixtures(observableDialMap, dialFixtures)

	// rainbow
	rainbowFixtures := fixture.NewIndividualObservableFixtures(universe.Clone())
	// A change to ANY individual observable fixtures SHUOLD notify all rainbow fixtures observers

	rainbowControls := patterns.NewRainbowControls(clock)
	rainbowControls.Rainbow(&rainbowFixtures.AddressableFixtures)

	// fall in
	fallInFixtures := fixture.NewIndividualObservableFixtures(universe.Clone())
	fallInControls := patterns.FallInControls{Clock: clock}
	fallInControls.FallIn(&fallInFixtures.AddressableFixtures)

	// mux
	mux := controls.NewMux[fixture.FixtureValues]()
	mux.Add("dials", dialFixtures)
	mux.Add("rainbow", rainbowFixtures)
	mux.Add("fall in", fallInFixtures)

	mux.SetSource("rainbow")

	// link mux emitter to universe fixture
	fixture.ConnectObservableValuesToFixtures(mux, universe)

	surface.SetItem(0, mux)
	surface.SetItem(1, dialList)
	surface.SetItem(2, rainbowControls)

	return surface
}
