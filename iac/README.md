# Observability Stack

Este repositório contém uma stack completa de ferramentas de observabilidade, configurada utilizando o Docker Compose (versão 3.9). Abaixo estão os serviços incluídos e suas respectivas funcionalidades.

## Serviços Disponíveis

### 1. **Jaeger**

**Imagem:** `jaegertracing/all-in-one:latest`

**Descrição:**
- Ferramenta de rastreamento distribuído para monitorar fluxos de requisições entre microserviços.
- Permite identificar gargalos e falhas na comunicação entre serviços.

**Portas Explicadas:**
- `5775/udp`: Para clientes que enviam spans via UDP.
- `5778`: Endpoint para configuração de cliente (HTTP).
- `16686`: Interface web para visualizar rastreamentos.
- `14250`: Endpoint gRPC para clientes OpenTelemetry.
- `14268`: Endpoint HTTP para clientes OpenTelemetry.
- `14269`: Endpoint para métricas.

### 2. **Prometheus**

**Imagem:** `prom/prometheus:latest`

**Descrição:**
- Ferramenta para coleta e consulta de métricas, amplamente utilizada para monitoramento de sistemas.
- Configurado para utilizar um arquivo customizado (`prometheus.yml`).

**Portas Explicadas:**
- `9090`: Interface web do Prometheus para consultas e monitoramento.

**Volumes Montados:**
- `./infra/docker/prometheus.yml:/etc/prometheus/prometheus.yml`: Arquivo de configuração do Prometheus.

### 3. **Alertmanager**

**Imagem:** `prom/alertmanager:latest`

**Descrição:**
- Sistema de gerenciamento de alertas do Prometheus, permitindo notificações via e-mail, Slack e outros canais.
- Configurado para utilizar um arquivo customizado (`alertmanager.yml`).

**Portas Explicadas:**
- `9093`: Interface web para gerenciamento de alertas.

**Volumes Montados:**
- `./infra/docker/alertmanager.yml:/etc/alertmanager/alertmanager.yml`: Arquivo de configuração do Alertmanager.

### 4. **Grafana**

**Imagem:** `grafana/grafana-oss:latest`

**Descrição:**
- Plataforma para visualização de métricas, logs e rastreamentos, com suporte para integrações como Prometheus, Loki e Jaeger.

**Portas Explicadas:**
- `3000`: Interface web do Grafana para criação de dashboards e visualização de dados.

**Ambiente Configurado:**
- Usuário administrador: `admin`.
- Senha: `admin`.

**Volumes Montados:**
- `grafana-data:/var/lib/grafana`: Diretório persistente para dados do Grafana.

### 5. **Loki**

**Imagem:** `grafana/loki:latest`

**Descrição:**
- Sistema de agregação e consulta de logs, projetado para trabalhar junto com o Grafana.
- Configurado para utilizar um arquivo customizado (`loki-config.yml`).

**Portas Explicadas:**
- `3100`: Endpoint de consulta de logs.

**Volumes Montados:**
- `./infra/docker/loki-config.yml:/etc/loki/local-config.yaml`: Arquivo de configuração do Loki.
- `loki-data:/loki-data`: A onde os dados são armazenados.

## Volumes Persistentes

Os seguintes volumes são utilizados para armazenar dados de forma persistente:
- `jaeger-data`: Dados do Jaeger.
- `prometheus-data`: Dados do Prometheus.
- `grafana-data`: Dados do Grafana.
-  `loki-data`: Dados do Loki.

## Como Utilizar

1. Certifique-se de ter o Docker e Docker Compose instalados em sua máquina.
2. Clone este repositório:
   ```bash
   git clone git@github.com:dosedetelemetria/projeto-otel-na-pratica.git
   cd iac
   ```
3. Suba os serviços com o comando:
   ```bash
   docker-compose up -d
   ```
4. Acesse as interfaces web dos serviços:
   - **Jaeger:** [http://localhost:16686](http://localhost:16686)
   - **Prometheus:** [http://localhost:9090](http://localhost:9090)
   - **Alertmanager:** [http://localhost:9093](http://localhost:9093)
   - **Grafana:** [http://localhost:3000](http://localhost:3000)

## Finalidade de Cada Serviço

- **Jaeger:** Rastreamento distribuído.
- **Prometheus:** Coleta de métricas.
- **Alertmanager:** Gerenciamento de alertas.
- **Grafana:** Visualização de métricas, logs e rastreamentos.
- **Loki:** Armazenamento e consulta de logs.

Com esta configuração, você terá uma stack de observabilidade completa para monitorar sistemas em desenvolvimento.

## Instalação e Uso do Kind (Kubernetes in Docker)

Para configurar um cluster Kubernetes local usando o Kind (Kubernetes in Docker), siga as instruções abaixo:

### Instalar o Kind
1. Certifique-se de ter o Docker instalado e em execução.
2. Baixe e instale o Kind executando o seguinte comando:
   ```bash
   curl -Lo ./kind https://kind.sigs.k8s.io/dl/latest/kind-linux-amd64
   chmod +x ./kind
   sudo mv ./kind /usr/local/bin/kind
   ```

### Criar um Cluster com o Kind
1. Certifique-se de que o arquivo `kind.yaml` está no diretório que será chamado
2. Execute o comando para criar um cluster Kubernetes com múltiplos nós:
   ```bash
   kind create cluster --name=multi-node-cluster --config=kind.yaml
   ```
3. Verifique se o cluster foi criado com sucesso:
   ```bash
   kubectl get nodes
   ```

