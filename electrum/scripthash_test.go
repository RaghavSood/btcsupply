package electrum

import "testing"

func TestScriptToElectrumScript(t *testing.T) {
	// https://electrumx.readthedocs.io/en/latest/protocol-basics.html#script-hashes
	// Genesis block P2PK script
	script := "4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac"
	expected := "740485f380ff6379d11ef6fe7d7cdd68aea7f8bd0d953d9fdf3531fb7d531833"

	result, _ := ScriptToElectrumScript(script)
	if result != expected {
		t.Errorf("ScriptToElectrumScript(%s) returned %s, expected %s", script, result, expected)
	}
}
