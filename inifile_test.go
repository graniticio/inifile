package inifile

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const testfiles_base = "testfiles"
const simple_file = "simple.ini"
const types_file = "types.ini"
const colon_file = "colons.ini"

func simplePath() string {
	return filepath.Join(testfiles_base, simple_file)
}

func typesPath() string {
	return filepath.Join(testfiles_base, types_file)
}

func colonPath() string {
	return filepath.Join(testfiles_base, colon_file)
}

func TestNewFunctions(t *testing.T) {

	path := simplePath()

	options := DefaultIniOptions()


	if ic, err := NewIniConfigFromPath(path); err != nil {
		t.Errorf("Error loading simple INI file with NewIniConfigFromPath: %s", err.Error())
		t.FailNow()
	} else if ic == nil {
		t.Error("Unexpected nil object from NewIniConfigFromPath")
	} else if !ic.SectionExists("Section1"){
		t.Error("Missing section when created from NewIniConfigFromPath")
	}

	if ic, err := NewIniConfigFromPathWithOptions(path, options); err != nil {
		t.Errorf("Error loading simple INI file with NewIniConfigFromPathWithOptions %s", err.Error())
		t.FailNow()
	} else if ic == nil {
		t.Error("Unexpected nil object from NewIniConfigFromPathWithOptions")
	} else if !ic.SectionExists("Section1"){
		t.Error("Missing section when created from NewIniConfigFromPathWithOptions")
	}

	if f, err := os.Open(path); err != nil {
		t.Errorf("Unable to open test file at %s: %s", path, err.Error())
	} else {

		defer f.Close()

		if ic, err := NewIniConfigFromFile(f); err != nil {
			t.Errorf("Error loading simple INI file with NewIniConfigFromFile: %s", err.Error())
			t.FailNow()
		} else if ic == nil {
			t.Error("Unexpected nil object from NewIniConfigFromFile")
		} else if !ic.SectionExists("Section1"){
			t.Error("Missing section when created from NewIniConfigFromFile")
		}

	}

	if f, err := os.Open(path); err != nil {
		t.Errorf("Unable to open test file at %s: %s", path, err.Error())
	} else {

		defer f.Close()

		if ic, err := NewIniConfigFromFileWithOptions(f, options); err != nil {
			t.Errorf("Error loading simple INI file with NewIniConfigFromFileWithOptions: %s", err.Error())
			t.FailNow()
		} else if ic == nil {
			t.Error("Unexpected nil object from NewIniConfigFromFileWithOptions")
		} else if !ic.SectionExists("Section1"){
			t.Error("Missing section when created from NewIniConfigFromFileWithOptions")
		}
	}

}

func TestAlternateComments(t *testing.T) {

	path := filepath.Join(testfiles_base, "alternate-comments.ini")

	options := DefaultIniOptions()

	_, err := NewIniConfigFromPathWithOptions(path, options)

	if err == nil {
		t.Errorf("Expected parse to fail")
		t.FailNow()
	}

	options.CommentStart = "#"
}

func TestBlankLines(t *testing.T) {

	path := simplePath()

	options := DefaultIniOptions()
	options.TolerateBlankLines = false

	_, err := NewIniConfigFromPathWithOptions(path, options)

	if err == nil {
		t.Errorf("Expected parse to fail")
		t.FailNow()
	}

}

func TestUnparseableLines(t *testing.T) {

	path := filepath.Join(testfiles_base, "unparseable-lines.ini")

	options := DefaultIniOptions()

	_, err := NewIniConfigFromPathWithOptions(path, options)

	if err == nil {
		t.Errorf("Expected parse to fail")
		t.FailNow()
	}

	options.IgnoreUnparseable = true

	_, err = NewIniConfigFromPathWithOptions(path, options)

	if err != nil {
		t.Errorf("Error loading INI file %s: %s", path, err.Error())
		t.FailNow()
	}



}


func TestGlobalProperties(t *testing.T) {

	path := filepath.Join(testfiles_base, "global-section.ini")

	options := DefaultIniOptions()
	options.AllowGlobalSection = false

	_, err := NewIniConfigFromPathWithOptions(path, options)

	if err == nil {
		t.Errorf("Expected parse to fail")
		t.FailNow()
	}

	options.AllowGlobalSection = true

	ic, err := NewIniConfigFromPathWithOptions(path, options)

	if err != nil {
		t.Errorf("Error loading INI file %s: %s", path, err.Error())
		t.FailNow()
	}

	if !ic.PropertyExists(GLOBAL_SECTION, "globalProp") {
		t.Errorf("Could not find property globalProp in global section")
	}

	if v, _ := ic.Value(GLOBAL_SECTION, "globalProp"); v != "A" {
		t.Errorf("Unexpected value %s", v)
	}

}

