package json_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	json "git.parallelcoin.io/pod/json"
)

// TestPodExtCmds tests all of the pod extended commands marshal and unmarshal into valid results include handling of optional fields being omitted in the marshalled command, while optional fields with defaults have the default assigned on unmarshalled commands.
func TestPodExtCmds(t *testing.T) {
	t.Parallel()
	testID := int(1)
	tests := []struct {
		name         string
		newCmd       func() (interface{}, error)
		staticCmd    func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "debuglevel",
			newCmd: func() (interface{}, error) {
				return json.NewCmd("debuglevel", "trace")
			},
			staticCmd: func() interface{} {
				return json.NewDebugLevelCmd("trace")
			},
			marshalled: `{"jsonrpc":"1.0","method":"debuglevel","params":["trace"],"id":1}`,
			unmarshalled: &json.DebugLevelCmd{
				LevelSpec: "trace",
			},
		},
		{
			name: "node",
			newCmd: func() (interface{}, error) {
				return json.NewCmd("node", json.NRemove, "1.1.1.1")
			},
			staticCmd: func() interface{} {
				return json.NewNodeCmd("remove", "1.1.1.1", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"node","params":["remove","1.1.1.1"],"id":1}`,
			unmarshalled: &json.NodeCmd{
				SubCmd: json.NRemove,
				Target: "1.1.1.1",
			},
		},
		{
			name: "node",
			newCmd: func() (interface{}, error) {
				return json.NewCmd("node", json.NDisconnect, "1.1.1.1")
			},
			staticCmd: func() interface{} {
				return json.NewNodeCmd("disconnect", "1.1.1.1", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"node","params":["disconnect","1.1.1.1"],"id":1}`,
			unmarshalled: &json.NodeCmd{
				SubCmd: json.NDisconnect,
				Target: "1.1.1.1",
			},
		},
		{
			name: "node",
			newCmd: func() (interface{}, error) {
				return json.NewCmd("node", json.NConnect, "1.1.1.1", "perm")
			},
			staticCmd: func() interface{} {
				return json.NewNodeCmd("connect", "1.1.1.1", json.String("perm"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"node","params":["connect","1.1.1.1","perm"],"id":1}`,
			unmarshalled: &json.NodeCmd{
				SubCmd:        json.NConnect,
				Target:        "1.1.1.1",
				ConnectSubCmd: json.String("perm"),
			},
		},
		{
			name: "node",
			newCmd: func() (interface{}, error) {
				return json.NewCmd("node", json.NConnect, "1.1.1.1", "temp")
			},
			staticCmd: func() interface{} {
				return json.NewNodeCmd("connect", "1.1.1.1", json.String("temp"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"node","params":["connect","1.1.1.1","temp"],"id":1}`,
			unmarshalled: &json.NodeCmd{
				SubCmd:        json.NConnect,
				Target:        "1.1.1.1",
				ConnectSubCmd: json.String("temp"),
			},
		},
		{
			name: "generate",
			newCmd: func() (interface{}, error) {
				return json.NewCmd("generate", 1)
			},
			staticCmd: func() interface{} {
				return json.NewGenerateCmd(1)
			},
			marshalled: `{"jsonrpc":"1.0","method":"generate","params":[1],"id":1}`,
			unmarshalled: &json.GenerateCmd{
				NumBlocks: 1,
			},
		},
		{
			name: "getbestblock",
			newCmd: func() (interface{}, error) {
				return json.NewCmd("getbestblock")
			},
			staticCmd: func() interface{} {
				return json.NewGetBestBlockCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getbestblock","params":[],"id":1}`,
			unmarshalled: &json.GetBestBlockCmd{},
		},
		{
			name: "getcurrentnet",
			newCmd: func() (interface{}, error) {
				return json.NewCmd("getcurrentnet")
			},
			staticCmd: func() interface{} {
				return json.NewGetCurrentNetCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getcurrentnet","params":[],"id":1}`,
			unmarshalled: &json.GetCurrentNetCmd{},
		},
		{
			name: "getheaders",
			newCmd: func() (interface{}, error) {
				return json.NewCmd("getheaders", []string{}, "")
			},
			staticCmd: func() interface{} {
				return json.NewGetHeadersCmd(
					[]string{},
					"",
				)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getheaders","params":[[],""],"id":1}`,
			unmarshalled: &json.GetHeadersCmd{
				BlockLocators: []string{},
				HashStop:      "",
			},
		},
		{
			name: "getheaders - with arguments",
			newCmd: func() (interface{}, error) {
				return json.NewCmd("getheaders", []string{"000000000000000001f1739002418e2f9a84c47a4fd2a0eb7a787a6b7dc12f16", "0000000000000000026f4b7f56eef057b32167eb5ad9ff62006f1807b7336d10"}, "000000000000000000ba33b33e1fad70b69e234fc24414dd47113bff38f523f7")
			},
			staticCmd: func() interface{} {
				return json.NewGetHeadersCmd(
					[]string{
						"000000000000000001f1739002418e2f9a84c47a4fd2a0eb7a787a6b7dc12f16",
						"0000000000000000026f4b7f56eef057b32167eb5ad9ff62006f1807b7336d10",
					},
					"000000000000000000ba33b33e1fad70b69e234fc24414dd47113bff38f523f7",
				)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getheaders","params":[["000000000000000001f1739002418e2f9a84c47a4fd2a0eb7a787a6b7dc12f16","0000000000000000026f4b7f56eef057b32167eb5ad9ff62006f1807b7336d10"],"000000000000000000ba33b33e1fad70b69e234fc24414dd47113bff38f523f7"],"id":1}`,
			unmarshalled: &json.GetHeadersCmd{
				BlockLocators: []string{
					"000000000000000001f1739002418e2f9a84c47a4fd2a0eb7a787a6b7dc12f16",
					"0000000000000000026f4b7f56eef057b32167eb5ad9ff62006f1807b7336d10",
				},
				HashStop: "000000000000000000ba33b33e1fad70b69e234fc24414dd47113bff38f523f7",
			},
		},
		{
			name: "version",
			newCmd: func() (interface{}, error) {
				return json.NewCmd("version")
			},
			staticCmd: func() interface{} {
				return json.NewVersionCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"version","params":[],"id":1}`,
			unmarshalled: &json.VersionCmd{},
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the command as created by the new static command creation function.
		marshalled, err := json.MarshalCmd(testID, test.staticCmd())
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}
		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}
		// Ensure the command is created without error via the generic new command creation function.
		cmd, err := test.newCmd()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}
		// Marshal the command as created by the generic new command creation function.
		marshalled, err = json.MarshalCmd(testID, cmd)
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}
		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}
		var request json.Request
		if err := json.Unmarshal(marshalled, &request); err != nil {
			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, err)
			continue
		}
		cmd, err = json.UnmarshalCmd(&request)
		if err != nil {
			t.Errorf("UnmarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}
		if !reflect.DeepEqual(cmd, test.unmarshalled) {
			t.Errorf("Test #%d (%s) unexpected unmarshalled command "+
				"- got %s, want %s", i, test.name,
				fmt.Sprintf("(%T) %+[1]v", cmd),
				fmt.Sprintf("(%T) %+[1]v\n", test.unmarshalled))
			continue
		}
	}
}
