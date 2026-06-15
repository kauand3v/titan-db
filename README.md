Aqui estão os dois arquivos `README.md` completos, estruturados com padrão de Engenharia Core para você copiar, colar e impressionar qualquer recrutador ou tech lead.

---

# 🇵🇹 Versão em Português (README.pt.md)

```markdown
# TitanDB 🚀

**TitanDB** é um banco de dados distribuído chave-valor (Key-Value Store) de alta performance, construído inteiramente do zero. O projeto foi desenhado para operar em clusters, utilizando um protocolo RPC proprietário focado em baixa latência e o algoritmo de consenso Raft para garantir a consistência e tolerância a falhas dos dados.

Este projeto não utiliza frameworks HTTP tradicionais ou soluções de mensageria prontas. Ele resolve os problemas complexos de redes na camada de transporte (TCP Sockets), gerencia concorrência extrema na memória e sobrevive a falhas parciais do cluster (como partições de rede e quedas de nós).

---

## 🏗️ Arquitetura do Sistema

O TitanDB é dividido em três camadas desacopladas através de interfaces rígidas:

```text
       +------------------+
       |   Client Script  |
       +------------------+
                 |
                 | (TitanRPC - TCP Binary)
                 v
       +------------------+
       |   Nó Líder (8001)|
       +------------------+
         /              \
        /                \  (TitanRPC - Heartbeats & Logs)
       v                  v
+------------------+  +------------------+
| Nó Seguidor 8002 |  | Nó Seguidor 8003 |
+------------------+  +------------------+

```

### Componentes Core

1. **TitanRPC (Camada de Rede):** Um framework RPC customizado rodando sobre TCP puro. Implementa multiplexação de mensagens, gerenciamento de conexões persistentes (*Connection Pooling*) e codificação binária eficiente (*Framing*) para eliminar o overhead de cabeçalhos HTTP/JSON.
2. **Raft Core (Camada de Consenso):** Implementação robusta do algoritmo de consenso Raft. Gerencia o ciclo de vida dos nós através de uma Máquina de Estados: *Leader*, *Follower* e *Candidate*. Responsável por eleições automáticas de líderes baseadas em *timeouts* aleatórios e replicação atômica de logs.
3. **KV Engine (Camada de Aplicação):** Um motor de armazenamento em memória ultraveloz baseado em estruturas de dados thread-safe. Ele é completamente isolado da rede; mudanças no estado só ocorrem quando o Raft confirma que o log foi comitado pela maioria do cluster.

---

## ✨ Funcionalidades Mecânicas

* **Consenso Forte:** Garantia de consistência linearizável (Linearizability) utilizando Raft.
* **Eleição Automática de Líder:** Se o nó líder falhar, os seguidores detectam a ausência de *heartbeats* e elegem um novo líder em milissegundos.
* **Protocolo Binário Customizado:** Mensagens de rede compactas estruturadas em `[Tamanho do Payload (4B)][ID da Mensagem (4B)][Tipo de Comando (1B)][Payload Data]`.
* **Tolerância a Falhas Dinâmica:** O cluster mantém a integridade contanto que a maioria absoluta dos nós ($N/2 + 1$) esteja ativa.

---

## 🛠️ Tecnologias Utilizadas

* **Linguagem Principal:** Go (ou Rust / C++ dependendo da sua escolha)
* **Rede:** Sockets TCP Nativos
* **Concorrência:** Mutexes, Canais e Primitivas de Sincronização Assíncrona

---

## 🚀 Como Executar

### Pré-requisitos

* Ter o runtime da linguagem instalado (ex: Go 1.22+ / Rust)

### 1. Inicializando o Cluster (3 Nós)

Abra três terminais diferentes para simular as três instâncias do cluster localmente:

**Terminal 1 (Nó 1):**

```bash
go run cmd/server/main.go --id=1 --address=":8001" --cluster="1@:8001,2@:8002,3@:8003"

```

**Terminal 2 (Nó 2):**

```bash
go run cmd/server/main.go --id=2 --address=":8002" --cluster="1@:8001,2@:8002,3@:8003"

```

**Terminal 3 (Nó 3):**

```bash
go run cmd/server/main.go --id=3 --address=":8003" --cluster="1@:8001,2@:8002,3@:8003"

```

### 2. Interagindo com o Cliente

Com o cluster estabilizado e o líder eleito, execute o script de cliente para enviar comandos de leitura e escrita:

```bash
# Gravar um dado
go run cmd/client/main.go --server=":8001" --cmd=PUT --key="infra:engineer" --val="elite"

# Ler um dado
go run cmd/client/main.go --server=":8001" --cmd=GET --key="infra:engineer"

```

---

## 🧪 Engenharia de Caos: Validando a Resiliência

Para provar que o TitanDB é um sistema de infraestrutura core real, execute os seguintes testes de estresse:

### Cenário 1: Queda do Líder (Failover)

1. Identifique qual nó assumiu o estado de `LEADER` nos logs do terminal.
2. Mate o processo desse nó (`Ctrl+C`).
3. Observe os logs dos dois nós restantes: eles entrarão em estado de `CANDIDATE`, solicitarão votos através do TitanRPC e um deles se tornará o novo líder automaticamente.
4. Faça uma consulta `GET` no novo líder e comprove que os dados inseridos anteriormente continuam lá.

### Cenário 2: Partição de Rede (Split-Brain)

Se você simular o isolamento de um nó impedindo a comunicação dele com os outros, o cluster remanescente continuará operando normalmente (pois mantém a maioria), enquanto o nó isolado recusará escritas para evitar a corrupção do estado do sistema, respeitando estritamente o **Teorema CAP** (Privilegiando Consistência sobre Disponibilidade).

```

