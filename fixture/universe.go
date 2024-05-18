package fixture

func GetUniverse() Fixtures {
	universe := NewFixtures()

	freedomPars := NewFixtures()
	for _, address := range []int{65, 81, 97, 113} {
		freedomPars.AddFixture(NewFreedomPar(), address)
	}
	universe.AddFixtureList(freedomPars)

	tomshine := NewFixtures()
	for _, address := range []int{1, 17, 33, 49} {
		tomshine.AddFixture(NewTomeshine(), address)
	}
	universe.AddFixtureList(tomshine)

	colorstripMinis := NewFixtures()
	for _, address := range []int{129} {
		colorstripMinis.AddFixture(NewColorstripMini(), address)
	}
	universe.AddFixtureList(colorstripMinis)

	parCans := NewFixtures()
	for _, address := range []int{140} {
		parCans.AddFixture(NewParCan(), address)
	}
	universe.AddFixtureList(parCans)

	return universe
}
