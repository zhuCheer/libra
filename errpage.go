package libra

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	ErrDefaultPage = `
<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
<html>
<head><title>{#title#}</title></head>
<body bgcolor="white">
<h1>{#title#}</h1>
<p>{#msg#}<br/>
Thank you very much!</p>
<table>
<tr>
<td>URL:</td>
<td>{#url#}</td>
</tr>
<tr>
<td>Server:</td>
<td>{#host#}</td>
</tr>
<tr>
<td>Date:</td>
<td>{#time#}</td>
</tr>
</table>
<hr/>Powered by libra/0.0.1</body>
</html>
`
)

// proxy not found page
func getDefaultErrorPage(statusCode int, msg string, req *http.Request) (resp *http.Response, err error) {
	errPageTemplate := ErrDefaultPage

	errPageTemplate = strings.Replace(errPageTemplate, "{#title#}",
		fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)), -1)
	errPageTemplate = strings.Replace(errPageTemplate, "{#title#}",
		fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)), -1)
	errPageTemplate = strings.Replace(errPageTemplate, "{#msg#}", msg, -1)
	errPageTemplate = strings.Replace(errPageTemplate, "{#url#}", req.URL.Path, -1)
	errPageTemplate = strings.Replace(errPageTemplate, "{#host#}", req.Host, -1)
	errPageTemplate = strings.Replace(errPageTemplate, "{#time#}", time.Now().String(), -1)

	resp = getResponsePage(statusCode, errPageTemplate, req)
	return resp, err
}

func getResponsePage(status int, msg string, req *http.Request) *http.Response {
	var resp *http.Response

	b, _ := ioutil.ReadAll(bytes.NewReader([]byte(msg)))
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp = &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(http.StatusOK)),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
		Body:       body,
		Request:    req,
	}
	return resp
}
