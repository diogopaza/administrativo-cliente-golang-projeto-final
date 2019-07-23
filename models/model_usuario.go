package models

type Usuario struct {
	Id      string `json:"id"`
	Nome    string `json:"nome"`
	Email   string `json:"email"`
	Senha   string `json:"-"`
	Cpf     string `json:"cpf"`
	Celular string `json:"celular"`
	Sexo    string `json:"sexo"`
	Perfil  string `json:"perfil"`
	Token   string `json: token `
}
