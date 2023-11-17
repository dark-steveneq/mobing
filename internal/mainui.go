package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dark-steveneq/mobiapi"
)

func messageRefresh(a fyne.App, w fyne.Window, api *mobiapi.MobiAPI, list *fyne.Container, info *widget.RichText, cont *widget.Entry) {
	list.RemoveAll()
	list.Add(widget.NewLabelWithStyle("Loading...", fyne.TextAlignCenter, widget.RichTextStyleHeading.TextStyle))
	list.Add(widget.NewProgressBarInfinite())
	messages, err := api.GetReceivedMessages(false)
	if err != nil {
		list.RemoveAll()
		list.Add(widget.NewLabelWithStyle("Couldn't get messages!", fyne.TextAlignCenter, widget.RichTextStyleHeading.TextStyle))
		return
	}
	list.RemoveAll()
	for _, rmessage := range messages {
		message := rmessage
		list.Add(newWidgetMessage(message.Title, message.Author, message.Read, func() {
			info.Segments = []widget.RichTextSegment{}
			cont.SetText("Loading message...")
			info.Refresh()
			cont.Refresh()
			messagecontent, err := api.GetMessageContent(message)
			if err != nil {
				info.Segments = []widget.RichTextSegment{}
				info.Refresh()
				cont.SetText("Couldn't read message, reason: " + err.Error())
				return
			}
			info.Segments = []widget.RichTextSegment{
				&widget.TextSegment{
					Text:  "Title:",
					Style: widget.RichTextStyleInline,
				},
				&widget.TextSegment{
					Text:  messagecontent.Info.Title,
					Style: widget.RichTextStyleParagraph,
				},
				&widget.TextSegment{
					Text:  "From:",
					Style: widget.RichTextStyleInline,
				},
				&widget.TextSegment{
					Text:  messagecontent.Info.Author,
					Style: widget.RichTextStyleParagraph,
				},
			}
			info.Refresh()
			cont.Refresh()
			cont.SetText(messagecontent.Content)
		}))
	}
}

func Mainui(a fyne.App, api *mobiapi.MobiAPI) {
	w := a.NewWindow("mobiNG")
	w.Resize(fyne.NewSize(812, 580))

	api.OnLostConnection = func() {
		a.SendNotification(fyne.NewNotification("Disconnected from mobiDziennik", "Sorry for the inconvinience!"))
		Loginui(a, api)
		w.Close()
	}

	// Messages
	messageListContainer := container.NewVBox()
	messageInfoWidget := widget.NewRichTextFromMarkdown("")
	messageContWidget := widget.NewMultiLineEntry()
	messageRSContainer := container.NewVSplit(
		messageInfoWidget,
		messageContWidget,
	)
	messageRootContainer := container.NewHSplit(
		container.NewVScroll(messageListContainer),
		messageRSContainer,
	)

	messageContWidget.Disable()
	messageRSContainer.SetOffset(0)
	messageRootContainer.SetOffset(0.5)

	// Tabs
	rootContainer := container.NewAppTabs(
		container.NewTabItem("Messages", messageRootContainer),
	)

	// Background
	go messageRefresh(a, w, api, messageListContainer, messageInfoWidget, messageContWidget)

	w.SetContent(rootContainer)
	w.Show()
}