---

# 🇺🇸 Versão em Inglês (README.md)

```markdown
# TitanDB 🚀

**TitanDB** is a high-performance distributed key-value store built entirely from scratch. The system is engineered to operate across multiple cluster nodes, leveraging a proprietary low-latency RPC framework and the Raft consensus algorithm to guarantee strong data consistency and high availability.

This project bypasses traditional HTTP frameworks and off-the-shelf messaging solutions. It solves core networking challenges at the transport layer (TCP Sockets), manages high-concurrency memory states, and survives arbitrary partial cluster failures (e.g., network partitions and node crashes).

---

## 🏗️ System Architecture

TitanDB is decoupled into three strict layers enforced by rigid interfaces:

```text
       +------------------+
       |   Client Script  |
       +------------------+
                 |
                 | (TitanRPC - TCP Binary)
                 v
       +------------------+
       |  Leader Node 8001|
       +------------------+
         /              \
        /                \  (TitanRPC - Heartbeats & Logs)
       v                  v
+------------------+  +------------------+
|Follower Node 8002|  |Follower Node 8003|
+------------------+  +------------------+

```

### Core Components

1. **TitanRPC (Networking Layer):** A custom RPC framework running on top of raw TCP sockets. It implements message multiplexing, connection pooling, and an efficient binary framing protocol to eliminate HTTP header and JSON encoding overhead.
2. **Raft Core (Consensus Layer):** A robust implementation of the Raft consensus algorithm. It manages node lifecycles via a strict Finite State Machine (FSM): *Leader*, *Follower*, and *Candidate*. It handles automated leader elections driven by randomized timers and atomic log replication.
3. **KV Engine (Application Layer):** An ultra-fast in-memory storage engine powered by thread-safe data structures. The engine is completely isolated from the network; state transitions only occur when the Raft layer signals that a log entry has been committed by a cluster majority.

---

## ✨ Mechanical Features

* **Strong Consistency:** Guarantees linearizable reads and writes using the Raft consensus protocol.
* **Automated Leader Election:** If the active leader goes offline, followers detect the missing heartbeats and elect a new leader within milliseconds.
* **Custom Binary Protocol:** Ultra-lean wire format structured as `[Payload Size (4B)][Message ID (4B)][Command Type (1B)][Payload Data]`.
* **Dynamic Fault Tolerance:** The cluster remains fully operational and consistent as long as a strict majority ($N/2 + 1$) of nodes are alive.

---

## 🛠️ Tech Stack

* **Core Language:** Go (or Rust / C++ depending on your choice)
* **Networking:** Native TCP Sockets
* **Concurrency:** Mutexes, Channels, and Async Synchronization Primitives

---

## 🚀 Getting Started

### Prerequisites

* Language runtime installed (e.g., Go 1.22+ / Rust)

### 1. Bootstrapping the Cluster (3 Nodes)

Open three distinct terminal windows to simulate separate infrastructure nodes locally:

**Terminal 1 (Node 1):**

```bash
go run cmd/server/main.go --id=1 --address=":8001" --cluster="1@:8001,2@:8002,3@:8003"

```

**Terminal 2 (Node 2):**

```bash
go run cmd/server/main.go --id=2 --address=":8002" --cluster="1@:8001,2@:8002,3@:8003"

```

**Terminal 3 (Node 3):**

```bash
go run cmd/server/main.go --id=3 --address=":8003" --cluster="1@:8001,2@:8002,3@:8003"

```

### 2. Interacting via Client

Once the cluster stabilizes and a leader is elected, execute the client script to run read/write workloads:

```bash
# Write a key-value pair
go run cmd/client/main.go --server=":8001" --cmd=PUT --key="infra:engineer" --val="elite"

# Read a key-value pair
go run cmd/client/main.go --server=":8001" --cmd=GET --key="infra:engineer"

```

---

## 🧪 Chaos Engineering: Validating Resilience

To demonstrate that TitanDB is an elite-tier core infrastructure system, execute the following stress-testing scenarios:

### Scenario 1: Leader Failover

1. Inspect the terminal logs to identify which node stepped up as the `LEADER`.
2. Kill that specific node process abruptly (`Ctrl+C`).
3. Monitor the surviving nodes' logs: they will swiftly transition to the `CANDIDATE` state, request votes via TitanRPC, and automatically elect a new leader.
4. Fire a `GET` request to the new leader and verify that data consistency is intact.

### Scenario 2: Network Partition (Split-Brain)

Simulating a network split that isolates a minority node will show that the majority partition keeps serving requests smoothly, while the isolated node rejects writes to prevent state drift—strictly adhering to the **CAP Theorem** by prioritizing Consistency over Availability.

```

```
