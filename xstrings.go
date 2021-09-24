package xstrings

import (
	"unicode"
	"unicode/utf8"
	"unsafe"
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

// ToUpper returns s with all Unicode letters mapped to their upper case.
func ToUpper(s string) string {
	isASCII, hasLower := true, false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= utf8.RuneSelf {
			isASCII = false
			break
		}
		hasLower = hasLower || ('a' <= c && c <= 'z')
	}

	if isASCII { // optimize for ASCII-only strings.
		if !hasLower {
			return s
		}
		b := make([]byte, len(s))
		for i := 0; i < len(s); i++ {
			c := s[i]
			if 'a' <= c && c <= 'z' {
				c -= 'a' - 'A'
			}
			b[i] = c
		}
		return *(*string)(unsafe.Pointer(&b))
	}
	return Map(unicode.ToUpper, s)
}

// ToLower returns s with all Unicode letters mapped to their lower case.
func ToLower(s string) string {
	isASCII, hasUpper := true, false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= utf8.RuneSelf {
			isASCII = false
			break
		}
		hasUpper = hasUpper || ('A' <= c && c <= 'Z')
	}

	if isASCII { // optimize for ASCII-only strings.
		if !hasUpper {
			return s
		}
		b := make([]byte, len(s))
		for i := 0; i < len(s); i++ {
			c := s[i]
			if 'A' <= c && c <= 'Z' {
				c += 'a' - 'A'
			}
			b[i] = c
		}
		return *(*string)(unsafe.Pointer(&b))
	}
	return Map(unicode.ToLower, s)
}


// TrimSpace returns a slice of the string s, with all leading
// and trailing white space removed, as defined by Unicode.
func TrimSpace(s string) string {
	// Fast path for ASCII: look for the first ASCII non-space byte
	start := 0
	for ; start < len(s); start++ {
		c := s[start]
		if c >= utf8.RuneSelf {
			// If we run into a non-ASCII byte, fall back to the
			// slower unicode-aware method on the remaining bytes
			return TrimFunc(s[start:], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	// Now look for the first ASCII non-space byte from the end
	stop := len(s)
	for ; stop > start; stop-- {
		c := s[stop-1]
		if c >= utf8.RuneSelf {
			return TrimRightFunc(s[start:stop], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	// At this point s[start:stop] starts and ends with an ASCII
	// non-space bytes, so we're done. Non-ASCII cases have already
	// been handled above.
	return s[start:stop]
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
