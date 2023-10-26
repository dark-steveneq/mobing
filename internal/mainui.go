package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dark-steveneq/mobiapi"
)

func messagerefresh(api *mobiapi.MobiAPI, messageContainer *fyne.Container, messageContents *fyne.Container) {
	messageContainer.RemoveAll()
	messageContainer.Add(container.NewCenter(container.NewVBox(widget.NewLabel("Loading Messages..."), widget.NewProgressBarInfinite())))
	messages, err := api.GetReadMessages(false)
	if err != nil {
		messageContainer.RemoveAll()
		messageContainer.Add(container.NewCenter(container.NewVBox(widget.NewLabel("Couldn't retrieve messages!"))))
	}
	messageContainer.RemoveAll()
	for i, rmessage := range messages {
		message := rmessage
		messageContainer.Add(container.NewVBox(widget.NewLabelWithStyle(message.Title, fyne.TextAlignCenter, widget.RichTextStyleBlockquote.TextStyle), widget.NewLabel("From: "+message.Author), widget.NewButton("Read", func() {
			messageContents.RemoveAll()
			messageContents.Add(widget.NewLabel("Loading..."))
			messageContents.Add(widget.NewProgressBarInfinite())
			messagecontent, err := api.GetMessageContent(message)
			if err != nil {
				messageContents.RemoveAll()
				messageContents.Add(widget.NewLabel("Couldn't read message, reason:"))
				messageContents.Add(widget.NewLabelWithStyle(err.Error(), fyne.TextAlignCenter, widget.RichTextStyleBlockquote.TextStyle))
			} else {
				messageContents.Add(widget.NewRichTextWithText("Title: " + message.Title + "\nFrom: " + message.Author))
				messageContents.Add(widget.NewRichTextWithText(messagecontent.RawContent))
			}
		})))
		if i != len(messages) {
			messageContainer.Add(widget.NewSeparator())
		}
	}
}

func Mainui(a fyne.App, api *mobiapi.MobiAPI) {
	w := a.NewWindow("mobiNG")
	w.Resize(fyne.NewSize(812, 620))

	searchWidget := widget.NewEntry()
	searchWidget.SetPlaceHolder("Search")
	commitSearchWidget := widget.NewButton("Search", func() {})
	searchContainer := container.NewHSplit(searchWidget, commitSearchWidget)
	searchContainer.SetOffset(1)
	messageContainer := container.NewVBox()
	messageContents := container.NewVBox()
	go messagerefresh(api, messageContainer, messageContents)

	tabContainer := container.NewAppTabs(
		container.NewTabItem("Messages", container.NewVBox(searchContainer, container.NewHSplit(container.NewVScroll(messageContainer), messageContents))),
	)
	tabContainer.OnSelected = func(ti *container.TabItem) {
		if ti.Text == "Messages" {
			go messagerefresh(api, messageContainer, messageContents)
		}
	}

	w.SetContent(tabContainer)
	w.Show()
}
