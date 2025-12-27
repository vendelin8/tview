// Demo code for the FormModal primitive.
package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/vendelin8/tview"
)

const (
	pageModal = "modal"
	pageBgd   = "background"
)

var b strings.Builder

func addText(txt string) {
	if len(txt) == 0 {
		return
	}

	if b.Len() > 0 {
		b.WriteByte(' ')
	}

	b.WriteString(txt)
}

func addProp(prepos, condition string, hasEnding bool, insertIndex int, text ...string) bool {
	if len(condition) == 0 {
		return hasEnding
	}

	if hasEnding {
		b.WriteString(".\n")
		b.WriteString(prepos)
	}

	for i, s := range text {
		b.WriteString(s)

		if insertIndex == i {
			b.WriteString(condition)
		}
	}

	return true
}

func main() {
	app := tview.NewApplication()

	maleTitles := []string{"Mr.", "Dr.", "Prof."}
	femaleTitles := []string{"Ms.", "Mrs.", "Dr.", "Prof."}
	preposBySex := map[int]string{0: "He", 1: "She"}
	adultMatcher := map[bool]string{false: "a child", true: "an adult"}

	pages := tview.NewPages()
	text := tview.NewTextView().SetText("Click on the button to edit user profile...")

	var (
		title     *tview.DropDown
		firstName *tview.InputField
		lastName  *tview.InputField
		address   *tview.TextArea
		sex       *tview.Radio
		notes     *tview.TextArea
		adult     *tview.Checkbox
		password  *tview.InputField
	)

	formModal := tview.NewFormModal(func(form *tview.Form) {
		form.AddDropDown("Title", maleTitles, 0, nil).
			AddInputField("First name", "", 20, nil, nil).
			AddInputField("Last name", "", 20, nil, nil).
			AddTextArea("Address", "", 40, 0, 0, nil).
			AddRadio("Sex", 0, true, func(newValue int) {
				var newOptions []string
				if newValue == 0 {
					newOptions = maleTitles
				} else {
					newOptions = femaleTitles
				}

				_, option := title.GetCurrentOption()
				newOption := 0

				for i, opt := range newOptions {
					if opt == option {
						newOption = i
						break
					}
				}

				title.SetOptions(newOptions, nil)
				title.SetCurrentOption(newOption)
			}, "male", "female").
			AddTextArea("Notes", "This is just a demo. You can enter whatever you wish. Mind how the radio changes title options", 40, 0, 0, nil).
			AddCheckbox("Age 18+", false, nil).
			AddPasswordField("Password", "", 10, '*', nil).
			AddButton("Save", func() {
				prepos := preposBySex[sex.Value()]

				b.Reset()
				addText(firstName.GetText())
				addText(lastName.GetText())
				if b.Len() == 0 {
					b.WriteString(prepos)
				} else {
					name := b.String()
					b.Reset()
					_, titleStr := title.GetCurrentOption()
					b.WriteString(titleStr)
					b.WriteByte(' ')
					b.WriteString(name)
				}

				hasEnding := addProp(prepos, address.GetText(), false, 0, " lives at ")
				hasEnding = addProp(prepos, notes.GetText(), hasEnding, 0, " has a note of '", "'")
				hasEnding = addProp(prepos, adultMatcher[adult.IsChecked()], hasEnding, 0, " is ")
				hasEnding = addProp(prepos, password.GetText(), hasEnding, 0, " has a password of '", "'")

				if hasEnding {
					b.WriteByte('.')
				}
				text.SetText(b.String())
				pages.HidePage("modal")
			}).
			AddButton("Cancel", func() {
				pages.HidePage("modal")
			}).AddButton("Some", func() {}).
			AddButton("More", func() {}).
			AddButton("Useless", func() {}).
			AddButton("Buttons", func() {}).
			AddButton("To", func() {}).
			AddButton("Extend", func() {}).
			SetTitle("Enter some data").
			SetTitleAlign(tview.AlignLeft)
	})

	title = formModal.GetFormItem(0).(*tview.DropDown)
	firstName = formModal.GetFormItem(1).(*tview.InputField)
	lastName = formModal.GetFormItem(2).(*tview.InputField)
	address = formModal.GetFormItem(3).(*tview.TextArea)
	sex = formModal.GetFormItem(4).(*tview.Radio)
	notes = formModal.GetFormItem(5).(*tview.TextArea).SetWrap(true)
	adult = formModal.GetFormItem(6).(*tview.Checkbox)
	password = formModal.GetFormItem(7).(*tview.InputField)

	menu := tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetWrap(false)
	fmt.Fprint(menu, ` F2 [""][green::b]Edit[white::-][""]  `)
	fmt.Fprint(menu, ` Esc [""][orange::b]Quit[white::-][""]  `)
	app.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		if ev.Key() == tcell.KeyEsc {
			app.Stop()
			return nil
		}

		if ev.Key() == tcell.KeyF2 {
			pages.ShowPage(pageModal)
			return nil
		}

		return ev
	})

	background := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(text, 0, 1, false).AddItem(menu, 1, 0, true)

	pages.AddPage(pageBgd, background, true, true).AddPage(pageModal, formModal, false, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
