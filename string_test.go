package null

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

var (
	stringJSON      = []byte(`"test"`)
	blankStringJSON = []byte(`""`)
	nullStringJSON  = []byte(`{"String":"test","Valid":true}`)

	nullJSON    = []byte(`null`)
	invalidJSON = []byte(`:)`)
)

type stringInStruct struct {
	Test String `json:"test,omitempty"`
}

func TestStringFrom(t *testing.T) {
	str := StringFrom("test")
	assertStr(t, str, "StringFrom() string")

	zero := StringFrom("")
	if !zero.Valid {
		t.Error("StringFrom(0)", "is invalid, but should be valid")
	}
}

func TestStringFromPtr(t *testing.T) {
	s := "test"
	sptr := &s
	str := StringFromPtr(sptr)
	assertStr(t, str, "StringFromPtr() string")

	null := StringFromPtr(nil)
	assertNullStr(t, null, "StringFromPtr(nil)")
}

func TestUnmarshalString(t *testing.T) {
	var str String
	err := json.Unmarshal(stringJSON, &str)
	maybePanic(err)
	assertStr(t, str, "string json")

	var ns String
	err = json.Unmarshal(nullStringJSON, &ns)
	maybePanic(err)
	assertStr(t, ns, "sql.NullString json")

	var blank String
	err = json.Unmarshal(blankStringJSON, &blank)
	maybePanic(err)
	if !blank.Valid {
		t.Error("blank string should be valid")
	}

	var null String
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullStr(t, null, "null json")

	var badType String
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullStr(t, badType, "wrong type json")

	var invalid String
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullStr(t, invalid, "invalid json")
}

func TestTextUnmarshalString(t *testing.T) {
	var str String
	err := str.UnmarshalText([]byte("test"))
	maybePanic(err)
	assertStr(t, str, "UnmarshalText() string")

	var null String
	err = null.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullStr(t, null, "UnmarshalText() empty string")
}

func TestMarshalString(t *testing.T) {
	str := StringFrom("test")
	data, err := json.Marshal(str)
	maybePanic(err)
	assertJSONEquals(t, data, `"test"`, "non-empty json marshal")
	data, err = str.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "test", "non-empty text marshal")

	// empty values should be encoded as an empty string
	zero := StringFrom("")
	data, err = json.Marshal(zero)
	maybePanic(err)
	assertJSONEquals(t, data, `""`, "empty json marshal")
	data, err = zero.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "string marshal text")

	null := StringFromPtr(nil)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, `null`, "null json marshal")
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "string marshal text")
}

// Tests omitempty... broken until Go 1.4
// func TestMarshalStringInStruct(t *testing.T) {
// 	obj := stringInStruct{Test: StringFrom("")}
// 	data, err := json.Marshal(obj)
// 	maybePanic(err)
// 	assertJSONEquals(t, data, `{}`, "null string in struct")
// }

func TestStringPointer(t *testing.T) {
	str := StringFrom("test")
	ptr := str.Ptr()
	if *ptr != "test" {
		t.Errorf("bad %s string: %#v ≠ %s\n", "pointer", ptr, "test")
	}

	null := NewString("", false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s string: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestStringIsZero(t *testing.T) {
	str := StringFrom("test")
	if str.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	blank := StringFrom("")
	if blank.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	empty := NewString("", true)
	if empty.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := StringFromPtr(nil)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestStringSetValid(t *testing.T) {
	change := NewString("", false)
	assertNullStr(t, change, "SetValid()")
	change.SetValid("test")
	assertStr(t, change, "SetValid()")
}

func TestStringScan(t *testing.T) {
	var str String
	err := str.Scan("test")
	maybePanic(err)
	assertStr(t, str, "scanned string")

	var null String
	err = null.Scan(nil)
	maybePanic(err)
	assertNullStr(t, null, "scanned null")
}

func TestStringValueOrZero(t *testing.T) {
	valid := NewString("test", true)
	if valid.ValueOrZero() != "test" {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := NewString("test", false)
	if invalid.ValueOrZero() != "" {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

type XmlNodeTest interface {
	Marshal() string
}

type NodeTest struct {
	Test String `xml:"test"`
}

func (n NodeTest) Marshal() string {
	data, err := xml.Marshal(n)
	maybePanic(err)
	return string(data)
}

type AttrTest struct {
	Test String `xml:"test,attr"`
}

func (a AttrTest) Marshal() string {
	data, err := xml.Marshal(a)
	maybePanic(err)
	return string(data)
}

type StrTest struct {
	Expected string
	Struct   XmlNodeTest
}

type AttrNullTest struct {
	Test String `xml:"test,attr,omitempty"`
}

func (a AttrNullTest) Marshal() string {
	data, err := xml.Marshal(a)
	maybePanic(err)
	return string(data)
}

type AttributesTest struct {
	Test1 String `xml:"test1,attr,omitempty"`
	Test2 String `xml:"test2,attr,omitempty"`
	Test3 String `xml:"test3,attr,omitempty"`
}

type NodesTest struct {
	Test1 String `xml:"test1,omitempty"`
	Test2 String `xml:"test2,omitempty"`
}

func (n NodesTest) Marshal() string {
	data, err := xml.Marshal(n)
	maybePanic(err)
	return string(data)
}

func (a AttributesTest) Marshal() string {
	data, err := xml.Marshal(a)
	maybePanic(err)
	return string(data)
}

var data = []StrTest{{
	Expected: "<NodeTest><test>123</test></NodeTest>",
	Struct: NodeTest{
		Test: NewString("123", true),
	},
}, {
	Expected: `<AttrTest test="123"></AttrTest>`,
	Struct: AttrTest{
		Test: NewString("123", true),
	},
}, {
	Expected: `<AttrNullTest></AttrNullTest>`,
	Struct:   AttrNullTest{Test: NewString("", false)},
}, {
	Expected: `<AttributesTest test1="1" test2="2"></AttributesTest>`,
	Struct: AttributesTest{
		Test1: NewString("1", true),
		Test2: NewString("2", true),
	},
}, {
	Expected: `<AttributesTest test1="1"></AttributesTest>`,
	Struct: AttributesTest{
		Test1: NewString("1", true),
		Test2: NewString("", false),
	},
}, {
	Expected: `<NodesTest><test1>1</test1></NodesTest>`,
	Struct: NodesTest{
		Test1: NewString("1", true),
		Test2: NewString("", false),
	},
}, {
	Expected: `<NodesTest></NodesTest>`,
	Struct: NodesTest{
		Test1: NewString("123", false),
	},
}}

func TestMarshalXML(t *testing.T) {
	for _, s := range data {
		t.Run(s.Expected, func(t *testing.T) {
			if s.Expected != s.Struct.Marshal() {
				t.Errorf(`got %s want %s`, s.Struct.Marshal(), s.Expected)
			}
		})
	}
}

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func assertStr(t *testing.T, s String, from string) {
	if s.String != "test" {
		t.Errorf("bad %s string: %s ≠ %s\n", from, s.String, "test")
	}
	if !s.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullStr(t *testing.T, s String, from string) {
	if s.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}
