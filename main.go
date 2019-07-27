package main

import (
	"fmt"
	"lupatini/controllers"
	"net/http"

	"lupatini/http_handlers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	allowedHeaders := handlers.AllowedHeaders([]string{"Cache-Control", "X-Requested-With", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	//endpoints aluno
	router.Handle("/alunos", controllers.ListAlunos).Methods("GET")
	router.Handle("/aluno/{id}", controllers.ListAluno).Methods("GET")
	router.Handle("/aluno/{id}", controllers.DeleteAluno).Methods("DELETE")
	router.Handle("/aluno", controllers.InsertAluno).Methods("POST")
	router.Handle("/aluno", controllers.AlterAluno).Methods("PUT")
	router.Handle("/aluno", controllers.AlterImagemAluno).Methods("PUT")

	//endpoints categoria
	router.Handle("/categorias", controllers.ListCategorias).Methods("GET")
	router.Handle("/categoria/{id}", controllers.ListCategoria).Methods("GET")
	router.Handle("/categoria/{id}", controllers.DeleteCategoria).Methods("DELETE")
	router.Handle("/categoria", controllers.InsertCategoria).Methods("POST")
	router.Handle("/categoria", controllers.AlterCategoria).Methods("PUT")

	//endpoints curso
	router.Handle("/cursos", controllers.ListCursos).Methods("GET")
	router.Handle("/curso/{id}", controllers.ListCurso).Methods("GET")
	router.Handle("/curso/{id}", controllers.DeleteCurso).Methods("DELETE")
	router.Handle("/curso", controllers.InsertCurso).Methods("POST")
	router.Handle("/curso", controllers.AlterCurso).Methods("PUT")

	//endpoints usuario
	router.Handle("/usuarios", controllers.ListUsuarios).Methods("GET")
	router.Handle("/usuario/{id}", controllers.ListUsuario).Methods("GET")
	router.Handle("/usuario-categoria/{id}", controllers.ListCategoriaUsuario).Methods("GET")
	router.Handle("/usuario/{id}", controllers.DeleteUsuario).Methods("DELETE")
	router.Handle("/usuario", controllers.InsertUsuario).Methods("POST")
	router.Handle("/usuario", controllers.AlterUsuario).Methods("PUT")

	//endpoints local
	router.Handle("/usuarios", controllers.ListUsuarios).Methods("GET")
	router.Handle("/usuario/{id}", controllers.ListUsuario).Methods("GET")
	router.Handle("/usuario-categoria/{id}", controllers.ListCategoriaUsuario).Methods("GET")
	router.Handle("/usuario/{id}", controllers.DeleteUsuario).Methods("DELETE")
	router.Handle("/usuario", controllers.InsertUsuario).Methods("POST")
	router.Handle("/usuario", controllers.AlterUsuario).Methods("PUT")

	//endpoints system
	router.Handle("/login", http_handlers.GetLogin).Methods("POST")

	fmt.Println("Rodando na 3000")
	http.ListenAndServe(":3000", handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router))

}
