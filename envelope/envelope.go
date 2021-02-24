package envelope

import (
	"encoding/json"
)

type Envelope struct {
	EnvelopeInterface

	Message string

	/// Some encoded details
	Details string
}

func (envelope *Envelope) GetMessage() string {
	return envelope.Message
}

func (envelope *Envelope) GetDetails() string {
	return envelope.Details
}

func (envelope *Envelope) GetDecodedDetails() (map[string]interface{}, error) {
	var details map[string]interface{}

	err := json.Unmarshal([]byte(envelope.Details), &details)
	if err != nil {
		return nil, err
	}

	return details, nil
}
