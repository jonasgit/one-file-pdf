// -----------------------------------------------------------------------------
// github.com/balacode/one-file-pdf                one-file-pdf/[pdf_encoding.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package pdf

// 8-bit text encoding for the built-in (standard 14) fonts and for
// document information strings.
//
// Go strings are UTF-8, but a PDF literal string used with a simple
// (Type1) font is a sequence of single-byte character codes. Text drawn
// with the built-in Latin fonts is therefore converted to WinAnsiEncoding
// (declared in the font dictionary), and the document information
// strings (Title, Author, ...) are converted to PDFDocEncoding.
//
// The code tables below come from the PDF Reference, Appendix D.1
// "Latin Character Set and Encodings" as transcribed in
// https://gist.github.com/zmughal/5130211 (WIN and PDF columns).
// Only non-ASCII entries are listed: codes 32-126 are identical to ASCII
// in both encodings. Characters that can't be encoded are replaced with '?'.

import (
	"bytes"
	"strconv"
)

// pdfIsSymbolicFont reports whether a built-in font uses its own
// built-in encoding instead of WinAnsiEncoding (Symbol, ZapfDingbats).
func pdfIsSymbolicFont(builtInIndex int) bool {
	return builtInIndex == 4 || builtInIndex == 9
} //                                                           pdfIsSymbolicFont

// encodeBuiltInFontText converts UTF-8 text to the byte codes of the
// current built-in font: WinAnsiEncoding for Latin text fonts, or the raw
// code (0-255) of each rune for the symbolic fonts.
func (p *PDF) encodeBuiltInFontText(s string) string {
	if p.page != nil && p.page.fontID > 0 && p.page.fontID <= len(p.fonts) &&
		pdfIsSymbolicFont(p.fonts[p.page.fontID-1].builtInIndex) {
		buf := make([]byte, 0, len(s))
		for _, r := range s {
			if r < 0 || r > 255 {
				r = '?'
			}
			buf = append(buf, byte(r))
		}
		return string(buf)
	}
	return p.encodeText(s, pdfWinAnsiEncoding)
} //                                                       encodeBuiltInFontText

// encodeText converts UTF-8 text to single-byte codes using 'table' for
// non-ASCII characters. Unencodable characters become '?'.
func (*PDF) encodeText(s string, table map[rune]byte) string {
	buf := make([]byte, 0, len(s))
	for _, r := range s {
		switch c, ok := table[r]; {
		case r < 128:
			buf = append(buf, byte(r))
		case ok:
			buf = append(buf, c)
		default:
			buf = append(buf, '?')
		}
	}
	return string(buf)
} //                                                                  encodeText

// escapeBytes formats already-encoded bytes for use inside a PDF literal
// string: '(', ')' and '\' are escaped, and bytes outside printable
// ASCII are written as octal escapes (\ddd) so the output stays 7-bit.
func (*PDF) escapeBytes(s string) string {
	var buf bytes.Buffer
	buf.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '(' || c == ')' || c == '\\':
			buf.WriteByte('\\')
			buf.WriteByte(c)
		case c < 32 || c > 126:
			oct := strconv.FormatInt(int64(c), 8)
			buf.WriteByte('\\')
			for j := len(oct); j < 3; j++ {
				buf.WriteByte('0')
			}
			buf.WriteString(oct)
		default:
			buf.WriteByte(c)
		}
	}
	return buf.String()
} //                                                                 escapeBytes

