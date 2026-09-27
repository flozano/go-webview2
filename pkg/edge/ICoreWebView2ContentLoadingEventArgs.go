package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// ContentLoading is when a document actually starts being this page,
// which is not the same moment as a navigation starting.
//
// A navigation that starts may never land: a download, a 204, a server
// that answers with nothing to render. The page on screen stays. So a
// host that treats «a navigation began» as «the document changed» tears
// down the old document's state while the old document is still there.

type _ICoreWebView2ContentLoadingEventArgsVtbl struct {
	_IUnknownVtbl
	GetIsErrorPage  ComProc
	GetNavigationId ComProc
}

type ICoreWebView2ContentLoadingEventArgs struct {
	vtbl *_ICoreWebView2ContentLoadingEventArgsVtbl
}

func (i *ICoreWebView2ContentLoadingEventArgs) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return r
}

// GetIsErrorPage reports whether what is loading is WebView2's own
// error page rather than anything a server sent.
func (i *ICoreWebView2ContentLoadingEventArgs) GetIsErrorPage() (bool, error) {
	var v int32
	_, _, err := i.vtbl.GetIsErrorPage.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&v)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return v != 0, nil
}

func (i *ICoreWebView2ContentLoadingEventArgs) GetNavigationId() (uint64, error) {
	var v uint64
	_, _, err := i.vtbl.GetNavigationId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&v)),
	)
	if err != windows.ERROR_SUCCESS {
		return 0, err
	}
	return v, nil
}
