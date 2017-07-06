package inifile

// IniSection provides access to an IniConfig object within the context of a single section.
//
// Call the Section(sectionName) function on your IniConfig to obtain an IniSection
type IniSection struct {
	name string
	ic *IniConfig
}

//Name returns the name of this section
func (is *IniSection) Name() string{
	return is.name
}

//See IniConfig.PropertyExists
func (is *IniSection) PropertyExists(propertyName string) bool {
	return is.ic.PropertyExists(is.name, propertyName)
}

//See IniConfig.Value
func (is *IniSection) Value(propertyName string) (string, error) {
	return is.ic.Value(is.name, propertyName)

}

//See IniConfig.ValueAsFloat64
func (is *IniSection) ValueAsFloat64(propertyName string) (float64, error) {
	return is.ic.ValueAsFloat64(is.name, propertyName)
}

//See IniConfig.ValueAsInt64
func (is *IniSection) ValueAsInt64(propertyName string) (int64, error) {
	return is.ic.ValueAsInt64(is.name, propertyName)
}

//See IniConfig.ValueAsUint64
func (is *IniSection) ValueAsUint64(propertyName string) (uint64, error) {
	return is.ic.ValueAsUint64(is.name, propertyName)
}

//See IniConfig.ValueAsBool
func (is *IniSection) ValueAsBool(propertyName string) (bool, error) {
	return is.ic.ValueAsBool(is.name, propertyName)
}

//See IniConfig.Add
func (is *IniSection) Add(propertyName string, value string) {
	is.ic.Add(is.name, propertyName, value)
}