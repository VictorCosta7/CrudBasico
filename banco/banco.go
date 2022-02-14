package banco

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //Driver para conetar com mysql
)

func Conectar() (*sql.DB, error) {
	stringDeConexao := "victor:Victor0099@/golangproj?charset=utf8&parseTime=True&loc=Local"

	//Abrindo conexão
	db, erro := sql.Open("mysql", stringDeConexao)
	if erro != nil {
		return nil, erro
	}
	//Verificando conexão com o banco de dados
	if erro = db.Ping(); erro != nil {
		return nil, erro
	}
	return db, nil
}
