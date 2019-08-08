package models

import "time"

type Curso struct {
	Id             string `json:"id"`
	Nome           string `json:"nome"`
	Vagas          string `json:"vagas"`
	Descricao      string `json:"descricao"`
	Valor          string `json:"valor"`
	Situacao       string `json:"situacao"`
	LocalId        string `json:"localId"`
	Imagem         string `json:"imagem"`
	ProfessorId    string `json:"professorId"`
	CategoriaId    string `json:"categoriaId"`
	FormaPagamento string `json:"formaPagamento"`
	Horarios       []Data `json:"horarios"`
	Professor      Usuario
}

type Data struct {
	Id       string    `json:"id"`
	Data     time.Time `json:"data"`
	HoraIni  string    `json:"horaIni"`
	HoraFim  string    `json:"horaFim"`
	Curso_id string    `json:"curso_id"`
}
