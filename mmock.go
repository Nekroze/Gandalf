package gandalf

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/jmartin82/mmock/pkg/mock"
)

// ToMMock exports Contract as MMock definitions to build a fake api endpoint
// with optional state via MMock scenarios. MMock
// (https://github.com/jmartin82/mmock) is an http mocking server.
type ToMMock struct {
	// The state(s) that the Scenario must be in to trigger this mock.
	TriggerStates []string
	// The Scenario to which state is stored.
	Scenario string
	// The state to transition the scenario to when this mock is triggered.
	NewState string
	// When set this is used for the request path definition instead of the path from the Contract's Requestor.
	Path string
	// Enables chaos testing by causing the mock, when triggered, may return a 5xx instead.
	ChaoticEvil bool
	// If true MMock will require the request headers to match exactly to trigger this mock.
	// This should be left false (the default ) for dynamic headers such as tokens/id's.
	MatchHeaders bool
	// If true MMock will require the request body to match exactly to trigger this mock.
	// This should be left false (the default ) for dynamic requests such as tokens/id's.
	MatchBody bool
	saved     bool
}

func headersToValues(h map[string][]string) mock.Values {
	return h
}

func cookiesToMap(cs []*http.Cookie) map[string]string {
	out := map[string]string{}
	for _, c := range cs {
		out[c.Name] = c.Value
	}
	return out
}

func (m *ToMMock) translateRequest(req *http.Request) mock.Request {
	out := mock.Request{
		Path:   req.URL.Path,
		Method: req.Method,
	}
	if m.Path != "" {
		out.Path = m.Path
	}
	if m.MatchHeaders {
		out.HttpHeaders = mock.HttpHeaders{
			Headers: headersToValues(req.Header),
			Cookies: cookiesToMap(req.Cookies()),
		}
	}
	if m.MatchBody {
		out.Body = GetRequestBody(req)
	}
	return out
}

func (m *ToMMock) translateResponse(resp *http.Response) mock.Response {
	return mock.Response{
		StatusCode: resp.StatusCode,
		HttpHeaders: mock.HttpHeaders{
			Headers: headersToValues(resp.Header),
			Cookies: cookiesToMap(resp.Cookies()),
		},
		Body: GetResponseBody(resp),
	}
}

func (m *ToMMock) translateMock() mock.Control {
	out := mock.Control{
		Crazy: OverrideChaos || m.ChaoticEvil,
	}
	if m.Scenario != "" {
		out.Scenario = mock.Scenario{
			Name: m.Scenario,
		}
		if len(m.TriggerStates) > 0 {
			out.Scenario.RequiredState = m.TriggerStates
		}
		if m.NewState != "" {
			out.Scenario.NewState = m.NewState
		}
	}
	return out
}

// Uses Requester.GetRequest and Checker.GetResponse as a basis to build
// an MMock definition.
func (m *ToMMock) contractToMock(c *Contract) mock.Definition {
	return mock.Definition{
		URI:         c.Name + ".json",
		Description: c.Name,
		Request:     m.translateRequest(c.Request.GetRequest()),
		Response:    m.translateResponse(c.Check.GetResponse()),
		Control:     m.translateMock(),
	}
}

func (m *ToMMock) saveMockToFile(mock mock.Definition) error {
	out, err := json.Marshal(mock)
	if err != nil {
		return err
	}
	out = mmockFixTimespansForParsing(out)
	err = ioutil.WriteFile(path.Join(MockSavePath, mock.Description+".json"), out, 0644)
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(MockDelay) * time.Millisecond)
	m.saved = true
	return err
}

// Save a valid MMock definition to a json file with the contract name as the filename.
// This incurs disk IO so is restricted to only saving once per instance.
func (m *ToMMock) Save(c *Contract) error {
	if m.saved || c.Tested || MockSkip {
		return nil
	}
	return m.saveMockToFile(m.contractToMock(c))
}

func mmockFixTimespansForParsing(in []byte) []byte {
	return []byte(strings.ReplaceAll(string(in), `{"Duration":0}`, `"0s"`))
}
