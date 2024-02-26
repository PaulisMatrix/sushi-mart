package tests

import (
	"encoding/json"
)

type GinResponse struct {
	Status string `json:"status"`
}

func (s *ServerSuite) TestPing() {
	path := "/testing/ping"
	resp, err := s.makeRequest(path)

	s.Nil(err)
	s.NotNil(resp)

	var serverResp GinResponse
	errUnMarshal := json.Unmarshal(resp, &serverResp)
	s.Nil(errUnMarshal)

	s.EqualValues("pong!!", serverResp.Status)
}
