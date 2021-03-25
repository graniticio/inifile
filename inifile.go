// Copyright 2017 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

/*
Package inifile provides a Go struct that can parse an INI-style file and make the configuration in that file
available via a series of type-safe data accessors. Parsing and data-access behaviour can be configured to support most
INI file variants.

Parsing

To parse an INI file and obtain an instance of IniConfig to access your configuration, call one of:
	inifile.NewIniConfigFromPath(string)
	inifile.NewIniConfigFromFile(*os.File)
	inifile.NewIniConfigFromPathWithOptions(string, *IniOptions)
	inifile.NewIniConfigFromFileWithOptions(*os.File, *IniOptions)


For example:

	import "github.com/graniticio/inifile"

	func main() {

		ic, err := inifile.NewIniConfigFromPath("/path/to/file.ini")
	}


Accessing properties

The key/value items in an INI file are referred to as properties. Each property is part of a section (properties in an INI file
defined before any [section] delimiters are considered to be part of the 'global' section). Values are stored as strings.

You can retrieve the value associated with a property by using one of:
	Value(sectionName, propertyName string)
	ValueAsFloat64(sectionName, propertyName string)
	ValueAsInt64(sectionName, propertyName string)
	ValueAsUint64(sectionName, propertyName string)
	ValueAsBool(sectionName, propertyName string)

These methods will return an error if the requested section or name does not exist or if the value associated with the
requested property could not be converted to the request data type.

To check that a section of property exists before you call one of these functions use:
	SectionExists(sectionName string)
	PropertyExists(sectionName, propertyName string)

Methods exist to return the zero value for a type instead of an error if the section/property didn't exist or if there
was a problem converting the value to the requested type:
	ValueOrZero(sectionName, propertyName string)
	ValueOrZeroAsFloat64(sectionName, propertyName string)
	ValueOrZeroAsInt64(sectionName, propertyName string)
	ValueOrZeroAsUint64(sectionName, propertyName string)
	ValueOrZeroAsBool(sectionName, propertyName string)

Accessing properties in the global section

Use the constant inifile.GLOBAL_SECTION as the sectionName when calling any of the above functions to work with properties that are not
attached to a named section

Accessing properties via an IniSection

If your code needs multiple property values from the same section:

	a, err := ic.Value("section1", "a")
	b, err := ic.Value("section1", "b")
	c, err := ic.Value("section1", "c")

it is recommended that you use an IniSection object, which offers the same functions as IniConfig but is bound to a single section:

	is, err := ic.Section("section1")
	a, err := is.Value("a")
	b, err := is.Value("b")
	c, err := is.Value("c")


Adding new properties

Properties can be added to an IniConfig at runtime by calling:
	Add(section, propertyName string, value string)


Customising parsing and configuration access

As INI files are not governed by an agreed standard, there are a number of variations in the structure and features
found in real-world INI files. To accommodate these variations, you can modify the IniOptions used when creating
an IniConfig object.

The recommended practice is to use the DefaultIniOptions() function to create a baseline set of options
and then only modify those you need to change (see the documentation for DefaultIniOptions below to see what the default values are). For example:

	opts := inifile.DefaultIniOptions()
	opts.CommentStart = "#"

	ic, err := NewIniConfigFromPathWithOptions("/path/to/file.ini", opts)

will parse a file using # instead of ; to identify comment lines.

Case-sensitivity for section and property names

By default look-ups of sections and properties are case sensitive - Value("mysection", "myproperty") would not match a property called myProperty in a section called [MYPROPERTY].

This is not the behaviour of many Windows implementations of INI file libraries, so you might find files that mix and match cases. To support this behaviour set:
	CaseSensitive = false
in your IniOptions.

Comment lines

Most INI files use the semi-colon symbol at the start of a line to indicate that the line is a comment. There are notable exceptions, including MySQL INI files, that use a different character, often #. To support this set:
	CommentStart = "#"
in your IniOptions.

Property names and values with spaces

For readability, some INI files align property names and values like:
	basedir	=	/usr
	datadir	=	/var/lib/mysql

By default the whitespace either side of the property name and value is discarded and not considered part of the property name or value. This behaviour can be disabled by setting:
	TrimProperties = false
in your IniOptions.

Blank lines

For readability, most INI files use blank lines to break up sections and properties. To disallow this and return an error if blank lines are encountered set:
	TolerateBlankLines = false
in your IniOptions.

Allow properties outside of a section

Some INI files have properties defined before the first [Section]. These properties are considered to be in the 'global' section and by default this behaviour is supported. To forbid
properties in the global section, set:
	AllowGlobalSection = false
in your IniOptions.

Use the inifile.GLOBAL_SECTION constant as the section name to access properties in the global section.

Unset properties

Some INI files allow a property to be defined without a value like
	[Reference]
	randomNumberSeed=

This can be problematic if, when set, that property is expected to hold a value of a specific data-type (a float in example above). Calling ValueAsFloat64 on the above would return an error.
By default, IniConfig ignores any unset properties, so in the above example calling
	PropertyExists("Reference", "randomNumberSeed")
would return false

To reverse this behaviour and have any unset properties stored with "" as their value, set:
	DiscardPropertiesWithNoValue = false
in your IniOptions.

Boolean values

Go's strconv.ParseBool is extremely permissive about the values it considers to represent true or false (see https://golang.org/pkg/strconv/#ParseBool) and by default
calling
	ValueAsBool(sectionName, propertyName string)
Will use those rules when trying to interpret a string value as a bool.

Many applications require bool representations to be more precise or use counter-intuitive values for true/false. This behaviour can supported by setting
	UseGoBoolRules = false
in your IniOptions and then providing values for
	StrictBoolTrue
	StrictBoolFalse

e.g:
	StrictBoolTrue = "0"
	StrictBoolFalse = "-1"

To support case-insensitve matching of the values in StrictBoolTrue and StrictBoolFalse set:
	StrictBoolCaseSensitive = false


Unparseable lines

By default, an INI file will not parse correctly if a line is encountered in the file that cannot be interpreted as a
section, comment, property or a blank line. To ignore unparseable lines, set:
	IgnoreUnparseable = true
in your IniOptions.


Inline comments

Some INI files allow a comment on the same line as a property or section:
	[database]	;Connection information
	host=localhost ;Default to localhost

By default this behaviour is not supported and the value of [database].host will be "host ;Default to localhost". To allow
inline comments set:
	AllowInlineComments true

When inline comments are allowed, you must escape any instances of the comment start character in your section names, property names and
property values. The default escape prefix is \ but you can change this by setting:
	CommentEscapePrefix
to whatever character you want.

For example, with CommentStart = ";" and CommentEscapePrefix="\" your INI file might contain

	[messages]
	message=Service is too busy\;load is high ;Message shown when site is down

The value returned by Value("messages", "message") would be "Service is too busy;load is high"

Quoted values

Some INI files surround their values with quotes like:
	inbox="/folder/with spaces"
	outbox='/another folder/sub-folder'

By default these enclosing quotes are preserved and are included when you retrieve the property's value. To change this behaviour, set:
	StripEnclosingQuotes = true
in your IniOptions.

The runes (characters) that will be recognised as enclosing quotes are defined in
	EnclosingQuoteSymbols
and default to the single (') and double (") quote symbols.


*/
package inifile

