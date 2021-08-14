package xstrings

import (
	"unicode/utf8"
)

type ASCIISet [8]uint32

func (as *ASCIISet) contains(c byte) bool {
	return (as[c>>5] & (1 << uint(c&31))) != 0
}

func makeASCIISet(chars string) (as ASCIISet, ok bool) {
	for i := 0; i < len(chars); i++ {
		c := chars[i]
		if c >= utf8.RuneSelf {
			return as, false
		}
		as[c>>5] |= 1 << uint(c&31)
	}
	return as, true
}


func MustMakeASCIISet(chars string) ASCIISet {
	as, ok := makeASCIISet(chars)
	if !ok {
		panic("non ascii")
	}
	return as
}

func indexRune(s string, f rune) int {
	for i, r := range s {
		if f != r {
			return i
		}
	}
	return -1
}

func indexASCIISet(s string, f ASCIISet) int {
	for i, r := range s {
		if !(r < utf8.RuneSelf && f.contains(byte(r))) {
			return i
		}
	}
	return -1
}

func lastIndexRune(s string, f rune) int {
	for i := len(s); i > 0; {
		r, size := utf8.DecodeLastRuneInString(s[0:i])
		i -= size
		if f != r {
			return i
		}
	}
	return -1
}

func lastIndexASCIISet(s string, f ASCIISet) int {
	for i := len(s); i > 0; {
		r, size := utf8.DecodeLastRuneInString(s[0:i])
		i -= size
		if !(r < utf8.RuneSelf && f.contains(byte(r))) {
			return i
		}
	}
	return -1
}

func TrimRightByte(s string, b rune) string {
	i := lastIndexRune(s, b)
	if i >= 0 && s[i] >= utf8.RuneSelf {
		_, wid := utf8.DecodeRuneInString(s[i:])
		i += wid
	} else {
		i++
	}
	return s[0:i]
}

func TrimRightASCIISet(s string, b ASCIISet) string {
	i := lastIndexASCIISet(s, b)
	if i >= 0 && s[i] >= utf8.RuneSelf {
		_, wid := utf8.DecodeRuneInString(s[i:])
		i += wid
	} else {
		i++
	}
	return s[0:i]
}

func TrimLeftASCIISet(s string, b ASCIISet) string {
	i := indexASCIISet(s, b)
	if i == -1 {
		return ""
	}
	return s[i:]
}

func TrimASCIISet(s string, b ASCIISet) string {
	return TrimLeftASCIISet(TrimRightASCIISet(s, b), b)
}

func TrimLeftByte(s string, b rune) string {
	i := indexRune(s, b)
	if i == -1 {
		return ""
	}
	return s[i:]
}

func TrimRune(s string, b rune) string {
	return TrimLeftByte(TrimRightByte(s, b), b)
}

//func TrimLeftByte(s string, b byte) string {
//	l := len(s)
//	for l > 0 {
//		l--
//		if s[0] != b {
//			return s
//		}
//		s = s[1:]
//	}
//	return s
//}

//func TrimRightByte(s string, b byte) string {
//	l := len(s)
//	for l > 0 {
//		l--
//		if s[l] != b {
//			return s
//		}
//		s = s[:l]
//	}
//
//
//	return s
//}
