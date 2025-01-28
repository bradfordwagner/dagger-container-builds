package flavors

type Flavor string

const (
	FlavorAnsibleRole     Flavor = "ansible_role"
	FlavorAnsiblePlaybook Flavor = "ansible_playbook"
	FlavorMirror          Flavor = "mirror"
	FlavorCustom          Flavor = "custom"
)

func (f Flavor) String() string {
	return string(f)
}

func FromString(s string) Flavor {
	return Flavor(s)
}
