package main

import (
	"context"
	"fmt"
	"log"
	"os"

	//"time"

	//gocqlastra "github.com/datastax/gocql-astra"
	//"github.com/gocql/gocql"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

func resetarSupabase(conn *pgx.Conn) {
	query := `
	DO $$
	DECLARE
	    r RECORD;
	BEGIN
	    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
	        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
	    END LOOP;
	END $$;
	`

	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		fmt.Println("erro ao dropar tabelas: %v", err)
	}
}

func criarSupabase(conn *pgx.Conn) {
	query := `
	create table usuario (
	cpf varchar(15)
	,email varchar(30)
	,senha varchar(10)
	,primary key (cpf)
	);

	create table professor (
	nome varchar(30)
	,cpf varchar(15)
	,primary key (cpf)
	,foreign key (cpf) references usuario(cpf)
	);

	create table aluno (
	nome varchar(30)
	,cpf varchar(15)
	,primary key (cpf)
	,foreign key (cpf) references usuario(cpf)
	);

	create table curso (
	nome varchar(20)
	,cpf_autor varchar(15)
	,id varchar(5)
	,primary key (id)
	,foreign key(cpf_autor) references professor(cpf)
	);

	create table certificado (
	id varchar(5)
	,horas varchar(6)
	,primary key (id)
	);

	create table aluno_curso (
	id_curso varchar(5)
	,cpf_aluno varchar(15)
	,data_in varchar(10)
	,data_fim varchar(10)
	,id_certificado varchar(5)
	,foreign key(id_curso) references curso(id)
	,foreign key(cpf_aluno) references aluno(cpf)
	,foreign key(id_certificado) references certificado(id) 
	);
	`

	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		fmt.Println("erro ao criar tabelas: %v", err)
	}
}

func main() {
	//Conexao com o Supabase
	godotenv.Load()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar arquivo .env")
	}
	USUARIO := os.Getenv("USUARIO")
	PASSWORD := os.Getenv("PASSWORD")
	HOST := os.Getenv("HOST")
	PORT := os.Getenv("PORT")
	DBNAME := os.Getenv("DBNAME")
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", USUARIO, PASSWORD, HOST, PORT, DBNAME)
	conn, err := pgx.Connect(context.Background(), url)

	if err != nil {
		log.Fatalf("Falha ao conectar ao Supabase: %v", err)
	}
	fmt.Println("Conectado ao SupaBase com sucesso!")
	resetarSupabase(conn)
	criarSupabase(conn)

	defer conn.Close(context.Background())
}
