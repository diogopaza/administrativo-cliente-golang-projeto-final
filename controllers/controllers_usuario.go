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

var InsertUsuario = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var usuario models.Usuario

	dataUsuario := r.FormValue("dataUsuario")

	json.Unmarshal([]byte(dataUsuario), &usuario)
	fmt.Println(usuario)

	//salva imagem do usuario
	file, handler, err := r.FormFile("selectedFile")
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

		usuario.Imagem = "/home/zaptec/img/" + handler.Filename
	}

	sqlQuery := "INSERT INTO public.usuario(nome,email,senha,sexo,perfil_id,token," +
		"imagem,cpf,celular) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)"
	row, err := connectingDB.Exec(sqlQuery, usuario.Nome, usuario.Email,
		usuario.Senha, usuario.Sexo, usuario.Perfil, usuario.Token, usuario.Imagem, &usuario.Cpf, &usuario.Celular)
	if err != nil {
		fmt.Println("Erro ao inserir dados")
		fmt.Println(err)
		return
	}
	_ = row
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(usuario)

})

var DeleteUsuario = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("id: ", id)
	sqlQuery := "DELETE FROM public.usuario WHERE id=$1"
	rows, err := connectingDB.Exec(sqlQuery, id)

	if err != nil {

		fmt.Println("Erro ao excluir usuario")
		w.WriteHeader(500)
		return

	}
	rowsDeleted, err := rows.RowsAffected()
	if err != nil {
		fmt.Println("Erro ao percorrer linha a ser excluida")
		w.WriteHeader(500)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(rowsDeleted)

})

var ListCategoriaUsuario = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	usuario := models.Usuario{}
	usuarios := []models.Usuario{}
	perfil := models.Perfil{}
	sql := "SELECT * from usuario" +
		" INNER JOIN perfil" +
		" ON perfil.id=usuario.perfil_id" +
		" WHERE perfil.id=$1"
	rows, err := connectingDB.Query(sql, id)
	if err != nil {
		fmt.Println("Não foi pesquisar usuários")
		w.WriteHeader(400)
		panic(err)
	}
	for rows.Next() {

		err = rows.Scan(&usuario.Id, &usuario.Nome, &usuario.Email, &usuario.Senha, &usuario.Sexo,
			&usuario.Perfil, &usuario.Token, &usuario.Imagem,
			&usuario.Cpf, &usuario.Celular, &perfil.Id, &perfil.Nome)
		if err != nil {
			fmt.Println("Erro ao listar usarios")
			w.WriteHeader(400)
			panic(err)
		}
		usuarios = append(usuarios, usuario)

	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(usuarios)

})

var ListUsuario = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entrei id")
	vars := mux.Vars(r)
	id := vars["id"]

	usuario := models.Usuario{}

	rows, err := connectingDB.Query("SELECT * FROM public.usuario WHERE id=$1", id)
	if err != nil {
		fmt.Println("Não foi pesquisar usuários")
		w.WriteHeader(400)
		panic(err)
	}
	for rows.Next() {

		err = rows.Scan(&usuario.Id, &usuario.Nome, &usuario.Email, &usuario.Senha,
			&usuario.Sexo, &usuario.Perfil, &usuario.Token, &usuario.Imagem, &usuario.Cpf, &usuario.Celular)
		if err != nil {
			fmt.Println("Erro ao listar usarios")
			w.WriteHeader(400)
			panic(err)
		}

	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(usuario)

})

var ListUsuarios = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	rows, err := connectingDB.Query("SELECT * FROM public.usuario")
	if err != nil {
		fmt.Println("Não foi pesquisar usuários")
	}

	usuario := models.Usuario{}
	usuarios := []models.Usuario{}

	for rows.Next() {
		fmt.Println(" next usuários")
		err = rows.Scan(&usuario.Id, &usuario.Nome, &usuario.Email, &usuario.Senha, &usuario.Sexo,
			&usuario.Perfil, &usuario.Token, &usuario.Imagem, &usuario.Cpf, &usuario.Celular)
		if err != nil {
			fmt.Println("Erro ao listar usarios")
			fmt.Println(err)
		}

		usuarios = append(usuarios, usuario)

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(usuarios)
	w.WriteHeader(http.StatusOK)

})

var AlterUsuario = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("estou alter")
	var usuario models.Usuario

	dataUsuario := r.FormValue("dataUsuario")

	json.Unmarshal([]byte(dataUsuario), &usuario)
	fmt.Println(usuario)

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

		usuario.Imagem = "/home/zaptec/img/" + handler.Filename
	}

	row, err := connectingDB.Prepare("UPDATE public.usuario SET nome=$1,email=$2,senha=$3,sexo=$4,perfil_id=$5,token=$6,imagem=$7,cpf=$8,celular=$9 WHERE id=$10")
	if err != nil {
		fmt.Println(err)
		return
	}
	row.Exec(usuario.Nome, usuario.Email, usuario.Senha, usuario.Sexo, usuario.Perfil,
		usuario.Token, usuario.Imagem, usuario.Cpf, usuario.Celular, usuario.Id)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(usuario)
	w.WriteHeader(http.StatusOK)

})
