package context

import (
	"public/message"
	"public/tools"
	"sync"

	"github.com/ant0ine/go-json-rest/rest"
)

//
var M_context Context

//
type Context struct {
}

var m sync.Map

/*
	更新
*/
func (a *Context) OnUpdate(w rest.ResponseWriter, r *rest.Request) {
	var req map[string]string
	tools.GetRequestJsonObj(r, &req)

	context := req["context"]
	access_token := req["access_token"]

	if !tools.CheckParam(access_token) {
		w.WriteJson(message.GetErrorMsg(message.ParameterInvalid))
		return
	}

	m.Store(access_token, context)

	w.WriteJson(message.GetSuccessMsg())
}

/*
	更新
*/
func (a *Context) OnGet(w rest.ResponseWriter, r *rest.Request) {
	var req map[string]string
	tools.GetRequestJsonObj(r, &req)

	access_token := req["access_token"]

	if !tools.CheckParam(access_token) {
		w.WriteJson(message.GetErrorMsg(message.ParameterInvalid))
		return
	}

	v, _ := m.Load(access_token)
	msg := message.GetSuccessMsg()
	msg.Data = v

	w.WriteJson(msg)
}
