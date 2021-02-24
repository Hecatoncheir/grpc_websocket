package envelope

import (
	"encoding/json"
	"testing"
)

func TestNew(t *testing.T) {
	details := map[string]string{"test key": "test value"}
	encodedDetails, err := json.Marshal(details)
	if err != nil {
		t.Error(err)
	}

	envelope := Envelope{
		Message: "test message",
		Details: string(encodedDetails),
	}

	if envelope.GetMessage() != "test message" {
		t.Fail()
	}

	if envelope.GetDetails() != "{\"test key\":\"test value\"}" {
		t.Fail()
	}

	decodedDetails, err := envelope.GetDecodedDetails()
	if err != nil {
		t.Error(err)
	}

	if decodedDetails["test key"] != "test value" {
		t.Fail()
	}
}
