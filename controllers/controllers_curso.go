package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"lupatini/models"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var ListCursos = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	data := models.Data{}
	professor := models.Usuario{}
	cursos := []models.Curso{}

	rows, err := connectingDB.Query("SELECT * FROM public.curso")
	if err != nil {
		fmt.Println("Não foi pesquisar as categorias dos cursos")
		w.WriteHeader(400)
		panic(err)
	}
	for rows.Next() {
		curso := models.Curso{}
		fmt.Println(" next cursos")
		err = rows.Scan(&curso.Id, &curso.Nome, &curso.Vagas,
			&curso.Descricao, &curso.Valor, &curso.Situacao, &curso.LocalId,
			&curso.Imagem, &curso.ProfessorId, &curso.CategoriaId, &curso.FormaPagamento)
		if err != nil {
			fmt.Println("Erro ao listar cursos")
			w.WriteHeader(400)
			panic(err)
		}

		//listar datas do curso
		rowsData, err := connectingDB.Query("SELECT * FROM public.curso_data WHERE curso_id=$1", curso.Id)
		if err != nil {
			fmt.Println("Não foi pesquisar as datas cursos")
			w.WriteHeader(400)
			panic(err)
		}
		var i = 0
		for rowsData.Next() {

			fmt.Println("next datas")
			err = rowsData.Scan(&data.Id, &data.Data, &data.HoraIni, &data.HoraFim, &data.Curso_id)
			if err != nil {
				fmt.Println("Erro ao listar cursos")
				w.WriteHeader(400)
				panic(err)
			}

			curso.Horarios = append(curso.Horarios, data)
			i++

		}
		//listar professor do curso
		var sql = "select usuario.id,usuario.nome,usuario.email," +
			" usuario.senha,usuario.sexo,usuario.perfil_id,usuario.token,usuario.imagem,usuario.cpf,usuario.celular" +
			" from curso" +
			" inner join usuario" +
			" on curso.professor_id=usuario.id" +
			" where curso.id = $1"
		rowsProfessor, err := connectingDB.Query(sql, curso.Id)
		if err != nil {
			fmt.Println("Não foi pesquisar o professor do curso")
			w.WriteHeader(400)
			return
		}
		for rowsProfessor.Next() {
			fmt.Println("next professor")
			err = rowsProfessor.Scan(&professor.Id, &professor.Nome, &professor.Email, &professor.Senha, &professor.Cpf, &professor.Celular,
				&professor.Sexo, &professor.Perfil, &professor.Token, &professor.Imagem)
			if err != nil {
				fmt.Println("Erro ao listar professor")
				fmt.Println(err)
				w.WriteHeader(400)
				return
			}

		}
		curso.Professor = professor

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
	data := models.Data{}
	professor := models.Usuario{}

	rows, err := connectingDB.Query("SELECT * FROM public.curso WHERE id=$1", id)
	if err != nil {
		fmt.Println("Não foi pesquisar o curso")
		w.WriteHeader(400)
		return
	}
	if rows.Next() == true {
		fmt.Println("true")
		err = rows.Scan(&curso.Id, &curso.Nome, &curso.Vagas, &curso.Descricao,
			&curso.Valor, &curso.Situacao, &curso.LocalId,
			&curso.Imagem, &curso.ProfessorId, &curso.CategoriaId, &curso.FormaPagamento)
		if err != nil {
			fmt.Println("Erro ao listar curso")
			w.WriteHeader(400)
			return
		}
		//listar datas do curso
		rowsData, err := connectingDB.Query("SELECT * FROM public.curso_data WHERE curso_id=$1", curso.Id)
		if err != nil {
			fmt.Println("Não foi pesquisar as datas cursos")
			w.WriteHeader(400)
			return
		}
		//	if rowsData.Next() {
		var i = 0
		for rowsData.Next() {
			fmt.Println("next datas")
			err = rowsData.Scan(&data.Id, &data.Data, &data.HoraIni, &data.HoraFim, &data.Curso_id)
			if err != nil {
				fmt.Println("Erro ao listar cursos")
				w.WriteHeader(400)
				return
			}
			curso.Horarios = append(curso.Horarios, data)
			i++
		}

		//	}
		//listar professor do curso
		var sql = "select usuario.id,usuario.nome,usuario.email," +
			" usuario.senha,usuario.sexo,usuario.perfil_id,usuario.token,usuario.imagem,usuario.cpf,usuario.celular" +
			" from curso" +
			" inner join usuario" +
			" on curso.professor_id=usuario.id" +
			" where curso.id = $1"
		fmt.Println(sql)
		rowsProfessor, err := connectingDB.Query(sql, curso.Id)
		if err != nil {
			fmt.Println("Não foi pesquisar o professor do curso")
			w.WriteHeader(400)
			return
		}
		for rowsProfessor.Next() {
			fmt.Println("next professor")
			err = rowsProfessor.Scan(&professor.Id, &professor.Nome, &professor.Email, &professor.Senha, &professor.Cpf, &professor.Celular,
				&professor.Sexo, &professor.Perfil, &professor.Token, &professor.Imagem)
			if err != nil {
				fmt.Println("Erro ao listar professor")
				fmt.Println(err)
				w.WriteHeader(400)
				return
			}

		}
		curso.Professor = professor

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(curso)
	} else {
		fmt.Println("Curso nao cadastrado")
		w.WriteHeader(400)
		return

	}

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
	var data models.Data
	dataCurso := r.FormValue("dataCurso")

	json.Unmarshal([]byte(dataCurso), &curso)

	fmt.Println(curso)

	//salva imagem do usuario
	file, handler, err := r.FormFile("imagem")
	if err != nil {
		fmt.Println("Arquivo vazio")
	} else {
		defer file.Close()
		f, err := os.OpenFile("/home/zaptec/img/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("Não é possível criar o arquivo")
		}
		defer f.Close()
		io.Copy(f, file)
		curso.Imagem = "/home/zaptec/img/" + handler.Filename
	}
	//var dataC int64
	sqlQuery := "INSERT INTO public.curso(nome,vagas,descricao,valor,situacao,local_id," +
		"imagem,professor_id,categoria_id,forma_pagamento)" +
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING ID"

	stmt, err := connectingDB.Prepare(sqlQuery)
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(curso.Nome, curso.Vagas, curso.Descricao, curso.Valor,
		curso.Situacao, curso.LocalId, curso.Imagem, curso.ProfessorId, curso.CategoriaId, curso.FormaPagamento).Scan(&curso.Id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("id:", curso.Id)
	//inseri datas do curso
	for index := range curso.Horarios {

		data.Data = curso.Horarios[index].Data
		data.HoraIni = curso.Horarios[index].HoraIni
		data.HoraFim = curso.Horarios[index].HoraFim

		sqlQuery := "INSERT INTO curso_data(data,curso_id,hora_inicio,hora_fim) VALUES($1,$2,$3,$4)"
		row, err := connectingDB.Exec(sqlQuery, data.Data, curso.Id, data.HoraIni, data.HoraFim)
		_ = row
		if err != nil {
			fmt.Println("Erro ao inserir data")
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(curso)

})

var AlterCurso = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	fmt.Println("estou alter")
	var curso models.Curso
	var data models.Data

	dataCurso := r.FormValue("dataCurso")
	json.Unmarshal([]byte(dataCurso), &curso)
	fmt.Println(curso)

	//salva imagem do usuario
	file, handler, err := r.FormFile("imagem")
	if err != nil {
		fmt.Println("Arquivo vazio")
	} else {
		defer file.Close()
		f, err := os.OpenFile("/home/zaptec/img/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("Não é possível criar o arquivo")
		}
		defer f.Close()
		io.Copy(f, file)
		curso.Imagem = "/home/zaptec/img/" + handler.Filename
	}

	row, err := connectingDB.Prepare("UPDATE public.curso SET nome=$1,vagas=$2,descricao=$3,valor=$4,situacao=$5,local_id=$6,imagem=$7,professor_id=$8,categoria_id=$9,forma_pagamento=$10 WHERE id=$11")
	if err != nil {
		fmt.Println(err)
		return
	}
	row.Exec(curso.Nome, curso.Vagas, curso.Descricao, curso.Valor, curso.Situacao,
		curso.LocalId, curso.Imagem, curso.ProfessorId, curso.CategoriaId, curso.FormaPagamento, curso.Id)
	if err != nil {
		fmt.Println(err)
		return
	}
	/*
		//converte string id_curso para int
		idCurso, err := strconv.ParseInt(curso.Id, 10, 64)
		if err == nil {
			fmt.Println(err)
		}*/
	//inseri datas do curso
	for index := range curso.Horarios {
		if index == 0 {
			sqlDelete := "DELETE FROM curso_data WHERE curso_id=$1"
			row, err := connectingDB.Exec(sqlDelete, curso.Id)
			if err != nil {
				fmt.Println("Erro ao excluir datas")
			}
			_ = row

		}

		data.Data = curso.Horarios[index].Data
		data.HoraIni = curso.Horarios[index].HoraIni
		data.HoraFim = curso.Horarios[index].HoraFim

		sqlQuery := "INSERT INTO curso_data(data,curso_id,hora_inicio,hora_fim) VALUES($1,$2,$3,$4)"
		row, err := connectingDB.Exec(sqlQuery, data.Data, curso.Id, data.HoraIni, data.HoraFim)
		_ = row
		if err != nil {
			fmt.Println("Erro ao inserir data")
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(curso)

})
