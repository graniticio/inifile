# inifile (v1.2.0)
A Go library for reading and parsing INI files.

Package inifile provides a Go struct that can parse an INI-style file and make the configuration in that file
available via a series of type-safe data accessors. Parsing and data-access behaviour can be configured to support most
INI file variants.

## Parsing

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


## Accessing properties

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


### Accessing properties via an IniSection

If your code needs multiple property values from the same section:

	a, err := ic.Value("section1", "a")
	b, err := ic.Value("section1", "b")
	c, err := ic.Value("section1", "c")

it is recommended that you use an IniSection object, which offers the same functions as IniConfig but is bound to a single section:

	is, err := ic.Section("section1")
	a, err := is.Value("a")
	b, err := is.Value("b")
	c, err := is.Value("c")


## Accessing properties in the global section

Use the constant <code>inifile.GLOBAL_SECTION</code> as the sectionName when calling any of the above functions to work with properties that are not
attached to a named section



## Adding new properties

Properties can be added to an IniConfig at runtime by calling:
	
	Add(section, propertyName string, value string)


## Customising parsing and configuration access

As INI files are not governed by an agreed standard, there are a number of variations in the structure and features
found in real-world INI files. To accommodate these variations, you can modify the IniOptions used when creating
an IniConfig object.

The recommended practice is to use the DefaultIniOptions() function to create a baseline set of options
and then only modify those you need to change (see the documentation for DefaultIniOptions below to see what the default values are). For example:

	opts := inifile.DefaultIniOptions()
	opts.CommentStart = "#"

	ic, err := NewIniConfigFromPathWithOptions("/path/to/file.ini", opts)

will parse a file using <code>#</code> instead of <code>;</code> to identify comment lines.

### Case-sensitivity for section and property names

By default look-ups of sections and properties are case sensitive - <code>Value("mysection", "myproperty")</code> would not match a property called myProperty in a section called [MYPROPERTY].

This is not the behaviour of many Windows implementations of INI file libraries, so you might find files that mix and match cases. To support this behaviour set:
	
	CaseSensitive = false
in your IniOptions.

### Comment lines

Most INI files use the semi-colon symbol at the start of a line to indicate that the line is a comment. There are notable exceptions, including MySQL INI files, that use a different character, often #. To support this set:
	
	CommentStart = "#"
in your IniOptions.

#### Property names and values with spaces

For readability, some INI files align property names and values like:
	
	basedir	=	/usr
	datadir	=	/var/lib/mysql

By default the whitespace either side of the property name and value is discarded and not considered part of the property name or value. This behaviour can be disabled by setting:
	
	TrimProperties = false
in your IniOptions.

### Blank lines

For readability, most INI files use blank lines to break up sections and properties. To disallow this and return an error if blank lines are encountered set:
	
	TolerateBlankLines = false
in your IniOptions.

### Allow properties outside of a section

Some INI files have properties defined before the first [Section]. These properties are considered to be in the 'global' section and by default this behaviour is supported. To forbid
properties in the global section, set:
	
	AllowGlobalSection = false
in your IniOptions.

Use the <code>inifile.GLOBAL_SECTION</code> constant as the section name to access properties in the global section.

### Unset properties

Some INI files allow a property to be defined without a value like
	
	[Reference]
	randomNumberSeed=

This can be problematic if, when set, that property is expected to hold a value of a specific data-type (a float in example above). Calling ValueAsFloat64 on the above would return an error.
By default, IniConfig ignores any unset properties, so in the above example calling
	
	PropertyExists("Reference", "randomNumberSeed")
would return false

To reverse this behaviour and have any unset properties stored with <code>""</code> as their value, set:
	
	DiscardPropertiesWithNoValue = false
in your IniOptions.

### Boolean values

Go's <code>strconv.ParseBool</code> function is extremely permissive about the values it considers to represent true or false (see https://golang.org/pkg/strconv/#ParseBool) and by default
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


### Unparseable lines

By default, an INI file will not parse correctly if a line is encountered in the file that cannot be interpreted as a
section, comment, property or a blank line. To ignore unparseable lines, set:
	
	IgnoreUnparseable = true
in your IniOptions.


### Inline comments

Some INI files allow a comment on the same line as a property or section:
	
	[database]	;Connection information
	host=localhost ;Default to localhost

By default this behaviour is not supported and the value of [database].host will be "host ;Default to localhost". To allow
inline comments set:
	
	AllowInlineComments true

When inline comments are allowed, you must escape any instances of the comment start character in your section names, property names and
property values. The default escape prefix is \ but you can change this by setting:
	
	CommentEscapePrefix
to whatever string you want.

For example, with 

    CommentStart = ";"
    CommentEscapePrefix="\" 
your INI file might contain:

	[messages]
	message=Service is too busy\;load is high ;Message shown when site is down

The value returned by 

    Value("messages", "message") 
    
would be "Service is too busy;load is high"

### Quoted values

Some INI files surround their values with quotes like:
	
	inbox="/folder/with spaces"
	outbox='/another folder/sub-folder'

By default these enclosing quotes are preserved and are included when you retrieve the property's value. To change this behaviour, set:
	
	StripEnclosingQuotes = true
in your IniOptions.

The runes (characters) that will be recognised as enclosing quotes are defined in
	
	EnclosingQuoteSymbols
and default to the single (') and double (") quote symbols.

### Assignment with Colon

Some INI files use the colon for assignment rather than the equals sign. To 
support this set:

    UseColonAssignment = true
in your IniOptions.