package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"lupatini/config"
	"lupatini/models"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" //postgresql
)

var connectingDB *sql.DB

func init() {
	connectingDB = config.DB()
}

var ListAlunos = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	rows, err := connectingDB.Query("SELECT * FROM public.aluno")
	if err != nil {
		fmt.Println("Não foi pesquisar alunos")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		panic(err)
	}

	aluno := models.Aluno{}
	alunos := []models.Aluno{}

	for rows.Next() {
		fmt.Println(" next alunos")
		err = rows.Scan(&aluno.Id, &aluno.Nome, &aluno.Email, &aluno.Senha, &aluno.Profissao,
			&aluno.Celular, &aluno.Telefone, &aluno.Sexo, &aluno.Cpf, &aluno.Imagem, &aluno.Nascimento)
		if err != nil {
			fmt.Println("Erro ao listar alunos")
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(400)
			panic(err)
		}

		alunos = append(alunos, aluno)

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(alunos)
	w.WriteHeader(http.StatusOK)

})

var ListAluno = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	aluno := models.Aluno{}

	rows, err := connectingDB.Query("SELECT * FROM public.aluno WHERE id=$1", id)
	if err != nil {
		fmt.Println("Não foi pesquisar o aluno")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		panic(err)
	}
	for rows.Next() {

		err = rows.Scan(&aluno.Id, &aluno.Nome, &aluno.Email, &aluno.Senha, &aluno.Profissao,
			&aluno.Celular, &aluno.Telefone, &aluno.Sexo, &aluno.Cpf, &aluno.Imagem)
		if err != nil {
			fmt.Println("Erro ao listar aluno")
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(400)
			panic(err)
		}

	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(aluno)

})

var AlterImagemAluno = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	//salva imagem do usuario

	file, handler, err := r.FormFile("imagem")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	f, err := os.OpenFile("/home/zaptec/img/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	fmt.Println(handler)
	imagemGravaBanco := "/home/zaptec/img/" + handler.Filename
	_ = imagemGravaBanco

})

var InsertAluno = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var aluno models.Aluno
	var dataNascimentoFinal string
	var imagemGravaBanco string
	nome := r.FormValue("nome")
	fmt.Println(nome)

	json.Unmarshal([]byte(nome), &aluno)
	fmt.Println(aluno)

	dataNascimento := strings.Split(aluno.Nascimento, "/")

	for i := range dataNascimento {

		dataNascimentoFinal += dataNascimento[i]
		if i < 2 {
			dataNascimentoFinal = dataNascimentoFinal + "-"
		}
	}
	fmt.Println(dataNascimentoFinal)

	file, handler, err := r.FormFile("selectedFile")
	if err != nil {
		fmt.Println("Erro ao abrir arquivo")

	}
	defer file.Close()
	f, err := os.OpenFile("/home/zaptec/img/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Erro ao criar arquivo", err)
	}
	defer f.Close()
	io.Copy(f, file)
	fmt.Println(handler.Filename)
	imagemGravaBanco = "/home/zaptec/img/" + handler.Filename

	sqlQuery := "INSERT INTO public.aluno(nome,email,senha,profissao,celular,telefone,sexo,cpf,imagem,data_nascimento) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)"
	row, err := connectingDB.Exec(sqlQuery, aluno.Nome, aluno.Email,
		aluno.Senha, aluno.Profissao, aluno.Celular, aluno.Telefone,
		aluno.Sexo, aluno.Cpf, imagemGravaBanco, dataNascimentoFinal)
	_ = row
	if err != nil {
		fmt.Println("Erro ao inserir aluno")
		fmt.Println(err)
		w.WriteHeader(400)
		return

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(aluno)

})

var AlterAluno = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var aluno models.Aluno

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &aluno)

	if err != nil {
		panic(err)
	}

	row, err := connectingDB.Prepare("UPDATE public.aluno SET nome=$1,email=$2,senha=$3,profissao=$4,celular=$5,telefone=$6,sexo=$7,cpf=$8,imagem=$9 WHERE id=$10")
	if err != nil {
		fmt.Println("Erro ao atualizar aluno")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		panic(err)
	}
	row.Exec(aluno.Nome, aluno.Email, aluno.Senha, aluno.Profissao, aluno.Celular, aluno.Telefone,
		aluno.Sexo, aluno.Cpf, aluno.Imagem, aluno.Id)
	if err != nil {
		fmt.Println("Erro ao atualizar aluno")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(aluno)

})

var DeleteAluno = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println(id)
	sqlQuery := "DELETE FROM public.aluno WHERE id = $1"
	rows, err := connectingDB.Exec(sqlQuery, id)

	if err != nil {
		fmt.Println("Erro ao excluir Aluno")
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteHeader(400)
		panic(err)

	}
	rowsDeleted, err := rows.RowsAffected()
	if err != nil {
		fmt.Println("Erro ao percorrer aluno a ser excluida")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(rowsDeleted)

})
