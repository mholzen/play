package home

import (
	"github.com/mholzen/play/controls"
	"github.com/mholzen/play/fixture"
	"github.com/mholzen/play/patterns"
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

func createObservableFixturesForChannels(universe fixture.Fixtures[fixture.Fixture], channelList []string) *fixture.ObservableFixtures {
	modelFixture := fixture.NewModelChannels(channelList...)
	addresses := universe.GetAddressesWithChannel(channelList[0])
	addressableFixtures := fixture.NewAddressableFixturesFromAddresses(modelFixture, addresses...)
	return fixture.NewIndividualObservableFixtures(addressableFixtures)
}

func GetRootSurface(universe fixture.Fixtures[fixture.Fixture], clock *controls.Clock) controls.Container {
	surface := controls.NewOrderedMap()

	// dial map
	dialFixtures := fixture.NewObservableFixtures(universe.Clone())
	dialMap := controls.NewDialMapFromChannels(controls.DefaultChannelList, controls.NewDialObservableNumeric)
	dialList := controls.NewDialListFromDialMap(dialMap)

	fixture.ConnectObservablesToFixtures(dialList.GetObservables(), dialFixtures)

	// TODO: have any changes to a dial of the list apply the entirety of the channel values to dialFixtures
	// observableDialMap := controls.NewObservableFromDialMap(dialMap)
	// fixture.ConnectObservableChannelValuesToFixtures(observableDialMap, dialFixtures)

	//
	// rainbow
	//
	rainbowFixtures := createObservableFixturesForChannels(universe, controls.ColorChannelList)
	rainbowControls := patterns.NewRainbowControls(clock)
	rainbowControls.Rainbow(&rainbowFixtures.AddressableFixtures)

	//
	// fall in
	//
	fallInFixtures := fixture.NewIndividualObservableFixtures(universe.Clone())
	fallInControls := patterns.FallInControls{Clock: clock}
	fallInControls.FallIn(&fallInFixtures.AddressableFixtures)

	// mux
	mux := controls.NewMux[fixture.FixtureValues]()
	mux.Add("dials", dialFixtures)
	mux.Add("rainbow", rainbowFixtures)
	mux.Add("fall in", fallInFixtures)

	mux.SetSource("dials")

	// link mux emitter to universe fixture
	fixture.ConnectObservableValuesToFixtures(mux, universe)

	surface.SetItem("mux", mux)
	surface.SetItem("dials", dialList)
	surface.SetItem("rainbow", rainbowControls)

	//
	// motion
	//
	motionDialMap := controls.NewDialMapFromChannels(controls.MotionChannelList, controls.NewDialObservableNumeric)
	motionDialList := controls.NewDialList(motionDialMap, controls.MotionChannelList)
	surface.SetItem("motion", motionDialList)

	// pan, ok := motionDialMap["pan"]
	// if !ok {
	// 	panic("pan not found")
	// }
	// discretizer := GetDiscretizer(clock)
	// controls.OnChange(discretizer, func(value int) {
	// 	log.Println("discretizer setting pan to", value)
	// 	pan.Set(byte(value))
	// })
	// surface.SetItem("discretizer", discretizer)

	clockDial := GetDialPattern(clock)
	surface.SetItem("dial_pattern", clockDial)

	motionFixtures := createObservableFixturesForChannels(universe, controls.MotionChannelList)

	controls.OnChange(clockDial.Mux, func(value byte) {
		motionFixtures.AddressableFixtures.SetChannelValue("tilt", value)
	})

	fixture.ConnectObservablesToFixtures(motionDialList.GetObservables(), motionFixtures)

	//
	// light
	//
	lightDialMap := controls.NewDialMapFromChannels(controls.LightChannelList, controls.NewDialObservableNumeric)
	lightDialList := controls.NewDialList(lightDialMap, controls.LightChannelList)
	surface.SetItem("light", lightDialList)

	lightFixtures := createObservableFixturesForChannels(universe, controls.LightChannelList)

	fixture.ConnectObservablesToFixtures(lightDialList.GetObservables(), lightFixtures)

	//
	// joiner
	//
	joiner := controls.NewJoiner[fixture.FixtureValues]()
	joiner.Add(motionFixtures)
	joiner.Add(lightFixtures)
	joiner.Add(mux)

	fixture.ConnectObservableValuesToFixtures(joiner, universe)

	return surface
}

func GetDialPattern(clock *controls.Clock) *patterns.DialPattern {
	return patterns.NewDialPattern(clock)
}
