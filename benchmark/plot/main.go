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

type BenchmarkResult struct {
	Library string
	P50     float64 // in microseconds
	P90     float64
	P99     float64
}

func main() {
	results := []BenchmarkResult{
		{"sahilm/fuzzy", 1005.083, 1071.5, 1418.417},
		{"lithammer/fuzzysearch", 129.375, 169.0, 192.25},
		{"raphamorim/fuzzy (BK-Tree)", 290.625, 426.667, 1127.584},
		{"raphamorim/fuzzy (NGram)", 23.708, 48.75, 107.208},
	}

	// Create the plot
	p := plot.New()
	p.Title.Text = "Fuzzy Search Library Latency Comparison (Log Scale)"
	p.Y.Label.Text = "Latency (microseconds)"
	p.X.Label.Text = "Libraries"
	p.Add(plotter.NewGrid())

	// Prepare data for grouped bars
	w := vg.Points(15)
	colors := []color.Color{
		color.RGBA{R: 77, G: 175, B: 74, A: 255},   // Green for P50
		color.RGBA{R: 255, G: 152, B: 0, A: 255},   // Orange for P90
		color.RGBA{R: 255, G: 82, B: 82, A: 255},   // Red for P99
	}

	// Create bars for each percentile
	for i, percentile := range []string{"P50", "P90", "P99"} {
		values := make(plotter.Values, len(results))

		switch i {
		case 0: // P50
			for j, r := range results {
				values[j] = math.Log10(r.P50)
			}
		case 1: // P90
			for j, r := range results {
				values[j] = math.Log10(r.P90)
			}
		case 2: // P99
			for j, r := range results {
				values[j] = math.Log10(r.P99)
			}
		}

		bars, err := plotter.NewBarChart(values, w)
		if err != nil {
			log.Fatal(err)
		}
		bars.Color = colors[i]
		bars.Offset = vg.Length(float64(i-1)) * w

		p.Add(bars)
		p.Legend.Add(percentile, bars)
	}

	// Set X axis labels
	labels := make([]string, len(results))
	for i, r := range results {
		labels[i] = r.Library
	}
	p.NominalX(labels...)

	// Customize Y axis to show actual values
	p.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{Value: math.Log10(10), Label: "10"},
		{Value: math.Log10(100), Label: "100"},
		{Value: math.Log10(1000), Label: "1000"},
		{Value: math.Log10(10000), Label: "10000"},
		{Value: math.Log10(100000), Label: "100000"},
	})

	// Save the plot
	if err := p.Save(12*vg.Inch, 8*vg.Inch, "latency_comparison.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Plot saved as latency_comparison.png")

	// Create a linear scale version
	createLinearPlot(results)

	// Create speedup plot
	createSpeedupPlot(results)
}

func createLinearPlot(results []BenchmarkResult) {
	p := plot.New()
	p.Title.Text = "Fuzzy Search Library Latency Comparison"
	p.Y.Label.Text = "Latency (microseconds)"
	p.X.Label.Text = "Libraries"
	p.Add(plotter.NewGrid())

	// Prepare data for grouped bars
	w := vg.Points(15)
	colors := []color.Color{
		color.RGBA{R: 77, G: 175, B: 74, A: 255},   // Green for P50
		color.RGBA{R: 255, G: 152, B: 0, A: 255},   // Orange for P90
		color.RGBA{R: 255, G: 82, B: 82, A: 255},   // Red for P99
	}

	// Create bars for each percentile
	for i, percentile := range []string{"P50", "P90", "P99"} {
		values := make(plotter.Values, len(results))

		switch i {
		case 0: // P50
			for j, r := range results {
				values[j] = r.P50
			}
		case 1: // P90
			for j, r := range results {
				values[j] = r.P90
			}
		case 2: // P99
			for j, r := range results {
				values[j] = r.P99
			}
		}

		bars, err := plotter.NewBarChart(values, w)
		if err != nil {
			log.Fatal(err)
		}
		bars.Color = colors[i]
		bars.Offset = vg.Length(float64(i-1)) * w

		p.Add(bars)
		p.Legend.Add(percentile, bars)
	}

	// Set X axis labels
	labels := make([]string, len(results))
	for i, r := range results {
		labels[i] = r.Library
	}
	p.NominalX(labels...)

	// Save the plot
	if err := p.Save(12*vg.Inch, 8*vg.Inch, "latency_comparison_linear.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Plot saved as latency_comparison_linear.png")
}

func createSpeedupPlot(results []BenchmarkResult) {
	p := plot.New()
	p.Title.Text = "Performance Speedup Relative to sahilm/fuzzy"
	p.Y.Label.Text = "Speedup Factor (X times faster)"
	p.X.Label.Text = "Libraries"
	p.Add(plotter.NewGrid())

	baseline := results[0] // sahilm/fuzzy as baseline
	w := vg.Points(15)

	colors := []color.Color{
		color.RGBA{R: 77, G: 175, B: 74, A: 255},   // Green for P50
		color.RGBA{R: 255, G: 152, B: 0, A: 255},   // Orange for P90
		color.RGBA{R: 255, G: 82, B: 82, A: 255},   // Red for P99
	}

	for i, percentile := range []string{"P50", "P90", "P99"} {
		values := make(plotter.Values, len(results))

		switch i {
		case 0: // P50
			for j, r := range results {
				values[j] = baseline.P50 / r.P50
			}
		case 1: // P90
			for j, r := range results {
				values[j] = baseline.P90 / r.P90
			}
		case 2: // P99
			for j, r := range results {
				values[j] = baseline.P99 / r.P99
			}
		}

		bars, err := plotter.NewBarChart(values, w)
		if err != nil {
			log.Fatal(err)
		}
		bars.Color = colors[i]
		bars.Offset = vg.Length(float64(i-1)) * w

		p.Add(bars)
		p.Legend.Add(percentile, bars)
	}

	// Set X axis labels
	labels := make([]string, len(results))
	for i, r := range results {
		labels[i] = r.Library
	}
	p.NominalX(labels...)

	// Add value labels on top of bars
	for i := range results {
		x := float64(i)

		// P50 speedup
		p50Speed := baseline.P50 / results[i].P50
		if p50Speed > 1.5 {
			label, _ := plotter.NewLabels(plotter.XYLabels{
				XYs: []plotter.XY{{X: x - 0.3, Y: p50Speed + 0.5}},
				Labels: []string{fmt.Sprintf("%.1fx", p50Speed)},
			})
			p.Add(label)
		}
	}

	// Save the plot
	if err := p.Save(12*vg.Inch, 8*vg.Inch, "speedup_comparison.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Plot saved as speedup_comparison.png")
}
