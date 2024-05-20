package main

import (
	"fmt"
	"strings"

	g "github.com/AllenDang/giu"
)

var s string
var path string
var x []string
var content []string = []string{"empty"}
var inter []interface{} = []interface{}{"empty for now"}
var case_count bool

func onDrop(names []string) {
	var sb strings.Builder
	for _, n := range names {
		sb.WriteString(n)
	}

	path = sb.String()
	g.Update()
}
func BuildRows() []*g.TableRowWidget {
	var rows []*g.TableRowWidget

	if len(content) > 0 {
		rows = append(rows, g.TableRow(
			g.Label("Binding"),
			g.Label("Shortcut"),
		))
		for i := range content {
			line := content[i]
			if strings.Contains(line, ":") {
				s := strings.Split(line, ":")
				rows = append(rows, g.TableRow(
					g.Markdown(&s[0]),
					g.Label(s[1]),
				))
			} else {
				rows = append(rows, g.TableRow(g.Label("no line found")))
			}
		}
	}

	return rows
}

func UpdateInterface() {
	inter = make([]interface{}, len(content))
	for i, v := range content {
		inter[i] = fmt.Sprint(v)
	}
}

func CaseInsensitiveContainsToggle(s, substr string, sensitive bool) bool {
	if !sensitive {
		s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	}
	return strings.Contains(s, substr)
}

func loop() {
	g.SingleWindow().Layout(
		g.Row(
			g.InputText(&path),
			g.Button("Parse").Disabled(path == "").OnClick(func() {
				content = parseXML(path)
				UpdateInterface()
				if s == "" {
					x = content
				}

			}),
		),
		g.Row(
			g.InputText(&s).Label("Search a pattern").OnChange(func() {
				if s == "" {
					content = x
					UpdateInterface()
					return
				}
				res := []string{}
				for _, l := range x {
					if CaseInsensitiveContainsToggle(l, s, case_count) {
						res = append(res, l)
					}
				}
				if len(res) == 0 {
					content = []string{"No results with this pattern"}
				} else {
					content = res
				}
				UpdateInterface()
			}),
			g.Checkbox("Case-sensitive", &case_count),
		),
		g.Table().Freeze(0, 1).FastMode(true).Rows(BuildRows()...),
	)
}
func main() {
	w := g.NewMasterWindow("Renoise Keys Util", 800, 1000, 0)
	g.Context.FontAtlas.SetDefaultFont("Arial.ttf", 14)
	w.SetDropCallback(onDrop)
	w.Run(loop)
}
