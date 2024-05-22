package stages

import "github.com/mholzen/play-go/fixture"

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
