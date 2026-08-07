
## Visão Geral

Repositório contendo a infraestrutura, aplicação e automação de deploy do Projeto Korp. A solução engloba uma API em Golang, monitoramento com Prometheus e Grafana, roteamento via Nginx, e uma esteira de automação baseada em Ansible.

O ambiente foi provisionado de forma isolada em Proxmox, utilizando o sistema operacional Ubuntu Server 24.04, com separação estrita entre a máquina de controle (`10.0.10.11`) e a máquina alvo (`10.0.10.10`).

---

## Atendimento aos Requisitos do Desafio

### 1. Serviço HTTP e Arquitetura Docker
* **Serviço Golang:** `http-server-projeto-korp` desenvolvido e configurado na porta interna `8080`.
* **Endpoint:** Rota `GET /projeto-korp` funcional, retornando JSON com as chaves `nome` ("Projeto Korp") e `horario` em formato UTC dinâmico.
* **Dockerfile:** Implementado em *multi-stage build* (base `golang:alpine` para compilação estática sem dependências CGO e imagem final otimizada baseada em `alpine`).
* **Docker Compose & Rede:** Orquestrado via rede bridge. A API não expõe portas diretamente ao host externo.
* **Proxy Reverso (Nginx):** Imagem oficial do Nginx mapeando a porta `80` do host para a `80` do container. O volume `/etc/nginx/conf.d/` injeta o arquivo de configuração para direcionar o tráfego de `localhost:80` para o serviço backend.

### 2. Monitoramento e Observabilidade
* **Métricas Expostas:** A aplicação utiliza a biblioteca oficial `client_golang` do Prometheus para exportar o contador de volume de requisições (`http_requests_total`) e o medidor de disponibilidade (`service_available` configurado como *Gauge* com estado binário).
* **Coleta (Scraping):** O Prometheus executa a coleta (*scrape*) direcionada ao target interno a cada 5 segundos.
* **Bônus (Dashboard as Code):** Provisionamento automático do Grafana estruturado declarativamente por meio de mapeamento de volumes para os diretórios `/etc/nginx/conf.d/`, `datasources.yml`, `dashboards.yml` e o arquivo JSON correspondente, eliminando intervenção manual na interface gráfica.

### 3. Automação (Ansible)
* **Playbooks Modulares:** A automação divide-se em:
  * `presetup.yaml`: Bootstrap de segurança, criação do usuário de automação (`sladmin`) e políticas de *Passwordless Sudo*.
  * `docker.yaml`: Instalação e configuração do Docker Engine e Docker Compose Plugin.
  * `playbook.yaml`: Pipeline principal que valida o ambiente local, sincroniza arquivos via `rsync`, executa o *build* e *up* dos containers, realiza validação de saúde (*healthcheck*) via HTTP e aplica rotina de rollback automático em caso de falha.

#### Algumas observações
* Foram mantidas as versões antigas em uma pasta "teste", mostrando que o projeto foi realizado passo a passo conforme as instruções.
---

## Estrutura do Repositório
```text
.
├── ansible
│   └── playbooks
│       ├── docker.yaml          # Instalação do Docker e Docker Compose
│       ├── hosts.ini            # Inventário de hosts (korp-stack)
│       ├── host_vars
│       │   └── korp-stack
│       │       └── main.yml     # Variáveis do host alvo
│       ├── playbook.yaml        # Playbook principal de deploy, healthcheck e rollback
│       └── presetup.yaml        # Playbook de bootstrap e configuração de sudoers
├── app
│   ├── DockerFile               # Build multi-stage da aplicação Go
│   ├── go.mod                   # Gerenciamento de módulos
│   ├── go.sum                   # Checksums de dependências
│   └── main.go                  # Código-fonte da API e registro de métricas
├── config
│   ├── grafana
│   │   └── provisioning
│   │       ├── dashboards       # Dashboards JSON e configuração yml
│   │       └── datasources      # Configuração do datasource Prometheus yml
│   ├── nginx
│   │   └── http-server-projeto-korp.conf  # Virtual host do proxy reverso
│   └── prometheus
│       └── prometheus.yml       # Configuração de scraping targets
├── docker-compose.yaml          # Orquestração de serviços em containers
└── README.md
```

## Guia de Implantação

### A execução ocorre inteiramente a partir do nó de controle (10.0.10.11).

* **Pode-se fazer todo o setup com o seguinte 'oneliner':** `git pull && ansible-playbook ansible/playbooks/presetup.yaml -i ansible/playbooks/hosts.ini --ask-become-pass && ansible-playbook ansible/playbooks/playbook.yaml -i ansible/playbooks/hosts.ini`

* **1. Setup Inicial (Bootstrap)**

 Executado uma única vez na máquina alvo para configurar o ambiente de privilégios:
`ansible-playbook ansible/playbooks/presetup.yaml -i ansible/playbooks/hosts.ini --ask-become-pass`

* **2. Deploy Completo via Ansible**

Executa a esteira automatizada (instalação de dependências, sincronização de arquivos, subida dos serviços, testes de saúde e salvamento de backup):
`ansible-playbook ansible/playbooks/playbook.yaml -i ansible/playbooks/hosts.ini`

<img width="1922" height="1053" alt="automacao-completa" src="https://github.com/user-attachments/assets/a89a0067-8c77-4ed1-bbf7-24ef6a4a8ad8" />

* **3. Validação do Serviço**

No nó alvo, verifique o retorno da API através do Nginx:
`curl http://localhost/projeto-korp`

Saída esperada:
```
{"horario":"2026-08-07T03:00:00Z","nome":"Projeto Korp"}
```

## Validação de Métricas e Monitoramento

O endpoint /metrics exposto pelo serviço Go pode ser inspecionado diretamente para validação das métricas:

```
curl http://localhost:8080/metrics
```

<img width="1922" height="1053" alt="endpoint-metrics" src="https://github.com/user-attachments/assets/8a259ad7-2dc0-4800-8610-eb55e952cda7" />
<img width="1922" height="1053" alt="win10-vm-dashboards" src="https://github.com/user-attachments/assets/eb56e99f-0a68-40c0-b65a-eb7fbc95e3bf" />"""
