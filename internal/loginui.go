package internal

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dark-steveneq/mobiapi"
	"github.com/zalando/go-keyring"
)

func Loginui(a fyne.App, api *mobiapi.MobiAPI) {
	w := a.NewWindow("Login - mobiNG")
	w.CenterOnScreen()
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(240, 320))

	// Settings
	/// Proxy
	proxyAddrWidget := widget.NewEntry()
	proxyAnySSLWidget := widget.NewCheck("Allow invalid SSL certificates", nil)
	proxyToggleWidget := widget.NewCheck("Use Proxy", func(b bool) {
		if b {
			proxyAddrWidget.Enable()
			proxyAnySSLWidget.Enable()
		} else {
			proxyAddrWidget.Disable()
			proxyAnySSLWidget.Disable()
		}
	})

	proxyAddrWidget.SetPlaceHolder("http://localhost:8080")
	proxyToggleWidget.OnChanged(proxyToggleWidget.Checked)
	settingsContainer := container.NewVBox(
		widget.NewLabelWithStyle("Proxy", fyne.TextAlignLeading, widget.RichTextStyleBlockquote.TextStyle),
		proxyToggleWidget,
		widget.NewLabel("Proxy URL"),
		proxyAddrWidget,
		proxyAnySSLWidget,
	)

	// Shared
	sharedDomWidget := widget.NewEntry()
	sharedAuthWidget := widget.NewButton("Login", nil)

	sharedDomWidget.SetPlaceHolder("Domain")

	// Password
	passUserWidget := widget.NewEntry()
	passPassWidget := widget.NewPasswordEntry()

	passUserWidget.SetPlaceHolder("Username")
	passPassWidget.SetPlaceHolder("Password")
	passContainer := container.NewVBox(
		sharedDomWidget,
		passUserWidget,
		passPassWidget,
		sharedAuthWidget,
	)

	// Token
	tokenTokenWidget := widget.NewPasswordEntry()
	tokenTokenWidget.SetPlaceHolder("Token")
	tokenContainer := container.NewVBox(
		sharedDomWidget,
		tokenTokenWidget,
		sharedAuthWidget,
	)

	// Tabs
	rootContainer := container.NewAppTabs(
		container.NewTabItem("Password", passContainer),
		container.NewTabItem("Token", tokenContainer),
		container.NewTabItem("Settings", settingsContainer),
	)

	// Functions
	funcLock := func() {
		for index := range rootContainer.Items {
			if index != rootContainer.SelectedIndex() {
				rootContainer.DisableIndex(index)
			}
		}
		sharedDomWidget.Disable()
		passUserWidget.Disable()
		passPassWidget.Disable()
		sharedAuthWidget.Disable()
	}

	funcUnlock := func() {
		for index := range rootContainer.Items {
			rootContainer.EnableIndex(index)
		}
		sharedDomWidget.Enable()
		passUserWidget.Enable()
		passPassWidget.Enable()
		sharedAuthWidget.Enable()
	}

	funcLoadKeyring := func() {
		var restoreStrings = map[string]func(text string){
			"proxyAddress": proxyAddrWidget.SetText,
			"domain":       sharedDomWidget.SetText,
			"username":     passUserWidget.SetText,
			"password":     passPassWidget.SetText,
		}
		var restoreBools = map[string]*bool{
			"useProxy": &proxyToggleWidget.Checked,
			"proxySSL": &proxyAnySSLWidget.Checked,
		}
		for key, ret := range restoreStrings {
			if value, err := keyring.Get("mobing", key); err == nil {
				ret(value)
			} else {
				return
			}
		}
		for key, ret := range restoreBools {
			if svalue, err := keyring.Get("mobing", key); err == nil {
				value, err := strconv.ParseBool(svalue)
				if err != nil {
					return
				}
				*ret = value
			} else {
				return
			}
		}
		proxyToggleWidget.OnChanged(proxyToggleWidget.Checked)
	}

	funcSaveKeyring := func() {
		// Proxy
		keyring.Set("mobing", "useProxy", strconv.FormatBool(proxyToggleWidget.Checked))
		keyring.Set("mobing", "proxyAddress", proxyAddrWidget.Text)
		keyring.Set("mobing", "proxySSL", strconv.FormatBool(proxyAnySSLWidget.Checked))

		// Shared
		keyring.Set("mobing", "domain", sharedDomWidget.Text)

		// Password
		keyring.Set("mobing", "username", passUserWidget.Text)
		keyring.Set("mobing", "password", passPassWidget.Text)
	}

	funcValAuthAccess := func(_ string) {
		if sharedDomWidget.Text == "" {
			sharedAuthWidget.Disable()
			return
		}
		switch rootContainer.Selected().Text {
		case "Password":
			if passUserWidget.Text == "" || passPassWidget.Text == "" {
				sharedAuthWidget.Disable()
				return
			}
		case "Token":
			if len(tokenTokenWidget.Text) != 26 {
				sharedAuthWidget.Disable()
				return
			}
		default:
			return
		}
		sharedAuthWidget.Enable()
	}

	funcAuthenticate := func(_ string) {
		funcSaveKeyring()
		funcLock()
		if err := api.SetupProxy(proxyAddrWidget.Text, proxyAnySSLWidget.Checked); !proxyToggleWidget.Checked || err != nil {
			api.SetupProxy("", false)
			if proxyToggleWidget.Checked && err != nil {
				funcUnlock()
				a.SendNotification(fyne.NewNotification("Couldn't setup proxy", "Reason: "+err.Error()))
				return
			}
		}

		if err := api.SetDomain(sharedDomWidget.Text); err != nil {
			funcUnlock()
			a.SendNotification(fyne.NewNotification("Couldn't set domain", "Reason: "+err.Error()))
			return
		}

		switch rootContainer.Selected().Text {
		case "Password":
			signedin, err := api.PasswordAuth(passUserWidget.Text, passPassWidget.Text)
			if err != nil {
				funcUnlock()
				a.SendNotification(fyne.NewNotification("Couldn't login", "Reason: "+err.Error()))
				return
			} else if !signedin {
				funcUnlock()
				a.SendNotification(fyne.NewNotification("Couldn't login", "Unknown reason"))
				return
			}
		case "Token":
			signedin, err := api.TokenAuth(tokenTokenWidget.Text)
			if err != nil {
				funcUnlock()
				a.SendNotification(fyne.NewNotification("Couldn't login", "Reason: "+err.Error()))
				return
			} else if !signedin {
				funcUnlock()
				a.SendNotification(fyne.NewNotification("Couldn't login", "Unknown reason"))
				return
			}
		default:
			funcUnlock()
			a.SendNotification(fyne.NewNotification("Authorize default switch case", rootContainer.Selected().Text))
			return
		}
		Mainui(a, api)
		w.Close()
	}

	// Function Calls
	sharedDomWidget.OnSubmitted = funcAuthenticate
	sharedDomWidget.OnChanged = funcValAuthAccess
	sharedAuthWidget.OnTapped = func() {
		funcAuthenticate("")
	}
	passUserWidget.OnSubmitted = funcAuthenticate
	passUserWidget.OnChanged = funcValAuthAccess
	passPassWidget.OnSubmitted = funcAuthenticate
	passPassWidget.OnChanged = funcValAuthAccess
	passPassWidget.OnChanged = funcValAuthAccess
	tokenTokenWidget.OnSubmitted = funcAuthenticate
	tokenTokenWidget.OnChanged = funcValAuthAccess
	rootContainer.OnSelected = func(ti *container.TabItem) {
		funcValAuthAccess("")
		funcSaveKeyring()
	}

	funcLoadKeyring()
	funcValAuthAccess("")
	w.SetContent(rootContainer)
	w.Show()
}
