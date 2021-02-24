package envelope

type EnvelopeInterface interface {
	GetMessage() string
	GetDetails() string
	GetDecodedDetails() (map[string]interface{}, error)
}
