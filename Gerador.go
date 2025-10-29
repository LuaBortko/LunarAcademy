//go get github.com/brianvoe/gofakeit/v6

package main

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"reflect"
)

type Usuario struct {
	cpf string
	email string
	senha string 
} 

type Professor struct {
	nome string
	cpf string
} 

type Aluno struct {
	nome string
	cpf string
} 

type Curso struct {
	nome string
	cpf_autor string
	id string 
} 

type Certificado struct {
	id string
	horas string
	} 

type Aluno_curso struct {
	cpf_aluno string
	id_curso string
	data_in string
	data_fim string 
	id_certificado string
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
}//Retorna verdadeiro se o valor ja existir no slice (vulgo lista)
//fmt.Println(valorExiste("123", usuarios, "cpf"))

func gerarUsuarios(n int) []Usuario{
	usuarios := []Usuario{}
	for len(usuarios) < n {
		cpf := gofakeit.Regex("[0-9]{11}")
		email := gofakeit.Email()
		if valorExiste(cpf, usuarios, "cpf") || valorExiste(email, usuarios, "email"){
			continue // já existe → gera outro
		}
		u := Usuario{cpf: cpf, email: email, senha: gofakeit.Password(true, true, true, false, false, 6)}
		usuarios = append(usuarios, u)
	}
	return usuarios
}

func alunos(lista []string) []Alunos{
	alunos := []Alunos{}
	n_alunos := len(lista)/2
	for len(alunos) < n_alunos{
		nome := gofakeit.Name()
		if valorExiste(nome, alunos, "nome"){
			continue
		}
		
		a:= Aluno{nome: nome, cpf: cpf}
		alunos = append(alunos, a)
	}
}

func professores(lista []string) []Professor{
	professores := []Professor{}
	n_professores := len(lista)
	for len(professores) < n_professores{
		nome := gofakeit.Name()
		if valorExiste(nome, professores, "nome"){
			continue
		}
		p:= Aluno{nome: nome, cpf: cpf}
		professores = append(professores, p)
	}
}

func main(){
	//usuarios := gerarUsuarios(10)
	//fmt.Println(usuarios)
}