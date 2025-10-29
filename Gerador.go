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
}

func gerarUsuarios(n int) []Usuario{
	usuarios := []Usuario{}
	for i := 0; i < n; i++{
		u := Usuario{cpf: gofakeit.Regex("[0-9]{11}"),
		email: gofakeit.Email(),
		senha: gofakeit.Password(true, true, true, true, false, 6)}
		usuarios = append(usuarios, u)
	}
	return usuarios
}

func main(){
	usuarios := gerarUsuarios(3)
	fmt.Println(usuarios)
}