func TestQuotedValues(t *testing.T) {

	path := filepath.Join(testfiles_base, "quoted-values.ini")

	options := DefaultIniOptions()
	options.StripEnclosingQuotes = true

	ic, err := NewIniConfigFromPathWithOptions(path, options)

	if err != nil {
		t.Errorf("Error loading INI file %s: %s", path, err.Error())
		t.FailNow()
	}

	if v, _ := ic.Value("section", "singleQuoted"); v != "quoted" {
		t.Errorf("Unexpected value")
	}

	if v, _ := ic.Value("section", "doubleQuoted"); v != "quoted" {
		t.Errorf("Unexpected value")
	}
}

func TestInlineAndEscapedComments(t *testing.T) {

	path := filepath.Join(testfiles_base, "inline-comments.ini")

	options := DefaultIniOptions()
	options.AllowInlineComments = true
	options.CommentStart = "#"

	ic, err := NewIniConfigFromPathWithOptions(path, options)

	if err != nil {
		t.Errorf("Error loading INI file %s: %s", path, err.Error())
		t.FailNow()
	}

	if !ic.SectionExists("uris") {
		t.Errorf("Expected to find section [uris]")
	}

	if !ic.SectionExists("#tags") {
		t.Errorf("Expected to find section [#tags]")
	}

	if !ic.PropertyExists("uris", "profile") {
		t.Errorf("Expected to find property [uris].profile")
		t.FailNow()
	}

	if v, _ := ic.Value("uris", "profile"); v != "http://example.com/profile#myprofile"{
		t.Errorf("Unexpected %s", v)
	}


	if !ic.PropertyExists("#tags", "latest#tag") {
		t.Errorf("Expected to find property [#tags].latest#tag")
		t.FailNow()
	}

	if v, _ := ic.Value("#tags", "latest#tag"); v != "#trending"{
		t.Errorf("Unexpected %s", v)
	}

}

func TestUnsetProperties(t *testing.T) {

	path := filepath.Join(testfiles_base, "unset.ini")

	options := DefaultIniOptions()

	ic, err := NewIniConfigFromPathWithOptions(path, options)

	if err != nil {
		t.Errorf("Error loading INI file %s: %s", path, err.Error())
		t.FailNow()
	}

	if ic.PropertyExists("section", "unset") {
		t.Errorf("Did not expect [section].unset to exist")
	}

	options.DiscardPropertiesWithNoValue = false

	ic, _ = NewIniConfigFromPathWithOptions(path, options)

	if !ic.PropertyExists("section", "unset") {
		t.Errorf("Did not expect [section].unset to exist")
	}

	if v, _ := ic.Value("section", "unset"); v != "" {
		t.Errorf("Did not expect [section].unset to have a values was >%s<", v)
	}

}

func TestSectionPropertyWhitespaceStripping(t *testing.T) {

	path := filepath.Join(testfiles_base, "whitespace-aligned.ini")

	options := DefaultIniOptions()

	ic, err := NewIniConfigFromPathWithOptions(path, options)

	if err != nil {
		t.Errorf("Error loading INI file %s: %s", path, err.Error())
		t.FailNow()
	}

	if v, err := ic.Value("section", "propA"); err != nil {
		t.Errorf("Expected to find value at [section].propA: %s", err.Error())
	} else if v != "Value1" {
		t.Errorf("Expected Value1 found >%s<", v)
	}

	options.TrimProperties = false

	ic, err = NewIniConfigFromPathWithOptions(path, options)

	if err != nil {
		t.Errorf("Error loading INI file %s: %s", path, err.Error())
		t.FailNow()
	}

	if ic.PropertyExists("section", "propA") {
		t.Errorf("Did not expect to find [section].propA")
	}

	if _, err := ic.Value("section", "propC"); err != nil {
		t.Errorf("Expected to find value at [section].propC: %s", err.Error())
	}
}