// pdfWinAnsiEncoding maps Unicode characters to WinAnsiEncoding codes
// (non-ASCII part, 123 codes). Same as Windows code page 1252.
var pdfWinAnsiEncoding = map[rune]byte{
	0x00A0: 0xA0, // no-break space (footnote 6)
	0x00AD: 0xAD, // soft hyphen (footnote 5)
	0x20AC: 0x80, // € Euro
	0x201A: 0x82, // ‚ quotesinglbase
	0x0192: 0x83, // ƒ florin
	0x201E: 0x84, // „ quotedblbase
	0x2026: 0x85, // … ellipsis
	0x2020: 0x86, // † dagger
	0x2021: 0x87, // ‡ daggerdbl
	0x02C6: 0x88, // ˆ circumflex
	0x2030: 0x89, // ‰ perthousand
	0x0160: 0x8A, // Š Scaron
	0x2039: 0x8B, // ‹ guilsinglleft
	0x0152: 0x8C, // Œ OE
	0x017D: 0x8E, // Ž Zcaron
	0x2018: 0x91, // ‘ quoteleft
	0x2019: 0x92, // ’ quoteright
	0x201C: 0x93, // “ quotedblleft
	0x201D: 0x94, // ” quotedblright
	0x2022: 0x95, // • bullet
	0x2013: 0x96, // – endash
	0x2014: 0x97, // — emdash
	0x02DC: 0x98, // ˜ tilde
	0x2122: 0x99, // ™ trademark
	0x0161: 0x9A, // š scaron
	0x203A: 0x9B, // › guilsinglright
	0x0153: 0x9C, // œ oe
	0x017E: 0x9E, // ž zcaron
	0x0178: 0x9F, // Ÿ Ydieresis
	0x00A1: 0xA1, // ¡ exclamdown
	0x00A2: 0xA2, // ¢ cent
	0x00A3: 0xA3, // £ sterling
	0x00A4: 0xA4, // ¤ currency
	0x00A5: 0xA5, // ¥ yen
	0x00A6: 0xA6, // ¦ brokenbar
	0x00A7: 0xA7, // § section
	0x00A8: 0xA8, // ¨ dieresis
	0x00A9: 0xA9, // © copyright
	0x00AA: 0xAA, // ª ordfeminine
	0x00AB: 0xAB, // « guillemotleft
	0x00AC: 0xAC, // ¬ logicalnot
	0x00AE: 0xAE, // ® registered
	0x00AF: 0xAF, // ¯ macron
	0x00B0: 0xB0, // ° degree
	0x00B1: 0xB1, // ± plusminus
	0x00B2: 0xB2, // ² twosuperior
	0x00B3: 0xB3, // ³ threesuperior
	0x00B4: 0xB4, // ´ acute
	0x00B5: 0xB5, // µ mu
	0x00B6: 0xB6, // ¶ paragraph
	0x00B7: 0xB7, // · periodcentered
	0x00B8: 0xB8, // ¸ cedilla
	0x00B9: 0xB9, // ¹ onesuperior
	0x00BA: 0xBA, // º ordmasculine
	0x00BB: 0xBB, // » guillemotright
	0x00BC: 0xBC, // ¼ onequarter
	0x00BD: 0xBD, // ½ onehalf
	0x00BE: 0xBE, // ¾ threequarters
	0x00BF: 0xBF, // ¿ questiondown
	0x00C0: 0xC0, // À Agrave
	0x00C1: 0xC1, // Á Aacute
	0x00C2: 0xC2, // Â Acircumflex
	0x00C3: 0xC3, // Ã Atilde
	0x00C4: 0xC4, // Ä Adieresis
	0x00C5: 0xC5, // Å Aring
	0x00C6: 0xC6, // Æ AE
	0x00C7: 0xC7, // Ç Ccedilla
	0x00C8: 0xC8, // È Egrave
	0x00C9: 0xC9, // É Eacute
	0x00CA: 0xCA, // Ê Ecircumflex
	0x00CB: 0xCB, // Ë Edieresis
	0x00CC: 0xCC, // Ì Igrave
	0x00CD: 0xCD, // Í Iacute
	0x00CE: 0xCE, // Î Icircumflex
	0x00CF: 0xCF, // Ï Idieresis
	0x00D0: 0xD0, // Ð Eth
	0x00D1: 0xD1, // Ñ Ntilde
	0x00D2: 0xD2, // Ò Ograve
	0x00D3: 0xD3, // Ó Oacute
	0x00D4: 0xD4, // Ô Ocircumflex
	0x00D5: 0xD5, // Õ Otilde
	0x00D6: 0xD6, // Ö Odieresis
	0x00D7: 0xD7, // × multiply
	0x00D8: 0xD8, // Ø Oslash
	0x00D9: 0xD9, // Ù Ugrave
	0x00DA: 0xDA, // Ú Uacute
	0x00DB: 0xDB, // Û Ucircumflex
	0x00DC: 0xDC, // Ü Udieresis
	0x00DD: 0xDD, // Ý Yacute
	0x00DE: 0xDE, // Þ Thorn
	0x00DF: 0xDF, // ß germandbls
	0x00E0: 0xE0, // à agrave
	0x00E1: 0xE1, // á aacute
	0x00E2: 0xE2, // â acircumflex
	0x00E3: 0xE3, // ã atilde
	0x00E4: 0xE4, // ä adieresis
	0x00E5: 0xE5, // å aring
	0x00E6: 0xE6, // æ ae
	0x00E7: 0xE7, // ç ccedilla
	0x00E8: 0xE8, // è egrave
	0x00E9: 0xE9, // é eacute
	0x00EA: 0xEA, // ê ecircumflex
	0x00EB: 0xEB, // ë edieresis
	0x00EC: 0xEC, // ì igrave
	0x00ED: 0xED, // í iacute
	0x00EE: 0xEE, // î icircumflex
	0x00EF: 0xEF, // ï idieresis
	0x00F0: 0xF0, // ð eth
	0x00F1: 0xF1, // ñ ntilde
	0x00F2: 0xF2, // ò ograve
	0x00F3: 0xF3, // ó oacute
	0x00F4: 0xF4, // ô ocircumflex
	0x00F5: 0xF5, // õ otilde
	0x00F6: 0xF6, // ö odieresis
	0x00F7: 0xF7, // ÷ divide
	0x00F8: 0xF8, // ø oslash
	0x00F9: 0xF9, // ù ugrave
	0x00FA: 0xFA, // ú uacute
	0x00FB: 0xFB, // û ucircumflex
	0x00FC: 0xFC, // ü udieresis
	0x00FD: 0xFD, // ý yacute
	0x00FE: 0xFE, // þ thorn
	0x00FF: 0xFF, // ÿ ydieresis
}

