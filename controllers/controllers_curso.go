package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lupatini/models"
	"net/http"

	"github.com/gorilla/mux"
)

var ListCursos = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	rows, err := connectingDB.Query("SELECT * FROM public.curso")
	if err != nil {
		fmt.Println("Não foi pesquisar as categorias dos cursos")
		w.WriteHeader(400)
		panic(err)
	}

	curso := models.Curso{}
	cursos := []models.Curso{}

	for rows.Next() {
		fmt.Println(" next cursos")
		err = rows.Scan(&curso.Id, &curso.Nome, &curso.Periodo, &curso.Vagas, &curso.Descricao, &curso.Valor, &curso.Situacao, &curso.Local)
		if err != nil {
			fmt.Println("Erro ao listar cursos")
			w.WriteHeader(400)
			panic(err)
		}

		cursos = append(cursos, curso)

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(cursos)
	w.WriteHeader(http.StatusOK)

})

var ListCurso = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	curso := models.Curso{}

	rows, err := connectingDB.Query("SELECT * FROM public.curso WHERE id=$1", id)
	if err != nil {
		fmt.Println("Não foi pesquisar o curso")
		w.WriteHeader(400)
		panic(err)
	}
	for rows.Next() {

		err = rows.Scan(&curso.Id, &curso.Nome, &curso.Periodo, &curso.Vagas, &curso.Descricao, &curso.Valor, &curso.Situacao, &curso.Local)
		if err != nil {
			fmt.Println("Erro ao listar curso")
			w.WriteHeader(400)
			panic(err)
		}

	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(curso)

})
var DeleteCurso = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	sqlQuery := "DELETE FROM public.curso WHERE id = $1"
	rows, err := connectingDB.Exec(sqlQuery, id)

	if err != nil {
		fmt.Println("Erro ao excluir curso")
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteHeader(400)
		panic(err)

	}
	rowsDeleted, err := rows.RowsAffected()
	if err != nil {
		fmt.Println("Erro ao percorrer categoria a ser excluida")
		w.WriteHeader(400)
		panic(err)
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(rowsDeleted)

})
var InsertCurso = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var curso models.Curso

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &curso)

	if err != nil {
		panic(err)
	}

	/*
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

		imagemGravaBanco := "/home/zaptec/img/"+handler.Filename
	*/
	sqlQuery := "INSERT INTO public.curso(nome,periodo,vagas,descricao,valor,situacao,local) VALUES($1,$2,$3,$4,$5,$6,$7)"
	row, err := connectingDB.Exec(sqlQuery, curso.Nome, curso.Periodo,
		curso.Vagas, curso.Descricao, curso.Valor,
		curso.Situacao, curso.Local)
	_ = row
	if err != nil {
		fmt.Println("Erro ao inserir curso")
		fmt.Println(err)
		w.WriteHeader(400)
		return

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(curso)

})

var AlterCurso = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var curso models.Curso

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &curso)

	if err != nil {
		panic(err)
	}

	row, err := connectingDB.Prepare("UPDATE public.curso SET nome=$1,periodo=$2,vagas=$3,descricao=$4,valor=$5,situacao=$6,local=$7 WHERE id=$8")
	if err != nil {
		fmt.Println("Erro ao atualizar curso")
		w.WriteHeader(400)
		panic(err)
	}
	row.Exec(curso.Nome, curso.Periodo, curso.Vagas, curso.Descricao, curso.Valor,
		curso.Situacao, curso.Local, curso.Id)
	if err != nil {
		fmt.Println("Erro ao atualizar curso")
		w.WriteHeader(400)
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(curso)

})