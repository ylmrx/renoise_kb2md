package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type (
	KeyboardBindings struct {
		XMLName    xml.Name   `xml:"KeyboardBindings"`
		Categories Categories `xml:"Categories"`
	}
	Categories struct {
		XMLName    xml.Name   `xml:"Categories"`
		Categories []Category `xml:"Category"`
	}
	Category struct {
		XMLName     xml.Name    `xml:"Category"`
		Identifier  string      `xml:"Identifier"`
		KeyBindings KeyBindings `xml:"KeyBindings"`
	}
	KeyBindings struct {
		XMLName     xml.Name     `xml:"KeyBindings"`
		KeyBindings []KeyBinding `xml:"KeyBinding"`
	}
	KeyBinding struct {
		XMLName xml.Name `xml:"KeyBinding"`
		Topic   string   `xml:"Topic"`
		Binding string   `xml:"Binding"`
		Key     string   `xml:"Key"`
	}
)

func parseXML(path string) {
	var out bytes.Buffer
	re := regexp.MustCompile(`\.xml$`)
	destination := re.ReplaceAllString(path, ".md")
	xmlFile, err := os.Open(path)
	if err != nil {
		log.Println("didn't open kb file")
	}
	defer xmlFile.Close()
	b, _ := io.ReadAll(xmlFile)
	var kbs KeyboardBindings
	xml.Unmarshal(b, &kbs)

	for _, kb := range kbs.Categories.Categories {
		fmt.Println("*** " + kb.Identifier)
		out.WriteString("\n## " + kb.Identifier + "\n")
		curr := ""
		for _, k := range kb.KeyBindings.KeyBindings {
			if len(k.Key) > 0 {
				if k.Topic != curr {
					curr = k.Topic
					out.WriteString("\n### " + curr + "\n\n")
				}
				out.WriteString("- " + k.Binding + " = `" + k.Key + "`\n")
			}
			// out.WriteString("\n")
		}
	}

	os.WriteFile(destination, out.Bytes(), os.FileMode(os.O_RDWR))
}

func main() {
	a := app.New()
	win := a.NewWindow("Keybindings")
	win.Resize(fyne.NewSize(800, 300))

	var path_to_renoise string

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter Renoise conf dir path...")

	i_wid := container.NewVBox(input, widget.NewButton("Save", func() {
		const kb_file = "KeyBindings.xml"
		if _, err := os.ReadDir(input.Text); err != nil {
			log.Println("failed to list dir")
		}
		path_to_renoise = filepath.Join(input.Text, kb_file)
		if _, err := os.Stat(path_to_renoise); err != nil {
			log.Println("Content was:", path_to_renoise)
			return
		}
		parseXML(path_to_renoise)
	}))

	win.SetContent(i_wid)
	win.ShowAndRun()
}
