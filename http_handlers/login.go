package http_handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"net/http"
	"strings"
	"lupatini/models"
	"lupatini/config"
	"github.com/dgrijalva/jwt-go"
)

var connectingDB *sql.DB
var myKey = []byte("secret")
func init(){
	connectingDB = config.DB()		
}

var GetLogin = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var usuario models.Usuario

	erro := r.ParseForm()
	if erro != nil {
		panic(erro)
	}
	usuario_login := r.FormValue("usuario")
	senha_login := r.FormValue("senha")
	
	sqlQuery := "SELECT id, nome, senha FROM public.usuario WHERE nome=$1"

	row := connectingDB.QueryRow(sqlQuery, usuario_login)

	err := row.Scan(&usuario.Id, &usuario.Nome, &usuario.Senha)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Usuario não encontrado")
			w.WriteHeader(401)
			
		}
		fmt.Println("Erros:%v", err)
	} else {
		fmt.Println("Achou o usuario")
		fmt.Println("ID:" + usuario.Id)
		
		if senha_login != usuario.Senha {
			w.WriteHeader(401)
			fmt.Println("Erro senha incorreta")
		} else {
			w.Header().Set("Content-Type", "application/json")

			myToken := GetToken(usuario)
			TokenSplit := strings.Split(myToken, ".")
			fmt.Println("Split: ", TokenSplit[2])

			var m = make(map[string]string)
			m["token"] = myToken

			SaveToken(TokenSplit[2], usuario.Id)
			json.NewEncoder(w).Encode(m)

		}

	}

})

func GetToken(u models.Usuario) string{

	token := jwt.New(jwt.SigningMethodHS256)
	
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"]=true
	claims["id"]= u.Id
	claims["name"]= u.Nome
	claims["password"]= u.Senha
	claims["exp"]=time.Now().Add(time.Minute * 30).Unix()
	
	tokenString, _ := token.SignedString(myKey)
	
	return tokenString

}

func SignedToken(name string) (token string){

	var userToken string

	sqlQuery := "SELECT token FROM public.user WHERE name=" + "'" + name + "'"
	rows, err := connectingDB.Query(sqlQuery)
	if err != nil{
		fmt.Println("Erro ao consultar usuario no banco de dados")
		return 
	}
	
	for rows.Next(){

		err = rows.Scan(&token)
		if err != nil{
			fmt.Println("Erro ao percorrer no banco de dados")
			return 
		}

		userToken = token	
	

	}
	
	return userToken

} 

func SaveToken( signatureToken string, id string){
	
	sqlQuery := "UPDATE public.usuario SET token=$2 WHERE id=$1"
	_ ,err := connectingDB.Exec(sqlQuery, id, signatureToken)
	if err != nil {
		panic(err)
	}
	fmt.Println("Token atualizado com sucesso")

}  