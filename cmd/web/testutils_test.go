package main

import (
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/golangcollege/sessions"
	"karolharasim.com/snippetbox/pkg/models/mock"
)

//Define a regex which captures the CSRF token value from the HTML for user signup page.
var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

func extractCSRFToken(t *testing.T, body []byte) string {
	//Use the FindSubmatch method to extract the token from the HTML body
	//Note that this returns an array with the entire matched pattern in the first position, and the values of any captured data in subsequent positions.
	matches := csrfTokenRX.FindSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	return html.UnescapeString(string(matches[1]))
}

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
    session.Lifetime = 12 * time.Hour
	session.Secure = true
	

	return &application{
		errorLog: log.New(ioutil.Discard, "", 0),
		infoLog: log.New(ioutil.Discard, "", 0),
		session: session,
		snippets: &mock.SnippetModel{},
		templateCache: templateCache,
		users: &mock.UserModel{},
	}
}

//Define a custom testServer which anonymously embeds a httptest.Server instance
type testServer struct {
	*httptest.Server
}

//Create a newTestServer helper which intializes and returns a new instance of our custom testServer type
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	//Add the cookie jar to the client, so that response cookies are stored and them semt woth subsequent requests
	ts.Client().Jar = jar

	//Disable redirect-following for the client. Essentially this funciton is called after a 3xx response is received by the client, and returning the 
	//http.ErrUseLastResponse error forces it to immediately return the received response
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

//Implement a get method on our custom testServer type. This makes a GET request to a given url path on the test server, and returns the response
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, []byte) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}