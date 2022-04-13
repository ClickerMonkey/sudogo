package sudogo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

type PuzzlePDF struct {
	Puzzles          []PuzzlePDFItem
	PuzzlesWide      int
	PuzzlesHigh      int
	PuzzleSpacing    float64
	MarginLeft       float64
	MarginRight      float64
	MarginTop        float64
	MarginBottom     float64
	Landscape        bool
	PageSize         string
	Font             string
	ValueFontColor   Color
	ValueBackColor   Color
	ValueFontScale   float64
	CandidateColor   Color
	CandidateScale   float64
	BorderThickColor Color
	BorderThickWidth float64
	BorderThinColor  Color
	BorderThinWidth  float64
}

type PuzzlePDFItem struct {
	Puzzle     *Puzzle
	Candidates bool
}

type Color struct {
	R int
	G int
	B int
}

func (c Color) Is(other Color) bool {
	return c.R == other.R && c.G == other.G && c.B == other.B
}

var (
	ColorWhite = Color{255, 255, 255}
	ColorBlack = Color{0, 0, 0}
	ColorGreen = Color{0, 255, 0}
	ColorRed   = Color{255, 0, 0}
	ColorBlue  = Color{0, 0, 255}
	ColorGray  = Color{127, 127, 127}
)

func NewPDF() PuzzlePDF {
	return PuzzlePDF{
		PageSize:         "A4",
		MarginBottom:     28.35,
		MarginTop:        28.35,
		MarginLeft:       28.35,
		MarginRight:      28.35,
		PuzzlesWide:      1,
		PuzzlesHigh:      1,
		PuzzleSpacing:    28.35,
		Font:             "Arial",
		CandidateColor:   ColorGray,
		CandidateScale:   0.2,
		ValueFontColor:   ColorBlack,
		ValueFontScale:   0.7,
		ValueBackColor:   ColorWhite,
		BorderThickColor: ColorBlack,
		BorderThickWidth: 2.0,
		BorderThinColor:  ColorBlack,
		BorderThinWidth:  0.5,
	}
}

func (pdf *PuzzlePDF) Add(puzzle *Puzzle, candidates bool) {
	pdf.Puzzles = append(pdf.Puzzles, PuzzlePDFItem{
		Puzzle:     puzzle,
		Candidates: candidates,
	})
}

func (pdf *PuzzlePDF) orientation() string {
	if pdf.Landscape {
		return "L"
	} else {
		return "P"
	}
}

func (pdf *PuzzlePDF) WriteFile(path string) {
	f := pdf.Generate()
	f.OutputFileAndClose(path)
}

func (pdf *PuzzlePDF) Write(writer io.Writer, direct bool) {
	f := pdf.Generate()
	if direct {
		f.Output(writer)
	} else {
		buf := bytes.Buffer{}
		f.Output(&buf)
		writer.Write(buf.Bytes())
	}
	f.Close()
}

func (pdf *PuzzlePDF) Send(w http.ResponseWriter, direct bool) (any, int) {
	w.Header().Set("Content-Type", "application/pdf")
	pdf.Write(w, direct)

	return nil, -1
}

func (pdf *PuzzlePDF) Download(w http.ResponseWriter, filename string, direct bool) (any, int) {
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	pdf.Send(w, direct)

	return nil, -1
}

func (pdf *PuzzlePDF) Generate() *gofpdf.Fpdf {
	p := gofpdf.New(pdf.orientation(), "pt", pdf.PageSize, "")
	p.SetMargins(pdf.MarginLeft, pdf.MarginTop, pdf.MarginRight)
	p.SetAutoPageBreak(false, pdf.MarginBottom)

	pageW, pageH, _ := p.PageSize(0)
	spaceW := pageW - pdf.MarginLeft - pdf.MarginRight
	spaceH := pageH - pdf.MarginTop - pdf.MarginBottom
	countW := float64(pdf.PuzzlesWide)
	countH := float64(pdf.PuzzlesHigh)
	spacingW := (countW - 1) * pdf.PuzzleSpacing
	spacingH := (countH - 1) * pdf.PuzzleSpacing
	puzzleSpaceW := (spaceW - spacingW) / countW
	puzzleSpaceH := (spaceH - spacingH) / countH
	puzzleSize := Min(puzzleSpaceW, puzzleSpaceH)
	puzzleSeparation := puzzleSize + pdf.PuzzleSpacing
	offsetX := pdf.MarginLeft + (spaceW-(puzzleSize*countW+spacingW))/2
	offsetY := pdf.MarginTop + (spaceH-(puzzleSize*countH+spacingH))/2
	perPage := pdf.PuzzlesWide * pdf.PuzzlesHigh

	for i, item := range pdf.Puzzles {
		puzzle := item.Puzzle
		pageIndex := i % perPage
		size := puzzle.Kind.Size()
		sizef := float64(size)
		cellSize := puzzleSize / sizef
		fontSize := cellSize * pdf.ValueFontScale
		candidateSize := cellSize * pdf.CandidateScale
		pageCol := pageIndex % pdf.PuzzlesWide
		pageRow := pageIndex / pdf.PuzzlesWide
		originX := offsetX + float64(pageCol)*puzzleSeparation
		originY := offsetY + float64(pageRow)*puzzleSeparation
		boxW := float64(puzzle.Kind.BoxSize.Width) * cellSize
		boxH := float64(puzzle.Kind.BoxSize.Height) * cellSize

		if pageIndex == 0 {
			p.AddPage()
		}

		p.SetLineWidth(pdf.BorderThinWidth)
		p.SetDrawColor(pdf.BorderThinColor.R, pdf.BorderThinColor.G, pdf.BorderThinColor.B)
		p.SetFillColor(pdf.ValueBackColor.R, pdf.ValueBackColor.G, pdf.ValueBackColor.B)

		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				cell := puzzle.Get(x, y)

				p.SetXY(originX+float64(x)*cellSize, originY+float64(y)*cellSize)
				if cell.HasValue() {
					cellValue := strconv.Itoa(cell.Value)
					fill := !pdf.ValueBackColor.Is(ColorWhite)

					p.SetFont(pdf.Font, "B", fontSize)
					p.SetTextColor(pdf.ValueFontColor.R, pdf.ValueFontColor.G, pdf.ValueFontColor.B)
					p.CellFormat(cellSize, cellSize, cellValue, "1", 0, "CM", fill, 0, "")
				} else {
					p.SetFont(pdf.Font, "", fontSize)
					p.CellFormat(cellSize, cellSize, "", "1", 0, "CM", false, 0, "")

					if item.Candidates {
						cand := fmt.Sprintf("%v", cell.Candidates())
						cand = strings.Trim(cand, "[]")

						p.SetFont(pdf.Font, "", candidateSize)
						p.SetTextColor(pdf.CandidateColor.R, pdf.CandidateColor.G, pdf.CandidateColor.B)
						p.SetXY(originX+float64(x)*cellSize, originY+float64(y)*cellSize+pdf.BorderThickWidth*2)
						p.CellFormat(cellSize, cellSize, cand, "0", 0, "CT", false, 0, "")
					}
				}
			}
		}

		p.SetLineWidth(pdf.BorderThickWidth)
		p.SetLineCapStyle("round")

		for y := 0.0; y <= puzzleSize+0.0001; y += boxH {
			p.Line(originX, originY+y, originX+puzzleSize, originY+y)
		}

		for x := 0.0; x <= puzzleSize+0.0001; x += boxW {
			p.Line(originX+x, originY, originX+x, originY+puzzleSize)
		}
	}

	return p
}
