package models

type Curso struct {
	Id        string `json:"id"`
	Nome      string `json:"nome"`
	Periodo   string `json:"periodo"`
	Vagas     int32 `json:"vagas"`
	Descricao string `json:"descricao"`
	Valor     float32 `json:"valor"`
	Situacao  string `json:"situacao"`
	Local     string `json:"local"`
}
