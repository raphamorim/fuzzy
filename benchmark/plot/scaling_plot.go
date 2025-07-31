package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	// Create scaling performance plot
	createScalingPlot()
}

func createScalingPlot() {
	p := plot.New()
	p.Title.Text = "Fuzzy Search Library Scaling Performance"
	p.X.Label.Text = "Dataset Size"
	p.Y.Label.Text = "Time per Operation (ns)"
	p.Y.Scale = plot.LogScale{}
	p.X.Scale = plot.LogScale{}
	p.Add(plotter.NewGrid())

	// Data from benchmark results (in nanoseconds)
	datasetSizes := []float64{100, 1000, 10000, 100000}
	
	// sahilm/fuzzy results
	sahilmTimes := []float64{7701, 97028, 1018647, 10352907}
	
	// lithammer/fuzzysearch results
	lithammerTimes := []float64{823.7, 10034, 123716, 1498196}
	
	// raphamorim/fuzzy BK-Tree results (new optimized version)
	raphamorimBKTimes := []float64{1003, 4975, 54623, 538290}
	
	// raphamorim/fuzzy NGram results (new optimized version)
	raphamorimNGramTimes := []float64{265.7, 495.0, 7581, 55022}

	// Create line plots
	sahilmPts := make(plotter.XYs, len(datasetSizes))
	lithammerPts := make(plotter.XYs, len(datasetSizes))
	raphamorimBKPts := make(plotter.XYs, len(datasetSizes))
	raphamorimNGramPts := make(plotter.XYs, len(datasetSizes))

	for i := range datasetSizes {
		sahilmPts[i].X = datasetSizes[i]
		sahilmPts[i].Y = sahilmTimes[i]
		
		lithammerPts[i].X = datasetSizes[i]
		lithammerPts[i].Y = lithammerTimes[i]
		
		raphamorimBKPts[i].X = datasetSizes[i]
		raphamorimBKPts[i].Y = raphamorimBKTimes[i]
		
		raphamorimNGramPts[i].X = datasetSizes[i]
		raphamorimNGramPts[i].Y = raphamorimNGramTimes[i]
	}

	// Add lines with different colors and styles
	sahilmLine, err := plotter.NewLine(sahilmPts)
	if err != nil {
		log.Fatal(err)
	}
	sahilmLine.Color = color.RGBA{R: 255, G: 87, B: 34, A: 255}
	sahilmLine.Width = vg.Points(2)
	sahilmLine.Dashes = []vg.Length{}

	lithammerLine, err := plotter.NewLine(lithammerPts)
	if err != nil {
		log.Fatal(err)
	}
	lithammerLine.Color = color.RGBA{R: 33, G: 150, B: 243, A: 255}
	lithammerLine.Width = vg.Points(2)

	raphamorimBKLine, err := plotter.NewLine(raphamorimBKPts)
	if err != nil {
		log.Fatal(err)
	}
	raphamorimBKLine.Color = color.RGBA{R: 76, G: 175, B: 80, A: 255}
	raphamorimBKLine.Width = vg.Points(2)

	raphamorimNGramLine, err := plotter.NewLine(raphamorimNGramPts)
	if err != nil {
		log.Fatal(err)
	}
	raphamorimNGramLine.Color = color.RGBA{R: 156, G: 39, B: 176, A: 255}
	raphamorimNGramLine.Width = vg.Points(2)
	raphamorimNGramLine.Dashes = []vg.Length{vg.Points(5), vg.Points(2)}

	// Add points
	sahilmPoints, err := plotter.NewScatter(sahilmPts)
	if err != nil {
		log.Fatal(err)
	}
	sahilmPoints.GlyphStyle.Color = sahilmLine.Color
	sahilmPoints.GlyphStyle.Radius = vg.Points(4)

	lithammerPoints, err := plotter.NewScatter(lithammerPts)
	if err != nil {
		log.Fatal(err)
	}
	lithammerPoints.GlyphStyle.Color = lithammerLine.Color
	lithammerPoints.GlyphStyle.Radius = vg.Points(4)

	raphamorimBKPoints, err := plotter.NewScatter(raphamorimBKPts)
	if err != nil {
		log.Fatal(err)
	}
	raphamorimBKPoints.GlyphStyle.Color = raphamorimBKLine.Color
	raphamorimBKPoints.GlyphStyle.Radius = vg.Points(4)

	raphamorimNGramPoints, err := plotter.NewScatter(raphamorimNGramPts)
	if err != nil {
		log.Fatal(err)
	}
	raphamorimNGramPoints.GlyphStyle.Color = raphamorimNGramLine.Color
	raphamorimNGramPoints.GlyphStyle.Radius = vg.Points(4)

	// Add to plot
	p.Add(sahilmLine, sahilmPoints)
	p.Add(lithammerLine, lithammerPoints)
	p.Add(raphamorimBKLine, raphamorimBKPoints)
	p.Add(raphamorimNGramLine, raphamorimNGramPoints)

	// Add legend
	p.Legend.Add("sahilm/fuzzy", sahilmLine)
	p.Legend.Add("lithammer/fuzzysearch", lithammerLine)
	p.Legend.Add("raphamorim/fuzzy (BK-Tree)", raphamorimBKLine)
	p.Legend.Add("raphamorim/fuzzy (NGram)", raphamorimNGramLine)
	p.Legend.Top = true
	p.Legend.Left = true

	// Custom tick marks for both axes
	p.X.Tick.Marker = plot.LogTicks{Prec: 0}
	p.Y.Tick.Marker = plot.LogTicks{Prec: 0}

	// Save the plot
	if err := p.Save(12*vg.Inch, 8*vg.Inch, "scaling_performance.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Scaling plot saved as scaling_performance.png")

	// Create speedup comparison plot
	createSpeedupComparisonPlot(datasetSizes, sahilmTimes, lithammerTimes, raphamorimBKTimes, raphamorimNGramTimes)
}

