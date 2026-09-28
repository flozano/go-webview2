package edge

import (
	"sync"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2CallDevToolsProtocolMethodCompletedHandler struct {
	vtbl *_ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerVtbl
	id   uintptr
}

// The handler has to outlive the call that made it, and it is COM on
// the other side of it: the browser keeps the pointer, calls Invoke
// whenever the answer is ready, and lets go. A Go closure handed over
// as a raw pointer is a closure the collector knows nothing about, so
// the live ones are held here, by a number the handler carries, and
// dropped when Invoke has run.
//
// One-shot on purpose: this completion fires exactly once per call, and
// a handler that stayed in the map after it would be a leak that grows
// with every use.
var devTools struct {
	sync.Mutex
	next uintptr
	live map[uintptr]func(errorCode uintptr, resultJSON string)
}

func _ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownQueryInterface(_ *iCoreWebView2CallDevToolsProtocolMethodCompletedHandler, _, _ uintptr) uintptr {
	return 0
}

func _ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownAddRef(_ *iCoreWebView2CallDevToolsProtocolMethodCompletedHandler) uintptr {
	return 1
}

func _ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownRelease(_ *iCoreWebView2CallDevToolsProtocolMethodCompletedHandler) uintptr {
	return 1
}

func _ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerInvoke(this *iCoreWebView2CallDevToolsProtocolMethodCompletedHandler, errorCode uintptr, returnObjectAsJson *uint16) uintptr {
	devTools.Lock()
	done := devTools.live[this.id]
	delete(devTools.live, this.id)
	devTools.Unlock()
	if done == nil {
		return 0
	}
	var result string
	if returnObjectAsJson != nil {
		result = windows.UTF16PtrToString(returnObjectAsJson)
	}
	done(errorCode, result)
	return 0
}

var _ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerFn = _ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerInvoke),
}

func newICoreWebView2CallDevToolsProtocolMethodCompletedHandler(done func(errorCode uintptr, resultJSON string)) *iCoreWebView2CallDevToolsProtocolMethodCompletedHandler {
	devTools.Lock()
	if devTools.live == nil {
		devTools.live = make(map[uintptr]func(uintptr, string))
	}
	devTools.next++
	id := devTools.next
	devTools.live[id] = done
	devTools.Unlock()
	return &iCoreWebView2CallDevToolsProtocolMethodCompletedHandler{
		vtbl: &_ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerFn,
		id:   id,
	}
}

// forget drops a handler whose call never got as far as the browser, so
// that a failed Call does not leave an entry behind for ever.
func (h *iCoreWebView2CallDevToolsProtocolMethodCompletedHandler) forget() {
	devTools.Lock()
	delete(devTools.live, h.id)
	devTools.Unlock()
}
