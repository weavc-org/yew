package configs

import (
	"testing"
)

type data struct {
	Payload string
}

func TestLoadConfig(t *testing.T) {
	var v *data = &data{}
	err := FetchConfig("./config.yaml", "examples.api", v)
	if err != nil {
		t.Errorf("%s", err)
	}

	if v.Payload != "foobar" {
		t.Errorf("incorrect payload binding. Expected %s, Got %s", "foobar", v.Payload)
	}
}
