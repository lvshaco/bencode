package bencode

import (
    "testing"
    "errors"
    "fmt"
)

type Case struct {
    stream string
    origin interface{}
}

var casesList = [...][]Case{
{
    {"0:", ""},
    {"1:a", "a"},
    {"5:hello", "hello"},
},
{
    {"i123e", 123},
    {"i0e", 0},
    {"i-1e", -1},
},
{
    {"li123ei-1ee", []interface{}{123, -1}},
    {"l5:helloe", []interface{}{"hello"}},
    {"ld5:hello5:worldee", []interface{}{map[string]interface{}{"hello": "world"}}},
    {"lli1ei2eee", []interface{}{[]interface{}{1, 2}}},
},
{
    {"d5:helloi100ee", map[string]interface{}{"hello": 100}},
    {"d3:foo3:bare", map[string]interface{}{"foo": "bar"}},
    {"d1:ad3:foo3:baree", map[string]interface{}{"a": map[string]interface{}{"foo": "bar"}}},
    {"d4:listli1eee", map[string]interface{}{"list": []interface{}{1}}},
},
}

func testEncode(cases []Case) (error) {
	for _, c := range cases {
		out, err := Encode(c.origin)
        if err != nil {
            return err
		}
        out2 := string(out)
        if out2 != c.stream {
            return errors.New(fmt.Sprintf("%s != %s", out2, c.stream))
        }
	}
    return nil
}

func TestEncode(t *testing.T) {
    for _, cases := range casesList {
        if err := testEncode(cases); err != nil {
            t.Error(err)
        }
    }
}

func BenchmarkEncode(b *testing.B) {
    for i:=0; i<b.N; i++ {
        for _, cases := range casesList {
            testEncode(cases)
        }
    }
}

func assertEqual(a interface{}, b interface{}) {
    switch v1 := a.(type) {
    case int:
        v2, ok := b.(int)
        if !ok {
            panic("Not int")
        }
        if v1 != v2 {
            panic("Not equal int")
        }
    case string:
        v2, ok := b.(string)
        if !ok {
            panic("Not string")
        }
        if v1 != v2 {
            panic("Not equal string")
        }
    case []interface{}:
        v2, ok := b.([]interface{})
        if !ok {
            panic("Not list")
        }
        if len(v1) != len(v2) {
            panic("Not equal list size")
        }
        for i := 0; i<len(v1); i++ {
            assertEqual(v1[i], v2[i])
        }
    case map[string]interface{}:
        v2, ok := b.(map[string]interface{})
        if !ok {
            panic("Not dict")
        }
        if len(v1) != len(v2) {
            panic("Not equal dict size")
        }
        for k, sub := range v1 {
            sub2, ok := v2[k]
            if !ok {
                panic("Not equal dict key")
            }
            assertEqual(sub, sub2)
        }
    }
}

func checkEqual(a interface{}, b interface{}) (err error) {
    defer func() {
        if e := recover(); e != nil {
            err = e.(error)
        }
    }()
    assertEqual(a, b)
    return nil
}

func testDecode(cases []Case) (error) {
    for _, c := range cases {
        out, err := Decode([]byte(c.stream))
        if err != nil {
            return err
        }
        if err := checkEqual(out, c.origin); err != nil {
            return err
        }
    }
    return nil
}

func TestDecode(t *testing.T) {
    for _, cases := range casesList {
        if err := testDecode(cases); err != nil {
            t.Error(err)
        }
    }
}

func BenchmarkDecode(b *testing.B) {
    for i:=0; i<b.N; i++ {
        for _, cases := range casesList {
            testDecode(cases)
        }
    }
}
