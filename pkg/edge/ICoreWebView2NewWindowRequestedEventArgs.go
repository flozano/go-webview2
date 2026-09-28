package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2NewWindowRequestedEventArgsVtbl struct {
	_IUnknownVtbl
	GetUri             ComProc
	PutNewWindow       ComProc
	GetNewWindow       ComProc
	PutHandled         ComProc
	GetHandled         ComProc
	GetIsUserInitiated ComProc
	GetDeferral        ComProc
	GetWindowFeatures  ComProc
}

type ICoreWebView2NewWindowRequestedEventArgs struct {
	vtbl *_ICoreWebView2NewWindowRequestedEventArgsVtbl
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return r
}

// GetUri is where the new window wanted to go: target=_blank, or
// window.open.
func (i *ICoreWebView2NewWindowRequestedEventArgs) GetUri() (string, error) {
	var _uri *uint16
	_, _, err := i.vtbl.GetUri.Call(
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

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetIsUserInitiated() (bool, error) {
	var v int32
	_, _, err := i.vtbl.GetIsUserInitiated.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&v)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return v != 0, nil
}

// PutHandled says the host dealt with it, so WebView2 opens nothing.
//
// Handled and cancelled are the same word here, and that is the API's
// doing, not ours: a host that claims a request and then does nothing
// with it has suppressed the window. Whoever claims it owes the person
// an answer -- opening the URL in their browser, say.
func (i *ICoreWebView2NewWindowRequestedEventArgs) PutHandled(handled bool) error {
	var v int32
	if handled {
		v = 1
	}
	_, _, err := i.vtbl.PutHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(v),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
