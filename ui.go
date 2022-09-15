package main

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os/exec"
)

var fieldTypes = []string{"string", "int", "float", "related", "dropdown"}

var modelNames = []string{"User", "Post", "Comment"}

var oreCommand = "ore"

type FieldResponse struct {
	FieldTypeNames []string `json:"field_types"`
}

func getFieldTypes() []string {
	var fieldResponse FieldResponse
	osCommand := exec.Command(oreCommand, "fields")
	jsonString, err := osCommand.Output()
	if err != nil {
		panic(err)
	}
	json.Unmarshal(jsonString, &fieldResponse)
	return fieldResponse.FieldTypeNames
}

func getModels() []Model {
	var models []Model
	osCommand := exec.Command(oreCommand, "models")
	jsonString, err := osCommand.Output()
	if err != nil {
		panic(err)
	}
	json.Unmarshal(jsonString, &models)
	return models
}

func getIndexOfFieldType(fieldType string) int {
	for index, value := range fieldTypes {
		if value == fieldType {
			return index
		}
	}
	return -1
}

// addField adds a new field to the currentModel and currentForm
func addField() *tview.Flex {

	ctx.CurrentPage.Model.Fields = append(ctx.CurrentPage.Model.Fields, JsonObject{
		"name":     "",
		"type":     "",
		"required": false,
	})

	itemIndex := ctx.CurrentPage.UIForm.GetItemCount()
	rowContainer := createFieldRow()

	addNameInput(rowContainer, itemIndex, "")
	addTypeSelector(rowContainer, itemIndex, -1)
	addFlagsLabel(rowContainer)

	ctx.CurrentPage.UIForm.AddItem(rowContainer, 1, 0, true)
	return rowContainer
}

func createFieldRow() *tview.Flex {
	rowContainer := tview.NewFlex().SetDirection(tview.FlexColumn)
	rowContainer.Box = tview.NewBox().SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	return rowContainer
}
func addFlagsLabel(rowContainer *tview.Flex) {
	flagsLabel := tview.NewTextView().SetText("Flags: ")
	flagsLabel.SetDynamicColors(true)
	rowContainer.AddItem(flagsLabel, 15, 0, false)
}

func addNameInput(rowContainer *tview.Flex, itemIndex int, presetNameValue string) {
	inputField := tview.NewInputField().
		SetLabel(" Name: ").
		SetText(presetNameValue).
		SetDoneFunc(func(key tcell.Key) {
			navigation(key, Start)
		}).SetChangedFunc(func(text string) {
		ctx.CurrentPage.Model.Fields[itemIndex]["name"] = text
	})
	rowContainer.AddItem(inputField, 0, 4, true)
}

func addOptionsInput(rowContainer *tview.Flex, itemIndex int, presetValue string) *tview.InputField {
	inputField := tview.NewInputField().
		SetLabel(" Options: ").
		SetText(presetValue).
		SetDoneFunc(func(key tcell.Key) {
			navigation(key, Middle)
		}).SetChangedFunc(func(text string) {
		ctx.CurrentPage.Model.Fields[itemIndex]["options"] = text
	})
	rowContainer.AddItem(inputField, 0, 4, true)
	return inputField
}

func addTypeSelector(rowContainer *tview.Flex, itemIndex int, presetValue int) *tview.DropDown {
	dropDown := tview.NewDropDown().
		SetLabel(" Type: ").
		SetOptions(fieldTypes, nil).
		SetDoneFunc(func(key tcell.Key) {
			navigation(key, End)
		})
	if presetValue >= 0 {
		dropDown.SetCurrentOption(presetValue)
	}
	dropDown.SetSelectedFunc(func(text string, index int) {
		ctx.CurrentPage.Model.Fields[itemIndex]["type"] = text
		changeRowType(itemIndex)
	})
	rowContainer.AddItem(dropDown, 25, 0, true)
	return dropDown
}

func addModelSelector(rowContainer *tview.Flex, itemIndex int, presetValue int) *tview.DropDown {
	dropDown := tview.NewDropDown().
		SetLabel(" Rel. Model: ").
		SetOptions(modelNames, nil).
		SetDoneFunc(func(key tcell.Key) {
			navigation(key, Middle)
		})
	if presetValue >= 0 {
		dropDown.SetCurrentOption(presetValue)
	}
	dropDown.SetSelectedFunc(func(text string, index int) {
		ctx.CurrentPage.Model.Fields[itemIndex]["related_model"] = text
	})
	rowContainer.AddItem(dropDown, 0, 4, true)
	return dropDown
}

