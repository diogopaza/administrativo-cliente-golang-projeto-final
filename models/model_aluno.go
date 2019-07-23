package models

type Aluno struct {
	Id        string `json:"id"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	Senha     string `json:"senha"`
	Profissao string `json:"profissao"`
	Celular   string `json:"celular"`
	Telefone  string `json:"telefone"`
	Sexo      string `json:"sexo"`
	Cpf       string `json:"cpf"`
	Imagem    string `json:"imagem"`
}
