# Documento de Visão

## Perfis dos Usuários

| Perfil        | Descrição                                     |
|---------------|-----------------------------------------------|
| Visitante     | Este usuário pode apenas visualizar snippets. |
| Usuário       | Pode criar snippets.                          |
| Administrador | Tem total controle sobre usuários e snippets. |

## Lista de Requisitos Funcionais

| Requisito                                 | Descrição                                                                             | Ator                  |
|-------------------------------------------|---------------------------------------------------------------------------------------|-----------------------|
| RF001 - Manter snippets                   | Uma snippet tem id, título, conteúdo, tempo de expiração e data de criação.           | Usuário               |
| RF002 - Manter usuários                   | Um usuário contem nome, email, senha e data de criação.                               | Usuário               |
| RF003 - Snippets adicionados recentemente | A aplicação deve ter uma página que lista as 10 últimas snippets inseridas            | -                     |
| RF004 - Formulário snippets               | A aplicação deve ter um formulário para criação de novas snippets pelos usuários.     | Usuário               |
| RF005 - Login                             | A aplicação deve ter um painel de login onde usuários e administradores possam logar. | Usuário/Administrador |

## Lista de Requisitos Não-Funcionais

| Requisito                   | Descrição                                                                                                                                  |
|-----------------------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| RNF001 - Logs e stack trace | O sistema deve logar todas as requisições e caso aconteça um erro interno no servidor, deve ser logado o erro, como também um stack-trace. |
| RNF002 - Izi deploy         | Todo o deploy da aplicação deve ser feito em apenas um comando utilizando docker e docker-compose.                                         |
| RNF002 - 70% testes         | O produto final deve apresentar 70% de cobertura em testes. |

## Riscos

Tabela com o mapeamento dos riscos do projeto, as possíveis soluções e os responsáveis.

| Data       | Risco                  | Prioridade | Responsável | Status | Providência/Solução                                      |
|------------|------------------------|------------|-------------|--------|----------------------------------------------------------|
| 21/08/2022 | Não terminar o projeto | Alta       | Eu          | Eu     | Deixar de ser preguiçoso e fazer oque tenho que fazer :) |

### Referências