// pdfDocEncoding maps Unicode characters to PDFDocEncoding codes
// (non-ASCII part). Used for document information strings.
var pdfDocEncoding = map[rune]byte{
	0x00A0: 0x20, // no-break space -> space
	0x00AD: 0x2D, // soft hyphen -> hyphen
	0x02D8: 0x18, // ˘ breve
	0x02C7: 0x19, // ˇ caron
	0x02C6: 0x1A, // ˆ circumflex
	0x02D9: 0x1B, // ˙ dotaccent
	0x02DD: 0x1C, // ˝ hungarumlaut
	0x02DB: 0x1D, // ˛ ogonek
	0x02DA: 0x1E, // ˚ ring
	0x02DC: 0x1F, // ˜ tilde
	0x2022: 0x80, // • bullet
	0x2020: 0x81, // † dagger
	0x2021: 0x82, // ‡ daggerdbl
	0x2026: 0x83, // … ellipsis
	0x2014: 0x84, // — emdash
	0x2013: 0x85, // – endash
	0x0192: 0x86, // ƒ florin
	0x2044: 0x87, // ⁄ fraction
	0x2039: 0x88, // ‹ guilsinglleft
	0x203A: 0x89, // › guilsinglright
	0x2212: 0x8A, // − minus
	0x2030: 0x8B, // ‰ perthousand
	0x201E: 0x8C, // „ quotedblbase
	0x201C: 0x8D, // “ quotedblleft
	0x201D: 0x8E, // ” quotedblright
	0x2018: 0x8F, // ‘ quoteleft
	0x2019: 0x90, // ’ quoteright
	0x201A: 0x91, // ‚ quotesinglbase
	0x2122: 0x92, // ™ trademark
	0xFB01: 0x93, // ﬁ fi
	0xFB02: 0x94, // ﬂ fl
	0x0141: 0x95, // Ł Lslash
	0x0152: 0x96, // Œ OE
	0x0160: 0x97, // Š Scaron
	0x0178: 0x98, // Ÿ Ydieresis
	0x017D: 0x99, // Ž Zcaron
	0x0131: 0x9A, // ı dotlessi
	0x0142: 0x9B, // ł lslash
	0x0153: 0x9C, // œ oe
	0x0161: 0x9D, // š scaron
	0x017E: 0x9E, // ž zcaron
	0x20AC: 0xA0, // € Euro
	0x00A1: 0xA1, // ¡ exclamdown
	0x00A2: 0xA2, // ¢ cent
	0x00A3: 0xA3, // £ sterling
	0x00A4: 0xA4, // ¤ currency
	0x00A5: 0xA5, // ¥ yen
	0x00A6: 0xA6, // ¦ brokenbar
	0x00A7: 0xA7, // § section
	0x00A8: 0xA8, // ¨ dieresis
	0x00A9: 0xA9, // © copyright
	0x00AA: 0xAA, // ª ordfeminine
	0x00AB: 0xAB, // « guillemotleft
	0x00AC: 0xAC, // ¬ logicalnot
	0x00AE: 0xAE, // ® registered
	0x00AF: 0xAF, // ¯ macron
	0x00B0: 0xB0, // ° degree
	0x00B1: 0xB1, // ± plusminus
	0x00B2: 0xB2, // ² twosuperior
	0x00B3: 0xB3, // ³ threesuperior
	0x00B4: 0xB4, // ´ acute
	0x00B5: 0xB5, // µ mu
	0x00B6: 0xB6, // ¶ paragraph
	0x00B7: 0xB7, // · periodcentered
	0x00B8: 0xB8, // ¸ cedilla
	0x00B9: 0xB9, // ¹ onesuperior
	0x00BA: 0xBA, // º ordmasculine
	0x00BB: 0xBB, // » guillemotright
	0x00BC: 0xBC, // ¼ onequarter
	0x00BD: 0xBD, // ½ onehalf
	0x00BE: 0xBE, // ¾ threequarters
	0x00BF: 0xBF, // ¿ questiondown
	0x00C0: 0xC0, // À Agrave
	0x00C1: 0xC1, // Á Aacute
	0x00C2: 0xC2, // Â Acircumflex
	0x00C3: 0xC3, // Ã Atilde
	0x00C4: 0xC4, // Ä Adieresis
	0x00C5: 0xC5, // Å Aring
	0x00C6: 0xC6, // Æ AE
	0x00C7: 0xC7, // Ç Ccedilla
	0x00C8: 0xC8, // È Egrave
	0x00C9: 0xC9, // É Eacute
	0x00CA: 0xCA, // Ê Ecircumflex
	0x00CB: 0xCB, // Ë Edieresis
	0x00CC: 0xCC, // Ì Igrave
	0x00CD: 0xCD, // Í Iacute
	0x00CE: 0xCE, // Î Icircumflex
	0x00CF: 0xCF, // Ï Idieresis
	0x00D0: 0xD0, // Ð Eth
	0x00D1: 0xD1, // Ñ Ntilde
	0x00D2: 0xD2, // Ò Ograve
	0x00D3: 0xD3, // Ó Oacute
	0x00D4: 0xD4, // Ô Ocircumflex
	0x00D5: 0xD5, // Õ Otilde
	0x00D6: 0xD6, // Ö Odieresis
	0x00D7: 0xD7, // × multiply
	0x00D8: 0xD8, // Ø Oslash
	0x00D9: 0xD9, // Ù Ugrave
	0x00DA: 0xDA, // Ú Uacute
	0x00DB: 0xDB, // Û Ucircumflex
	0x00DC: 0xDC, // Ü Udieresis
	0x00DD: 0xDD, // Ý Yacute
	0x00DE: 0xDE, // Þ Thorn
	0x00DF: 0xDF, // ß germandbls
	0x00E0: 0xE0, // à agrave
	0x00E1: 0xE1, // á aacute
	0x00E2: 0xE2, // â acircumflex
	0x00E3: 0xE3, // ã atilde
	0x00E4: 0xE4, // ä adieresis
	0x00E5: 0xE5, // å aring
	0x00E6: 0xE6, // æ ae
	0x00E7: 0xE7, // ç ccedilla
	0x00E8: 0xE8, // è egrave
	0x00E9: 0xE9, // é eacute
	0x00EA: 0xEA, // ê ecircumflex
	0x00EB: 0xEB, // ë edieresis
	0x00EC: 0xEC, // ì igrave
	0x00ED: 0xED, // í iacute
	0x00EE: 0xEE, // î icircumflex
	0x00EF: 0xEF, // ï idieresis
	0x00F0: 0xF0, // ð eth
	0x00F1: 0xF1, // ñ ntilde
	0x00F2: 0xF2, // ò ograve
	0x00F3: 0xF3, // ó oacute
	0x00F4: 0xF4, // ô ocircumflex
	0x00F5: 0xF5, // õ otilde
	0x00F6: 0xF6, // ö odieresis
	0x00F7: 0xF7, // ÷ divide
	0x00F8: 0xF8, // ø oslash
	0x00F9: 0xF9, // ù ugrave
	0x00FA: 0xFA, // ú uacute
	0x00FB: 0xFB, // û ucircumflex
	0x00FC: 0xFC, // ü udieresis
	0x00FD: 0xFD, // ý yacute
	0x00FE: 0xFE, // þ thorn
	0x00FF: 0xFF, // ÿ ydieresis
}

//end
