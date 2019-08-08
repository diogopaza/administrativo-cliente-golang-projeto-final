package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lupatini/models"
	"net/http"

	"github.com/gorilla/mux"
)

var ListDatas = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	rows, err := connectingDB.Query("SELECT * FROM public.curso_data")
	if err != nil {
		fmt.Println("Não foi pesquisar as datas")
		w.WriteHeader(500)
		panic(err)
	}

	data := models.Data{}
	datas := []models.Data{}

	for rows.Next() {
		fmt.Println(" next datas")
		err = rows.Scan(&data.Id, &data.Data, &data.Curso_id, &data.HoraIni, &data.HoraFim)

		if err != nil {
			fmt.Println("Erro ao listar cursos")
			w.WriteHeader(400)
			panic(err)
		}

		datas = append(datas, data)

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(datas)
	w.WriteHeader(http.StatusOK)

})

var ListData = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	data := models.Data{}

	rows, err := connectingDB.Query("SELECT * FROM public.curso_data WHERE id=$1", id)
	if err != nil {
		fmt.Println("Não foi pesquisar o data")
		w.WriteHeader(500)
		panic(err)
	}
	for rows.Next() {

		err = rows.Scan(&data.Id, &data.Data, &data.Curso_id, &data.HoraIni, &data.HoraFim)
		if err != nil {
			fmt.Println("Erro ao listar data")
			w.WriteHeader(500)
			panic(err)
		}

	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(data)

})

var InsertData = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var data models.Data

	dadosDatas := r.FormValue("dadosData")

	json.Unmarshal([]byte(dadosDatas), &data)
	fmt.Println(data)

	//data.Curso_id = "1"
	sqlQuery := "INSERT INTO public.curso_data(data,curso_id,hora_inicio,hora_fim) VALUES($1,$2,$3,$4)"
	row, err := connectingDB.Exec(sqlQuery, data.Data, data.Curso_id, data.HoraIni, data.HoraFim)
	_ = row
	if err != nil {
		fmt.Println("Erro ao inserir data")
		fmt.Println(err)
		w.WriteHeader(500)
		return

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(data)

})

var DeleteData = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	sqlQuery := "DELETE FROM public.curso_data WHERE id = $1"
	rows, err := connectingDB.Exec(sqlQuery, id)

	if err != nil {
		fmt.Println("Erro ao excluir data")
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteHeader(400)
		panic(err)

	}
	rowsDeleted, err := rows.RowsAffected()
	if err != nil {
		fmt.Println("Erro ao percorrer a data a ser excluida")
		w.WriteHeader(500)
		panic(err)
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(rowsDeleted)

})

var AlterData = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	fmt.Println("estou alter")
	var data models.Data

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	row, err := connectingDB.Prepare("UPDATE public.curso_data SET data=$1,curso_id=$2,hora_inicio=$3,hora_fim=$4 WHERE id=$5")
	if err != nil {
		fmt.Println(err)
		return
	}
	row.Exec(data.Data, data.Curso_id, data.HoraIni, data.HoraFim, data.Id)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
	w.WriteHeader(http.StatusOK)

})
