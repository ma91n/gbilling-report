package main

import (
	"log"
	"os"
	"bytes"

	"github.com/wcharczuk/go-chart"
	"github.com/laqiiz/gbilling-report/billing"
	"github.com/laqiiz/gbilling-report/storage"
	"strconv"
	"time"
)

func main() {
	log.Print("start gbiling-report")

	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	if len(projectID) == 0 {
		log.Fatal("env GOOGLE_PROJECT_ID is not found")
	}

	billingTable := os.Getenv("BILLING_TABLE")
	if len(billingTable) == 0 {
		log.Fatal("env BILLING_TABLE is not found")
	}

	billings, err := billing.New(projectID, billingTable)
	if err != nil {
		log.Fatal(err)
	}

	costs, err := billings.FetchCost(30)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("load biling")

	var series []chart.Series
	for k, v := range costs.DividePerProject() {
		times, values := v.Timeseries()
		perSeries := chart.TimeSeries{Name: k, XValues: times, YValues: values}
		series = append(series, perSeries)
	}

	// c.f. https://github.com/wcharczuk/go-chart
	graph := chart.Chart{
		Title:      "gcp billing",
		TitleStyle: chart.StyleShow(),
		XAxis: chart.XAxis{
			Name:      "Day",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string { // change format YYYY-MM-DD to MM/DD
				return time.Unix(0, int64(v.(float64))).Format("01/02")
			},
		},
		YAxis: chart.YAxis{
			Name:      costs.Currency(),
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string { // round under decimal point
				return strconv.FormatInt(int64(v.(float64)), 10)
			},
			Range:     &chart.ContinuousRange{Min: 0, Max: costs.RoundMaxCost()},
			//GridLines: []chart.GridLine{{Value: 200, Style: chart.StyleShow()}, {Value: 400, Style: chart.StyleShow()}, {Value: 600, Style: chart.StyleShow()}, {Value: 800, Style: chart.StyleShow()}},
		},
		Series: series,
	}

	// line description
	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := graph.Render(chart.PNG, buffer); err != nil {
		log.Fatal(err)
	}

	store := storage.New()
	if err := store.Save(buffer.Bytes(), "output.png"); err != nil {
		log.Fatal(err)
	}

	log.Print("done")
}
