package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"lupatini/models"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var InsertUsuario = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var usuario models.Usuario

	var imagemGravaBanco string
	dataUsuario := r.FormValue("dataUsuario")
	fmt.Println(dataUsuario)

	json.Unmarshal([]byte(dataUsuario), &usuario)
	fmt.Println(usuario)

	//salva imagem do usuario
	file, handler, err := r.FormFile("selectedFile")
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

	imagemGravaBanco = "/home/zaptec/img/" + handler.Filename

	sqlQuery := "INSERT INTO public.usuario(nome,email,senha,cpf,celular,sexo,perfil_id,token,imagem) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)"
	row, err := connectingDB.Exec(sqlQuery, usuario.Nome, usuario.Email,
		usuario.Senha, usuario.Cpf, usuario.Celular, usuario.Sexo, usuario.Perfil, &usuario.Token, imagemGravaBanco)
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

	sqlQuery := "DELETE FROM public.usuario WHERE id = $1"
	rows, err := connectingDB.Exec(sqlQuery, id)

	if err != nil {
		fmt.Println("Erro ao excluir usuario")
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteHeader(400)

	}
	rowsDeleted, err := rows.RowsAffected()
	if err != nil {
		fmt.Println("Erro ao percorrer linha a ser excluida")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(rowsDeleted)

})

var ListCategoriaUsuario = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println(id)
	usuario := models.Usuario{}
	usuarios := []models.Usuario{}
	perfil := models.Perfil{}
	sql := "SELECT * from usuario" +
		" INNER JOIN perfil" +
		" ON perfil.id=usuario.perfil_id" +
		" WHERE perfil.id=$1"
	fmt.Println(sql)
	rows, err := connectingDB.Query(sql, id)
	if err != nil {
		fmt.Println("Não foi pesquisar usuários")
		w.WriteHeader(400)
		panic(err)
	}
	for rows.Next() {

		err = rows.Scan(&usuario.Id, &usuario.Nome, &usuario.Email, &usuario.Senha, &usuario.Cpf, &usuario.Celular,
			&usuario.Sexo, &usuario.Perfil, &usuario.Token, &perfil.Id, &perfil.Nome)
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
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		panic(err)
	}
	for rows.Next() {

		err = rows.Scan(&usuario.Id, &usuario.Nome, &usuario.Email, &usuario.Senha,
			&usuario.Cpf, &usuario.Celular, &usuario.Sexo, &usuario.Perfil, &usuario.Token)
		if err != nil {
			fmt.Println("Erro ao listar usarios")
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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
		err = rows.Scan(&usuario.Id, &usuario.Nome, &usuario.Email, &usuario.Senha,
			&usuario.Cpf, &usuario.Celular, &usuario.Sexo, &usuario.Perfil, &usuario.Token)
		if err != nil {
			fmt.Println("Erro ao listar usarios")
		}

		usuarios = append(usuarios, usuario)

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(usuarios)
	w.WriteHeader(http.StatusOK)

})

var AlterUsuario = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var usuario models.Usuario

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &usuario)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(usuario)
	row, err := connectingDB.Prepare("UPDATE public.usuario SET nome=$1,email=$2,senha=$3,cpf=$4,celular=$5,sexo=$6,perfil_id=$7,token=$8 WHERE id=$9")
	if err != nil {
		fmt.Println(err)
		return
	}
	row.Exec(usuario.Nome, usuario.Email, usuario.Senha, usuario.Cpf, usuario.Celular,
		usuario.Sexo, usuario.Perfil, usuario.Token, usuario.Id)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(usuario)
	w.WriteHeader(http.StatusOK)

})
