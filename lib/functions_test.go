package lib

import "testing"

func TestLogs(t *testing.T) {
	_, _ = GetFunctionNodes(ParseRequest{
		Path: "/Users/das/Baylor/RA/bpm-go/testbed",
	})
}
