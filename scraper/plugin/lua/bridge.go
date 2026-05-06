package lua

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dokidokikoi/go-common/tools"
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

// Bridge provides Go functions to Lua scripts
type Bridge struct {
	logger *zap.Logger
}

// NewBridge creates a new Lua bridge
func NewBridge(logger *zap.Logger) *Bridge {
	return &Bridge{
		logger: logger,
	}
}

// SetupSandbox configures the Lua VM with a secure sandboxed environment
func (b *Bridge) SetupSandbox(L *lua.LState) {
	// In gopher-lua, we need to manually remove dangerous globals
	// Remove dangerous modules by setting them to nil
	L.SetGlobal("os", lua.LNil)
	L.SetGlobal("io", lua.LNil)
	L.SetGlobal("package", lua.LNil)
	L.SetGlobal("debug", lua.LNil)
	L.SetGlobal("loadfile", lua.LNil)
	L.SetGlobal("dofile", lua.LNil)
	L.SetGlobal("load", lua.LNil)

	// Create safe modules
	httpModule := b.createHTTPModule(L)
	htmlModule := b.createHTMLModule(L)
	logModule := b.createLogModule(L)
	jsonModule := b.createJSONModule(L)
	stringModule := L.GetGlobal("string")
	tableModule := L.GetGlobal("table")
	mathModule := L.GetGlobal("math")

	// Set safe globals
	L.SetGlobal("http", httpModule)
	L.SetGlobal("html", htmlModule)
	L.SetGlobal("log", logModule)
	L.SetGlobal("json", jsonModule)
	L.SetGlobal("string", stringModule)
	L.SetGlobal("table", tableModule)
	L.SetGlobal("math", mathModule)
}

// createHTTPModule creates the http module for Lua
func (b *Bridge) createHTTPModule(L *lua.LState) *lua.LTable {
	module := L.NewTable()

	L.SetField(module, "get", L.NewFunction(b.httpGet))
	L.SetField(module, "post", L.NewFunction(b.httpPost))

	return module
}

// createHTMLModule creates the html module for Lua
func (b *Bridge) createHTMLModule(L *lua.LState) *lua.LTable {
	module := L.NewTable()

	L.SetField(module, "parse", L.NewFunction(b.htmlParse))

	return module
}

// createLogModule creates the log module for Lua
func (b *Bridge) createLogModule(L *lua.LState) *lua.LTable {
	module := L.NewTable()

	L.SetField(module, "info", L.NewFunction(b.logInfo))
	L.SetField(module, "error", L.NewFunction(b.logError))
	L.SetField(module, "warn", L.NewFunction(b.logWarn))
	L.SetField(module, "debug", L.NewFunction(b.logDebug))

	return module
}

// createJSONModule creates the json module for Lua
func (b *Bridge) createJSONModule(L *lua.LState) *lua.LTable {
	module := L.NewTable()

	L.SetField(module, "encode", L.NewFunction(b.jsonEncode))
	L.SetField(module, "decode", L.NewFunction(b.jsonDecode))

	return module
}

// HTTP functions

func (b *Bridge) httpGet(L *lua.LState) int {
	url := L.CheckString(1)

	headers := make(map[string]string)
	if L.GetTop() >= 2 {
		if tbl, ok := L.CheckAny(2).(*lua.LTable); ok {
			headers = luaTableToMap(L, tbl)
		}
	}

	proxy := ""
	if L.GetTop() >= 3 {
		proxy = L.CheckString(3)
	}

	rsp, err := b.doRequest(http.MethodGet, url, headers, "", proxy)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LNumber(rsp.StatusCode))
	L.Push(lua.LString(rsp.Body))
	return 2
}

func (b *Bridge) httpPost(L *lua.LState) int {
	url := L.CheckString(1)

	headers := make(map[string]string)
	if L.GetTop() >= 2 {
		if tbl, ok := L.CheckAny(2).(*lua.LTable); ok {
			headers = luaTableToMap(L, tbl)
		}
	}

	body := ""
	if L.GetTop() >= 3 {
		body = L.CheckString(3)
	}

	proxy := ""
	if L.GetTop() >= 4 {
		proxy = L.CheckString(4)
	}

	rsp, err := b.doRequest(http.MethodPost, url, headers, body, proxy)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LNumber(rsp.StatusCode))
	L.Push(lua.LString(rsp.Body))
	return 2
}

type httpResponse struct {
	StatusCode int
	Body       string
}

func (b *Bridge) doRequest(method, url string, headers map[string]string, body, proxy string) (*httpResponse, error) {
	rsp, err := tools.ReqWithProxy(method, url, body, proxy, tools.SetHeadersWithOption(headers))
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode()/100 != 2 {
		return nil, errors.Errorf("status code: %d, body: %s", rsp.StatusCode(), rsp.String())
	}

	return &httpResponse{
		StatusCode: rsp.StatusCode(),
		Body:       rsp.String(),
	}, nil
}

// HTML functions

func (b *Bridge) htmlParse(L *lua.LState) int {
	body := L.CheckString(1)

	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(body))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// Create a userdata to hold the document
	ud := L.NewUserData()
	ud.Value = doc

	// Set metatable with __index method
	mt := L.NewTable()
	L.SetField(mt, "__index", L.NewFunction(b.htmlDocumentIndex))
	L.SetMetatable(ud, mt)

	L.Push(ud)
	return 1
}

