// Copyright (c) 2021 Proton Technologies AG
//
// This file is part of ProtonMail Bridge.
//
// ProtonMail Bridge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// ProtonMail Bridge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with ProtonMail Bridge.  If not, see <https://www.gnu.org/licenses/>.

package pmapi

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"testing"

	"github.com/hashicorp/go-multierror"
)

var (
	colRed     = "\033[1;31m"
	colNon     = "\033[0;39m"
	reHTTPCode = regexp.MustCompile(`(HTTP|get|post|put|delete)_(\d{3}).*.json`)
)

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		vv := []interface{}{filepath.Base(file), line, colRed}
		vv = append(vv, v...)
		vv = append(vv, colNon)
		fmt.Printf("%s:%d: %s"+msg+"%s\n\n", vv...)
		tb.FailNow()
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: %sunexpected error: %s%s\n\n", filepath.Base(file), line, colRed, err.Error(), colNon)
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n%s\texp: %#v\n\n\tgot: %#v%s\n\n", filepath.Base(file), line, colRed, exp, act, colNon)
		tb.FailNow()
	}
}

// newTestServer is old function and should be replaced everywhere by newTestServerCallbacks.
func newTestServer(h http.Handler) (*httptest.Server, *client) {
	s := httptest.NewServer(h)

	serverURL, err := url.Parse(s.URL)
	if err != nil {
		panic(err)
	}

	cm := newTestClientManager(testClientConfig)
	cm.host = serverURL.Host
	cm.scheme = serverURL.Scheme

	return s, newTestClient(cm)
}

func newTestServerCallbacks(tb testing.TB, callbacks ...func(testing.TB, http.ResponseWriter, *http.Request) string) (func(), *client) {
	reqNum := 0
	_, file, line, _ := runtime.Caller(1)
	file = filepath.Base(file)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqNum++
		if reqNum > len(callbacks) {
			fmt.Printf(
				"%s:%d: %sServer was requeted %d times which is more requests than expected %d%s\n\n",
				file, line, colRed, reqNum, len(callbacks), colNon,
			)
			tb.FailNow()
		}
		response := callbacks[reqNum-1](tb, w, r)
		if response != "" {
			writeJSONResponsefromFile(tb, w, response, reqNum-1)
		}
	}))

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		panic(err)
	}

	finish := func() {
		server.CloseClientConnections() // Closing without waiting for finishing requests.
		if reqNum != len(callbacks) {
			fmt.Printf(
				"%s:%d: %sServer was requested %d times but expected to be %d times%s\n\n",
				file, line, colRed, reqNum, len(callbacks), colNon,
			)
			tb.Error("server failed")
		}
	}

	cm := newTestClientManager(testClientConfig)
	cm.host = serverURL.Host
	cm.scheme = serverURL.Scheme

	return finish, newTestClient(cm)
}

func checkMethodAndPath(r *http.Request, method, path string) error {
	var result *multierror.Error
	if err := checkHeader(r.Header, "x-pm-appversion", "GoPMAPI_1.0.14"); err != nil {
		result = multierror.Append(result, err)
	}
	if r.Method != method {
		err := fmt.Errorf("Invalid request method expected %v, got %v", method, r.Method)
		result = multierror.Append(result, err)
	}
	if r.URL.RequestURI() != path {
		err := fmt.Errorf("Invalid request path expected %v, got %v", path, r.URL.RequestURI())
		result = multierror.Append(result, err)
	}
	return result.ErrorOrNil()
}

func httpResponse(code int) string {
	return fmt.Sprintf("HTTP_%d.json", code)
}

func writeJSONResponsefromFile(tb testing.TB, w http.ResponseWriter, response string, reqNum int) {
	if match := reHTTPCode.FindAllSubmatch([]byte(response), -1); len(match) != 0 {
		httpCode, err := strconv.Atoi(string(match[0][len(match[0])-1]))
		Ok(tb, err)
		w.WriteHeader(httpCode)
	}
	f, err := os.Open("./testdata/routes/" + response)
	Ok(tb, err)
	w.Header().Set("content-type", "application/json;charset=utf-8")
	w.Header().Set("x-test-pmapi-response", fmt.Sprintf("%s:%d", tb.Name(), reqNum))
	_, err = io.Copy(w, f)
	Ok(tb, err)
}

func checkHeader(h http.Header, field, exp string) error {
	val := h.Get(field)
	if val != exp {
		msg := "wrong field %s expected %q but have %q"
		return fmt.Errorf(msg, field, exp, val)
	}
	return nil
}

func isAuthReq(r *http.Request, uid, token string) error { //nolint[unparam] always retrieves testUID
	if err := checkHeader(r.Header, "x-pm-uid", uid); err != nil {
		return err
	}
	if err := checkHeader(r.Header, "authorization", "Bearer "+token); err != nil { //nolint[revive] can return the error right away but this is easier to read
		return err
	}
	return nil
}
