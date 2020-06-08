package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

const (
	pagew  = 210.
	pageh  = 297.
	margin = 5.
	innerw = pagew - (margin * 2)

	headerh = margin + 22.
	footerh = margin + 15.

	fontSize = 12.
	lineHt   = 4.1
)

const (
	contact  = "(814) 977-7556\njon@calhoun.io\ngophercises.com"
	address  = "123 Fake St\nSome Town, PA\n12345"
	billTo   = "Client Name\n1 Client Address\nCity, State, Country\n98476"
	pdfName  = "hello.pdf"
	dataFile = "invoice-data.json"
)

// Invoice holds the parts that make up an invoice.
type Invoice struct {
	pdf     *gofpdf.Fpdf
	data    InvoiceData
	details InvoiceDetails
}

// InvoiceData holds the data items that make up an invoice.
type InvoiceData []struct {
	UnitName       string  `json: "UnitName"`
	PricePerUnit   float64 `json:"PricePerUnit"`
	UnitsPurchased float64 `json:"UnitsPurchased"`
}

// InvoiceDetails holds the info details that make up an invoice.
type InvoiceDetails struct {
	contact       []string
	address       []string
	billTo        []string
	invoiceNumber string
	invoiceTotal  string
	issueDate     string
}

func main() {
	var err error

	var bytes []byte
	if bytes, err = ioutil.ReadFile(dataFile); err != nil {
		log.Fatal(err)
	}

	var invdt InvoiceData
	if err = json.Unmarshal(bytes, &invdt); err != nil {
		log.Fatal(err)
	}

	inv := New(invdt, contact, address, billTo)

	if err = inv.Generate(pdfName); err != nil {
		log.Fatal(err)
	}
}

// New creates a new instance of an invoice.
func New(
	d InvoiceData,
	contact, address, billTo string,
) *Invoice {
	var inv Invoice

	inv.data = d

	inv.details = InvoiceDetails{
		contact:       strings.Split(contact, "\n"),
		address:       strings.Split(address, "\n"),
		billTo:        strings.Split(billTo, "\n"),
		invoiceNumber: fmt.Sprintf("%v", 777000123),
		invoiceTotal:  fmt.Sprintf("%v", calculateSubtotal(d)),
		issueDate:     time.Now().Format("1/2/2006"),
	}

	inv.pdf = gofpdf.New("P", "mm", "A4", "")
	inv.pdf.SetFont("Arial", "", fontSize)
	inv.pdf.SetMargins(margin, margin, margin)

	return &inv
}

