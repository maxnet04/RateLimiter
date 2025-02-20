# Rate Limiter em Go 

## Objetivo 

Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

## Descrição

Este projeto implementa um rate limiter que controla o tráfego de requisições para um serviço web. O rate limiter pode limitar o número de requisições com base em dois critérios: 

**Endereço IP**: Restringe o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.

**Token de Acesso**: Limita as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O token deve ser informado no header no seguinte formato: •API_KEY: <TOKEN>. As configurações de limite do token de acesso se sobrepõem às do IP. As configurações de limite do token de acesso se sobrepõem às do IP.



## Requisitos 

O rate limiter funciona como um middleware injetado no servidor web. 

- Permite a configuração do número máximo de requisições permitidas por segundo. 
- Oferece a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições seja excedida. 
- As configurações de limite são realizadas via variáveis de ambiente ou em um arquivo env na pasta raiz. Configurável tanto para limitação por IP quanto por token de acesso. 
- Responde adequadamente quando o limite é excedido com o código HTTP 429 e a mensagem: "you have reached the maximum number of requests or actions allowed within a certain time frame". 
- Todas as informações de "limiter" são armazenadas e consultadas de um banco de dados Redis.

 Implementa uma estratégia que permite trocar facilmente o Redis por outro mecanismo de persistência. A lógica do limiter está separada do middleware. 
 
 ## Exemplos 
 
 **Limitação por IP**: Se configurado para permitir no máximo 5 requisições por segundo por IP, e o IP `192.168.11` enviar 6 requisições em um segundo, a sexta requisição será bloqueada. 
 
 **Limitação por Token**: Se um token `FCYCLE` tiver um limite de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira será bloqueada. 
 
 ## Configuração 
 
 ### Arquivo `config.env`
 
 Crie um arquivo `config.env` na raiz do projeto com as seguintes configurações: 
 
 
```
API_KEY=FCYCLE
RATE_LIMIT_IP=8
RATE_LIMIT_TOKEN=10
BLOCK_DURATION_SECONDS=60
```
 
#### Requisitos:

- [GO](https://golang.org/doc/insttall) 1.17 ou superior
- [Docker](https://docs.docker.com/get-docker/)



### Como Utilizar localmente:

  1. Clonar o Repositório:~
  ```git clone https://github.com/maxnet04/RateLimiter.git```


  2. Acesse a pasta do app:
  ```cd ReteLimiter```

  3. Rode o docker para buildar a imagem gerando o container com a aplicação e o redis:
   ```docker-compose up```

Porta: HTTP server on port :8080

#### Execute o curl abaixo ou use um aplicação client REST para realizar a requisição:

    curl --request GET \
    --url http://localhost:8080/ \
    --header 'API_KEY: FCYCLE'

OBs:  O desafio especifica que token tem prioridade em relação a IP sendo assim solicitações enviadas com um API_KEY diferente do configurado no arquivo config.env o rate limit a ser considerado sera o por IP