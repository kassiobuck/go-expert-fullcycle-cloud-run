# Google cloud run

## Desafio:
Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual 
(temperatura em graus celsius, fahrenheit e kelvin). Esse sistema deverá ser publicado no Google Cloud Run.

## Pré requisitos
- [Go](https://golang.org/doc/install)
- [Docker](https://www.docker.com/get-started)
- [WeatherApi KEY](https://www.weatherapi.com)


## Executando o projeto em modo de desenvolvimento

Abrir o arquivo .env na raiz do projeto e preencher a constante WEATHER_API_KEY com uma WeatherApi KEY válida;
Executar o comando:
```shell
    go run main.go
``` 

## Testando API localmente

Abra o arquivo api.http, nele é possivel acessar a API e executar os testes já realizados. Para novos testes basta substituir o valor de "cep".

## Testando API no Google cloud run

Abra o endereço e insira um cep válido após `cep=`

https://cloud-run-i5qc44ni2q-uc.a.run.app/?cep=
