# get_weather_services_OTEL
# Get Weather Services OTEL

## 📌 Descrição
Este projeto consiste em dois serviços escritos em Go que trabalham juntos para receber um CEP, identificar a cidade correspondente e retornar as temperaturas em Celsius, Fahrenheit e Kelvin. A arquitetura inclui OpenTelemetry (OTEL) e Zipkin para rastreamento distribuído.

- **Service A**: Recebe o CEP via POST e encaminha para o Service B.
- **Service B**: Consulta a API ViaCEP para obter a cidade e depois a WeatherAPI para buscar a temperatura.
- **Zipkin**: Utilizado para rastreamento distribuído das requisições.

---

## 🛠 Tecnologias Utilizadas
- **Go** 1.20+
- **Docker & Docker Compose**
- **OpenTelemetry** (OTEL)
- **Zipkin**
- **API ViaCEP**
- **API WeatherAPI**

---

## 🚀 Como Rodar o Projeto em Ambiente de Desenvolvimento

### 1️⃣ **Pré-requisitos**
Certifique-se de ter instalado:
- **Docker** e **Docker Compose**
- **Go** (caso queira rodar os serviços sem Docker)

### 2️⃣ **Clonar o repositório**
```bash
 git clone https://github.com/SEU_USUARIO/SEU_REPOSITORIO.git
 cd get_weather_services_OTEL
```

### 3️⃣ **Configurar variáveis de ambiente**
Crie um arquivo `.env` na raiz do projeto e adicione sua chave da WeatherAPI:
```bash
WEATHER_API_KEY=INSIRA_SUA_CHAVE_AQUI
```

---

## 🐳 **Rodando com Docker Compose**
### 🔹 **Passo 1: Build das imagens**
```bash
docker compose build
```

### 🔹 **Passo 2: Subir os serviços**
```bash
docker compose up -d
```
Isso iniciará os seguintes serviços:
- `service-a`: disponível em [http://localhost:8010](http://localhost:8010)
- `service-b`: disponível em [http://localhost:8091](http://localhost:8091)
- `zipkin`: disponível em [http://localhost:9411](http://localhost:9411) (para visualização dos spans)

Para ver os logs:
```bash
docker logs -f service-a
```
```bash
docker logs -f service-b
```

Para parar os serviços:
```bash
docker compose down
```

---

## 🛠 **Rodando sem Docker**
Caso queira rodar os serviços manualmente sem Docker, siga os passos:

### 🔹 **Rodar o Service B**
```bash
cd service_B
export WEATHER_API_KEY=INSIRA_SUA_CHAVE_AQUI
go run main.go
```

### 🔹 **Rodar o Service A**
Em outra aba do terminal:
```bash
cd service_A
go run main.go
```

Agora, os serviços estarão rodando localmente nas portas 8010 e 8091.

---

## 📡 **Testando a API**

### 🔹 **Enviar um CEP para o Service A**
```bash
curl -X POST http://localhost:8010/service-A \
     -H "Content-Type: application/json" \
     -d '{"cep": "29902555"}'
```
Resposta esperada:
```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

Se algo não funcionar, verifique os logs:
```bash
docker logs service-a --tail=50
```
```bash
docker logs service-b --tail=50
```

---

## 🔍 **Monitoramento com Zipkin**
Para visualizar os traces, acesse o Zipkin em:
[http://localhost:9411](http://localhost:9411)

Pesquise por `service-A` ou `service-B` para ver o fluxo de requisições entre os serviços.

---

## ✅ **Checklist de Solução de Problemas**

### ❌ **Erro de conexão entre os serviços?**
- Certifique-se de que `service-A` está chamando `service-B` com `http://service-b:8091`, e **não** `localhost:8091`.
- Confirme se os serviços estão rodando:
  ```bash
  docker ps
  ```

### ❌ **Erro ao conectar ao Zipkin?**
- Troque `http://localhost:9411/api/v2/spans` por `http://zipkin:9411/api/v2/spans` nos serviços.
- Verifique os logs:
  ```bash
  docker logs zipkin --tail=50
  ```

### ❌ **Resposta vazia do Service B?**
- Verifique se a API **WeatherAPI** está acessível:
  ```bash
  curl "http://api.weatherapi.com/v1/current.json?key=SEU_API_KEY&q=São Paulo&aqi=no"
  ```
- Veja os logs do `service-b`.

---

## 📄 **Licença**
Este projeto é open-source e está sob a licença MIT.

---

Feito com ❤️ por [Seu Nome] 🚀

