package main

import (
	"fmt"
	"image/color"
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type BenchmarkResult struct {
	Name       string
	BuildTime  float64 // seconds
	MemoryMB   int
	QueryTime  float64 // microseconds
}

func main() {
	// Results from the 10GB benchmark
	results := []BenchmarkResult{
		{
			Name:      "BK-Tree (1M items)",
			BuildTime: 51.25,
			MemoryMB:  399,
			QueryTime: 114.472, // 114472 ns = 114.472 μs
		},
	}

	// Create build time and memory plot
	createBuildMemoryPlot(results)

	// Create query performance plot
	createQueryPerformancePlot(results)

	// Create comparison with theoretical limits
	createScalingPlot()
}

func createBuildMemoryPlot(results []BenchmarkResult) {
	p := plot.New()
	p.Title.Text = "10GB Dataset: Index Build Time and Memory Usage"
	p.X.Label.Text = "Index Type"

	// Create two Y axes
	p.Y.Label.Text = "Build Time (seconds)"
	
	// Build time bars
	buildTimes := make(plotter.Values, len(results))
	for i, r := range results {
		buildTimes[i] = r.BuildTime
	}

	buildBars, err := plotter.NewBarChart(buildTimes, vg.Points(40))
	if err != nil {
		log.Fatal(err)
	}
	buildBars.Color = color.RGBA{R: 77, G: 175, B: 74, A: 255}

	p.Add(buildBars)
	p.Legend.Add("Build Time", buildBars)

	// Set X axis labels
	labels := make([]string, len(results))
	for i, r := range results {
		labels[i] = r.Name
	}
	p.NominalX(labels...)

	// Add value labels
	for i, r := range results {
		label, _ := plotter.NewLabels(plotter.XYLabels{
			XYs:    []plotter.XY{{X: float64(i), Y: r.BuildTime + 2}},
			Labels: []string{fmt.Sprintf("%.1fs", r.BuildTime)},
		})
		p.Add(label)

		// Memory label
		memLabel, _ := plotter.NewLabels(plotter.XYLabels{
			XYs:    []plotter.XY{{X: float64(i), Y: r.BuildTime / 2}},
			Labels: []string{fmt.Sprintf("%d MB", r.MemoryMB)},
		})
		p.Add(memLabel)
	}

	if err := p.Save(10*vg.Inch, 6*vg.Inch, "10gb_build_memory.png"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Plot saved as 10gb_build_memory.png")
}

func createQueryPerformancePlot(results []BenchmarkResult) {
	p := plot.New()
	p.Title.Text = "10GB Dataset: Query Performance"
	p.Y.Label.Text = "Query Time (microseconds)"
	p.X.Label.Text = "Index Type"
	p.Add(plotter.NewGrid())

	queryTimes := make(plotter.Values, len(results))
	for i, r := range results {
		queryTimes[i] = r.QueryTime
	}

	bars, err := plotter.NewBarChart(queryTimes, vg.Points(40))
	if err != nil {
		log.Fatal(err)
	}
	bars.Color = color.RGBA{R: 255, G: 82, B: 82, A: 255}

	p.Add(bars)

	// Set X axis labels
	labels := make([]string, len(results))
	for i, r := range results {
		labels[i] = r.Name
	}
	p.NominalX(labels...)

	// Add value labels
	for i, r := range results {
		label, _ := plotter.NewLabels(plotter.XYLabels{
			XYs:    []plotter.XY{{X: float64(i), Y: r.QueryTime + 5}},
			Labels: []string{fmt.Sprintf("%.1f μs", r.QueryTime)},
		})
		p.Add(label)
	}

	if err := p.Save(10*vg.Inch, 6*vg.Inch, "10gb_query_performance.png"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Plot saved as 10gb_query_performance.png")
}

func createScalingPlot() {
	p := plot.New()
	p.Title.Text = "BK-Tree Performance Scaling"
	p.X.Label.Text = "Dataset Size (millions of items)"
	p.Y.Label.Text = "Query Time (microseconds)"
	p.Add(plotter.NewGrid())

	// Data points from various benchmarks
	sizes := []float64{0.01, 0.1, 1.0}  // 10K, 100K, 1M items
	times := []float64{2.9, 29.0, 114.472}  // Approximate values

	pts := make(plotter.XYs, len(sizes))
	for i := range pts {
		pts[i].X = sizes[i]
		pts[i].Y = times[i]
	}

	// Create the line
	line, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatal(err)
	}
	line.Color = color.RGBA{R: 77, G: 175, B: 74, A: 255}
	line.Width = vg.Points(2)

	// Create scatter points
	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatal(err)
	}
	scatter.GlyphStyle.Color = color.RGBA{R: 77, G: 175, B: 74, A: 255}
	scatter.GlyphStyle.Radius = vg.Points(4)

	p.Add(line, scatter)
	p.Legend.Add("BK-Tree", line)

	// Add theoretical O(log n) line for comparison
	theoreticalPts := make(plotter.XYs, 50)
	for i := range theoreticalPts {
		x := float64(i+1) * 0.02
		theoreticalPts[i].X = x
		theoreticalPts[i].Y = 10 * (1 + 2*x) // Simplified scaling
	}

	theoryLine, _ := plotter.NewLine(theoreticalPts)
	theoryLine.Color = color.RGBA{R: 255, G: 152, B: 0, A: 255}
	theoryLine.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	
	p.Add(theoryLine)
	p.Legend.Add("Linear Scaling", theoryLine)

	if err := p.Save(10*vg.Inch, 6*vg.Inch, "bktree_scaling.png"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Plot saved as bktree_scaling.png")
}