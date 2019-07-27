package models

type Local struct {
	Id          int32  `json:"id"`
	Rua         string `json:"rua"`
	Numero      int32  `json:numero`
	Complemento string `json:complemento`
	Bairro      string `json:bairro`
	Cidade      string `json:cidade`
	Estado      string `json:estado`
}
