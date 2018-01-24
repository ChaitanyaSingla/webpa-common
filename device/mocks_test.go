package device

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockReader struct {
	mock.Mock
}

func (m *mockReader) Read(b []byte) (int, error) {
	arguments := m.Called(b)
	return arguments.Int(0), arguments.Error(1)
}

type mockDevice struct {
	mock.Mock
}

func (m *mockDevice) String() string {
	return m.Called().String(0)
}

func (m *mockDevice) MarshalJSON() ([]byte, error) {
	arguments := m.Called()
	return arguments.Get(0).([]byte), arguments.Error(1)
}

func (m *mockDevice) ID() ID {
	return m.Called().Get(0).(ID)
}

func (m *mockDevice) Pending() int {
	return m.Called().Int(0)
}

func (m *mockDevice) RequestClose() {
	m.Called()
}

func (m *mockDevice) Closed() bool {
	arguments := m.Called()
	return arguments.Bool(0)
}

func (m *mockDevice) Statistics() Statistics {
	arguments := m.Called()
	first, _ := arguments.Get(0).(Statistics)
	return first
}

func (m *mockDevice) Send(request *Request) (*Response, error) {
	arguments := m.Called(request)
	first, _ := arguments.Get(0).(*Response)
	return first, arguments.Error(1)
}

type mockDialer struct {
	mock.Mock
}

func (m *mockDialer) DialDevice(deviceName, url string, extra http.Header) (*websocket.Conn, *http.Response, error) {
	var (
		arguments = m.Called(deviceName, url, extra)
		first, _  = arguments.Get(0).(*websocket.Conn)
		second, _ = arguments.Get(1).(*http.Response)
	)

	return first, second, arguments.Error(2)
}

type mockWebsocketDialer struct {
	mock.Mock
}

func (m *mockWebsocketDialer) Dial(url string, requestHeader http.Header) (*websocket.Conn, *http.Response, error) {
	var (
		arguments = m.Called(url, requestHeader)
		first, _  = arguments.Get(0).(*websocket.Conn)
		second, _ = arguments.Get(1).(*http.Response)
	)

	return first, second, arguments.Error(2)
}

// deviceSet is a convenient map type for capturing visited devices
// and asserting expectations.
type deviceSet map[*device]bool

func (s deviceSet) len() int {
	return len(s)
}

func (s deviceSet) add(d Interface) {
	s[d.(*device)] = true
}

func (s *deviceSet) reset() {
	*s = make(deviceSet)
}

// managerCapture returns a high-level visitor for Manager testing
func (s deviceSet) managerCapture() func(Interface) {
	return func(d Interface) {
		s.add(d)
	}
}

// registryCapture returns a low-level visitor for registry testing
func (s deviceSet) registryCapture() func(*device) {
	return func(d *device) {
		s[d] = true
	}
}

func (s deviceSet) assertSameID(assert *assert.Assertions, expected ID) {
	for d := range s {
		assert.Equal(expected, d.ID())
	}
}

func (s deviceSet) assertDistributionOfIDs(assert *assert.Assertions, expected map[ID]int) {
	actual := make(map[ID]int, len(expected))
	for d := range s {
		actual[d.ID()] += 1
	}

	assert.Equal(expected, actual)
}

// drain copies whatever is available on the given channel into this device set
func (s deviceSet) drain(source <-chan Interface) {
	for d := range source {
		s.add(d)
	}
}

func expectsDevices(devices ...*device) deviceSet {
	result := make(deviceSet, len(devices))
	for _, d := range devices {
		result[d] = true
	}

	return result
}

type mockRouter struct {
	mock.Mock
}

func (m *mockRouter) Route(request *Request) (*Response, error) {
	arguments := m.Called(request)
	first, _ := arguments.Get(0).(*Response)
	return first, arguments.Error(1)
}

type mockConnector struct {
	mock.Mock
}

func (m *mockConnector) Connect(response http.ResponseWriter, request *http.Request, header http.Header) (Interface, error) {
	arguments := m.Called(response, request, header)
	first, _ := arguments.Get(0).(Interface)
	return first, arguments.Error(1)
}

func (m *mockConnector) Disconnect(id ID) bool {
	return m.Called().Bool(0)
}

func (m *mockConnector) DisconnectIf(predicate func(ID) bool) int {
	return m.Called(predicate).Int(0)
}

type mockRegistry struct {
	mock.Mock
}

func (m *mockRegistry) Get(id ID) (Interface, bool) {
	arguments := m.Called(id)
	first, _ := arguments.Get(0).(Interface)
	return first, arguments.Bool(1)
}

func (m *mockRegistry) VisitIf(predicate func(ID) bool, visitor func(Interface)) int {
	return m.Called(predicate, visitor).Int(0)
}

func (m *mockRegistry) VisitAll(visitor func(Interface)) int {
	return m.Called(visitor).Int(0)
}
