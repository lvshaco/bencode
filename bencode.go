package bencode

import (
    "strconv"
    "bytes"
    "errors"
)

func encodeString(dst []byte, s string) ([]byte, error) {
    dst = strconv.AppendInt(dst, int64(len(s)), 10)
    dst = append(dst, ':')
    dst = append(dst, s...)
    return dst, nil
}

func encodeInt(dst []byte, i int) ([]byte, error) {
    dst = append(dst, 'i')
    dst = strconv.AppendInt(dst, int64(i), 10)
    dst = append(dst, 'e')
    return dst, nil
}

func encodeList(dst []byte, l []interface{}) ([]byte, error) {
    var e error
    dst = append(dst, 'l')
    for _, v := range l {
        if dst, e = encodeItem(dst, v); e != nil {
            return nil, e
        }
    }
    dst = append(dst, 'e')
    return dst, nil
}

func encodeDict(dst []byte, d map[string]interface{}) ([]byte, error) {
    var e error
    dst = append(dst, 'd')
    for k, v := range d {
        if dst, e = encodeItem(dst, k); e != nil {
            return nil, e
        }
        if dst, e = encodeItem(dst, v); e != nil {
            return nil, e
        }
    }
    dst = append(dst, 'e')
    return dst, nil
}

func encodeItem(dst []byte, t interface{}) ([]byte, error) {
    switch v := t.(type) {
    case string:
        return encodeString(dst, v)
    case int:
        return encodeInt(dst, v)
    case []interface{}:
        return encodeList(dst, v)
    case map[string]interface{}:
        return encodeDict(dst, v)
    default:
        return nil, errors.New("Encode invalid type")
    }
}

func Encode(it interface{}) ([]byte, error) {
    dst := make([]byte, 0, 16) // 16 enough for int
    return encodeItem(dst, it)
}

func decodeString(s []byte) (string, []byte, error) {
    split := bytes.IndexByte(s, ':')
    if split == -1 {
        return "", nil, errors.New("Decode invalid string")
    }
    size, e := strconv.Atoi(string(s[:split]))
    if e != nil {
        return "", nil, errors.New("Decode illegal string length")
    }
    last := split + size + 1
    if last > len(s) {
        return "", nil, errors.New("Decode invalid string length")
    }
    return string(s[split+1:last]), s[last:], nil
}

func decodeInt(s []byte) (int, []byte, error) {
    last := bytes.IndexByte(s, 'e')
    if last == -1 {
        return 0, nil, errors.New("Decode int no `e`")
    }
    v, e := strconv.Atoi(string(s[:last]))
    if e != nil {
        return 0, nil, e
    }
    return v, s[last+1:], nil
}

func decodeList(s []byte) ([]interface{}, []byte, error) {
    var e error
    var v interface{}
    r := make([]interface{}, 0)
    for len(s) > 0 {
        if s[0] == 'e' {
            return r, s[1:], nil
        }
        if v, s, e = decodeItem(s); e != nil {
            return nil, nil, e
        }
        r = append(r, v)
    }
    return nil, nil, errors.New("Decode list no `e`")
}

func decodeDict(s []byte) (map[string]interface{}, []byte, error) {
    var e error
    var k string
    var v interface{}
    r := make(map[string]interface{})
    for len(s) > 0 {
        if s[0] == 'e' {
            return r, s[1:], nil
        } else {
            if k, s, e = decodeString(s); e != nil {
                return nil, nil, e
            }
        }
        if v, s, e = decodeItem(s); e != nil {
            return nil, nil, e
        }
        r [k] = v
    }
    return nil, nil, errors.New("Decode dict no `e`")
}

func decodeItem(s []byte) (interface{}, []byte, error) {
    switch s[0] {
    case 'i':
        return decodeInt(s[1:])
    case 'l':
        return decodeList(s[1:])
    case 'd':
        return decodeDict(s[1:])
    default:
        return decodeString(s)
    }
}

func Decode(s []byte) (interface{}, error){
    r, s, e := decodeItem(s)
    if len(s) > 0 {
        e = errors.New("Decode no end")
    }
    return r, e
}
