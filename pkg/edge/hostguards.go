package edge

// What a host needs to be a host, rather than a frame that shows
// whatever it is handed.
//
// Everything in this file is a Go wrapper over a slot that the vtables
// in this package already declare. Nothing new is mapped: the COM was
// always there, it simply could not be reached from outside the package
// because the vtbl fields are unexported and no method called them.
//
// Three questions, and a host that cannot ask them is not one:
//
//   - where is this navigation going, and may I refuse it;
//   - somebody asked for a new window, may I open it myself instead;
//   - what page is on screen right now, when something calls in.
//
// Plus the one that is not about trust: the runtime has been replaced
// under a process that never closes.

import (
	"errors"
	"unsafe"

	"github.com/jchv/go-webview2/internal/w32"

	"golang.org/x/sys/windows"
)

// GetSource is the URL of the document currently loaded.
//
// Read at the moment a call arrives, and not remembered from the last
// navigation: what a page says it is and what it is are the same thing
// only if nobody is trying.
func (i *ICoreWebView2) GetSource() (string, error) {
	var _uri *uint16
	_, _, err := i.vtbl.GetSource.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_uri)),
	)
	if err != windows.ERROR_SUCCESS {
		return "", err
	}
	uri := windows.UTF16PtrToString(_uri)
	windows.CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, nil
}