func changeRowType(index int) {
	fieldDefinition := ctx.CurrentPage.Model.Fields[index]
	rowContainer := ctx.CurrentPage.UIForm.GetItem(index).(*tview.Flex)
	var dropdown *tview.DropDown
	rowContainer.Clear()
	if fieldDefinition["type"] == "related" {
		addNameInput(rowContainer, index, fieldDefinition["name"].(string))
		addModelSelector(rowContainer, index, getIndexOfModelName(keyValueOrNothing(fieldDefinition, "related_model")))
		dropdown = addTypeSelector(rowContainer, index, getIndexOfFieldType(fieldDefinition["type"].(string)))
		addFlagsLabel(rowContainer)
		updateFlags(index)
	} else if fieldDefinition["type"] == "dropdown" {
		addNameInput(rowContainer, index, fieldDefinition["name"].(string))
		addOptionsInput(rowContainer, index, listToString(keyValueOrEmptyList(fieldDefinition, "options")))
		dropdown = addTypeSelector(rowContainer, index, getIndexOfFieldType(fieldDefinition["type"].(string)))
		addFlagsLabel(rowContainer)
		updateFlags(index)
	} else {
		addNameInput(rowContainer, index, fieldDefinition["name"].(string))
		dropdown = addTypeSelector(rowContainer, index, getIndexOfFieldType(fieldDefinition["type"].(string)))
		addFlagsLabel(rowContainer)
		updateFlags(index)
	}
	ctx.App.SetFocus(dropdown)
}

func getIndexOfModelName(modelName string) int {
	for index, value := range modelNames {
		if value == modelName {
			return index
		}
	}
	return -1
}

func keyValueOrFalse(object JsonObject, key string) bool {
	value, ok := object[key]
	if !ok {
		return false
	}
	return value.(bool)
}

func keyValueOrNothing(object JsonObject, key string) string {
	value, ok := object[key]
	var options string
	if !ok {
		options = ""
	} else {
		options = value.(string)
	}
	return options
}

func listToString(list []string) string {
	var result string
	for index, value := range list {
		if index > 0 {
			result += ", " + value
		}
	}
	return result
}

func keyValueOrEmptyList(object JsonObject, key string) []string {
	value, ok := object[key]
	var options []string
	if !ok {
		options = []string{}
	} else {
		for _, optionValue := range value.([]interface{}) {
			options = append(options, optionValue.(string))
		}
	}
	return options
}

func removeField() {
	itemCount := ctx.CurrentPage.UIForm.GetItemCount()
	if itemCount > 1 {
		lastFieldItem := ctx.CurrentPage.UIForm.GetItem(itemCount - 1)
		ctx.CurrentPage.UIForm.RemoveItem(lastFieldItem)
		ctx.CurrentPage.Model.Fields = ctx.CurrentPage.Model.Fields[:len(ctx.CurrentPage.Model.Fields)-1]
	}
}

func toggleRequired() {
	ctx.CurrentPage.Model.Fields[focusIndex]["required"] = !keyValueOrFalse(ctx.CurrentPage.Model.Fields[focusIndex], "required")
	updateFlags(focusIndex)
}

func toggleSearchable() {
	fieldName := ctx.CurrentPage.Model.Fields[focusIndex]["name"].(string)
	if foundIndex := contains(ctx.CurrentPage.Model.SearchFields, fieldName); foundIndex > -1 {
		ctx.CurrentPage.Model.SearchFields = append(ctx.CurrentPage.Model.SearchFields[:foundIndex], ctx.CurrentPage.Model.SearchFields[foundIndex+1:]...)
	} else {
		ctx.CurrentPage.Model.SearchFields = append(ctx.CurrentPage.Model.SearchFields, fieldName)
	}
	updateFlags(focusIndex)
}

func updateFlags(index int) {
	rowContainer := ctx.CurrentPage.UIForm.GetItem(index).(*tview.Flex)

	flagsLabel := rowContainer.GetItem(rowContainer.GetItemCount() - 1).(*tview.TextView)
	flagsLabel.SetText(fmt.Sprintf("Flags: %s", getFlagsString(ctx.CurrentPage.Model, ctx.CurrentPage.Model.Fields[index])))
}

func getFlagsString(model Model, fieldDefinition JsonObject) string {
	fieldName := fieldDefinition["name"].(string)
	flags := ""
	isRequired := keyValueOrFalse(fieldDefinition, "required")
	isSearchField := contains(model.SearchFields, fieldName) > -1
	if isRequired {
		flags += "R"
	} else {
		flags += "[gray]-[-]"
	}
	if isSearchField {
		flags += "S"
	} else {
		flags += "[gray]-[-]"
	}
	return flags
}

func contains(list []string, value string) int {
	for index, item := range list {
		if item == value {
			return index
		}
	}
	return -1
}