import (
	"os"
	"bufio"
	"regexp"
	"strings"
	"errors"
	"fmt"
	"strconv"
)

type sectionPropertyMap map[string]map[string]*nilableString

// If your INI file contains properties outside of a named section, use this constant as the 'section name' when
// looking up property values. For example:
//
//		val, err := ic.Value(inifile.GLOBAL_SECTION, "propertyName")
const GLOBAL_SECTION = ""


// DefaultIniOptions returns an IniOptions object populated with default values useful for working with most INI files.
//
// Default values are:
//		CaseSensitive 					true
//		CommentStart 					";"
//		TrimProperties 					true
//		TolerateBlankLines				true
//		AllowGlobalSection				true
//		DiscardPropertiesWithNoValue	true
//		UseGoBoolRules					true
//		StrictBoolTrue					""
//		StrictBoolFalse					""
//		StrictBoolCaseSensitive			true
//		IgnoreUnparseable				false
// 		AllowInlineComments				false
//		CommentEscapePrefix				"\"
//		StripEnclosingQuotes			false
//		EnclosingQuoteSymbols			[]rune{'\'','"'}
//      UseColonAssignment              false
//
func DefaultIniOptions() *IniOptions {
	io := new(IniOptions)

	io.CaseSensitive = true
	io.CommentStart = ";"
	io.TrimProperties = true
	io.TolerateBlankLines = true
	io.AllowGlobalSection = true
	io.DiscardPropertiesWithNoValue = true
	io.UseGoBoolRules = true
	io.StrictBoolCaseSensitive = true
	io.IgnoreUnparseable = false
	io.AllowInlineComments = false
	io.CommentEscapePrefix = "\\"
	io.StripEnclosingQuotes = false
	io.EnclosingQuoteSymbols = []rune{'\'','"'}
    io.UseColonAssignment = false

	return io
}

