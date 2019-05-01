package billing

import (
	"time"
)

// Cost is gcp billing per project and day
type Cost struct {
	ProjectName string  `bigquery:"project_name"`
	Day         string  `bigquery:"export_time"`
	Cost        float64 `bigquery:"cost"`
	Currency    string  `bigquery:"currency"`
}

// ProjectCost is timeseries cost per project
type TotalCost []Cost

// Divide total cost to per project.
func (t TotalCost) DividePerProject() map[string]ProjectCost {

	m := map[string]ProjectCost{}

	for _, cost := range t {
		if costs, ok := m[cost.ProjectName]; ok {
			m[cost.ProjectName] = append(costs, cost)
			continue
		}
		m[cost.ProjectName] = append(make([]Cost, 0), cost)
	}

	return m
}

func (t TotalCost) Currency() string {
	if len(t) == 0 {
		return ""
	}
	return t[0].Currency
}

// ProjectCost is timeseries cost per project
type ProjectCost []Cost

// Timeseries is convert ProjectCost to X,Y axis.
func (p ProjectCost) Timeseries() ([]time.Time, []float64) {

	times := make([]time.Time, 0, len(p))
	costs := make([]float64, 0, len(p))

	for _, c := range p {
		times = append(times, MustParse("2006-01-02", c.Day))
		costs = append(costs, c.Cost)
	}

	return times, costs
}

func MustParse(layout, value string) time.Time {
	parse, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return parse
}
