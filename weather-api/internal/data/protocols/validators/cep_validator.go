package validators_protocols

type CepValidator interface {
	Validate(cep string) error
}
