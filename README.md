# desafio-go



Para executar a aplicação, utilize:

```
docker-compose up --build
```

## Endpoints disponiveis:

### Criar planetas
#### URL
POST http://localhost:8080/planets

Parametros Obrigatórios:

- name 
- climate
- terrain

Status Possiveis:

- 201(Created) : Sucesso;
- 422(422 Unprocessable Entity): Erro de validação;

Exemplo:
```
curl --request POST \
  --url http://localhost:8080/planets \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "tatooine",
	"climate": "arid",
	"terrain": "teste"
}'
```

### Listar planetas
#### URL
GET http://localhost:8080/planets

Status Possiveis:
- 200(OK) : Sucesso;

Exemplo:
```
curl --request GET \
  --url 'http://127.0.0.1:8080/planets
```

É possivel também filtar por id:
```
curl --request GET \
  --url 'http://127.0.0.1:8080/planets?id=612d93ebe4e7e3aa6d9377a6
```

E por nome:

```
curl --request GET \
  --url 'http://127.0.0.1:8080/planets?nome=tatooine
```
### Remover planetas
#### URL
DELETE http://localhost:8080/planets/<id>

Status Possiveis:
- 204(No Content) : Sucesso;
- 404(Not Found) : Planeta não encontrado;  

Exemplo:
```
curl --request DELETE \
  --url http://localhost:8080/planets/612d93ebe4e7e3aa6d9377a6 \
  --header 'Content-Type: application/json'
```