//IniOptions allows you to alter the behaviour of parsing and subsequent access to parsed configuration.
type IniOptions struct {
	//Set to true if section and variable names should be treated as-case sensitive.
	CaseSensitive bool

	//The string, which if found at the start of a line, indicates a comment line
	CommentStart string

	//Removes leading and trailing and spaces from property names and values
	TrimProperties bool

	//Ignore blank lines if true; if false return an error when parsing if a blank line found
	TolerateBlankLines bool

	//Allow properties to be defined before the first section is encountered
	AllowGlobalSection bool

	//How to handle properties with a name but no value
	DiscardPropertiesWithNoValue bool

	//Use Go's standard string-to-bool rules https://golang.org/pkg/strconv/#ParseBool
	//If set to false, StrictBoolTrue and StrictBoolFalse must be set.
	UseGoBoolRules bool

	//A string which must be matched exactly to consider a property value a 'true' boolean
	//Only used if UseGoBoolRules = false
	StrictBoolTrue string

	//A string which must be matched exactly to consider a property value a 'false' boolean
	//Only used if UseGoBoolRules = false
	StrictBoolFalse string

	//Use case sensitive matching when in StrictBool mode.
	//Only used if UseGoBoolRules = false
	StrictBoolCaseSensitive bool

	//Ignore lines that cannot be parsed as a section, property, comment or blank
	IgnoreUnparseable bool

	//Allow comments on the same line as a section or property
	AllowInlineComments bool

	//The string which, when found before the comment symbol, escapes that symbol
	//Only used when AllowInlineComments = true
	CommentEscapePrefix string

	//Remove leading and trailing quotes from property values before storing them
	StripEnclosingQuotes bool

	//The symbols that are used as enclosing quotes
	EnclosingQuoteSymbols []rune

    //Assignment uses colon not equals
    UseColonAssignment bool
}

// NewIniConfigFromPath loads the INI file at the supplied path into a new IniConfig object.
// The IniOptions used will be those returned from DefaultIniOptions()
//
// An error will be returned if there was a problem accessing the specified file or parsing it as an INI file.
func NewIniConfigFromPath(path string) (*IniConfig, error) {
	return NewIniConfigFromPathWithOptions(path, DefaultIniOptions())
}

// NewIniConfigFromPathWithOptions loads the INI file at the supplied path into a new IniConfig object using the supplied
// options.
//
// An error will be returned if there was a problem accessing the specified file or parsing it as an INI file.
func NewIniConfigFromPathWithOptions(path string, options *IniOptions) (*IniConfig, error) {

	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		defer f.Close()

		return NewIniConfigFromFileWithOptions(f, options)
	}
}

// NewIniConfigFromFile loads the INI file behind the supplied file handle into a new IniConfig object
// using the default options returned from DefaultIniOptions(). Caller is responsible for closing the supplied file.
//
// An error will be returned if there was a problem using the supplied file or parsing it as an INI file.
func NewIniConfigFromFile(file *os.File) (*IniConfig, error) {
	return NewIniConfigFromFileWithOptions(file, DefaultIniOptions())
}

