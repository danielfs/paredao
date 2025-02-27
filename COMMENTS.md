# Desafio técnico globo.com
## Projeto Paredão

### Objetivo
Resolver o problema de votação do BBB usando a linguagem Go.

### Etapas Realizadas
#### 1. Desenho do Projeto
Para esclarecer os requisitos do projeto e obter uma visão geral, o primeiro passo foi criar um desenho no draw.io.

[Acesse aqui](https://drive.google.com/file/d/1A8ZBGVKZbtWmE02lHDirFH7fRujAYFGb/view?usp=sharing)

- Utilizei uma Arquitetura em Camadas para separar a apresentação da lógica de negócio.
- Contém uma camada frontend para interação com os usuários.
- Contém uma camada backend para registro das informações e agregações de dados.
- A camada backend conta com dois componentes: um para votação e outro para estatísticas.
- O componente de estatísticas possui uma camada de cache, para evitar consultar os dados no banco com muita frequência.

#### 2. Escolhas das Tecnologias
- **Linguagem:** Go
- **Banco de Dados:** MySQL
- **Cache:** Redis

Porque são amplamente adotadas pela comunidade de desenvolvimento de software.

## Documentação do Sistema

### Visão Geral
O Paredão é uma aplicação full-stack para gerenciar participantes, sessões de votação e votos. Consiste em uma API backend em Go e um frontend em JavaScript.

### Estrutura do Projeto
- `backend/` - API backend em Go
- `frontend/` - Frontend em JavaScript
- `docker-compose.yml` - Configuração do Docker Compose
- `load-tests/` - Testes de carga para a API

### Executando a Aplicação
A maneira mais fácil de executar a aplicação é usando o Docker Compose:

```bash
docker-compose up -d
```

Isso iniciará os seguintes serviços:
- Frontend: http://localhost:3000
- API Backend: http://localhost:8080
- Administrador de Banco de Dados (Adminer): http://localhost:8081

### Funcionalidades

#### Interface de Administração
- Gerenciar participantes (criar, ler, atualizar, excluir)
- Gerenciar sessões de votação (criar, ler, atualizar, excluir)
- Adicionar participantes às sessões de votação

#### Interface de Votação do Usuário
- Selecionar uma sessão de votação
- Visualizar participantes na sessão de votação selecionada
- Votar em um participante

### API Backend

#### Endpoints da API

##### Participantes
- **GET /participantes** - Listar todos os participantes
- **GET /participantes/{id}** - Obter um participante específico por ID
- **POST /participantes** - Criar um novo participante
- **PUT /participantes/{id}** - Atualizar um participante
- **DELETE /participantes/{id}** - Excluir um participante

##### Votações
- **GET /votacoes** - Listar todas as sessões de votação
- **GET /votacoes/{id}** - Obter uma sessão de votação específica por ID
- **POST /votacoes** - Criar uma nova sessão de votação
- **PUT /votacoes/{id}** - Atualizar uma sessão de votação
- **DELETE /votacoes/{id}** - Excluir uma sessão de votação
- **GET /votacoes/{id}/participantes** - Obter todos os participantes de uma sessão de votação específica
- **POST /votacoes/{id}/participantes** - Adicionar um participante a uma sessão de votação

##### Votos
- **GET /votos** - Listar todos os votos
- **GET /votos/{participanteId}/{votacaoId}** - Obter um voto específico
- **POST /votos** - Criar um novo voto
- **PUT /votos/{participanteId}/{votacaoId}** - Atualizar um voto (redefine o timestamp)
- **DELETE /votos/{participanteId}/{votacaoId}** - Excluir um voto

##### Estatísticas
- **GET /estatisticas/votacoes/{id}/total** - Obter o número total de votos para uma sessão de votação
- **GET /estatisticas/votacoes/{id}/participantes** - Obter o número total de votos por participante para uma sessão de votação
- **GET /estatisticas/votacoes/{id}/hourly** - Obter o número total de votos por hora para uma sessão de votação

#### Modelos de Dados

##### Participante
```go
type Participante struct {
    Id      int64
    Nome    string
    UrlFoto string
}
```

##### Votacao
```go
type Votacao struct {
    Id        int64
    Descricao string
}
```

##### Voto
```go
type Voto struct {
    Participante *Participante
    Votacao      *Votacao
    DataHora     time.Time
}
```

### Frontend

#### Estrutura de Arquivos
- `index.html` - Interface de votação do usuário
- `admin.html` - Interface de administração
- `css/styles.css` - Estilos para ambas as interfaces
- `js/voting.js` - JavaScript para a interface de votação do usuário
- `js/admin.js` - JavaScript para a interface de administração

### Testes de Carga

#### Executando o Teste de Carga
O programa Go fornece uma maneira programática de executar testes de carga:

1. Compile o programa Go:
   ```
   cd load-tests
   go build -o load_test vegeta_runner.go
   ```
2. Execute o teste com opções:
   ```
   ./load_test --rate 2000 --duration 30s --endpoint votos
   ```

Opções disponíveis:
- `--rate`: Requisições por segundo (padrão: 2000)
- `--duration`: Duração do teste (padrão: 30s)
- `--endpoint`: Endpoint da API para testar (padrão: votos)
- `--output`: Diretório de saída (padrão: ./results)
- `--threshold`: Mínimo aceitável de requisições por segundo (padrão: 2000)

#### Interpretando Resultados
O teste é considerado bem-sucedido se:
1. A API conseguir lidar com pelo menos 2000 requisições por segundo
2. Todas as requisições forem bem-sucedidas (código de status HTTP 2xx)
