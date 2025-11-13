package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	//"time"

	gocqlastra "github.com/datastax/gocql-astra"
	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type cursoMongoLeitura struct {
	Nome        string `bson:"nome"`
	Autor       string `bson:"autor"`
	Descricao   string `bson:"descricao"`
	Requisitos  string `bson:"requisitos"`
	Preco       string `bson:"preco"`
	Dificuldade string `bson:"dificuldade"`
	Avaliacao   string `bson:"avaliacao"`
}

type professorMongoLeitura struct {
	Nome            string `bson:"nome"`
	Cpf             string `bson:"cpf"`
	Formacao        string `bson:"formacao"`
	TempoPlataforma string `bson:"tempo_plataforma"`
	QtdeCursos      string `bson:"qtde_cursos"`
}

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
		fmt.Print("erro ao dropar tabelas: %v", err)
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
		fmt.Print("erro ao criar tabelas: %v", err)
	}
}

func inserirUsuario(conn *pgx.Conn, usuarios []Usuario) {
	for i := 0; i < len(usuarios); i++ {
		u := usuarios[i]
		query := `INSERT INTO usuario (cpf, email, senha) VALUES ($1, $2, $3)`
		_, err := conn.Exec(context.Background(), query, u.cpf, u.email, u.senha)
		if err != nil {
			fmt.Print("erro ao inserir o usuario: ", u.cpf)
		} else {
			fmt.Print("Usuario ", u.cpf, " inserido com sucesso!\n")
		}
	}
	fmt.Print("\n")
}

func inserirProfessor(conn *pgx.Conn, professores []Professor) {
	for i := 0; i < len(professores); i++ {
		u := professores[i]
		query := `INSERT INTO professor (nome, cpf) VALUES ($1, $2)`
		_, err := conn.Exec(context.Background(), query, u.nome, u.cpf)
		if err != nil {
			fmt.Print("erro ao inserir o professor: ", u.nome, "\n")
		} else {
			fmt.Print("Professor ", u.nome, " inserido com sucesso!\n")
		}
	}
	fmt.Print("\n")
}

func inserirAluno(conn *pgx.Conn, alunos []Aluno) {
	for i := 0; i < len(alunos); i++ {
		u := alunos[i]
		query := `INSERT INTO aluno (nome, cpf) VALUES ($1, $2)`
		_, err := conn.Exec(context.Background(), query, u.nome, u.cpf)
		if err != nil {
			fmt.Print("erro ao inserir o aluno: ", u.nome, "\n")
		} else {
			fmt.Print("Aluno ", u.nome, " inserido com sucesso!\n")
		}
	}
	fmt.Print("\n")
}

func inserirCurso(conn *pgx.Conn, cursos []Curso) {
	for i := 0; i < len(cursos); i++ {
		u := cursos[i]
		query := `INSERT INTO curso (nome, cpf_autor, id) VALUES ($1, $2, $3)`
		_, err := conn.Exec(context.Background(), query, u.nome, u.cpf_autor, u.id)
		if err != nil {
			fmt.Print("erro ao inserir o curso: ", u.nome, "\n")
		} else {
			fmt.Print("Curso ", u.nome, " inserido com sucesso!\n")
		}
	}
	fmt.Print("\n")
}

func inserirCertificado(conn *pgx.Conn, certificados []Certificado) {
	for i := 0; i < len(certificados); i++ {
		u := certificados[i]
		query := `INSERT INTO certificado (id, horas) VALUES ($1, $2)`
		_, err := conn.Exec(context.Background(), query, u.id, u.horas)
		if err != nil {
			fmt.Print("erro ao inserir o certificado: ", u.id, "\n")
		} else {
			fmt.Print("Certificado ", u.id, " inserido com sucesso!\n")
		}
	}
	fmt.Print("\n")
}

