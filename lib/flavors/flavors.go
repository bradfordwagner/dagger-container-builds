package flavors

type Flavor string

const (
	FlavorMirror Flavor = "mirror"
)

func (f Flavor) String() string {
	return string(f)
}
