package edge

type _ICoreWebView2DocumentTitleChangedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2DocumentTitleChangedEventHandler struct {
	vtbl *_ICoreWebView2DocumentTitleChangedEventHandlerVtbl
	impl _ICoreWebView2DocumentTitleChangedEventHandlerImpl
}

func _ICoreWebView2DocumentTitleChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2DocumentTitleChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2DocumentTitleChangedEventHandlerIUnknownAddRef(this *ICoreWebView2DocumentTitleChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2DocumentTitleChangedEventHandlerIUnknownRelease(this *ICoreWebView2DocumentTitleChangedEventHandler) uintptr {
	return this.impl.Release()
}

// The event carries no arguments of its own: the title is read back off
// the webview that sent it. That is how the API declares it, and it is
// also why this handler's Invoke takes an args pointer it never uses.
func _ICoreWebView2DocumentTitleChangedEventHandlerInvoke(this *ICoreWebView2DocumentTitleChangedEventHandler, sender *ICoreWebView2, args uintptr) uintptr {
	return this.impl.DocumentTitleChanged(sender)
}

type _ICoreWebView2DocumentTitleChangedEventHandlerImpl interface {
	_IUnknownImpl
	DocumentTitleChanged(sender *ICoreWebView2) uintptr
}

var _ICoreWebView2DocumentTitleChangedEventHandlerFn = _ICoreWebView2DocumentTitleChangedEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2DocumentTitleChangedEventHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2DocumentTitleChangedEventHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2DocumentTitleChangedEventHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2DocumentTitleChangedEventHandlerInvoke),
}

func newICoreWebView2DocumentTitleChangedEventHandler(impl _ICoreWebView2DocumentTitleChangedEventHandlerImpl) *ICoreWebView2DocumentTitleChangedEventHandler {
	return &ICoreWebView2DocumentTitleChangedEventHandler{
		vtbl: &_ICoreWebView2DocumentTitleChangedEventHandlerFn,
		impl: impl,
	}
}
