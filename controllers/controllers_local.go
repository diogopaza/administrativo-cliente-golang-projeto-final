package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lupatini/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var ListLocais = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	rows, err := connectingDB.Query("SELECT * FROM public.local_curso")
	if err != nil {
		fmt.Println("Não foi pesquisar os locais")
		w.WriteHeader(400)
		panic(err)
	}

	local := models.Local{}
	locais := []models.Local{}

	for rows.Next() {
		fmt.Println(" next locais")
		err = rows.Scan(&local.Id, &local.Rua, &local.Numero, &local.Complemento, &local.Bairro, &local.Cidade,
			&local.Estado)
		if err != nil {
			fmt.Println("Erro ao listar locais")
			w.WriteHeader(400)
			panic(err)
		}

		locais = append(locais, local)

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(locais)
	w.WriteHeader(http.StatusOK)

})

var ListLocal = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	local := models.Local{}

	rows, err := connectingDB.Query("SELECT * FROM public.local_curso WHERE id=$1", id)
	if err != nil {
		fmt.Println("Não foi pesquisar o curso")
		w.WriteHeader(400)
		panic(err)
	}
	for rows.Next() {

		err = rows.Scan(&local.Id, &local.Rua, &local.Numero, &local.Complemento, &local.Bairro, &local.Cidade,
			&local.Estado)
		if err != nil {
			fmt.Println("Erro ao listar curso")
			w.WriteHeader(400)
			panic(err)
		}

	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(local)

})
var DeleteLocal = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println(id)
	sqlQuery := "DELETE FROM public.local_curso WHERE id = $1"
	rows, err := connectingDB.Exec(sqlQuery, id)

	if err != nil {
		fmt.Println("Erro ao excluir local")
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteHeader(400)
		panic(err)

	}
	rowsDeleted, err := rows.RowsAffected()
	if err != nil {
		fmt.Println("Erro ao percorrer local a ser excluido")
		w.WriteHeader(400)
		panic(err)
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(rowsDeleted)

})

var InsertLocal = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var local models.Local
	var numero int64

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &local)

	if err != nil {
		panic(err)
	}
	//converte Numero de string para int
	if local.Numero != "" {
		numero, err = strconv.ParseInt(local.Numero, 10, 32)
		if err != nil {
			fmt.Println(err)
		}
	}
	sqlQuery := "INSERT INTO public.local_curso(rua,numero,complemento,bairro,cidade,estado) VALUES($1,$2,$3,$4,$5,$6)"
	row, err := connectingDB.Exec(sqlQuery, &local.Rua, numero, &local.Complemento, &local.Bairro, &local.Cidade,
		&local.Estado)
	_ = row
	if err != nil {
		fmt.Println("Erro ao inserir curso")
		fmt.Println(err)
		w.WriteHeader(400)
		return

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(local)

})

var AlterLocal = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var local models.Local
	var numero int64

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &local)

	if err != nil {
		panic(err)
	}

	//converte Numero de string para int
	if local.Numero != "" {
		numero, err = strconv.ParseInt(local.Numero, 10, 32)
		if err != nil {
			fmt.Println(err)
		}
	}

	row, err := connectingDB.Prepare("UPDATE local_curso SET rua=$1,numero=$2,complemento=$3,bairro=$4,cidade=$5,estado=$6 WHERE id=$7")
	if err != nil {
		fmt.Println("Erro ao atualizar local")
		w.WriteHeader(400)
		panic(err)
	}
	row.Exec(local.Rua, numero, local.Complemento, local.Bairro, local.Cidade, local.Estado, local.Id)
	if err != nil {
		fmt.Println("Erro ao atualizar local")
		w.WriteHeader(400)
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(local)

})
