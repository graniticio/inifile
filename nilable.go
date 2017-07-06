package inifile


// Create a new nilableString with the supplied value.
func newNilableString(v string) *nilableString {
	ns := new(nilableString)
	ns.Set(v)

	return ns
}

// A string where it can be determined if "" is an explicitly set value, or just the default zero value
type nilableString struct {
	val string
	set bool
}

// Set sets the contained value to the supplied value and makes IsSet true even if the supplied value is the empty
// string.
func (ns *nilableString) Set(v string) {
	ns.val = v
	ns.set = true
}

// The currently stored value (whether or not it has been explicitly set).
func (ns *nilableString) String() string {
	return ns.val
}

// See Nilable.IsSet
func (ns *nilableString) IsSet() bool {
	return ns.set
}