func TestSectionPropertyCaseSensitivity(t *testing.T) {

	path := filepath.Join(testfiles_base, "case-sensitivity.ini")

	options := DefaultIniOptions()

	ic, err := NewIniConfigFromPathWithOptions(path, options)

	if err != nil {
		t.Errorf("Error loading INI file %s: %s", path, err.Error())
		t.FailNow()
	}

	if !ic.SectionExists("ABC"){
		t.Errorf("Expected section ABC to exist")
	}

	if !ic.SectionExists("aBc"){
		t.Errorf("Expected section aBc to exist")
	}

	if ic.SectionExists("abc"){
		t.Errorf("Did not expect section abc to exist")
	}

	if v, _ := ic.Value("ABC", "value1"); v != "123" {
		t.Errorf("Expected [ABC].value1=123 was %s", v)
	}

	options.CaseSensitive = false

	ic, err = NewIniConfigFromPathWithOptions(path, options)

	if err != nil {
		t.Errorf("Error loading INI file %s: %s", path, err.Error())
		t.FailNow()
	}

	if !ic.SectionExists("ABC"){
		t.Errorf("Expected section ABC to exist")
	}

	if !ic.SectionExists("aBc"){
		t.Errorf("Expected section aBc to exist")
	}

	if !ic.SectionExists("abc"){
		t.Errorf("Expected section abc to exist")
	}

	if v, _ := ic.Value("ABC", "value1"); v != "890" {
		t.Errorf("Expected [ABC].value1=890 was %s", v)
	}
}

func TestSimpleParse(t *testing.T) {

	path := simplePath()

	ic, err := NewIniConfigFromPath(path)

	if err != nil {
		t.Errorf("Error loading simple INI file with NewIniConfigFromPath: %s", err.Error())
		t.FailNow()
	}

	if !ic.SectionExists("Section1") {
		t.Errorf("Expected Section1 to exist")
	}

	if ic.SectionExists("Section2") {
		t.Errorf("Did not expect Section2 to exist")
	}

	if !ic.PropertyExists("Section1", "name1") {
		t.Errorf("Expected [Section1].name1 to exist")
	}

	if ic.PropertyExists("Section2", "name1") {
		t.Errorf("Did not expect [Section2].name1 to exist")
	}

	v, _ := ic.Value("Section1", "name1")

	if v != "value1" {
		t.Errorf("Unexpected value %s for [Section1].name1. Was expecting value1", v)
	}

}

func TestSectionView(t *testing.T) {

	path := typesPath()

	ic, err := NewIniConfigFromPath(path)

	if err != nil {
		t.Fatalf("Problem loading test file %s", err.Error())
	}

	if s, err := ic.Section("uint"); err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	} else {

		if v, _ := s.Value("positive"); v != "4" {
			t.Errorf("Unexpected value %s", v)
		}

	}

	if s, _ := ic.Section("xxx"); s != nil {
		t.Errorf("Expected nil section")
	}
}

func TestFloat64Handling(t *testing.T) {

	path := typesPath()

	ic, err := NewIniConfigFromPath(path)

	if err != nil {
		t.Fatalf("Problem loading test file %s", err.Error())
	}

	if v, err := ic.ValueAsFloat64("float", "positive"); err != nil {
		t.Errorf("Unexpected error with valid float64: %s ", err.Error())
	} else if v != float64(4) {
		t.Errorf("Unexpected value with valid float64: %v", v)
	}

	if v, err := ic.ValueAsFloat64("float", "negative"); err != nil {
		t.Errorf("Unexpected error with valid float64: %s ", err.Error())
	} else if v != float64(-2.3333) {
		t.Errorf("Unexpected value with valid float64: %v", v)
	}

	if _, err := ic.ValueAsFloat64("float", "string"); err == nil {
		t.Errorf("Expected string to fail")
	}

}

func TestUintHandling(t *testing.T) {

	path := typesPath()

	ic, err := NewIniConfigFromPath(path)

	if err != nil {
		t.Fatalf("Problem loading test file %s", err.Error())
	}

	if v, err := ic.ValueAsUint64("uint", "positive"); err != nil {
		t.Errorf("Unexpected error with valid uint: %s ", err.Error())
	} else if v != uint64(4) {
		t.Errorf("Unexpected value with valid uint: %v", v)
	}

	if _, err := ic.ValueAsUint64("uint", "negative"); err == nil {
		t.Errorf("Expected negative number to fail")
	}

	if _, err := ic.ValueAsUint64("uint", "float"); err == nil {
		t.Errorf("Expected float to fail")
	}

	if _, err := ic.ValueAsUint64("uint", "string"); err == nil {
		t.Errorf("Expected string to fail")
	}

}

