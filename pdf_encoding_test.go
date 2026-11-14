// -----------------------------------------------------------------------------
// github.com/balacode/one-file-pdf           one-file-pdf/[pdf_encoding_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package pdf

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func Test_PDF_encodeText_WinAnsi_(t *testing.T) {
	var p PDF
	for _, tc := range []struct{ in, want string }{
		{"Hello", "Hello"},
		{"räksmörgås", "r\xe4ksm\xf6rg\xe5s"},
		{"RÄKSMÖRGÅS", "R\xc4KSM\xd6RG\xc5S"},
		{"€ … – — “x” ‘y’ • ™ Œœ Šš Žž Ÿ ƒ", "\x80 \x85 \x96 \x97 \x93x\x94 \x91y\x92 \x95 \x99 \x8c\x9c \x8a\x9a \x8e\x9e \x9f \x83"},
		{"© ® ° ± µ ¹²³ ß ÿ", "\xa9 \xae \xb0 \xb1 \xb5 \xb9\xb2\xb3 \xdf \xff"},
		{"ẞ日本", "???"}, // not in WinAnsiEncoding
	} {
		if got := p.encodeText(tc.in, pdfWinAnsiEncoding); got != tc.want {
			t.Errorf("encodeText(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func Test_PDF_encodeText_PDFDoc_(t *testing.T) {
	var p PDF
	got := p.encodeText("Åsa € • ˘", pdfDocEncoding)
	if want := "\xc5sa \xa0 \x80 \x18"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func Test_PDF_escapeBytes_(t *testing.T) {
	var p PDF
	got := p.escapeBytes("a(b)c\\d\xe4\x80\n")
	if want := `a\(b\)c\\d\344\200\012`; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func Test_PDF_DrawText_8bit_(t *testing.T) {
	doc := NewPDF("A4")
	doc.SetCompression(false).SetUnits("cm").SetXY(1, 1).
		SetFont("Helvetica", 12).DrawText("räksmörgås €").
		SetFont("ZapfDingbats", 12).DrawText("a").
		SetDocTitle("Smörgåsbord €")
	out := string(doc.Bytes())
	for _, s := range []string{
		`(r\344ksm\366rg\345s \200) Tj`,
		"/BaseFont/Helvetica\n/Encoding/WinAnsiEncoding>>",
		"/BaseFont/ZapfDingbats\n/Encoding/StandardEncoding>>",
		`/Title (Sm\366rg\345sbord \240)`,
	} {
		if !strings.Contains(out, s) {
			t.Errorf("output does not contain %q:\n%s", s, out)
		}
	}
	if errs := doc.Errors(); len(errs) > 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
	// width uses the WinAnsi code: "ä" is as wide as "a" in Helvetica
	// (widths use the last font applied by a draw call)
	doc.SetFont("Helvetica", 12).DrawText(" ")
	if a, b := doc.TextWidth("ä"), doc.TextWidth("a"); a != b {
		t.Errorf("TextWidth(ä)=%v, TextWidth(a)=%v", a, b)
	}
	if w := doc.TextWidth("€"); w <= 0 {
		t.Errorf("TextWidth(€)=%v", w)
	}
}

func Test_PDF_WrapTextLines_UTF8_(t *testing.T) {
	doc := NewPDF("A4")
	doc.SetUnits("cm").SetXY(1, 1).SetFont("Helvetica", 12).DrawText(" ")
	for _, w := range []float64{0.3, 0.5, 1, 2, 3} {
		for _, line := range doc.WrapTextLines(w, "åäöåäöåäö ÅÄÖÅÄÖ €€€€ smörgås") {
			if !utf8.ValidString(line) {
				t.Errorf("width %v: line %q split inside a character", w, line)
			}
		}
	}
}

//end