func inserirAluno_Curso(conn *pgx.Conn, alunos_curso []Aluno_curso) {
	for i := 0; i < len(alunos_curso); i++ {
		u := alunos_curso[i]

		if u.id_certificado == "" {
			query := `INSERT INTO aluno_curso (id_curso, cpf_aluno, data_in, data_fim, id_certificado) VALUES ($1, $2, $3, $4, $5)`
			_, err := conn.Exec(context.Background(), query, u.id_curso, u.cpf_aluno, u.data_in, u.data_fim, nil)
			if err != nil {
				fmt.Print(err)
				fmt.Print("erro ao inserir o Aluno_curso: ", u.cpf_aluno, " em ", u.id_curso, "\n")
			} else {
				fmt.Print("Aluno_curso ", u.cpf_aluno, " em ", u.id_curso, " inserido com sucesso!\n")
			}
		} else {
			query := `INSERT INTO aluno_curso (id_curso, cpf_aluno, data_in, data_fim, id_certificado) VALUES ($1, $2, $3, $4, $5)`
			_, err := conn.Exec(context.Background(), query, u.id_curso, u.cpf_aluno, u.data_in, u.data_fim, u.id_certificado)
			if err != nil {
				fmt.Print(err)
				fmt.Print("erro ao inserir o Aluno_curso: ", u.cpf_aluno, " em ", u.id_curso, "\n")
			} else {
				fmt.Print("Aluno_curso ", u.cpf_aluno, " em ", u.id_curso, " inserido com sucesso!\n")
			}
		}
	}
	fmt.Print("\n")
}

func lerUsuarios(conn *pgx.Conn) []Usuario {
	rows, err := conn.Query(context.Background(), "SELECT cpf, email, senha FROM usuario")
	if err != nil {
		fmt.Errorf("erro ao executar SELECT: %v", err)
		return nil
	}
	defer rows.Close()

	var usuarios []Usuario

	for rows.Next() {
		var u Usuario
		err := rows.Scan(&u.cpf, &u.email, &u.senha)
		if err != nil {
			fmt.Errorf("erro ao ler linha: %v", err)
			return nil
		}
		usuarios = append(usuarios, u)
	}
	return usuarios
}

func lerProfessores(conn *pgx.Conn) []Professor {
	rows, err := conn.Query(context.Background(), "SELECT nome, cpf FROM professor")
	if err != nil {
		fmt.Errorf("erro ao executar SELECT: %v", err)
		return nil
	}
	defer rows.Close()

	var professores []Professor

	for rows.Next() {
		var u Professor
		err := rows.Scan(&u.nome, &u.cpf)
		if err != nil {
			fmt.Errorf("erro ao ler linha: %v", err)
			return nil
		}
		professores = append(professores, u)
	}
	return professores
}

func lerAlunos(conn *pgx.Conn) []Aluno {
	rows, err := conn.Query(context.Background(), "SELECT nome, cpf FROM aluno")
	if err != nil {
		fmt.Errorf("erro ao executar SELECT: %v", err)
		return nil
	}
	defer rows.Close()

	var alunos []Aluno

	for rows.Next() {
		var u Aluno
		err := rows.Scan(&u.nome, &u.cpf)
		if err != nil {
			fmt.Errorf("erro ao ler linha: %v", err)
			return nil
		}
		alunos = append(alunos, u)
	}
	return alunos
}

func lerCursos(conn *pgx.Conn) []Curso {
	rows, err := conn.Query(context.Background(), "SELECT nome, cpf_autor, id FROM curso")
	if err != nil {
		fmt.Errorf("erro ao executar SELECT: %v", err)
		return nil
	}
	defer rows.Close()

	var cursos []Curso

	for rows.Next() {
		var u Curso
		err := rows.Scan(&u.nome, &u.cpf_autor, &u.id)
		if err != nil {
			fmt.Errorf("erro ao ler linha: %v", err)
			return nil
		}
		cursos = append(cursos, u)
	}
	return cursos
}

func lerCertificados(conn *pgx.Conn) []Certificado {
	rows, err := conn.Query(context.Background(), "SELECT id, horas FROM certificado")
	if err != nil {
		fmt.Errorf("erro ao executar SELECT: %v", err)
		return nil
	}
	defer rows.Close()

	var certificados []Certificado

	for rows.Next() {
		var u Certificado
		err := rows.Scan(&u.id, &u.horas)
		if err != nil {
			fmt.Errorf("erro ao ler linha: %v", err)
			return nil
		}
		certificados = append(certificados, u)
	}
	return certificados
}

func lerAlunos_Curso(conn *pgx.Conn) []Aluno_curso {
	rows, err := conn.Query(context.Background(),
		"SELECT id_curso, cpf_aluno, data_in, data_fim, id_certificado FROM aluno_curso")
	if err != nil {
		log.Printf("Erro ao executar SELECT: %v", err)
		return nil
	}
	defer rows.Close()

	var alunos_curso []Aluno_curso

	for rows.Next() {
		var u Aluno_curso
		var idCertificado sql.NullString // 👈 suporta valores NULL

		err := rows.Scan(&u.id_curso, &u.cpf_aluno, &u.data_in, &u.data_fim, &idCertificado)
		if err != nil {
			log.Printf("Erro ao ler linha: %v", err)
			continue
		}
		if idCertificado.Valid {
			u.id_certificado = idCertificado.String
		} else {
			u.id_certificado = ""
		}

		alunos_curso = append(alunos_curso, u)
	}
	return alunos_curso
}