// NewIniConfigFromFileWithOptions loads the INI file behind the supplied file handle into a new IniConfig object
// using the supplied options. Caller is responsible for closing the supplied file.
//
// An error will be returned if there was a problem using the supplied file or parsing it as an INI file.
func NewIniConfigFromFileWithOptions(file *os.File, options *IniOptions) (*IniConfig, error) {

	if file == nil {
		return nil, errors.New("Nil file provided")
	}

	if options == nil {
		return nil, errors.New("Nil IniOptions provided")
	}

	if len(strings.TrimSpace(options.CommentStart)) == 0 {
		return nil, errors.New("CommentStart field in IniOptions cannot be empty")
	}

	ic := new(IniConfig)
	ic.options = options
	ic.sections = make(sectionPropertyMap)

	if err := ic.parse(file); err != nil {
		return nil, err
	} else {
		return ic, nil
	}

}

const rx_section = "\\[(.*)\\]"
const rx_property = "([^=]*)=(.*)"
const rx_colon_property = "([^=]*):(.*)"

// IniConfig provides access to configuration loaded in from an INI file. Functions exist to
// check whether a section or property exists; to recover the raw string value of a property or
// to try and interpret a property's value as a Go type
//
// The various PropertyValueAsXXX methods are generally convenience functions over the builtin strconv.Parse functions.
type IniConfig struct {
	sections sectionPropertyMap
	options  *IniOptions
}

//SectionExists returns true if a section with the supplied name was found and parsed.
func (ic *IniConfig) SectionExists(sectionName string) bool {

	return ic.findSection(sectionName) != nil
}

//Section returns a view on the IniConfig with the same methods but constrained to a single section
func (ic *IniConfig) Section(sectionName string) (*IniSection, error) {

	if ic.SectionExists(sectionName) {
		is := new(IniSection)
		is.name = sectionName
		is.ic = ic

		return is, nil
	} else {

		return nil, errorf("Section %s does not exist", sectionName)

	}

}

//PropertyExists returns true if the section exists and it contains a property with the requested name
func (ic *IniConfig) PropertyExists(sectionName, propertyName string) bool {
	propertyName = ic.normalise(propertyName)

	if foundSection := ic.findSection(sectionName); foundSection == nil {
		return false
	} else {
		return foundSection[propertyName] != nil
	}

}

// Value returns the value of the specified property in the specified section.
//
// Returns an error if the section or property does not exist.
func (ic *IniConfig) Value(sectionName, propertyName string) (string, error) {

	section := ic.findSection(sectionName)
	propertyName = ic.normalise(propertyName)

	if section == nil {
		return "", errorf("No such section %s", sectionName)
	}

	if value := section[propertyName]; value == nil {
		return "",  errorf("No such property [%s].%s", sectionName, propertyName)
	} else {
		return value.String(), nil
	}

}

// ValueOrDefault returns the value of the specified property in the specified section or the string zero value
// (empty string) if the value could not be found
func (ic *IniConfig) ValueOrZero(sectionName, propertyName string) (string) {

	if v, err := ic.Value(sectionName, propertyName); err == nil {
		return v
	} else {
		return ""
	}

}

// ValueAsFloat64 attempts to convert the specified property to a float64.
//
// Returns an error if the section or property does not exist or if the value could not be converted to a float64
func (ic *IniConfig) ValueAsFloat64(sectionName, propertyName string) (float64, error) {

	origSectionName := sectionName
	origPropName := propertyName

	sv, err := ic.Value(sectionName, propertyName)


	if err != nil {
		//Value not found
		return 0, err
	}

	if v, err := strconv.ParseFloat(sv, 64); err == nil {
		return v, nil
	} else {

		return 0, errorf("Unable to interpret [%s].%s (%s) as a float64.", origSectionName, origPropName, sv)

	}

}

// ValueOrZeroAsFloat64 returns the value of the specified property in the specified section as a float64 or
// the float64 zero value (0) if the value could not be found
func (ic *IniConfig) ValueOrZeroAsFloat64(sectionName, propertyName string) (float64) {

	if v, err := ic.ValueAsFloat64(sectionName, propertyName); err == nil {
		return v
	} else {
		return 0
	}

}


