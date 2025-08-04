/*
 * Copyright 2022-present Kuei-chun Chen. All rights reserved.
 * template.go
 */

package hatchet

import (
	"fmt"
	"html/template"
	"sort"
)

// GetTablesTemplate returns HTML
func GetTablesTemplate() (*template.Template, error) {
	html := headers + getContentHTML() + getMainPage() + "</body>"
	return template.New("hatchet").Funcs(template.FuncMap{
		"getHatchetImage": func() string {
			return HATCHET_PNG
		},
		"add": func(a int, b int) int {
			return a + b
		}}).Parse(html)
}

const headers = `<!DOCTYPE html>
<html lang="en">
<head>
  <title>Ken Chen's Hatchet</title>
	<meta http-equiv="Cache-Control" content="no-cache, no-store, must-revalidate" />
	<meta http-equiv="Pragma" content="no-cache" />
	<meta http-equiv="Expires" content="0" />

  <script src="/assets/jquery.min.js"></script>
  <script src="/assets/chart.min.js"></script>
  <link href="/favicon.ico" rel="icon" type="image/x-icon" />
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
  <link rel="stylesheet" href="/assets/main.css">
</head>
<body>
  <div id="loading">
    <div class="spinner"></div>
  </div>
  <script src="/assets/main.js"></script>
`

func getContentHTML() string {
	html := `
	<div class="navbar">
		<a href="/"><i class="fa fa-home"></i> Hatchet</a>
		<a href="#" onClick="loadData('/hatchets/{{.Hatchet}}/stats/audit'); return false;"><i class="fa fa-shield"></i> Audit</a>
		<a href="#" onClick="loadData('/hatchets/{{.Hatchet}}/stats/slowops'); return false;"><i class="fa fa-info"></i> Stats</a>
		<a href="#" onClick="loadData('/hatchets/{{.Hatchet}}/logs/slowops'); return false;"><i class="fa fa-list"></i> Top N</a>
		<a href="#" onClick="loadData('/hatchets/{{.Hatchet}}/logs/all?component=NONE'); return false;"><i class="fa fa-search"></i> Search</a>
		<div class="dropdown">
			<button class="dropbtn"><i class="fa fa-bar-chart"></i> Charts
			<i class="fa fa-caret-down"></i>
			</button>
			<div class="dropdown-content">
	`
	items := []Chart{}
	for _, chart := range charts {
		items = append(items, chart)
	}
	sort.Slice(items, func(i int, j int) bool {
		return items[i].Index < items[j].Index
	})

	for i, item := range items {
		if i == 0 {
			continue
		}
		html += fmt.Sprintf("<a href='#' onClick=\"loadData('/hatchets/{{.Hatchet}}/charts%v'); return false;\">%v</a>", item.URL, item.Title)
	}
	html += `
			</div>
		</div>
	</div>
	<div id="content"></div>
	`
	return html
}

func getMainPage() string {
	template := `
<div align='center'>
	<h2><img class='rotate23' width='60' valign="middle" src='data:image/png;base64,{{ getHatchetImage }}'>Hatchet - MongoDB JSON Log Analyzer</img></h2>
	<select id='table' class='hatchet-sel' onchange='javascript:redirect(); return false'>
		<option value=''>select a hatchet</option>
{{range $n, $value := .Hatchets}}
		<option value='{{$value}}'>{{$value}}</option>
{{end}}
	</select>
</div>
<hr/>
<div align='center'>
<iframe width="560" height="315" src="https://www.youtube.com/embed/WavOyaFTDE8" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>
</div>
<h3>Reports</h3>
    <table width='100%'>
      <tr><th></th><th>Title</th><th>Description</th></tr>
      <tr><td align=center><i class="fa fa-shield"></i></td><td>Audit</td><td>Display information on security audits and performance metrics</td></tr>
      <tr><td align=center><i class="fa fa-bar-chart"></i></td><td>Charts</td><td>A number of charts are available for security audits and performance metrics</td></tr>
      <tr><td align=center><i class="fa fa-search"></i></td><td>Search</td><td>Powerful log searching function with key metrics highlighted</td></tr>
      <tr><td align=center><i class="fa fa-info"></i></td><td>Stats</td><td>Summary of slow operational query patterns and duration</td></tr>
      <tr><td align=center><i class="fa fa-list"></i></td><td>TopN</td><td>Display the slowest 23 operation logs</td></tr>
    </table>
<h3>Charts</h3>
    <table width='100%'>
      <tr><th></th><th>Title</th><th>Description</th></tr>`
	size := len(charts) - 1
	tables := make([]Chart, size)
	for k, chart := range charts {
		if k == "instruction" {
			continue
		}
		tables[chart.Index-1] = chart
	}
	for _, chart := range tables {
		template += fmt.Sprintf("<tr><td align=right>%d</td><td>%v</td><td>%v</td></tr>\n",
			chart.Index, chart.Title, chart.Descr)
	}
	template += "</table>"
	template += `<h3>URL</h3>
<ul class="api">
	<li>/</li>
	<li>/hatchets/{hatchet}/charts/{chart}[?type={str}]</li>
	<li>/hatchets/{hatchet}/logs/all[?component={str}&context={str}&duration={date},{date}&severity={str}&limit=[{offset},]{int}]</li>
	<li>/hatchets/{hatchet}/logs/slowops[?topN={int}]</li>
	<li>/hatchets/{hatchet}/stats/slowops[?COLLSCAN={bool}&orderBy={str}]</li>
</ul>

<h3>API</h3>
<ul class="api">
	<li>/api/hatchet/v1.0/hatchets/{hatchet}/logs/all[?component={str}&context={str}&duration={date},{date}&severity={str}&limit=[{offset},]{int}]</li>
	<li>/api/hatchet/v1.0/hatchets/{hatchet}/logs/slowops[?topN={int}]</li>
	<li>/api/hatchet/v1.0/hatchets/{hatchet}/stats/audit</li>
	<li>/api/hatchet/v1.0/hatchets/{hatchet}/stats/slowops[?COLLSCAN={bool}&orderBy={str}]</li>
	<li>/api/hatchet/v1.0/mongodb/{version}/drivers/{driver}?compatibleWith={driver version}</li>
</ul>
<h4 align='center'><hr/>{{.Version}}</h4>
`
	template += fmt.Sprintf(`<div class="footer"><img valign="middle" src='data:image/png;base64,%v'/> Ken Chen</div>`, CHEN_ICO)
	return template
}
