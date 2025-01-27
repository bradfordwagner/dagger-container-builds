package flavors

type Flavor string

const (
	FlavorMirror Flavor = "mirror"
	FlavorCustom Flavor = "custom"
)

func (f Flavor) String() string {
	return string(f)
}

func FromString(s string) Flavor {
	return Flavor(s)
}
