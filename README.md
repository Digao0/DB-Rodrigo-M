# Estrutura do Projeto

## Objetivo Inicial

Implementar um banco de dados chave-valor simples utilizando:

* armazenamento em memória (`map`)
* persistência em disco
* append-only log (WAL)
* recuperação após reinicialização

Comandos suportados inicialmente:

```text
SET chave valor
GET chave
DEL chave
```

---

## cmd/db/main.go

### Responsabilidade

Ponto de entrada da aplicação.

### Deve fazer

* Inicializar o banco.
* Carregar dados persistidos.
* Iniciar a interface de linha de comando.
* Receber comandos do usuário.
* Chamar os métodos do banco.

### Não deve fazer

* Manipular arquivos diretamente.
* Implementar lógica de persistência.
* Manipular estruturas internas do banco.

---

## internal/storage/db.go

### Responsabilidade

Implementar a lógica principal do banco.

### Deve conter

Estrutura principal:

```go
type DB struct {
    mem *MemTable
    wal *WAL
}
```

Operações:

```go
Set(key, value)
Get(key)
Delete(key)
```

### Fluxo esperado

SET:

```text
Recebe operação
↓
Atualiza memória
↓
Escreve no WAL
↓
fsync
```

GET:

```text
Consulta memória
↓
Retorna valor
```

DELETE:

```text
Remove da memória
↓
Registra no WAL
↓
fsync
```

---

## internal/storage/mem.go

### Responsabilidade

Gerenciar o estado em memória.

### Deve conter

Estrutura:

```go
type MemTable struct {
    data map[string]string
}
```

Operações:

```go
Set()
Get()
Delete()
```

### Observação

Inicialmente será apenas um map.

No futuro poderá ser substituído por:

* Skip List
* B+Tree
* LSM Tree

sem alterar a lógica principal do banco.

---

## internal/storage/wal.go (criar futuramente)

### Responsabilidade

Persistência em disco.

### Deve conter

Estrutura:

```go
type WAL struct {
    file *os.File
}
```

Operações:

```go
Append(record)
Load()
Sync()
```

### Exemplo de conteúdo do arquivo

```text
SET nome Rodrigo
SET idade 22
DEL nome
```

### Objetivo

Garantir durabilidade e recuperação após falhas.

---

## internal/types/record.go

### Responsabilidade

Definir como uma operação é representada.

### Deve conter

Tipos de operação:

```go
SET
DELETE
```

Estrutura:

```go
type Record struct {
    Operation string
    Key       string
    Value     string
}
```

### Objetivo

Padronizar o formato usado pelo banco e pelo WAL.

---

## tests/

### Responsabilidade

Garantir que o comportamento do banco continue correto.

### Testes iniciais

#### Set/Get

```text
SET nome Rodrigo
GET nome

Resultado:
Rodrigo
```

#### Delete

```text
SET nome Rodrigo
DEL nome
GET nome

Resultado:
não encontrado
```

#### Recovery

```text
SET nome Rodrigo
Fechar banco

Abrir banco novamente

GET nome

Resultado:
Rodrigo
```

---

# Próximo Marco

Quando esta etapa estiver concluída, o banco deverá:

1. Abrir normalmente.
2. Aceitar SET, GET e DEL.
3. Salvar todas as operações em disco.
4. Recuperar os dados após reiniciar.
5. Utilizar append-only log para persistência.

Nesse ponto o projeto já será um banco de dados persistente básico.
