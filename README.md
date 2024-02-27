## API para Buscar Clima de Cidades a Partir do CEP
Este projeto é uma API para buscar informações climáticas de cidades com base no CEP. A seguir, você encontrará instruções sobre como executar e testar o projeto.

### Pré-requisitos
Certifique-se de ter o [Docker](https://www.docker.com/) instalado em sua máquina.

**Obs**: A ferramenta zipkin fornece a funcionalidade de exportação de `tracing` em formato json, caso você tenha algum problema em usar o projeto por favor me encaminhe o/os trancings para que eu possa analisar.

### Executando os Serviços 
Para iniciar os serviços, execute o seguinte comando no terminal:

```sh
    docker compose up -d
```

**Observação**: Os projetos possuem variáveis de ambiente definidas no arquivo `docker-compose.yaml`. Caso deseje modificá-las, você pode fazer isso antes de executar o comando `docker compose up -d`.

### Testando a API
Utilize o arquivo `api/request.http` em conjunto com a extensão [Rest Client](https://github.com/Huachao/vscode-restclient/tree/master) no VSCode para realizar testes.

Exemplo de request:

```sh
curl -X POST http://localhost:8000/weather -H "Content-Type: application/json" -d '{"cep": "80010000"}'
```
