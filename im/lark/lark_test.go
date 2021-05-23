package lark

import "testing"

func TestLark(t *testing.T) {
	larkCli := New()

	// TODO: fix test id
	logger := larkCli.Logger(map[string]interface{}{
		"chat_id": "oc_843dd49366635fc3fd33c29b4b4ece00",
		"open_id": "ou_2988d05af705b066d80c14293649659c",
	})
	logger.Infof("unit-test-info")
}
