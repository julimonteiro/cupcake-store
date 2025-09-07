# 🧁 Cupcake Store

Uma aplicação web para gerenciamento de cupcakes gourmet, desenvolvida em Go com arquitetura MVC e API REST.

## 📋 Características

- **Backend**: Go 1.24.3 com Chi Router
- **Banco de Dados**: SQLite (desenvolvimento) e PostgreSQL (produção)
- **ORM**: GORM
- **Frontend**: HTML + JavaScript puro
- **Arquitetura**: MVC com separação clara de responsabilidades
- **Testes**: Unitários com testify
- **Containerização**: Docker e Docker Compose

## 🏗️ Arquitetura

```
cupcake-store/
├── cmd/                    # Ponto de entrada da aplicação
│   └── main.go
├── internal/               # Código interno da aplicação
│   ├── config/            # Configurações
│   ├── database/          # Conexão com banco de dados
│   ├── handler/           # Handlers HTTP
│   ├── models/            # Modelos de dados
│   ├── repository/        # Camada de acesso a dados
│   ├── router/            # Configuração de rotas
│   └── service/           # Lógica de negócio
├── web/                   # Frontend
│   └── index.html
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

## 🚀 Quick Start

### Pré-requisitos

- Go 1.24.3 ou superior
- Docker e Docker Compose (opcional)

### Execução Local

1. **Clone o repositório**
   ```bash
   git clone <repository-url>
   cd cupcake-store
   ```

2. **Instale as dependências**
   ```bash
   make deps
   ```

3. **Configure as variáveis de ambiente**
   ```bash
   cp env.example .env
   # Edite o arquivo .env conforme necessário
   ```

4. **Execute a aplicação**
   ```bash
   make run
   ```

5. **Acesse a aplicação**
   - Frontend: http://localhost:8080
   - API Health Check: http://localhost:8080/health

### Execução com Docker

1. **Inicie os containers**
   ```bash
   make docker-up
   ```

2. **Acesse a aplicação**
   - Frontend: http://localhost:8080
   - API Health Check: http://localhost:8080/health

3. **Para os containers**
   ```bash
   make docker-down
   ```

## 📚 API Endpoints

### Health Check
- `GET /health` - Verifica o status da aplicação

### Cupcakes
- `GET /api/v1/cupcakes` - Lista todos os cupcakes
- `POST /api/v1/cupcakes` - Cria um novo cupcake
- `GET /api/v1/cupcakes/{id}` - Obtém um cupcake específico
- `PUT /api/v1/cupcakes/{id}` - Atualiza um cupcake
- `DELETE /api/v1/cupcakes/{id}` - Remove um cupcake

### Exemplo de Requisição POST
```json
{
  "name": "Chocolate Especial",
  "flavor": "Chocolate Belga",
  "price_cents": 1500
}
```

### Exemplo de Resposta
```json
{
  "id": 1,
  "name": "Chocolate Especial",
  "flavor": "Chocolate Belga",
  "price_cents": 1500,
  "is_available": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## 🗄️ Modelo de Dados

### Cupcake
- `id` (uint, auto increment) - Identificador único
- `name` (string, obrigatório, min 2 chars) - Nome do cupcake
- `flavor` (string, obrigatório) - Sabor do cupcake
- `price_cents` (int, obrigatório > 0) - Preço em centavos
- `is_available` (bool, default true) - Status de disponibilidade
- `created_at` (timestamp) - Data de criação
- `updated_at` (timestamp) - Data de atualização

## 🧪 Testes

### Executar todos os testes
```bash
make test
```

### Executar testes com cobertura
```bash
make test-coverage
```

### Executar testes específicos
```bash
go test -v ./internal/service
```

## 🛠️ Comandos Make

```bash
make help          # Mostra todos os comandos disponíveis
make run           # Executa a aplicação localmente
make build         # Compila a aplicação
make test          # Executa os testes
make clean         # Remove arquivos temporários
make docker-up     # Inicia containers com Docker Compose
make docker-down   # Para e remove containers
make docker-build  # Constrói imagem Docker
make deps          # Instala dependências
make fmt           # Formata o código
make lint          # Executa linter
make check         # Executa todos os checks (fmt, lint, test)
```

## 🔧 Configuração

### Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `PORT` | Porta do servidor | `8080` |
| `DB_DIALECT` | Tipo de banco (`sqlite` ou `postgres`) | `sqlite` |
| `DB_DSN` | String de conexão com banco | `cupcake_store.db` |
| `LOG_LEVEL` | Nível de log | `info` |

### Exemplo de .env
```env
PORT=8080
DB_DIALECT=sqlite
DB_DSN=cupcake_store.db
LOG_LEVEL=info
```

## 🐳 Docker

### Construir imagem
```bash
make docker-build
```

### Executar container
```bash
make docker-run
```

### Logs dos containers
```bash
make logs          # Todos os logs
make logs-app      # Apenas aplicação
make logs-db       # Apenas banco de dados
```

## 📦 Estrutura para Expansão

O projeto está estruturado para facilitar futuras expansões:

- **Pedidos**: Adicionar `internal/models/order.go`, `internal/service/order_service.go`, etc.
- **Entregas**: Adicionar `internal/models/delivery.go`, `internal/service/delivery_service.go`, etc.
- **Pagamentos**: Adicionar `internal/models/payment.go`, `internal/service/payment_service.go`, etc.
- **Usuários**: Adicionar autenticação e autorização
- **Relatórios**: Adicionar endpoints para relatórios de vendas

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

## 👨‍💻 Autor

Desenvolvido como continuação do PIT I - Loja de Cupcakes Gourmet.

---

**Cupcake Store** - Gerenciando cupcakes com sabor e qualidade! 🧁✨

