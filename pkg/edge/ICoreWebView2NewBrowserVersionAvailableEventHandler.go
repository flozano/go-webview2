package edge

// The runtime updates itself; an application that never closes keeps
// using the version it started with. This is how Windows says so, and
// the only way a till that is open all day ever moves off a browser
// whose replacement is already on disk.

type _ICoreWebView2NewBrowserVersionAvailableEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2NewBrowserVersionAvailableEventHandler struct {
	vtbl *_ICoreWebView2NewBrowserVersionAvailableEventHandlerVtbl
	impl _ICoreWebView2NewBrowserVersionAvailableEventHandlerImpl
}

func _ICoreWebView2NewBrowserVersionAvailableEventHandlerIUnknownQueryInterface(this *ICoreWebView2NewBrowserVersionAvailableEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2NewBrowserVersionAvailableEventHandlerIUnknownAddRef(this *ICoreWebView2NewBrowserVersionAvailableEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2NewBrowserVersionAvailableEventHandlerIUnknownRelease(this *ICoreWebView2NewBrowserVersionAvailableEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2NewBrowserVersionAvailableEventHandlerInvoke(this *ICoreWebView2NewBrowserVersionAvailableEventHandler, sender *ICoreWebView2Environment, args uintptr) uintptr {
	return this.impl.NewBrowserVersionAvailable(sender)
}

type _ICoreWebView2NewBrowserVersionAvailableEventHandlerImpl interface {
	_IUnknownImpl
	NewBrowserVersionAvailable(sender *ICoreWebView2Environment) uintptr
}

var _ICoreWebView2NewBrowserVersionAvailableEventHandlerFn = _ICoreWebView2NewBrowserVersionAvailableEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2NewBrowserVersionAvailableEventHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2NewBrowserVersionAvailableEventHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2NewBrowserVersionAvailableEventHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2NewBrowserVersionAvailableEventHandlerInvoke),
}

func newICoreWebView2NewBrowserVersionAvailableEventHandler(impl _ICoreWebView2NewBrowserVersionAvailableEventHandlerImpl) *ICoreWebView2NewBrowserVersionAvailableEventHandler {
	return &ICoreWebView2NewBrowserVersionAvailableEventHandler{
		vtbl: &_ICoreWebView2NewBrowserVersionAvailableEventHandlerFn,
		impl: impl,
	}
}