func TestIntHandling(t *testing.T) {

	path := typesPath()

	ic, err := NewIniConfigFromPath(path)

	if err != nil {
		t.Fatalf("Problem loading test file %s", err.Error())
	}

	if v, err := ic.ValueAsInt64("int", "positive"); err != nil {
		t.Errorf("Unexpected error with valid int: %s ", err.Error())
	} else if v != int64(4) {
		t.Errorf("Unexpected value with valid int: %v", v)
	}

	if v, err := ic.ValueAsInt64("int", "negative"); err != nil {
		t.Errorf("Unexpected error with valid int: %s ", err.Error())
	} else if v != int64(-1) {
		t.Errorf("Unexpected value with valid int: %v", v)
	}


	if _, err := ic.ValueAsInt64("int", "float"); err == nil {
		t.Errorf("Expected float to fail")
	}

	if _, err := ic.ValueAsInt64("int", "string"); err == nil {
		t.Errorf("Expected string to fail")
	}

}


func TestBooleanHandling(t *testing.T) {

	path := typesPath()
	options := DefaultIniOptions()

	ic, err := NewIniConfigFromPathWithOptions(path, options)

	if err != nil {
		t.Fatalf("Problem loading test file %s", err.Error())
	}

	expectBool(t, ic, "Boolean", "value1", true) //True
	expectBool(t, ic, "Boolean", "value2", true)
	expectBool(t, ic, "Boolean", "value3", true)
	expectBool(t, ic, "Boolean", "value4", false) //False
	expectBool(t, ic, "Boolean", "value5", false)
	expectBool(t, ic, "Boolean", "value6", false)
	expectBool(t, ic, "Boolean", "value7", false)


	//Enable strict bool matching
	options.UseGoBoolRules = false
	options.StrictBoolTrue = "True"
	options.StrictBoolFalse = "False"

	expectBool(t, ic, "Boolean", "value1", true) //True
	expectBool(t, ic, "Boolean", "value4", false) //False

	if _, err := ic.ValueAsBool("Boolean", "value2"); err == nil {
		//1
		t.Errorf("Expected bool conversion to fail")
	}

	if _, err := ic.ValueAsBool("Boolean", "value3"); err == nil {
		//TRUE
		t.Errorf("Expected bool conversion to fail")
	}

	if _, err := ic.ValueAsBool("Boolean", "value7"); err == nil {
		//FALSE
		t.Errorf("Expected bool conversion to fail")
	}

	options.StrictBoolCaseSensitive = false

	expectBool(t, ic, "Boolean", "value3", true) //TRUE
	expectBool(t, ic, "Boolean", "value7", false) //FALSE

	if _, err := ic.ValueAsBool("Boolean", "value5"); err == nil {
		//0
		t.Errorf("Expected bool conversion to fail")
	}

}

func TestColonParse(t *testing.T) {
	path := colonPath()
	options := DefaultIniOptions()
	options.UseColonAssignment = true

	ic, err := NewIniConfigFromPathWithOptions(path, options)
	if err != nil {
		t.Errorf("Unexpected Error: %s", err)
	}
	val, err := ic.Value("Section1", "name1")
	if err != nil {
		t.Errorf("Can't find expected string '%s': error: %s", "value1", err)
	}
	if val != "value1" {
		t.Errorf("Got %s, expected %s", val, "value1")
	}
}

func expectBool(t *testing.T, ic *IniConfig, sectionName, propertyName string, expected bool) {

	if v, err := ic.ValueAsBool(sectionName, propertyName); err != nil {
		t.Errorf("Unexpected error at line %s: %s", determineLine(), err.Error())
	} else if v != expected {
		t.Errorf("Expected %v found %v at line %s", expected, v, determineLine())
	}

}

func determineLine() string {
	trace := make([]byte, 2048)
	runtime.Stack(trace, false)

	splitTrace := strings.SplitN(string(trace), "\n", -1)

	if len(splitTrace) < 7 {
		return "?"
	}

	l := splitTrace[6]
	trimmed := strings.TrimSpace(l)
	p := strings.SplitN(trimmed, " +", -1)[0]

	f := strings.SplitN(p, string(os.PathSeparator), -1)

	return f[len(f)-1]

}