func createSpeedupComparisonPlot(sizes, sahilm, lithammer, raphamorimBK, raphamorimNGram []float64) {
	p := plot.New()
	p.Title.Text = "Performance Comparison Across Dataset Sizes"
	p.X.Label.Text = "Dataset Size"
	p.Y.Label.Text = "Speedup vs sahilm/fuzzy"
	p.Add(plotter.NewGrid())

	// Calculate speedups
	lithammerSpeedup := make(plotter.XYs, len(sizes))
	raphamorimBKSpeedup := make(plotter.XYs, len(sizes))
	raphamorimNGramSpeedup := make(plotter.XYs, len(sizes))

	for i := range sizes {
		lithammerSpeedup[i].X = math.Log10(sizes[i])
		lithammerSpeedup[i].Y = sahilm[i] / lithammer[i]
		
		raphamorimBKSpeedup[i].X = math.Log10(sizes[i])
		raphamorimBKSpeedup[i].Y = sahilm[i] / raphamorimBK[i]
		
		raphamorimNGramSpeedup[i].X = math.Log10(sizes[i])
		raphamorimNGramSpeedup[i].Y = sahilm[i] / raphamorimNGram[i]
	}

	// Create bar chart
	w := vg.Points(15)
	
	lithammerBars, err := plotter.NewBarChart(plotter.Values{
		lithammerSpeedup[0].Y,
		lithammerSpeedup[1].Y,
		lithammerSpeedup[2].Y,
		lithammerSpeedup[3].Y,
	}, w)
	if err != nil {
		log.Fatal(err)
	}
	lithammerBars.Color = color.RGBA{R: 33, G: 150, B: 243, A: 255}
	lithammerBars.Offset = -w

	raphamorimBKBars, err := plotter.NewBarChart(plotter.Values{
		raphamorimBKSpeedup[0].Y,
		raphamorimBKSpeedup[1].Y,
		raphamorimBKSpeedup[2].Y,
		raphamorimBKSpeedup[3].Y,
	}, w)
	if err != nil {
		log.Fatal(err)
	}
	raphamorimBKBars.Color = color.RGBA{R: 76, G: 175, B: 80, A: 255}
	raphamorimBKBars.Offset = 0

	raphamorimNGramBars, err := plotter.NewBarChart(plotter.Values{
		raphamorimNGramSpeedup[0].Y,
		raphamorimNGramSpeedup[1].Y,
		raphamorimNGramSpeedup[2].Y,
		raphamorimNGramSpeedup[3].Y,
	}, w)
	if err != nil {
		log.Fatal(err)
	}
	raphamorimNGramBars.Color = color.RGBA{R: 156, G: 39, B: 176, A: 255}
	raphamorimNGramBars.Offset = w

	p.Add(lithammerBars, raphamorimBKBars, raphamorimNGramBars)
	p.Legend.Add("lithammer/fuzzysearch", lithammerBars)
	p.Legend.Add("raphamorim/fuzzy (BK-Tree)", raphamorimBKBars)
	p.Legend.Add("raphamorim/fuzzy (NGram)", raphamorimNGramBars)

	// Set X axis labels
	p.NominalX("100", "1K", "10K", "100K")

	// Add value labels on bars
	for i := range sizes {
		x := float64(i)
		
		// Lithammer speedup label
		label1, _ := plotter.NewLabels(plotter.XYLabels{
			XYs: []plotter.XY{{X: x - 0.2, Y: lithammerSpeedup[i].Y + 0.5}},
			Labels: []string{fmt.Sprintf("%.1fx", lithammerSpeedup[i].Y)},
		})
		p.Add(label1)
		
		// Raphamorim BK speedup label
		label2, _ := plotter.NewLabels(plotter.XYLabels{
			XYs: []plotter.XY{{X: x, Y: raphamorimBKSpeedup[i].Y + 0.5}},
			Labels: []string{fmt.Sprintf("%.1fx", raphamorimBKSpeedup[i].Y)},
		})
		p.Add(label2)
		
		// Raphamorim NGram speedup label
		label3, _ := plotter.NewLabels(plotter.XYLabels{
			XYs: []plotter.XY{{X: x + 0.2, Y: raphamorimNGramSpeedup[i].Y + 0.5}},
			Labels: []string{fmt.Sprintf("%.1fx", raphamorimNGramSpeedup[i].Y)},
		})
		p.Add(label3)
	}

	// Save the plot
	if err := p.Save(12*vg.Inch, 8*vg.Inch, "speedup_by_size.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Speedup comparison plot saved as speedup_by_size.png")
}