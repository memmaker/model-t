package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SubItemPosition int64

const (
	Start SubItemPosition = iota
	Middle
	End
)

var focusIndex = 0

func navigation(key tcell.Key, subPositionHint SubItemPosition) {
	if key == tcell.KeyTab && subPositionHint == End {
		focusForward(ctx.CurrentPage.UIForm)
	} else if key == tcell.KeyBacktab && subPositionHint == Start {
		focusBackwardForFields(ctx.CurrentPage.UIForm)
	} else if key == tcell.KeyTab && (subPositionHint == Start || subPositionHint == Middle) {
		focusForwardInRow(ctx.CurrentPage.UIForm)
	} else if key == tcell.KeyBacktab && (subPositionHint == End || subPositionHint == Middle) {
		focusBackwardInRow(ctx.CurrentPage.UIForm)
	} else if key == tcell.KeyDown {
		focusRowDown(ctx.CurrentPage.UIForm)
	} else if key == tcell.KeyUp {
		focusRowUp(ctx.CurrentPage.UIForm)
	}
}

func modalNavigation(key tcell.Key, subPositionHint SubItemPosition) {
	if key == tcell.KeyTab && subPositionHint == End {
		focusForward(ctx.Modal)
	} else if key == tcell.KeyBacktab && subPositionHint == Start {
		focusBackward(ctx.Modal)
	} else if key == tcell.KeyTab && (subPositionHint == Start || subPositionHint == Middle) {
		focusForward(ctx.Modal)
	} else if key == tcell.KeyBacktab && (subPositionHint == End || subPositionHint == Middle) {
		focusBackward(ctx.Modal)
	} else if key == tcell.KeyDown {
		focusForward(ctx.Modal)
	} else if key == tcell.KeyUp {
		focusBackward(ctx.Modal)
	}
}

func focusRowUp(container *tview.Flex) {
	currentlyFocusedItem := ctx.App.GetFocus()
	field := container.GetItem(focusIndex).(*tview.Flex)
	index := getIndexOfChild(field, currentlyFocusedItem)
	focusIndex = focusIndex - 1
	if focusIndex < 0 {
		focusIndex = container.GetItemCount() - 1
	}
	field = container.GetItem(focusIndex).(*tview.Flex)
	nextSubItem := field.GetItem(index)
	ctx.App.SetFocus(nextSubItem)
}

func focusRowDown(container *tview.Flex) {
	currentlyFocusedItem := ctx.App.GetFocus()
	field := container.GetItem(focusIndex).(*tview.Flex)
	index := getIndexOfChild(field, currentlyFocusedItem)
	focusIndex = (focusIndex + 1) % container.GetItemCount()
	field = container.GetItem(focusIndex).(*tview.Flex)
	nextSubItem := field.GetItem(index)
	ctx.App.SetFocus(nextSubItem)
}

func focusForwardInRow(container *tview.Flex) {
	currentlyFocusedItem := ctx.App.GetFocus()
	field := container.GetItem(focusIndex).(*tview.Flex)

	index := getIndexOfChild(field, currentlyFocusedItem)
	nextSubItem := field.GetItem(index + 1)
	ctx.App.SetFocus(nextSubItem)
}

func focusBackwardInRow(container *tview.Flex) {
	currentlyFocusedItem := ctx.App.GetFocus()
	field := container.GetItem(focusIndex).(*tview.Flex)

	index := getIndexOfChild(field, currentlyFocusedItem)
	nextSubItem := field.GetItem(index - 1)
	ctx.App.SetFocus(nextSubItem)
}

func getIndexOfChild(container *tview.Flex, child tview.Primitive) int {
	for i := 0; i < container.GetItemCount(); i++ {
		if container.GetItem(i) == child {
			return i
		}
	}
	return -1
}

func focusForward(container *tview.Flex) {
	focusIndex = (focusIndex + 1) % container.GetItemCount()
	field := container.GetItem(focusIndex)
	ctx.App.SetFocus(field)
}

func focusBackward(container *tview.Flex) {
	focusIndex = focusIndex - 1
	if focusIndex < 0 {
		focusIndex = container.GetItemCount() - 1
	}
	field := container.GetItem(focusIndex)
	ctx.App.SetFocus(field)
}

func focusBackwardForFields(container *tview.Flex) {
	focusIndex = focusIndex - 1
	if focusIndex < 0 {
		focusIndex = container.GetItemCount() - 1
	}
	field := container.GetItem(focusIndex).(*tview.Flex)
	lastItemInField := field.GetItem(field.GetItemCount() - 2)
	ctx.App.SetFocus(lastItemInField)
}