func (b *Bridge) htmlDocumentIndex(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if ud.Value == nil {
		L.Push(lua.LNil)
		return 1
	}

	doc, ok := ud.Value.(*goquery.Document)
	if !ok {
		L.Push(lua.LNil)
		return 1
	}

	key := L.CheckString(2)

	if key == "select" {
		selector := L.CheckString(3)
		selection := doc.Find(selector)
		return b.pushSelection(L, selection)
	}

	L.Push(lua.LNil)
	return 1
}

func (b *Bridge) pushSelection(L *lua.LState, sel *goquery.Selection) int {
	ud := L.NewUserData()
	ud.Value = sel

	mt := L.NewTable()
	L.SetField(mt, "__index", L.NewFunction(b.htmlSelectionIndex))
	L.SetMetatable(ud, mt)

	L.Push(ud)
	return 1
}

func (b *Bridge) htmlSelectionIndex(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if ud.Value == nil {
		L.Push(lua.LNil)
		return 1
	}

	sel, ok := ud.Value.(*goquery.Selection)
	if !ok {
		L.Push(lua.LNil)
		return 1
	}

	key := L.CheckString(2)

	switch key {
	case "select":
		selector := L.CheckString(3)
		selection := sel.Find(selector)
		return b.pushSelection(L, selection)
	case "text":
		L.Push(lua.LString(sel.Text()))
		return 1
	case "html":
		html, _ := sel.Html()
		L.Push(lua.LString(html))
		return 1
	case "attr":
		name := L.CheckString(3)
		attr, _ := sel.Attr(name)
		L.Push(lua.LString(attr))
		return 1
	case "each":
		callback := L.CheckFunction(3)
		idx := 1
		sel.Each(func(i int, s *goquery.Selection) {
			L.Push(callback)
			L.Push(lua.LNumber(i + 1))
			b.pushSelection(L, s)
			L.Call(2, 0)
			idx++
		})
		return 0
	case "length":
		L.Push(lua.LNumber(sel.Length()))
		return 1
	case "first":
		return b.pushSelection(L, sel.First())
	case "last":
		return b.pushSelection(L, sel.Last())
	case "eq":
		idx := L.CheckInt(3)
		return b.pushSelection(L, sel.Eq(idx))
	case "next":
		return b.pushSelection(L, sel.Next())
	case "prev":
		return b.pushSelection(L, sel.Prev())
	case "parent":
		return b.pushSelection(L, sel.Parent())
	case "children":
		return b.pushSelection(L, sel.Children())
	case "siblings":
		return b.pushSelection(L, sel.Siblings())
	default:
		L.Push(lua.LNil)
	}

	return 1
}

// Log functions

func (b *Bridge) logInfo(L *lua.LState) int {
	msg := L.CheckString(1)
	b.logger.Info(msg)
	return 0
}

func (b *Bridge) logError(L *lua.LState) int {
	msg := L.CheckString(1)
	b.logger.Error(msg)
	return 0
}

func (b *Bridge) logWarn(L *lua.LState) int {
	msg := L.CheckString(1)
	b.logger.Warn(msg)
	return 0
}

func (b *Bridge) logDebug(L *lua.LState) int {
	msg := L.CheckString(1)
	b.logger.Debug(msg)
	return 0
}

// JSON functions

func (b *Bridge) jsonEncode(L *lua.LState) int {
	tbl := L.CheckTable(1)

	result := b.tableToJSON(L, tbl)
	L.Push(lua.LString(result))
	return 1
}

func (b *Bridge) tableToJSON(L *lua.LState, tbl *lua.LTable) string {
	if tbl == nil {
		return "null"
	}

	// Check if it's an array (sequential keys)
	isArray := true
	maxIdx := 0
	tbl.ForEach(func(key, value lua.LValue) {
		if numKey, ok := key.(lua.LNumber); ok {
			idx := int(numKey)
			if idx > maxIdx {
				maxIdx = idx
			}
		} else {
			isArray = false
		}
	})

	var builder strings.Builder

	if isArray && maxIdx > 0 {
		builder.WriteString("[")
		for i := 1; i <= maxIdx; i++ {
			if i > 1 {
				builder.WriteString(",")
			}
			value := tbl.RawGetInt(i)
			builder.WriteString(b.valueToJSON(L, value))
		}
		builder.WriteString("]")
	} else {
		builder.WriteString("{")
		first := true
		tbl.ForEach(func(key, value lua.LValue) {
			if !first {
				builder.WriteString(",")
			}
			first = false

			if strKey, ok := key.(lua.LString); ok {
				builder.WriteString(fmt.Sprintf("\"%s\":", string(strKey)))
			}
			builder.WriteString(b.valueToJSON(L, value))
		})
		builder.WriteString("}")
	}

	return builder.String()
}

func (b *Bridge) valueToJSON(L *lua.LState, value lua.LValue) string {
	switch v := value.(type) {
	case lua.LString:
		return fmt.Sprintf("\"%s\"", string(v))
	case lua.LNumber:
		return fmt.Sprintf("%v", float64(v))
	case lua.LBool:
		if bool(v) {
			return "true"
		}
		return "false"
	case *lua.LTable:
		return b.tableToJSON(L, v)
	default:
		return "null"
	}
}

func (b *Bridge) jsonDecode(L *lua.LState) int {
	jsonStr := L.CheckString(1)

	// Simple JSON decoder - for complex JSON, consider using a proper library
	// For now, return the string as-is
	L.Push(lua.LString(jsonStr))
	return 1
}