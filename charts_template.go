// Copyright 2022-present Kuei-chun Chen. All rights reserved.
package hatchet

import (
	"fmt"
	"html/template"
	"time"
)

func getFooter() string {
	summary := "{{.Summary}}"
	return fmt.Sprintf(`<div class="footer"><img valign="middle" src='data:image/png;base64,%v'> %v</img></div>`,
		CHEN_ICO, summary)
}

// GetChartTemplate returns HTML
func GetChartTemplate(chartType string) (*template.Template, error) {
	var html string
	if chartType == BUBBLE_CHART {
		html = getOpStatsChart()
	} else if chartType == PIE_CHART {
		html = getPieChart()
	} else if chartType == BAR_CHART {
		html = getConnectionsChart()
	}
	html += `
	<div style="float: left; width: 100%; clear: left;">
		<input type='datetime-local' id='start' value='{{.Start}}'></input>
		<input type='datetime-local' id='end' value='{{.End}}'></input>
		<button onClick="refreshChart(); return false;" class="button">Refresh</button>
  	</div>
	<div id='hatchetChart' class='chart' style="clear: left;"></div>`

	return template.New("hatchet").Funcs(template.FuncMap{
		"descr": func(v OpCount) template.HTML {
			if v.Filter == "" {
				return template.HTML(v.Namespace)
			}
			str := fmt.Sprintf("%v, QP: %v", v.Namespace, v.Filter)
			return template.HTML(str)
		},
		"toSeconds": func(n float64) float64 {
			return n / 1000
		},
		"substr": func(str string, n int) string {
			return str[:n]
		},
		"epoch": func(d string, s string) int64 {
			dfmt := "2016-01-02T23:59:59"
			sdt, _ := time.Parse("2006-01-02T15:04:05", s+dfmt[len(s):])
			dt, _ := time.Parse("2006-01-02T15:04:05", d+dfmt[len(d):])
			return dt.Unix() - sdt.Unix()
		}}).Parse(html)
}

func getOpStatsChart() string {
	return `
{{ if .OpCounts }}
<canvas id="hatchetChart"></canvas>
<script>
	setChartType();
	const ctx = document.getElementById('hatchetChart').getContext('2d');
	const chart = new Chart(ctx, {
		type: 'bubble',
		data: {
			datasets: [{
				label: '{{.Chart.Title}}',
				data: [
				{{range $i, $v := .OpCounts}}
					{
						x: new Date("{{$v.Date}}"),
						y: {{toSeconds $v.Milli}},
						r: {{$v.Count}}
					},
				{{end}}
				],
				backgroundColor: 'rgba(255, 99, 132, 0.2)',
				borderColor: 'rgba(255, 99, 132, 1)',
				borderWidth: 1
			}]
		},
		options: {
			responsive: true,
			plugins: {
				title: {
					display: true,
					text: '{{.Chart.Title}}'
				}
			},
			scales: {
				x: {
					type: 'time',
					time: {
						unit: 'day'
					}
				},
				y: {
					title: {
						display: true,
						text: '{{.VAxisLabel}}'
					}
				}
			}
		}
	});
</script>
{{else}}
<div align='center' class='btn'><span style='color: red'>no data found</span></div>
{{end}}`
}

func getPieChart() string {
	return `
{{ if .NameValues }}
<canvas id="hatchetChart"></canvas>
<script>
	setChartType();
	const ctx = document.getElementById('hatchetChart').getContext('2d');
	const chart = new Chart(ctx, {
		type: 'pie',
		data: {
			labels: [
				{{range $i, $v := .NameValues}}
					'{{$v.Name}}',
				{{end}}
			],
			datasets: [{
				label: '{{.Chart.Title}}',
				data: [
					{{range $i, $v := .NameValues}}
						{{$v.Value}},
					{{end}}
				],
				backgroundColor: [
					'rgba(255, 99, 132, 0.2)',
					'rgba(54, 162, 235, 0.2)',
					'rgba(255, 206, 86, 0.2)',
					'rgba(75, 192, 192, 0.2)',
					'rgba(153, 102, 255, 0.2)',
					'rgba(255, 159, 64, 0.2)'
				],
				borderColor: [
					'rgba(255, 99, 132, 1)',
					'rgba(54, 162, 235, 1)',
					'rgba(255, 206, 86, 1)',
					'rgba(75, 192, 192, 1)',
					'rgba(153, 102, 255, 1)',
					'rgba(255, 159, 64, 1)'
				],
				borderWidth: 1
			}]
		},
		options: {
			responsive: true,
			plugins: {
				title: {
					display: true,
					text: '{{.Chart.Title}}'
				}
			}
		}
	});
</script>
{{else}}
<div align='center' class='btn'><span style='color: red'>no data found</span></div>
{{end}}`
}

func getConnectionsChart() string {
	return `
{{ if .Remote }}
<canvas id="hatchetChart"></canvas>
<script>
	setChartType();
	const ctx = document.getElementById('hatchetChart').getContext('2d');
	const chart = new Chart(ctx, {
		type: 'bar',
		data: {
			labels: [
				{{range $i, $v := .Remote}}
					'{{$v.IP}}',
				{{end}}
			],
			datasets: [
				{
					label: 'Accepted',
					data: [
						{{range $i, $v := .Remote}}
							{{$v.Accepted}},
						{{end}}
					],
					backgroundColor: 'rgba(54, 162, 235, 0.2)',
					borderColor: 'rgba(54, 162, 235, 1)',
					borderWidth: 1
				},
				{
					label: 'Ended',
					data: [
						{{range $i, $v := .Remote}}
							{{$v.Ended}},
						{{end}}
					],
					backgroundColor: 'rgba(255, 99, 132, 0.2)',
					borderColor: 'rgba(255, 99, 132, 1)',
					borderWidth: 1
				}
			]
		},
		options: {
			responsive: true,
			plugins: {
				title: {
					display: true,
					text: '{{.Chart.Title}}'
				}
			},
			scales: {
				x: {
					stacked: true,
				},
				y: {
					stacked: true
				}
			}
		}
	});
</script>
{{else}}
<div align='center' class='btn'><span style='color: red'>no data found</span></div>
{{end}}`
}
