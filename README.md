# LunarAcademy
## Tema: Plataforma online de cursos
O tema escolhido foi o de uma plataforma de cursos online. Nós escolhemos este tema para dar continuidade ao tema que foi escolhido pelo grupo para o desenvolvimento do projeto de Engenharia de Software no semestre anterior.
## Bancos de dados
Neste projeto, utilizamos três bancos de dados diferentes: Supabase, MongoDB e Cassandra. Escolhemos o primeiro por causa de dois semestres de experiências prévias(3° e 5° semestre), consequentemente o grupo sente mais facilidade em prosseguir o projeto com o PostgreSQL, utilizando a plataforma do Supabase, pois foi a plataforma com experiência mais recente. O MongoDB foi escolhido pois ele é um banco não estruturado e os dados que guardaremos podem alterar sem deixar de ser da mesma tabela. Utilizaremos para guardar as informações dos professores, por exemplo, não é necessário ter uma formação ou certificados, mas serão salvas se tiver. Outro exemplos, os cursos podem ou não ter pré-requisitos para ingressão. E por último, escolhemos o Cassandra pois ele é bom para guardar histórico, retorna mais rápido e mantém consistência porque quando algo é atualizado, o resto não é afetado. Assim, escolhemos o Cassandra para guardar  histórico de cursos feitos pelos alunos da plataforma.
## Implementação do S2
Cada banco tem certas tabelas dentro dele:
### Supabase
O banco do Supabase guarda as seguintes tabelas: 
<img width="978" height="359" alt="TabelaSupaBase" src="https://github.com/user-attachments/assets/2825aae7-610d-4876-bdc9-2abdfd1435cc" />

#### curso
Onde são guardados todos os cursos cadastrados na plataforma. Nesta tabela, é guardado o nome do cursos, o autor e um ID para identificação.
#### usuario
Onde são guardados todos os usuários cadastrados na plataforma. Nesta tabela, são guardados o email,o CPF e a senha dos usuários. 
#### professor
Onde são guardadas todas as informações de todos os professores cadastrados na plataforma. É guardado apenas o nome e o CPF do professor.
#### certificado
Onde são guardadas as informações do certificado do aluno após o término do curso. Nela é guardado o id do certificado e o número de horas gastos no curso.


