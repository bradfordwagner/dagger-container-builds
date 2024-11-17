package flavors

type Flavor string

const (
	FlavorMirror Flavor = "mirror"
)

func (f Flavor) String() string {
	return string(f)
}

func FromString(s string) Flavor {
	return Flavor(s)
}
