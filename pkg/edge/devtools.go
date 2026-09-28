package edge

// The DevTools protocol, which is how a host asks the browser for the
// things the WebView2 API itself does not offer.
//
// Like hostguards.go, nothing new is mapped here: the slot is already
// in the vtable and was simply unreachable from outside the package.
//
// WHAT IT IS FOR, in the caller this fork exists for: a screenshot of
// the WHOLE page. CapturePreview, which the API does offer, captures
// the viewport -- what fits on the screen right now. A person at a
// counter reporting that a list came out wrong needs the list, not the
// part of it that fitted, and `Page.captureScreenshot` with
// `captureBeyondViewport` is the only way to get it.

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// CallDevToolsProtocolMethod runs one DevTools protocol method and
// hands the answer to done.
//
// ASYNCHRONOUS, and deliberately not wrapped into something that looks
// synchronous. The answer arrives on the thread that owns the browser,
// through its message loop; a wrapper that blocked waiting for it would
// have to pump that loop from inside a call the loop is making, which
// is how this package's caller once hung a window for good.
//
// done is called exactly once, unless this returns an error -- in which
// case the call never reached the browser and done is called not at
// all.
//
// parametersAsJSON is a JSON object, and "{}" when the method takes no
// parameters. An empty string is not valid JSON and the browser would
// answer with a protocol error.
func (i *ICoreWebView2) CallDevToolsProtocolMethod(method, parametersAsJSON string, done func(err error, resultJSON string)) error {
	if i == nil {
		return ErrNoBrowser
	}
	if done == nil {
		return fmt.Errorf("go-webview2: CallDevToolsProtocolMethod needs somewhere to put the answer")
	}
	_method, err := windows.UTF16PtrFromString(method)
	if err != nil {
		return fmt.Errorf("go-webview2: the method name is not valid UTF-16: %w", err)
	}
	if parametersAsJSON == "" {
		parametersAsJSON = "{}"
	}
	_params, err := windows.UTF16PtrFromString(parametersAsJSON)
	if err != nil {
		return fmt.Errorf("go-webview2: the parameters are not valid UTF-16: %w", err)
	}

	handler := newICoreWebView2CallDevToolsProtocolMethodCompletedHandler(
		func(errorCode uintptr, resultJSON string) {
			if errorCode != 0 {
				done(fmt.Errorf("go-webview2: %s failed: HRESULT 0x%08x", method, uint32(errorCode)), resultJSON)
				return
			}
			done(nil, resultJSON)
		})

	_, _, callErr := i.vtbl.CallDevToolsProtocolMethod.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_method)),
		uintptr(unsafe.Pointer(_params)),
		uintptr(unsafe.Pointer(handler)),
	)
	if callErr != windows.ERROR_SUCCESS {
		// The browser never took it, so the completion will never fire
		// and the handler would sit in the map for the life of the
		// process.
		handler.forget()
		return callErr
	}
	return nil
}

// CallDevToolsProtocolMethod on the browser this Chromium is driving,
// or ErrNoBrowser when there is none -- a degraded mode, a window that
// never came up -- rather than a nil dereference.
func (e *Chromium) CallDevToolsProtocolMethod(method, parametersAsJSON string, done func(err error, resultJSON string)) error {
	if e.webview == nil {
		return ErrNoBrowser
	}
	return e.webview.CallDevToolsProtocolMethod(method, parametersAsJSON, done)
}
