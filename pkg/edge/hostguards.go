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
	"unsafe"

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