func resetarBancos(conn *pgx.Conn, session *gocql.Session, ctx context.Context, client *mongo.Client, dbName string) {
	resetarSupabase(conn)
	criarSupabase(conn)
	resetarCassandra(session)
	resetarMongoDB(ctx, client, dbName)
}

// MongoDB

func resetarMongoDB(ctx context.Context, client *mongo.Client, dbName string) {
	err := client.Database(dbName).Drop(ctx)
	if err != nil {
		fmt.Print("Erro ao dropar MongoDB")
		return
	}
}

func inserirCurso_Mongo(ctx context.Context, client *mongo.Client, dbName string, cursos []Curso_Mongo) {
	coll := client.Database(dbName).Collection("cursos")

	var docs []interface{}
	for _, c := range cursos {
		docs = append(docs, bson.M{
			"nome":        c.nome,
			"autor":       c.autor,
			"descricao":   c.descricao,
			"requisitos":  c.requisitos,
			"preco":       c.preco,
			"dificuldade": c.dificuldade,
			"avaliacao":   c.avaliacao,
		})
	}
	_, err := coll.InsertMany(ctx, docs)
	if err != nil {
		fmt.Println("Erro ao inserir cursos:", err)
	} else {
		fmt.Printf("\n%d cursos inseridos com sucesso\n", len(cursos))
	}
}

func inserirProfessor_MongoDB(ctx context.Context, client *mongo.Client, dbName string, professores []Professor_Mongodb) {
	coll := client.Database(dbName).Collection("professores")

	var docs []interface{}
	for _, p := range professores {
		docs = append(docs, bson.M{
			"nome":             p.nome,
			"cpf":              p.cpf,
			"formacao":         p.formacao,
			"tempo_plataforma": p.tempo_plataforma,
			"qtde_cursos":      p.qtde_cursos,
		})
	}
	_, err := coll.InsertMany(ctx, docs)
	if err != nil {
		fmt.Println("\nErro ao inserir professores:", err)
	}
}

func lerCursos_Mongo(ctx context.Context, client *mongo.Client, dbName string) []Curso_Mongo {
	coll := client.Database(dbName).Collection("cursos")

	cursor, err := coll.Find(ctx, struct{}{})
	if err != nil {
		fmt.Println("Erro ao buscar cursos:", err)
		return nil
	}
	defer cursor.Close(ctx)

	var cursos []Curso_Mongo
	for cursor.Next(ctx) {
		var tmp cursoMongoLeitura
		if err := cursor.Decode(&tmp); err != nil {
			fmt.Println("Erro ao decodificar curso:", err)
			continue
		}

		// converte para seu tipo original
		cursos = append(cursos, Curso_Mongo{
			nome:        tmp.Nome,
			autor:       tmp.Autor,
			descricao:   tmp.Descricao,
			requisitos:  tmp.Requisitos,
			preco:       tmp.Preco,
			dificuldade: tmp.Dificuldade,
			avaliacao:   tmp.Avaliacao,
		})
	}

	if err := cursor.Err(); err != nil {
		fmt.Println("Erro no cursor:", err)
	}
	return cursos
}

func lerProfessores_MongoDB(ctx context.Context, client *mongo.Client, dbName string) []Professor_Mongodb {
	coll := client.Database(dbName).Collection("professores")

	cursor, err := coll.Find(ctx, struct{}{})
	if err != nil {
		fmt.Println("Erro ao buscar professores:", err)
		return nil
	}
	defer cursor.Close(ctx)

	var professores []Professor_Mongodb
	for cursor.Next(ctx) {
		var tmp professorMongoLeitura
		if err := cursor.Decode(&tmp); err != nil {
			fmt.Println("Erro ao decodificar professor:", err)
			continue
		}

		professores = append(professores, Professor_Mongodb{
			nome:             tmp.Nome,
			cpf:              tmp.Cpf,
			formacao:         tmp.Formacao,
			tempo_plataforma: tmp.TempoPlataforma,
			qtde_cursos:      tmp.QtdeCursos,
		})
	}

	if err := cursor.Err(); err != nil {
		fmt.Println("Erro no cursor:", err)
	}
	return professores
}