// ValueAsInt64 attempts to convert the specified property to an int64.
//
// Returns an error if the section or property does not exist or if the value could not be converted to an int64
func (ic *IniConfig) ValueAsInt64(sectionName, propertyName string) (int64, error) {

	origSectionName := sectionName
	origPropName := propertyName

	sv, err := ic.Value(sectionName, propertyName)


	if err != nil {
		//Value not found
		return 0, err
	}

	if v, err := strconv.ParseInt(sv, 10, 64); err == nil {
		return v, nil
	} else {

		return 0, errorf("Unable to interpret [%s].%s (%s) as an int64.", origSectionName, origPropName, sv)

	}

}

// ValueOrZeroAsInt64 returns the value of the specified property in the specified section as an int64 or
// the int64 zero value (0) if the value could not be found
func (ic *IniConfig) ValueOrZeroAsInt64(sectionName, propertyName string) (int64) {

	if v, err := ic.ValueAsInt64(sectionName, propertyName); err == nil {
		return v
	} else {
		return 0
	}

}

// ValueAsUint64 attempts to convert the specified property to a uint64.
//
// Returns an error if the section or property does not exist or if the value could not be converted to a uint64
func (ic *IniConfig) ValueAsUint64(sectionName, propertyName string) (uint64, error) {

	origSectionName := sectionName
	origPropName := propertyName

	sv, err := ic.Value(sectionName, propertyName)


	if err != nil {
		//Value not found
		return 0, err
	}

	if v, err := strconv.ParseUint(sv, 10, 64); err == nil {
		return v, nil
	} else {

		return 0, errorf("Unable to interpret [%s].%s (%s) as a uint64.", origSectionName, origPropName, sv)

	}

}

// ValueOrZeroAsUint64 returns the value of the specified property in the specified section as a uint64 or
// the uint64 zero value (0) if the value could not be found
func (ic *IniConfig) ValueOrZeroAsUint64(sectionName, propertyName string) (uint64) {

	if v, err := ic.ValueAsUint64(sectionName, propertyName); err == nil {
		return v
	} else {
		return 0
	}

}


// ValueAsBool attempts to convert the specified property to a uint64.
//
// Behaviour is affected by the IniOptions  supplied when creating this IniConfig. If the UseGoBoolRules field is set
// to true, conversion behaviour is as defined by strconv.ParseBool
//
// If UseGoBoolRules is set to false, the property value must be equal to the StrictBoolTrue field to be considered 'true' or
// match StrictBoolFalse to be considered 'false'. This matching can be made case insensitive by setting StrictBoolCaseSensitive to false
func (ic *IniConfig) ValueAsBool(sectionName, propertyName string) (bool, error) {

	sv, err := ic.Value(sectionName, propertyName)

	origSv := sv

	options := ic.options

	if err != nil {
		//Value not found
		return false, err
	}


	if options.UseGoBoolRules {
		//Allow any value Go would normally interpret as a bool
		if bv, err := strconv.ParseBool(sv); err == nil {
			return bv, nil
		} else {
			return false, errorf("Unable to interpret [%s].%s as a Go bool.", sectionName, propertyName)
		}

	}

	//Require that specific values for true or false be matched
	strictTrue := options.StrictBoolTrue
	strictFalse := options.StrictBoolFalse

	if !options.StrictBoolCaseSensitive {
		strictTrue = strings.ToUpper(strictTrue)
		strictFalse = strings.ToUpper(strictFalse)
		sv = strings.ToUpper(sv)
	}


	if sv == strictTrue {
		return true, nil
	} else if sv == strictFalse {
		return false, nil
	} else {

		return false, errorf("Value of [%s].%s (%s) could not be matched to %s or %s", sectionName, propertyName, origSv, options.StrictBoolTrue, options.StrictBoolFalse)

	}



	return false, nil
}

