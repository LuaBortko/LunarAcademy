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
	,email varchar(50)
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
	nome varchar(50)
	,cpf_autor varchar(15)
	,id varchar(8)
	,primary key (id)
	,foreign key(cpf_autor) references professor(cpf)
	);

	create table certificado (
	id varchar(8)
	,horas varchar(6)
	,primary key (id)
	);

	create table aluno_curso (
	id_curso varchar(8)
	,cpf_aluno varchar(15)
	,data_in varchar(15)
	,data_fim varchar(15)
	,id_certificado varchar(8)
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

func inserirUsuario(conn *pgx.Conn, usuarios []Usuario){
	for i:= 0; i < len(usuarios); i++{
		u := usuarios[i]
		query := `INSERT INTO usuario (cpf, email, senha) VALUES ($1, $2, $3)`
		_, err := conn.Exec(context.Background(), query, u.cpf, u.email, u.senha)
		if err != nil {
			fmt.Print("erro ao inserir o usuario: ", u.cpf)
		}else{
			fmt.Print("Usuario ", u.cpf, " inserido com sucesso!\n")
		}
	}
	fmt.Print("\n")
}

func inserirProfessor(conn *pgx.Conn, professores []Professor){
	for i:= 0; i < len(professores); i++{
		u := professores[i]
		query := `INSERT INTO professor (nome, cpf) VALUES ($1, $2)`
		_, err := conn.Exec(context.Background(), query, u.nome, u.cpf)
		if err != nil {
			fmt.Print("erro ao inserir o professor: ", u.nome, "\n")
		}else{
			fmt.Print("Professor ", u.nome, " inserido com sucesso!\n")
		}
	}
	fmt.Print("\n")
}

func inserirAluno(conn *pgx.Conn, alunos []Aluno){
	for i:= 0; i < len(alunos); i++{
		u := alunos[i]
		query := `INSERT INTO aluno (nome, cpf) VALUES ($1, $2)`
		_, err := conn.Exec(context.Background(), query, u.nome, u.cpf)
		if err != nil {
			fmt.Print("erro ao inserir o aluno: ", u.nome, "\n")
		}else{
			fmt.Print("Aluno ", u.nome, " inserido com sucesso!\n")
		}
	}
	fmt.Print("\n")
}

func inserirCurso(conn *pgx.Conn, cursos []Curso){
	for i:= 0; i < len(cursos); i++{
		u := cursos[i]
		query := `INSERT INTO curso (nome, cpf_autor, id) VALUES ($1, $2, $3)`
		_, err := conn.Exec(context.Background(), query, u.nome, u.cpf_autor, u.id)
		if err != nil {
			fmt.Print("erro ao inserir o curso: ", u.nome, "\n")
		}else{
			fmt.Print("Curso ", u.nome, " inserido com sucesso!\n")
		}
	}
	fmt.Print("\n")
}

func inserirCertificado(conn *pgx.Conn, certificados []Certificado){
	for i:= 0; i < len(certificados); i++{
		u := certificados[i]
		query := `INSERT INTO certificado (id, horas) VALUES ($1, $2)`
		_, err := conn.Exec(context.Background(), query, u.id, u.horas)
		if err != nil {
			fmt.Print("erro ao inserir o certificado: ", u.id, "\n")
		}else{
			fmt.Print("Certificado ", u.id, " inserido com sucesso!\n")
		}
	}
	fmt.Print("\n")
}

func inserirAluno_Curso(conn *pgx.Conn, alunos_curso []Aluno_curso){
	for i:= 0; i < len(alunos_curso); i++{
		u := alunos_curso[i]

		if u.id_certificado == ""{
			query := `INSERT INTO aluno_curso (id_curso, cpf_aluno, data_in, data_fim, id_certificado) VALUES ($1, $2, $3, $4, $5)`
			_, err := conn.Exec(context.Background(), query, u.id_curso, u.cpf_aluno, u.data_in, u.data_fim, nil)
			if err != nil {
				fmt.Print(err)
				fmt.Print("erro ao inserir o Aluno_curso: ", u.cpf_aluno, " em ", u.id_curso, "\n")
			}else{
				fmt.Print("Aluno_curso ", u.cpf_aluno, " em ", u.id_curso, " inserido com sucesso!\n")
			}
		}else{
			query := `INSERT INTO aluno_curso (id_curso, cpf_aluno, data_in, data_fim, id_certificado) VALUES ($1, $2, $3, $4, $5)`
			_, err := conn.Exec(context.Background(), query, u.id_curso, u.cpf_aluno, u.data_in, u.data_fim, u.id_certificado)
			if err != nil {
				fmt.Print(err)
				fmt.Print("erro ao inserir o Aluno_curso: ", u.cpf_aluno, " em ", u.id_curso, "\n")
			}else{
				fmt.Print("Aluno_curso ", u.cpf_aluno, " em ", u.id_curso, " inserido com sucesso!\n")
			}
		}
	}
	fmt.Print("\n")
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

	//Criando as tabelas do supa
	resetarSupabase(conn)
	criarSupabase(conn)
	//gerando as listas
	usuarios := gerarUsuarios(10)
	//fmt.Println(usuarios)
	lista := []int{}
	tam_usuario := len(usuarios)
	for len(lista) < tam_usuario {
		num := randInt(0, tam_usuario)
		if numExiste(num, lista) {
			continue
		}
		lista = append(lista, num)
	}
	lista1 := lista[:(len(lista)/2)-1]
	lista2 := lista[(len(lista1)):]
	professores := gerarProfessores(usuarios, lista1)
	alunos := gerarAlunos(usuarios, lista2)
	cursos := gerarCursos(professores)
	alunos_curso := gerarAlunos_curso(alunos, cursos)
	certificados := []Certificado{}
	for i := range alunos_curso {
		if alunos_curso[i].id_certificado == "a" {
			certificado := gerarCertificado()
			alunos_curso[i].id_certificado = certificado.id
			certificados = append(certificados, certificado)
		}
	}
	fmt.Println(certificados)
	fmt.Print(alunos_curso)
	//cursos_mongo := gerarCurso_Mongo(cursos, professores)
	//professores_mongodb := gerarProfessores_Mongodb(professores, cursos)
	//historico := gerarHistorico(cursos, professores, alunos_curso)
	
	//inserindo no bamco
	inserirUsuario(conn,usuarios)
	inserirProfessor(conn, professores)
	inserirAluno(conn, alunos)
	inserirCurso(conn, cursos)
	inserirCertificado(conn, certificados)
	inserirAluno_Curso(conn, alunos_curso)


	defer conn.Close(context.Background())
}
