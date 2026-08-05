# desafio-korp
Este repositório contém a resolução do desafio técnico para provisionamento de uma infraestrutura baseada em containers, monitoramento e automação. O ambiente foi inteiramente projetado para rodar em um servidor Proxmox utilizando redes isoladas e automação via Ansible.

## Arquitetura do Ambiente (Proxmox)

A infraestrutura foi montada no Proxmox com isolamento de rede, garantindo segurança e resolução de nomes local. O diagrama abaixo ilustra o fluxo da arquitetura:

```mermaid
graph TD
    Internet((Internet)) -->|WAN| pfSense[VM pfSense - Gateway/Firewall]
    
    subgraph Rede_Isolada [Rede Isolada do Projeto]
        pfSense -->|LAN| SwitchVirtual[Switch Virtual / Bridge]
        
        SwitchVirtual -->|IP: <IP_DA_VM_UBUNTU>| VM_Ubuntu[VM Ubuntu 24.04<br>Host Docker & Ansible Target]
    end

    subgraph VM_Ubuntu_Docker [Containers na VM Ubuntu]
        Nginx[NGINX Proxy Reverso<br>Porta 80]
        GoApp[Serviço Golang<br>Porta 8080]
        Prometheus[Prometheus<br>Porta 9090]
        Grafana[Grafana<br>Porta 3000]

        Nginx -->|Encaminha tráfego| GoApp
        Prometheus -->|Coleta métricas via /metrics| GoApp
        Grafana -->|Lê dados e gera Dashboard| Prometheus
    end
```
