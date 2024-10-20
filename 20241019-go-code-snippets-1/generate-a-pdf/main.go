package main

import (
	"github.com/jung-kurt/gofpdf"
)

// https://www.golangcode.com/generate-a-pdf/
func main() {
	err := GeneratePdf("hello.pdf")
	if err != nil {
		panic(err)
	}
}

// GeneratePdf generates our pdf by adding text and images to the page
// then saving it to a file (name specified in params).
func GeneratePdf(filename string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// CellFormat(width, height, text, border, position after, align, fill, link, linkStr)
	pdf.CellFormat(190, 7, "Welcome to golangcode.com", "0", 0, "CM", false, 0, "")

	// ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
	// pdf.ImageOptions(
	// 	"avatar.jpg",
	// 	80, 20,
	// 	0, 0,
	// 	false,
	// 	gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true},
	// 	0,
	// 	"",
	// )

	return pdf.OutputFileAndClose(filename)
}
