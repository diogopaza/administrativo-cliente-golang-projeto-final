package http_handlers

import (
	"encoding/json"
	"fmt"
	"lupatini/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var GetPaginacao = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tabela := vars["tabela"]
	pagina := vars["pagina"]
	curso := models.Curso{}
	cursos := []models.Curso{}

	//converte pagina para int
	paginaInt, err := strconv.ParseInt(pagina, 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	valor_offset := 10 * (paginaInt - 1)

	sql := "SELECT * from " + tabela + " ORDER BY id offset $1 limit 10"

	rows, err := connectingDB.Query(sql, valor_offset)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		fmt.Println("estou paginando")
		rows.Scan(&curso.Id, &curso.Nome, &curso.Vagas, &curso.Descricao, &curso.Valor, &curso.Situacao, &curso.LocalId,
			&curso.Imagem, &curso.ProfessorId, &curso.CategoriaId, &curso.FormaPagamento)

		cursos = append(cursos, curso)
	}

	json.NewEncoder(w).Encode(cursos)
})