// Generate builds the invoice.
func (inv *Invoice) Generate(fileName string) error {
	log.Printf("generating invoice %v issued on %v", inv.details.invoiceNumber, inv.details.issueDate)

	inv.pdf.SetHeaderFuncMode(func() {
		inv.pdf.SetFillColor(255, 127, 0)
		inv.pdf.Rect(0, 0, pagew, headerh, "F")

		// Page title.
		inv.pdf.SetFont("Arial", "B", 28)
		inv.pdf.Text(15, 15, "INVOICE")

		// Contact details.
		inv.pdf.SetXY(pagew/2, margin)
		inv.writeBlockOfCells("", inv.details.contact, pagew/4, lineHt)

		inv.pdf.SetXY((innerw/4)*3, margin)
		inv.writeBlockOfCells("", inv.details.address, pagew/4, lineHt)

	}, false)

	inv.pdf.SetFooterFunc(func() {
		inv.pdf.SetFillColor(255, 127, 0)
		inv.pdf.Rect(0, pageh-footerh, pagew, footerh, "F")
	})

	inv.pdf.AddPage()
	inv.pdf.SetY(headerh + margin)

	// Bill to
	inv.writeBlockOfCells(
		"Bill to", inv.details.billTo, innerw/3, lineHt+2,
	)

	// Store the greatest Y value as the 3 columns
	// are created so the table can be created below.
	maxY := inv.pdf.GetY()

	// Invoice number
	inv.pdf.SetXY(innerw/3, headerh+margin)
	inv.writeBlockOfCells(
		"Invoice number", []string{inv.details.invoiceNumber}, innerw/3, lineHt+2,
	)

	// Date of issue
	inv.pdf.SetXY(innerw/3, inv.pdf.GetY()+5.)
	inv.writeBlockOfCells(
		"Date of issue", []string{inv.details.issueDate}, innerw/3, lineHt+2,
	)

	// End of col, so check for a height larger than current max height.
	if inv.pdf.GetY() > maxY {
		maxY = inv.pdf.GetY()
	}

	// Invoice Total
	inv.pdf.SetXY((innerw/3)*2, headerh+margin)
	inv.writeBlockOfCells(
		"Invoice total", []string{inv.details.invoiceTotal}, innerw/3, lineHt+2,
	)

	// End of col, so check for a height larger than current max height.
	if inv.pdf.GetY() > maxY {
		maxY = inv.pdf.GetY()
	}

	log.Printf("generating line item table for invoice %v with inner width %.2f", inv.details.invoiceNumber, innerw)

	// Table should begin below the longest column.
	inv.pdf.SetY(maxY + 5)

	// Set the column widths and set the first row of the table, the headings.
	halfcol, sixthcol := innerw/2, innerw/6
	colWidths := []float64{halfcol, sixthcol, sixthcol, sixthcol}

	rows := [][]string{}
	rows = append(rows, []string{"Description", "Price per unit", "Quantity", "Amount"})

	// Build the table data, and calculate the subtotal as we go.
	for _, item := range inv.data {
		price := fmt.Sprintf("$%.2f", item.PricePerUnit/100.)
		numUnits := fmt.Sprintf("%.0f", item.UnitsPurchased)
		amount := fmt.Sprintf("$%.2f", item.PricePerUnit*item.UnitsPurchased/100.)

		rows = append(rows, []string{item.UnitName, price, numUnits, amount})
	}

	// Subtotal is the final row.
	rows = append(rows, []string{"", "Subtotal", "", "$" + inv.details.invoiceTotal})

	for rowIdx, row := range rows {
		curx, cury := inv.pdf.GetXY()

		x, y := inv.pdf.GetXY()

		// Get the row height by getting the maximum of all cell heights.
		rowh := lineHt + margin
		for i, txt := range row {
			lines := inv.pdf.SplitLines([]byte(txt), colWidths[i])
			cellh := (float64(len(lines)) * lineHt) + margin

			if cellh > rowh {
				rowh = cellh
			}
		}

		// Start on the next page, if necessary.
		if y+rowh+margin+footerh > pageh {
			inv.pdf.AddPage()
			x, y = inv.pdf.GetXY()
			inv.pdf.SetY(headerh + margin)
		}

		// Bold the first row.
		if rowIdx == 0 {
			inv.pdf.SetFont("Arial", "B", fontSize)
		} else {
			inv.pdf.SetFont("Arial", "", fontSize)
		}

		// Write the text into the cells.
		for i, cell := range row {
			cellw := colWidths[i]

			// SplitLines will break up the text by cell width.
			txtLns := inv.pdf.SplitLines([]byte(cell), cellw)

			inv.pdf.SetXY(x, y)
			for _, ln := range txtLns {
				inv.pdf.Cell(cellw, lineHt, string(ln))
				inv.pdf.Ln(-1)
			}
			x += cellw
		}
		inv.pdf.SetXY(curx, cury+rowh)
	}

	// Write the file contents to hello.pdf.
	var err error
	if err = inv.pdf.OutputFileAndClose(fileName); err != nil {
		return err
	}

	return nil
}

// writeBlockOfCells takes an optional title and lines and writes
// them into cells at the current coordinates.
func (inv *Invoice) writeBlockOfCells(title string, lines []string, w, h float64) {
	x := inv.pdf.GetX()

	// Write the title, if there is one.
	if title != "" {
		inv.pdf.SetFont("Arial", "B", fontSize)
		inv.pdf.Cell(w, h, title)
		inv.pdf.Ln(-1)
	}

	// Write lines at the x pos., line by line.
	inv.pdf.SetFont("Arial", "", fontSize)
	for _, ln := range lines {
		inv.pdf.SetX(x)
		inv.pdf.Cell(w, h, ln)
		inv.pdf.Ln(-1)
	}
}

// calculateSubtotal sums the values of all invoice items.
func calculateSubtotal(d InvoiceData) float64 {
	var subtotal float64
	for _, item := range d {
		subtotal += (float64(item.PricePerUnit) * float64(item.UnitsPurchased)) / 100.
	}

	return subtotal
}
