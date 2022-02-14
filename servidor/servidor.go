package servidor

import (
	"crud/banco"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `josn:"email"`
}

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	//Pacote iotiul lerá o corpo da requisição e retornar um erro casjo não tenha requisição
	CorpoDaRequiscao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Erro ao ler corpo da requisição!"))
		return
	}

	var usuario usuario //(Varialvel e Estrutura)
	//Converte usuario em inserido na rquisição em json para struct
	if erro = json.Unmarshal(CorpoDaRequiscao, &usuario); erro != nil {
		w.Write([]byte("Erro ao converter usuario para struct!"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados!"))
		return
	}
	defer db.Close() //fechando banco de dados

	//Prepara dados que serão inseridos no banco dedados
	//Evita o sql injection
	statement, erro := db.Prepare("insert into usuarios (nome, email) values (?, ?)")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close() //fechando statement

	inserindoDados, erro := statement.Exec(usuario.Nome, usuario.Email) //Obedecer ordem de inserção
	if erro != nil {
		w.Write([]byte("Erro ao inserir valores de usuários!"))
		return
	}

	//Pegando id inserido chamando "inserindoDados"
	idInserido, erro := inserindoDados.LastInsertId()
	if erro != nil {
		w.Write([]byte("Erro ao coletar ID da struct!"))
		return
	}

	//Status Codes e retorno de Resultado da inseção
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuario inserido com sucesso! Id: %d", idInserido)))
}

//Busca todos os usuários do banco de dados
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados!"))
		return
	}

	//Query que retorna linhas de todos os usuários
	linhas, erro := db.Query("Select * from usuarios")
	if erro != nil {
		w.Write([]byte("Erro ao buscar usuarios!"))
		return
	}
	defer linhas.Close()

	//Cria uma slice de usuários para ser alimentada
	//.Next vai iterar pelas linhas de usuarios
	var DadosUsuarios []usuario
	for linhas.Next() {
		var usuario usuario

		//.Scan define oq deseja pegar das linhas em ordem  e retorna um erro
		if erro := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro so escanear o usuários!"))
			return
		}
		//Após popular a variavel usuário devo transpostar esses dados para o slice usuários
		DadosUsuarios = append(DadosUsuarios, usuario)
	}
	w.WriteHeader(http.StatusOK)
	//decodificando dados pata json
	if erro := json.NewEncoder(w).Encode(DadosUsuarios); erro != nil {
		w.Write([]byte("Erro ao converter ususários para json!"))
		return
	}
}

//Busca de um ID de usuário especifico do banco de dados
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	//Usando o pacote mux
	//Retorna parametros da requisição
	parametros := mux.Vars(r)

	//convertendo parametro qu vem em string para uint32
	//strconv.ParceUint recebe três parametros: referencia que quer modificar,  base do numero e o tamanho dos bits
	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter parametro para Uint!"))
		return
	}

	//Posso abrir a coneção com o banco após converter o ID, para otimizar o processo de busca
	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados!"))
		return
	}

	//Coleta toda linha de um determinado ID
	linha, erro := db.Query("select * from usuarios where id = ?", ID)
	if erro != nil {
		w.Write([]byte("Erro ao coletar informações do ID!"))
		return
	}

	var BuscarUsuario usuario
	//Uso linha diretamente com o next e Scan SEM O FOR, porque so preciso de um valor
	if linha.Next() {
		if erro := linha.Scan(&BuscarUsuario.ID, &BuscarUsuario.Nome, &BuscarUsuario.Email); erro != nil {
			w.Write([]byte("Erro ao escanear o Usuário!"))
			return
		}
	}

	if erro := json.NewEncoder(w).Encode(BuscarUsuario); erro != nil {
		w.Write([]byte("Erro ao decodificar para json"))
		return
	}

}

//AtuualizarUsusario altera os dados de um usuário no banco
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter json para Uint32!"))
		return
	}

	corpoDaRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Erro ao ler corpo da rquisição"))
		return
	}

	//A Requisição está como json será passada para um struct após ter sido modificada pelo "Unmarshal"
	var auterarDado usuario
	if erro := json.Unmarshal(corpoDaRequisicao, &auterarDado); erro != nil {
		w.Write([]byte("Erro ao converter corpo da requisição!"))
		return
	}

	//Abrindo conexão com banco de dados
	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados!"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("update usuarios set nome = ?, email = ? where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(auterarDado.Nome, auterarDado.Email, ID); erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados!"))
		return
	}

	//Não utilizaremnos de retorno no update!
	w.WriteHeader(http.StatusNoContent)
}

//DeletarUsuario Remove usuario do banco de dados
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter para Uint32!"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com obanco de dados!"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("delete from usuarios where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao preparar statement"))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(ID); erro != nil {
		w.Write([]byte("Erro ao Deletar Usuário!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