func (i *ICoreWebView2) AddNavigationStarting(eventHandler *ICoreWebView2NavigationStartingEventHandler, token *_EventRegistrationToken) error {
	_, _, err := i.vtbl.AddNavigationStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2) AddNewWindowRequested(eventHandler *ICoreWebView2NewWindowRequestedEventHandler, token *_EventRegistrationToken) error {
	_, _, err := i.vtbl.AddNewWindowRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2) AddContentLoading(eventHandler *ICoreWebView2ContentLoadingEventHandler, token *_EventRegistrationToken) error {
	_, _, err := i.vtbl.AddContentLoading.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (e *ICoreWebView2Environment) AddNewBrowserVersionAvailable(eventHandler *ICoreWebView2NewBrowserVersionAvailableEventHandler, token *_EventRegistrationToken) error {
	_, _, err := e.vtbl.AddNewBrowserVersionAvailable.Call(
		uintptr(unsafe.Pointer(e)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

// Source is the current URL, asked of the webview this Chromium drives.
//
// Empty before the browser exists, which is not an error and is not
// permission either: a caller that gets "" has learned nothing and must
// refuse, like any other URL that is not the one expected.
func (e *Chromium) Source() (string, error) {
	if e.webview == nil {
		return "", nil
	}
	return e.webview.GetSource()
}

// NavigationStarting is the COM callback. The host's answer lives in
// NavigationStartingCallback, which may cancel through the args.
func (e *Chromium) NavigationStarting(sender *ICoreWebView2, args *ICoreWebView2NavigationStartingEventArgs) uintptr {
	if e.NavigationStartingCallback != nil {
		e.NavigationStartingCallback(sender, args)
	}
	return 0
}

func (e *Chromium) NewWindowRequested(sender *ICoreWebView2, args *ICoreWebView2NewWindowRequestedEventArgs) uintptr {
	if e.NewWindowRequestedCallback != nil {
		e.NewWindowRequestedCallback(sender, args)
	}
	return 0
}

// ContentLoading is a document actually becoming this page. See the
// comment on its event args for why that is not NavigationStarting.
func (e *Chromium) ContentLoading(sender *ICoreWebView2, args *ICoreWebView2ContentLoadingEventArgs) uintptr {
	if e.ContentLoadingCallback != nil {
		e.ContentLoadingCallback(sender, args)
	}
	return 0
}

func (e *Chromium) NewBrowserVersionAvailable(sender *ICoreWebView2Environment) uintptr {
	if e.NewBrowserVersionAvailableCallback != nil {
		e.NewBrowserVersionAvailableCallback(sender)
	}
	return 0
}

// Reload is the button a browser has and a host window does not.
//
// A page that will not draw is the commonest thing that goes wrong in
// front of a person, and in a browser they press F5 without thinking.
// A host that embeds a webview and offers no way to reload leaves them
// with nothing to try -- and the keyboard is not a way, because the
// accelerators only reach the browser when it has the focus, which a
// host window does not hand over by itself.
func (i *ICoreWebView2) Reload() error {
	_, _, err := i.vtbl.Reload.Call(uintptr(unsafe.Pointer(i)))
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

// Stop is Reload's pair: what a person presses when a page is taking
// too long and they would rather have the old one back.
func (i *ICoreWebView2) Stop() error {
	_, _, err := i.vtbl.Stop.Call(uintptr(unsafe.Pointer(i)))
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

// Reload and Stop, on the Chromium, for callers that never see the
// ICoreWebView2.
//
// NOT no-ops before the browser exists: they say so. A button that
// answers «fine» and does nothing is the failure that gets reported as
// «I pressed it and nothing happened», and the person reporting it is
// right and has nothing else to say. ErrNoBrowser is something a
// caller can log.
func (e *Chromium) Reload() error {
	if e.webview == nil {
		return ErrNoBrowser
	}
	return e.webview.Reload()
}

func (e *Chromium) Stop() error {
	if e.webview == nil {
		return ErrNoBrowser
	}
	return e.webview.Stop()
}

// ErrNoBrowser means the webview has not been created yet, or has gone.
var ErrNoBrowser = errors.New("go-webview2: there is no browser to ask")

// Controller is the webview's controller, for a host that has to place
// it somewhere other than the whole client area -- beside a toolbar,
// say -- or hand it the keyboard.
func (e *Chromium) Controller() *ICoreWebView2Controller { return e.controller }

// PutBoundsRect is PutBounds in numbers anybody can pass.
//
// PutBounds takes a w32.Rect, and w32 is an INTERNAL package: a caller
// outside this module cannot name the type, so the method it belongs
// to might as well not exist. Four coordinates are the whole argument,
// and they cross a package boundary without help.
func (i *ICoreWebView2Controller) PutBoundsRect(left, top, right, bottom int32) error {
	return i.PutBounds(w32.Rect{Left: left, Top: top, Right: right, Bottom: bottom})
}

// PlaceBelow is what a host with a strip of buttons on top actually
// wants: the page fills the client area except the first `top` pixels.
func (e *Chromium) PlaceBelow(left, top, right, bottom int32) error {
	if e.controller == nil {
		return nil
	}
	return e.controller.PutBoundsRect(left, top, right, bottom)
}

// GetSource on the browser this Chromium is driving, or ErrNoBrowser
// when there is none.
//
// ASK, RATHER THAN REMEMBER. A host that recorded the last URL it let
// through NavigationStarting would be right until the page changed
// route without navigating, which is what a single-page application
// does all day. The browser knows; nothing else does.
func (e *Chromium) GetSource() (string, error) {
	if e.webview == nil {
		return "", ErrNoBrowser
	}
	return e.webview.GetSource()
}

// GetDocumentTitle is what the page calls itself: its <title>.
//
// WHAT IT IS FOR. A host window that is not a browser still has to say
// which application is inside it, and the honest place to read that
// from is the page, not the host's own configuration -- one binary
// serves several tenants, and each should show its own name without
// anybody compiling it in.
//
// Read at the moment it is wanted, like GetSource, and re-read when
// DocumentTitleChanged says so: a single-page application changes its
// title without navigating.
func (i *ICoreWebView2) GetDocumentTitle() (string, error) {
	var _title *uint16
	_, _, err := i.vtbl.GetDocumentTitle.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_title)),
	)
	if err != windows.ERROR_SUCCESS {
		return "", err
	}
	title := windows.UTF16PtrToString(_title)
	windows.CoTaskMemFree(unsafe.Pointer(_title))
	return title, nil
}

// GetDocumentTitle on the browser this Chromium is driving, or
// ErrNoBrowser when there is none.
func (e *Chromium) GetDocumentTitle() (string, error) {
	if e.webview == nil {
		return "", ErrNoBrowser
	}
	return e.webview.GetDocumentTitle()
}

func (i *ICoreWebView2) AddDocumentTitleChanged(eventHandler *ICoreWebView2DocumentTitleChangedEventHandler, token *_EventRegistrationToken) error {
	_, _, err := i.vtbl.AddDocumentTitleChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

// DocumentTitleChanged is the page saying it now calls itself something
// else. The title itself is not in the event: it is read back with
// GetDocumentTitle, which is how the API declares it.
func (e *Chromium) DocumentTitleChanged(sender *ICoreWebView2) uintptr {
	if e.DocumentTitleChangedCallback != nil {
		e.DocumentTitleChangedCallback(sender)
	}
	return 0
}
