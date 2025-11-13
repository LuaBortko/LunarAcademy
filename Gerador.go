//go get github.com/brianvoe/gofakeit/v6

package main

import (
	//"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"github.com/brianvoe/gofakeit/v6"
)

//=====================================Supabase=====================================

type Usuario struct {
	cpf   string
	email string
	senha string
}

type Professor struct {
	nome string
	cpf  string
}

type Aluno struct {
	nome string
	cpf  string
}

type Curso struct {
	nome      string
	cpf_autor string
	id        string
}

type Certificado struct {
	id    string
	horas string
}

type Aluno_curso struct {
	cpf_aluno      string
	id_curso       string
	data_in        string
	data_fim       string
	id_certificado string
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func valorExiste(valor string, lista interface{}, campo string) bool {
	v := reflect.ValueOf(lista) // retorna um objeto reflect, permitindo inspecionar o tipo, os campos e até o conteúdo de forma dinâmica

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i) // retorna o i elemento da slice, ex: um item do struct, mas ainda sendo um reflect.value
		// Acessar o campo
		c := item.FieldByName(campo) // campo que estou comparando (ex: cpf)
		// Comparar se é string
		if c.String() == valor {
			return true
		}
	}
	return false
} //Retorna verdadeiro se o valor ja existir no slice (vulgo lista)
//fmt.Println(valorExiste("123", usuarios, "cpf"))

func numExiste(num int, lista []int) bool {
	for i := 0; i < len(lista); i++ {
		if num == lista[i] {
			return true
		}
	}
	return false
}

func stringExiste(palavra string, lista []string) bool {
	for i := 0; i < len(lista); i++ {
		if palavra == lista[i] {
			return true
		}
	}
	return false
}

func gerarUsuarios(n int) []Usuario {
	usuarios := []Usuario{}
	for len(usuarios) < n {
		cpf := gofakeit.Regex("[0-9]{11}")
		email := gofakeit.Email()
		if valorExiste(cpf, usuarios, "cpf") || valorExiste(email, usuarios, "email") {
			continue // já existe → gera outro
		}
		u := Usuario{cpf: cpf, email: email, senha: gofakeit.Password(true, true, true, false, false, 6)}
		usuarios = append(usuarios, u)
	}
	return usuarios
}

func gerarAlunos(usuarios []Usuario, lista []int) []Aluno {
	alunos := []Aluno{}
	n_alunos := len(lista)
	for len(alunos) < n_alunos {
		nome := gofakeit.Name()
		if valorExiste(nome, alunos, "nome") {
			continue
		}
		index := lista[0]
		cpf := usuarios[index].cpf
		lista = lista[1:]
		a := Aluno{nome: nome, cpf: cpf}
		alunos = append(alunos, a)
	}
	return alunos
}

func gerarProfessores(usuarios []Usuario, lista []int) []Professor {
	professores := []Professor{}
	n_professores := len(lista)
	for len(professores) < n_professores {
		nome := gofakeit.Name()
		if valorExiste(nome, professores, "nome") {
			continue
		}
		index := lista[0]
		cpf := usuarios[index].cpf
		lista = lista[1:]
		p := Professor{nome: nome, cpf: cpf}
		professores = append(professores, p)
	}
	return professores
}

func gerarCursos(professores []Professor) []Curso {
	cursos := []Curso{}
	n_cursos := len(professores) + (len(professores) / 2)
	cpfs := []string{}
	for i := 0; i < len(professores); i++ {
		cpf := professores[i].cpf
		cpfs = append(cpfs, cpf)
	}
	lista := make([]int, len(cpfs)) // lista que vai guardar a quant de cursos de cada professor
	nomes := []string{}
	adjectives := []string{"Introdução à ", " Avançado"}
	for i := 0; i < n_cursos; i++ {
		id := "CU-" + strconv.Itoa(i+1)

		cpf := ""
		index := 0
		aux := 0
		for aux == 0 {
			if len(cpfs) != 1 {
				index = randInt(0, len(cpfs))
			}
			lista[index] += 1
			if lista[index] <= 2 {
				cpf = cpfs[index]
				aux = 1
			}
		}
		aux1 := 0
		nome := ""
		for aux1 == 0 {
			nome = gofakeit.ProgrammingLanguage()
			if stringExiste(nome, nomes) {
				continue
			}
			nomes = append(nomes, nome)
			aux1 = 1
		}
		n := randInt(0, 2)
		if n == 0 {
			nome = adjectives[0] + nome
		} else {
			nome = nome + adjectives[1]
		}

		c := Curso{nome: nome, cpf_autor: cpf, id: id}
		cursos = append(cursos, c)
	}
	return cursos

}

func gerarCertificado() Certificado {
	id := gofakeit.Regex("[0-9]{5}")
	id = "CE-" + id
	horas := strconv.Itoa(randInt(1, 50)) + ":00"

	return Certificado{id: id, horas: horas}
}

