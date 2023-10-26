package internal

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dark-steveneq/mobiapi"
)

func Loginui(a fyne.App, api *mobiapi.MobiAPI) {
	w := a.NewWindow("Login - mobiNG")
	w.CenterOnScreen()

	// Info
	infoWidget := widget.NewLabel("")

	// Proxy
	proxyToggleWidget := widget.NewCheck("Proxy", nil)
	proxyAddrWidget := widget.NewEntry()
	proxySSLChkWidget := widget.NewCheck("Any SSL", nil)
	proxySettingsContainer := container.NewHSplit(proxyAddrWidget, proxySSLChkWidget)
	proxyContainer := container.NewHSplit(proxyToggleWidget, proxySettingsContainer)

	proxyToggleWidget.SetChecked(false)
	proxyToggleWidget.OnChanged = func(b bool) {
		if b {
			proxyAddrWidget.Enable()
			proxySSLChkWidget.Enable()
		} else {
			proxyAddrWidget.Disable()
			proxySSLChkWidget.Disable()
		}
	}
	proxyAddrWidget.SetPlaceHolder("http://localhost:8080")
	proxyAddrWidget.Disable()
	proxySSLChkWidget.Disable()
	proxySettingsContainer.SetOffset(1)
	proxyContainer.SetOffset(0.3)

	// Token
	tokenDomWidget := widget.NewEntry()
	tokenTokenWidget := widget.NewEntry()
	tokenLoginWidget := widget.NewButton("Login with token", nil)

	tokenDomWidget.SetPlaceHolder("Domain")
	tokenDomWidget.OnChanged = func(_ string) {
		if len(tokenDomWidget.Text) != 0 && len(tokenTokenWidget.Text) == 26 {
			tokenLoginWidget.Enable()
		} else {
			tokenLoginWidget.Disable()
		}
	}
	tokenTokenWidget.SetPlaceHolder("Token")
	tokenTokenWidget.OnChanged = tokenDomWidget.OnChanged
	tokenLoginWidget.Disable()

	tokenTokenWidget.OnSubmitted = func(_ string) {
		tokenDomWidget.Disable()
		tokenTokenWidget.Disable()
		tokenLoginWidget.Disable()
		proxyToggleWidget.Disable()
		proxyAddrWidget.Disable()
		proxySSLChkWidget.Disable()
		infoWidget.SetText("")
		if proxyToggleWidget.Checked {
			if err := api.SetupProxy(proxyAddrWidget.Text, proxySSLChkWidget.Checked); err != nil {
				tokenDomWidget.Enable()
				tokenTokenWidget.Enable()
				tokenLoginWidget.Enable()
				proxyToggleWidget.Enable()
				proxyAddrWidget.Enable()
				proxySSLChkWidget.Enable()
				infoWidget.SetText("Couldn't setup proxy!")
				log.Println("Couldn't setup proxy, reason:", err)
				return
			}
		}
		if err := api.SetDomain(tokenDomWidget.Text); err != nil {
			tokenDomWidget.Enable()
			tokenTokenWidget.Enable()
			tokenLoginWidget.Enable()
			proxyToggleWidget.Enable()
			proxyAddrWidget.Enable()
			proxySSLChkWidget.Enable()
			infoWidget.SetText("Couldn't find instance")
			log.Println("Couldn't instance at provided domain, reason:", err)
			return
		}
		if signedin, err := api.TokenAuth(tokenTokenWidget.Text); err != nil || !signedin {
			infoWidget.SetText("Couldn't authenticate")
			tokenDomWidget.Enable()
			tokenTokenWidget.Enable()
			tokenLoginWidget.Enable()
			proxyToggleWidget.Enable()
			proxyAddrWidget.Enable()
			proxySSLChkWidget.Enable()
			if err != nil {
				log.Println("Couldn't authenticate, reason:", err)
			} else {
				log.Println("Couldn't authenticate, unknown reason")
			}
			return
		}
		Mainui(a, api)
		w.Close()
	}
	tokenDomWidget.OnSubmitted = tokenTokenWidget.OnSubmitted
	tokenLoginWidget.OnTapped = func() {
		tokenTokenWidget.OnSubmitted("")
	}

	tokenContainer := container.NewVBox(infoWidget, widget.NewLabel("Login with token"), tokenDomWidget, tokenTokenWidget, tokenLoginWidget)

	// Login
	loginDomWidget := widget.NewEntry()
	loginUserWidget := widget.NewEntry()
	loginPassWidget := widget.NewPasswordEntry()
	loginLoginWidget := widget.NewButton("Login", nil)

	loginDomWidget.SetPlaceHolder("Domain")
	loginDomWidget.OnChanged = func(_ string) {}
	loginUserWidget.SetPlaceHolder("Login")
	loginPassWidget.SetPlaceHolder("Password")
	loginLoginWidget.Disable()

	loginContainer := container.NewVBox(infoWidget, widget.NewLabel("Login to mobiDziennik"), loginDomWidget, loginUserWidget, loginPassWidget, loginLoginWidget)

	loginLoginWidget.OnTapped = func() {
		loginDomWidget.Disable()
		loginUserWidget.Disable()
		loginPassWidget.Disable()
		loginLoginWidget.Disable()
		proxyToggleWidget.Disable()
		proxyAddrWidget.Disable()
		proxySSLChkWidget.Disable()
		infoWidget.SetText("")
		if proxyToggleWidget.Checked {
			if err := api.SetupProxy(proxyAddrWidget.Text, proxySSLChkWidget.Checked); err != nil {
				loginDomWidget.Enable()
				loginUserWidget.Enable()
				loginPassWidget.Enable()
				loginPassWidget.SetText("")
				loginLoginWidget.Disable()
				proxyToggleWidget.Enable()
				proxyAddrWidget.Enable()
				proxySSLChkWidget.Enable()
				infoWidget.SetText("Couldn't setup proxy!")
				log.Println("Couldn't setup proxy, reason:", err)
				return
			}
		}

		if signedin, err := api.PasswordAuth(loginUserWidget.Text, loginPassWidget.Text); err != nil || !signedin {
			loginDomWidget.Enable()
			loginUserWidget.Enable()
			loginPassWidget.Enable()
			loginPassWidget.SetText("")
			loginLoginWidget.Disable()
			proxyToggleWidget.Enable()
			proxyAddrWidget.Enable()
			proxySSLChkWidget.Enable()
			infoWidget.SetText("Couldn't login!")
			if err != nil {
				log.Println("Couldn't login, reason:", err)
			} else {
				log.Println("Couldn't login, unknown reason")
			}
			return
		}

		Mainui(a, api)
		w.Close()
	}
	loginDomWidget.OnChanged = func(_ string) {
		if loginDomWidget.Text != "" && loginUserWidget.Text != "" && loginPassWidget.Text != "" {
			loginLoginWidget.Enable()
		} else {
			loginLoginWidget.Disable()
		}
	}
	loginDomWidget.OnSubmitted = func(_ string) {
		loginLoginWidget.OnTapped()
	}
	loginUserWidget.OnChanged = loginDomWidget.OnChanged
	loginUserWidget.OnSubmitted = func(_ string) {
		loginLoginWidget.OnTapped()
	}
	loginPassWidget.OnChanged = loginDomWidget.OnChanged
	loginPassWidget.OnSubmitted = func(_ string) {
		loginLoginWidget.OnTapped()
	}

	w.SetContent(container.NewVBox(container.NewVSplit(container.NewAppTabs(container.NewTabItem("Password", container.NewCenter(loginContainer)), container.NewTabItem("Token", container.NewCenter(tokenContainer))), proxyContainer)))
	w.SetFixedSize(true)
	w.Show()
}
