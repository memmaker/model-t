package main

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

type MemoryLayout struct {
	RootContainer *tview.Flex
	TopBar        *tview.Flex
	BottomBar     *tview.Flex
	App           *tview.Application
	Modal         *tview.Flex
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
	RootContainer: tview.NewFlex().SetDirection(tview.FlexRow),
	TopBar:        tview.NewFlex().SetDirection(tview.FlexColumn),
	BottomBar:     tview.NewFlex().SetDirection(tview.FlexColumn),
	App:           tview.NewApplication(),
	PageContainer: tview.NewPages(),
}

func main() {
	loadDataFromServer()

	if len(ctx.Pages) == 0 {
		addPage(Model{})
		addField()
	}
	ctx.RootContainer.SetInputCapture(inputHandler)
	ctx.RootContainer.AddItem(ctx.TopBar, 1, 0, false)
	ctx.RootContainer.AddItem(ctx.PageContainer, 0, 1, true)
	ctx.RootContainer.AddItem(ctx.BottomBar, 1, 0, false)
	ctx.BottomBar.AddItem(tview.NewTextView().SetText("F2: Add Field, F3: Remove Field, F4: Toggle Required, F7: Page Backward, F8: Page Forward, F10: Print and Quit"), 0, 1, false)

	if err := ctx.App.SetRoot(ctx.RootContainer, true).SetFocus(ctx.CurrentPage.UIForm).Run(); err != nil {
		panic(err)
	}
}

func toggleModal() {
	if ctx.Modal == nil {
		showModal()
	} else {
		closeModal()
	}
}

func makeModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func showModal() {
	createModelUI(ctx.CurrentPage.Model.Name)
	ctx.PageContainer.AddAndSwitchToPage(".modal.", makeModal(ctx.Modal, 50, 25), true)
	ctx.PageContainer.ShowPage(ctx.CurrentPage.Title)
	//ctx.App.SetFocus(ctx.Modal)
}

func createModelUI(modelName string) {
	ctx.Modal = tview.NewFlex().SetDirection(tview.FlexRow)
	ctx.Modal.Box = tview.NewBox().SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	ctx.Modal.SetTitle("Model: " + modelName)
	ctx.Modal.SetBorder(true)
	ctx.Modal.AddItem(tview.NewInputField().SetLabel(" Name: ").SetText(modelName), 1, 0, true)
	ctx.Modal.AddItem(tview.NewInputField().SetLabel(" Disp. Template: ").SetText(""), 1, 0, true)
	// add checkbox
	ctx.Modal.AddItem(tview.NewCheckbox().SetLabel(" Hidden: "), 1, 0, true)
}

func closeModal() {
	ctx.PageContainer.RemovePage(".modal.")
	ctx.App.SetFocus(ctx.CurrentPage.UIForm)
	ctx.Modal = nil
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
	pageForm.SetBorder(true).SetTitleAlign(tview.AlignCenter)
	setFormTitle(pageForm, name)
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

func setFormTitle(form *tview.Flex, givenModelName string) {
	var title string
	for _, modelName := range modelNames {
		if modelName == givenModelName {
			title += "[white:blue:b]|" + modelName + "|[-:-:-]"
		} else {
			title += " " + modelName + " "
		}
	}
	form.SetTitle(" " + title + " ")
}

func loadDataFromServer() {
	fieldTypes = getFieldTypes()
	allModels := getModels()
	setModelNames(allModels)
	for _, model := range allModels {
		addPage(model)
		for rowIndex, _ := range model.Fields {
			addField()
			changeRowType(rowIndex)
		}
	}
}

func setModelNames(allModels []Model) {
	modelNames = make([]string, len(allModels))
	for index, model := range allModels {
		modelNames[index] = model.Name
	}
}

func inputHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyF2 {
		addField()
	} else if event.Key() == tcell.KeyF3 {
		removeField()
	} else if event.Key() == tcell.KeyF4 {
		toggleRequired()
	} else if event.Key() == tcell.KeyF5 {
		toggleSearchable()
	} else if event.Key() == tcell.KeyF7 {
		pageBackward()
	} else if event.Key() == tcell.KeyF8 {
		pageForward()
	} else if event.Key() == tcell.KeyF9 {
		toggleModal()
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
			return
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
			return
		}
	}
}

func printAndQuit() {
	ctx.App.Stop()
	// model to json
	jsonString, _ := json.MarshalIndent(ctx.CurrentPage.Model, "", "  ")
	fmt.Println(string(jsonString))
}