func gerarAlunos_curso(alunos []Aluno, cursos []Curso) []Aluno_curso {
	alunos_curso := []Aluno_curso{}
	for i := 0; i < len(alunos); i++ {
		cpf := alunos[i].cpf
		curso := ""
		id_cer := ""
		n_cursos := randInt(1, len(cursos)+1)
		cursados := []int{}
		for len(cursados) < n_cursos {
			index := randInt(0, len(cursos))
			if numExiste(index, cursados) {
				continue
			}
			cursados = append(cursados, index)

			curso = cursos[index].id

			formou := randInt(0, 2) // 0 -> nao terminou, 1 -> terminou o curso;
			dia := gofakeit.Day()
			mes := gofakeit.Month()
			ano := randInt(2020, 2026)
			data_in := strconv.Itoa(dia) + "/" + strconv.Itoa(mes) + "/" + strconv.Itoa(ano)
			data_fim := "--/--/--"
			if formou == 1 {
				id_cer = "a"
				aux_mes := 0
				aux_ano := 0
				dia_f := gofakeit.Day()
				if dia_f <= dia {
					aux_mes = 1
				}
				mes_f := gofakeit.Month() + aux_mes
				if mes_f > 12 {
					mes_f = 1
					aux_ano += 1
				}
				if mes_f < mes {
					aux_ano += 1
				}
				ano_f := ano + aux_ano
				data_fim = strconv.Itoa(dia_f) + "/" + strconv.Itoa(mes_f) + "/" + strconv.Itoa(ano_f)
			}
			a := Aluno_curso{cpf_aluno: cpf, id_curso: curso, data_in: data_in, data_fim: data_fim, id_certificado: id_cer}
			alunos_curso = append(alunos_curso, a)
			id_cer = ""
		}
	}
	return alunos_curso
}

//=====================================MongoDB=====================================

type Curso_Mongo struct {
	nome        string
	autor       string
	descricao   string
	requisitos  string
	preco       string
	dificuldade string
	avaliacao   string
}

type Professor_Mongodb struct {
	nome             string
	cpf              string
	formacao         string
	tempo_plataforma string
	qtde_cursos      string
}

func gerarCurso_Mongo(cursos []Curso, professores []Professor) []Curso_Mongo {
	cursos_mongo := []Curso_Mongo{}
	for i := 0; i < len(cursos); i++ {
		nome_autor := ""
		nome_curso := cursos[i].nome
		cpf_autor := cursos[i].cpf_autor
		for j := 0; j < len(professores); j++ {
			if cpf_autor == professores[j].cpf {
				nome_autor = professores[j].nome
			}
		}
		descricao := "Curso de " + nome_curso
		possui := randInt(0, 2)
		requisitos := ""
		if possui == 0 {
			requisitos = ""
		} else {
			requisitos = "Requisito minimo: O Basico"
		}
		preco := strconv.FormatFloat(gofakeit.Price(10.0, 100.0), 'f', -1, 64)
		dificuldade := ""
		nivel := randInt(1, 4)
		if nivel == 1 {
			dificuldade = "Iniciante"
		} else if nivel == 2 {
			dificuldade = "Intermediario"
		} else {
			dificuldade = "Avancado"
		}
		nota1 := randInt(0, 5)
		nota2 := randInt(0, 10)
		avaliacao := strconv.Itoa(nota1) + "." + strconv.Itoa(nota2)

		c := Curso_Mongo{nome: nome_curso, autor: nome_autor, descricao: descricao, requisitos: requisitos, preco: preco, dificuldade: dificuldade, avaliacao: avaliacao}
		cursos_mongo = append(cursos_mongo, c)
	}
	return cursos_mongo
}

func gerarProfessores_Mongodb(professores []Professor, cursos []Curso) []Professor_Mongodb {
	professores_mongodb := []Professor_Mongodb{}
	cursos_formacao := []string{"Ciência da Computação", "Ciência de dados", "Matemática", "Engenharia de Computação", "Sistemas de Informação", "Análise e Desenvolvimento de Sistemas"}
	for i := 0; i < len(professores); i++ {
		nome := professores[i].nome
		cpf := professores[i].cpf
		formado := randInt(0, 2)
		formacao := ""
		if formado == 1 {
			curso_formacao := gofakeit.RandomString(cursos_formacao)
			faculdade := gofakeit.School()
			formacao = "Formado na " + faculdade + " no curso de " + curso_formacao
		}
		anos := randInt(0, 5)
		meses := randInt(0, 13)
		tempo_plataforma := strconv.Itoa(anos) + " anos e " + strconv.Itoa(meses) + " meses"
		qtde_cursos := 0
		for j := 0; j < len(cursos); j++ {
			if cursos[j].cpf_autor == cpf {
				qtde_cursos = qtde_cursos + 1
			}
		}
		p := Professor_Mongodb{nome: nome, cpf: cpf, formacao: formacao, tempo_plataforma: tempo_plataforma, qtde_cursos: strconv.Itoa(qtde_cursos)}
		professores_mongodb = append(professores_mongodb, p)
	}
	return professores_mongodb
}

//=====================================Cassandra=====================================

type Historico struct {
	cpf_aluno   string
	id_curso    string
	avaliacao   string
	data_inicio string
	data_final  string
}

func gerarHistorico(alunos_curso []Aluno_curso) []Historico {
	historico := []Historico{}
	for i := 0; i < len(alunos_curso); i++ {
		cpf := alunos_curso[i].cpf_aluno
		id_curso := alunos_curso[i].id_curso
		
		nota1 := randInt(0, 5)
		nota2 := randInt(0, 10)
		avaliacao := strconv.Itoa(nota1) + "." + strconv.Itoa(nota2)
		data_inicio := alunos_curso[i].data_in
		data_final := alunos_curso[i].data_fim

		h := Historico{cpf_aluno: cpf, id_curso: id_curso, avaliacao: avaliacao, data_inicio: data_inicio, data_final: data_final}
		id_certificado := alunos_curso[i].id_certificado
		if id_certificado != "" {
			historico = append(historico, h)
		}
	}
	return historico
}

