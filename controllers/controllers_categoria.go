package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lupatini/models"
	"net/http"

	"github.com/gorilla/mux"
)

var ListCategorias = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	rows, err := connectingDB.Query("SELECT * FROM public.categoria_curso")
	if err != nil {
		fmt.Println("Não foi pesquisar as categorias dos cursos")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		panic(err)
	}

	categoria := models.Categoria{}
	categorias := []models.Categoria{}

	for rows.Next() {
		fmt.Println(" next categorias")
		err = rows.Scan(&categoria.Id, &categoria.Nome)
		if err != nil {
			fmt.Println("Erro ao listar categorias")
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(400)
			panic(err)
		}

		categorias = append(categorias, categoria)

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(categorias)
	w.WriteHeader(http.StatusOK)

})
var ListCategoria = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	categoria := models.Categoria{}

	rows, err := connectingDB.Query("SELECT * FROM public.categoria_curso WHERE id=$1", id)
	if err != nil {
		fmt.Println("Não foi pesquisar a categoria")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		panic(err)
	}
	for rows.Next() {

		err = rows.Scan(&categoria.Id, &categoria.Nome)
		if err != nil {
			fmt.Println("Erro ao listar categoria")
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(400)
			panic(err)
		}

	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(categoria)

})
var DeleteCategoria = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	sqlQuery := "DELETE FROM public.categoria_curso WHERE id = $1"
	rows, err := connectingDB.Exec(sqlQuery, id)

	if err != nil {
		fmt.Println("Erro ao excluir categoria")
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteHeader(400)
		panic(err)

	}
	rowsDeleted, err := rows.RowsAffected()
	if err != nil {
		fmt.Println("Erro ao percorrer categoria a ser excluida")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(rowsDeleted)

})
var InsertCategoria = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var categoria models.Categoria

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &categoria)
	if err != nil {
		panic(err)
	}

	sqlQuery := "INSERT INTO public.categoria_curso(nome) VALUES($1)"
	row, err := connectingDB.Exec(sqlQuery, categoria.Nome)
	_ = row
	if err != nil {
		fmt.Println("Erro ao inserir categoria")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		return

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(categoria)

})

var AlterCategoria = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var categoria models.Categoria

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &categoria)

	if err != nil {
		panic(err)
	}

	row, err := connectingDB.Prepare("UPDATE public.categoria_curso SET nome=$1 WHERE id=$2")
	if err != nil {
		fmt.Println("Erro ao atualizar aluno")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		panic(err)
	}
	row.Exec(categoria.Nome, categoria.Id)
	if err != nil {
		fmt.Println("Erro ao atualizar categoria")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categoria)
})