// Cassandra

func resetarCassandra(session *gocql.Session) {
	session.Query(`DROP TABLE IF EXISTS historico`).Exec()

	err := session.Query(`
		CREATE TABLE IF NOT EXISTS historico (
		id UUID PRIMARY KEY,
		nome text,
		autor text,
		avaliacao text,
		data_inicio text,
		data_final text)
	`).Exec()

	if err != nil {
		log.Fatalf("Erro ao recriar tabela: %v", err)
	}
}

func inserirHistorico_Cassandra(session *gocql.Session, historico []Historico) {
	for _, h := range historico {
		id := gocql.TimeUUID()
		err := session.Query(`
		INSERT INTO historico (id, nome, autor, avaliacao, data_inicio, data_final)
		VALUES (?, ?, ? , ?, ?, ?)`,
			id, h.nome, h.autor, h.avaliacao, h.data_inicio, h.data_final,
		).Exec()
		if err != nil {
			log.Printf("Erro ao inserir historico: %v", err)
		}
	}
}

func lerHistorico_Cassandra(session *gocql.Session) []Historico {
	iter := session.Query(`SELECT nome, autor, avaliacao, data_inicio, data_final FROM historico`).Iter()

	var historicos []Historico
	var h Historico

	for iter.Scan(&h.nome, &h.autor, &h.avaliacao, &h.data_inicio, &h.data_final) {
		historicos = append(historicos, h)
	}

	if err := iter.Close(); err != nil {
		log.Printf("Erro ao ler historico: %v", err)
	}
	return historicos
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
	//fmt.Println("Conectado ao SupaBase com sucesso!")

	// // MongoDB

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Println("MONGO_URI não configurada, pulando conexão MongoDB")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Falha ao conectar ao MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	dbName := "CursosProfessores"

	// Cassandra

	if os.Getenv("ASTRA_DB_ID") == "" || os.Getenv("APPLICATION_TOKEN") == "" {
		log.Println("Variáveis de ambiente do Astra DB não configuradas, pulando conexão Cassandra.")
		return
	}

	cluster, err := gocqlastra.NewClusterFromURL(
		"https://api.astra.datastax.com",
		os.Getenv("ASTRA_DB_ID"),
		os.Getenv("APPLICATION_TOKEN"),
		10*time.Second,
	)
	if err != nil {
		log.Fatalf("Erro ao carregar cluster Astra: %v", err)
	}

	cluster.Timeout = 30 * time.Second
	cluster.Keyspace = "default_keyspace"
	session, err := gocql.NewSession(*cluster)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Cassandra: %v", err)
	}

	var i int
	fmt.Print("Resetar bancos?(1-sim,2-não)")
	fmt.Scan(&i)
	if i == 1 {
		resetarBancos(conn, session, ctx, client, dbName)
	}

	var j int
	fmt.Print("\nGerar mais dados?(1-sim,2-não)")
	fmt.Scan(&j)
	if j == 1 {
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

		//inserindo no bamco
		inserirUsuario(conn, usuarios)
		inserirProfessor(conn, professores)
		inserirAluno(conn, alunos)
		inserirCurso(conn, cursos)
		inserirCertificado(conn, certificados)
		inserirAluno_Curso(conn, alunos_curso)

		//resetarMongoDB(ctx, client, dbName)

		cursos_mongo := gerarCurso_Mongo(cursos, professores)
		professores_mongodb := gerarProfessores_Mongodb(professores, cursos)
		inserirCurso_Mongo(ctx, client, dbName, cursos_mongo)
		inserirProfessor_MongoDB(ctx, client, dbName, professores_mongodb)
		fmt.Print("\nDados do MongoDB inseridos\n")

		//resetarCassandra(session)
		historico := gerarHistorico(cursos, professores, alunos_curso)
		inserirHistorico_Cassandra(session, historico)
		fmt.Print("\nDados do Cassandra inseridos")
	}

	usuarios_banco := lerProfessores_MongoDB(ctx, client, dbName)
	fmt.Println("\n")
	fmt.Println(usuarios_banco)

	defer conn.Close(context.Background())
	defer session.Close()

	// Queries
	//Mostrar nome e senha de todos os alunos (supa)
	//Mostrar todas as notas acima de 3 de um certo aluno (cassandra)
	//Mostrar o nome e se o professor tem alguma formação
}