// ValueOrZeroAsBool returns the value of the specified property in the specified section as a bool or
// the bool zero value (false) if the value could not be found
func (ic *IniConfig) ValueOrZeroAsBool(sectionName, propertyName string) (bool) {

	if v, err := ic.ValueAsBool(sectionName, propertyName); err == nil {
		return v
	} else {
		return false
	}

}

// Add stores a property in the named section. If the property already exists, its value is overwritten.
func (ic *IniConfig) Add(section, propertyName string, value string) {

	section = ic.normalise(section)
	propertyName = ic.normalise(propertyName)

	storedSection := ic.sections[section]

	if storedSection == nil {
		storedSection = make(map[string]*nilableString)
		ic.sections[section] = storedSection
	}

	storedSection[propertyName] = newNilableString(value)

}

//parse scans the supplied file line by line according to the rules defined in the IniOptions
func (ic *IniConfig) parse(cf *os.File) error {
	s := bufio.NewScanner(cf)
	section := GLOBAL_SECTION

	options := ic.options

    var propRx *regexp.Regexp
	sectionRx := regexp.MustCompile(rx_section)
    if options.UseColonAssignment == true {
	    propRx = regexp.MustCompile(rx_colon_property)
    } else {
	    propRx = regexp.MustCompile(rx_property)
    }

	lineNumber := 0

	for s.Scan() {

		lineNumber++

		l := strings.TrimSpace(s.Text())
		lineLength := len(l)

		if lineLength == 0 && !options.TolerateBlankLines {
			return errorf("Blank line on line %d (forbidden in IniOptions)", lineNumber)
		} else if lineLength == 0 || strings.HasPrefix(l, options.CommentStart) {
			//Blank line or comment - ignore
			continue
		}

		l = ic.stripInlineComments(l)

		if sectionRx.MatchString(l) {
			matches := sectionRx.FindStringSubmatch(l)

			if len(matches) != 2 {
				return errorf("Unparseable section line in file at line %d", lineNumber)
			}

			section = matches[1]

		} else if propRx.MatchString(l) {

			if section == GLOBAL_SECTION && !options.AllowGlobalSection {
				return errorf("Property on line %d is outside of a named section (forbidden in IniOptions)", lineNumber)
			}


			matches := propRx.FindStringSubmatch(l)

			if len(matches) != 3{
				return errorf("Unparseable property line in file at line %d", lineNumber)
			}

			key := matches[1]
			value := matches[2]

			if options.TrimProperties {
				key = strings.TrimSpace(key)
				value = strings.TrimSpace(value)
			}

			value = ic.stripQuotes(value)

			if len(value) > 0 || !options.DiscardPropertiesWithNoValue {
				ic.Add(section, key, value)
			}

		} else {

			if !options.IgnoreUnparseable {
				return errorf("Unparseable line in file at line %d", lineNumber)
			}
		}
	}

	return nil
}

func (ic *IniConfig) stripQuotes(value string) string {

	options := ic.options

	vLength := len(value)

	if !options.StripEnclosingQuotes || vLength < 2{
		return value
	}

	for _, r := range options.EnclosingQuoteSymbols {

		if rune(value[0]) == r && rune(value[vLength-1]) == r {

			stripped := value[1:vLength-1]

			return stripped
		}

	}

	return value
}

func (ic *IniConfig) stripInlineComments(line string) string {

	options := ic.options

	if !options.AllowInlineComments {
		return line
	}

	ph := "[ESC_PH?]"
	escapeSeq := options.CommentEscapePrefix + options.CommentStart

	line = strings.Replace(line, escapeSeq, ph, -1)

	line = strings.Split(line, options.CommentStart)[0]

	line = strings.Replace(line, ph, options.CommentStart, -1)

	return line

}


func (ic *IniConfig) findSection(sectionName string) map[string]*nilableString {
	sectionName = ic.normalise(sectionName)

	return ic.sections[sectionName]
}



func (ic *IniConfig) normalise(s string) string {
	if ic.options.CaseSensitive {
		return s
	} else {
		return strings.ToLower(s)
	}
}

func errorf(template string, args ...interface{}) error {
	m := fmt.Sprintf(template, args...)

	return errors.New(m)
}
