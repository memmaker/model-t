package main

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

type MemoryLayout struct {
	TopLabel      *tview.TextView
	App           *tview.Application
	PageContainer *tview.Pages
	Pages         []Page
	CurrentPage   Page
}

type Page struct {
	Title  string
	UIForm *tview.Flex
	Model  Model
}

type JsonObject = map[string]interface{}
type StringMap = map[string]string

type MongoIndex struct {
	Keys   []map[string]int `json:"keys" bson:"keys"`
	Unique bool             `json:"unique" bson:"unique"`
}

type ForeignRelation struct {
	Name         string `json:"name" bson:"name"`
	RelatedModel string `json:"related_model" bson:"related_model"`
	RelatedField string `json:"related_field" bson:"related_field"`
}
type Model struct {
	Name             string            `json:"name" bson:"name"`
	DisplayTemplate  string            `json:"display_template" bson:"display_template,omitempty"`
	DisplayFields    []string          `json:"display_fields" bson:"display_fields,omitempty"`
	Indexes          []MongoIndex      `json:"indexes" bson:"indexes"`
	SearchFields     []string          `json:"search_fields" bson:"search_fields,omitempty"`
	Fields           []JsonObject      `json:"fields" bson:"fields"`
	ForeignRelations []ForeignRelation `json:"foreign_relations" bson:"foreign_relations"`
	Hidden           bool              `json:"hidden" bson:"hidden"`
}

var ctx = MemoryLayout{
	TopLabel:      tview.NewTextView(),
	App:           tview.NewApplication(),
	PageContainer: tview.NewPages(),
}

func main() {
	loadDataFromServer()

	if len(ctx.Pages) == 0 {

		addPage(Model{})
		addField()
	}
	ctx.PageContainer.SetInputCapture(inputHandler)

	if err := ctx.App.SetRoot(ctx.PageContainer, true).SetFocus(ctx.CurrentPage.UIForm).Run(); err != nil {
		panic(err)
	}
}

// addPage creates a new Page struct with a new UIForm and adds the given model to it.
// then the page struct is added to the ctx.Pages list and the ui page itself gets added to
// the ctx.PageContainer
func addPage(model Model) {
	var name = "model #" + strconv.Itoa(len(ctx.Pages)+1)
	if model.Name != "" {
		name = model.Name
	}
	pageForm := tview.NewFlex().SetDirection(tview.FlexRow)
	pageForm.SetBorder(true).SetTitle(name).SetTitleAlign(tview.AlignLeft)

	page := Page{
		Title:  name,
		UIForm: pageForm,
		Model:  model,
	}

	content := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageForm, 0, 1, true)
	content.Box = tview.NewBox().SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	ctx.PageContainer.AddAndSwitchToPage(name, content, true)
	ctx.Pages = append(ctx.Pages, page)
	ctx.CurrentPage = page
}

func loadDataFromServer() {
	fieldTypes = getFieldTypes()
	allModels := getModels()
	modelNames = make([]string, len(allModels))
	for index, model := range allModels {
		modelNames[index] = model.Name
		addPage(model)
		for rowIndex, _ := range model.Fields {
			addField()
			changeRowType(rowIndex)
		}
	}
}

func inputHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyF2 {
		addField()
	} else if event.Key() == tcell.KeyF3 {
		removeField()
	} else if event.Key() == tcell.KeyF4 {
		toggleRequired()
	} else if event.Key() == tcell.KeyF7 {
		pageBackward()
	} else if event.Key() == tcell.KeyF8 {
		pageForward()
	} else if event.Key() == tcell.KeyF10 {
		printAndQuit()
	}
	//topLabel.SetText(fmt.Sprintf("Focus index: %d of %d", focusIndex, currentForm.GetItemCount()))
	return event
}

func pageForward() {
	for index, page := range ctx.Pages {
		if page.Title == ctx.CurrentPage.Title {
			if index < len(ctx.Pages)-1 {
				ctx.PageContainer.SwitchToPage(ctx.Pages[index+1].Title)
				ctx.CurrentPage = ctx.Pages[index+1]
			}
		}
	}
}

func pageBackward() {
	for index, page := range ctx.Pages {
		if page.Title == ctx.CurrentPage.Title {
			if index > 0 {
				ctx.PageContainer.SwitchToPage(ctx.Pages[index-1].Title)
				ctx.CurrentPage = ctx.Pages[index-1]
			}
		}
	}
}

func printAndQuit() {
	ctx.App.Stop()
	// model to json
	jsonString, _ := json.MarshalIndent(ctx.CurrentPage.Model, "", "  ")
	fmt.Println(string(jsonString))
}
