// fsb project fsb.go
package fsb

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const ( // iota is reset to 0
	JNodeTypeRaw  = iota // c0 == 0
	JNodeTypeHash = iota // c1 == 1
	JNodeTypeStr  = iota // c2 == 2
	JNodeTypeByte = iota // c3==3
)

type JNode struct {
	str    string
	binary []byte
	raw    json.RawMessage
	hash   map[string]interface{}
	jtype  int
}

func NewJNodeMap(m map[string]interface{}) *JNode {
	return &JNode{hash: m, jtype: JNodeTypeHash}
}
func NewJNodeRaw(r json.RawMessage) *JNode {
	return &JNode{raw: r, jtype: JNodeTypeRaw}
}

func NewJNodeString(s string) *JNode {
	return &JNode{binary: []byte(s), jtype: JNodeTypeByte}
}
func NewJNodeByte(b []byte) *JNode {
	return &JNode{binary: b, jtype: JNodeTypeByte}
}

func NewJNodeInterface(x interface{}) (*JNode, error) {
	b, err := json.Marshal(x)
	if err != nil {
		return nil, err
	}
	return &JNode{binary: b, jtype: JNodeTypeByte}, nil
}

func MapToRaw(m map[string]interface{}) (raw json.RawMessage, err error) {
	b, err := json.Marshal(m)
	if err != nil {
		return raw, err
	}

	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	err = decoder.Decode(&raw)
	//  err = json.Unmarshal(b,&raw)
	return raw, err
}

func (r *JNode) GetType() int {
	return r.jtype
}
func (r *JNode) GetBytes() (b []byte) {
	b, _ = r.MarshalJSON()
	return b
}
func (r *JNode) GetMap() (m map[string]interface{}, err error) {
	decoder := json.NewDecoder(bytes.NewReader(r.GetBytes()))
	decoder.UseNumber()
	err = decoder.Decode(&m)
	//   err = json.Unmarshal(r.GetBytes(),&m)
	return m, err
}

func (r *JNode) UnmarshalJSON(b []byte) (err error) {
	/*
		switch r.jtype {
			case JNodeTypeRaw:
			   return r.raw.UnmarshalJSON(b)
			case JNodeTypeHash:
		       decoder := json.NewDecoder(bytes.NewReader(b))
		       decoder.UseNumber()
		       return decoder.Decode(r.hash)
			case JNodeTypeByte:
			   r.binary=b
			   return nil
		} */
	r.jtype = JNodeTypeByte
	r.binary = b
	return nil
}
func (r *JNode) MarshalJSON() (b []byte, err error) {
	switch r.jtype {
	case JNodeTypeRaw:
		// return r.raw.UnmarshalJSON(b)
		b, err = r.raw.MarshalJSON()
		return b, err
	case JNodeTypeHash:
		b, err = json.Marshal(&r.hash)
		return b, err
	case JNodeTypeByte:
		return r.binary, nil
	}
	return r.binary, nil
}

func (r *JNode) String() string {
	b, _ := r.MarshalJSON()
	return string(b)
}

type Attachment struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Hash string `json:"hash,omitempty"`
}

/*
type GenericMessage struct {
	From        string                 `json:"from,omitempty"`
	Kargs    *JNode `json:"kargs"`
	//Kargs       map[string]interface{} `json:"kargs"`
	Attachments []*Attachment          `json:"attachments,omitempty"`
}
*/

/*
	type CallRequest struct {
		Function    string                 `json:"function"`
		Kargs       map[string]interface{} `json:"kargs"`
		Attachments []*Attachment          `json:"attachments,omitempty"`
		//	Attachments  []map[string]interface{} `json:"attachments,omitempty"`
	}

	type CallResponse struct {
		Verdict     string                 `json:"verdict"`
		Response    map[string]interface{} `json:"response,omitempty"`
		Attachments []*Attachment           `json:"attachments,omitempty"`
		Error string `json:"error,omitempty"`
		Exception string `json:"exception,omitempty"`
	}
*/
type CallRequest struct {
	Function    string        `json:"function"`
	Kargs       *JNode        `json:"kargs"`
	Attachments []*Attachment `json:"attachments,omitempty"`
	//	Attachments  []map[string]interface{} `json:"attachments,omitempty"`
}

type CallResponse struct {
	Verdict string `json:"verdict"`
	//Response    map[string]interface{} `json:"response,omitempty"`
	Response    *JNode        `json:"response"`
	Attachments []*Attachment `json:"attachments,omitempty"`
	Error       string        `json:"error,omitempty"`
	Exception   string        `json:"exception,omitempty"`
}

func (r *CallResponse) IsOK() bool {
	return r.Verdict == "OK"
}
func (r *CallResponse) GetError() string {
	if r.Verdict == "OK" {
		return ""
	} else if r.Verdict == "ERROR" {
		return r.Error
	} else if r.Verdict == "EXCEPTION" {
		return r.Exception
	}
	return "unknown verdict!"
}
func NewErrorResponse(err string) *CallResponse {
	return &CallResponse{
		Verdict: "ERROR",
		Error:   err,
	}
}
func NewExceptionResponse(exception string) *CallResponse {
	return &CallResponse{
		Verdict:   "EXCEPTION",
		Exception: exception,
	}
}
func NewOKResponse(res *JNode) *CallResponse {
	return &CallResponse{
		Verdict:  "OK",
		Response: res,
	}
}
func NewOKMapResponse(m map[string]interface{}) *CallResponse {
	return &CallResponse{
		Verdict:  "OK",
		Response: NewJNodeMap(m),
	}
}
func NewOKByteResponse(b []byte) *CallResponse {
	return &CallResponse{
		Verdict:  "OK",
		Response: NewJNodeByte(b),
	}
}
func NewOKStringResponse(s string) *CallResponse {
	return &CallResponse{
		Verdict:  "OK",
		Response: NewJNodeString(s),
	}
}
func NewOKInterfaceResponse(x interface{}) (*CallResponse, error) {
	m, err := NewJNodeInterface(x)
	if err != nil {
		return nil, err
	}
	return &CallResponse{
		Verdict:  "OK",
		Response: m,
	}, nil
}

func NewOKEmptyResponse() *CallResponse {
	return &CallResponse{
		Verdict:  "OK",
		Response: NewJNodeString("{}"),
	}
}
func NewEmptyMap() *JNode {
	return NewJNodeString("{}")
}

type GenericServiceWithContext func(req *CallRequest, context map[string]interface{}) (res *CallResponse, err error)
type ServiceCall func(req *CallRequest, context interface{}) (*CallResponse, error)
type GenericService func(req *CallRequest) (res *CallResponse, err error)
type ServiceAction func(kargs *JNode, attachments []*Attachment, context interface{}) (*CallResponse, error)
type InternalServiceAction func(kargs *JNode, context interface{}) (*CallResponse, error)


type ApiHandle interface {
	CheckAPI(functionName string) bool
	Call(functionName string, request *CallRequest, context interface{}) (*CallResponse, error)
}

type OpenAPIHandle interface {
	CheckAPI(functionName string) bool
	Call(functionName string, request json.RawMessage, context interface{}) (*CallResponse, error)
}

type ApiError struct {
	Code int
	Info string
}

func (p ApiError) Error() string {
	return fmt.Sprintf("code %d:%s", p.Code, p.Info)
}
