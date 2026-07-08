
# TitanDB 

# 🇵🇹 Versão em Português
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



---

# 🇺🇸 Versão em Inglês 


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



---


Aqui está o arquivo **README.md** contendo exclusivamente a seção do **Ruby Connector**, totalmente isolada em ambas as línguas.

---

# 💎 TitanDB - Ruby Connector

🇵🇹 Versão em Português

Esta seção documenta o **TitanConnector**, uma implementação nativa de cliente escrita em **Ruby** para o TitanDB. Este conector não utiliza requisições HTTP ou camadas REST; ele se comunica via sockets TCP puros e serializa as estruturas utilizando mapeamento binário direto (`Array#pack`), garantindo comunicação de ultra-baixa latência com o cluster.

## 🏗️ O Protocolo Binário (Wire Format)

O cliente Ruby empacota as mensagens seguindo estritamente o layout de memória do TitanDB:
`[Tamanho do Payload (4B)][ID da Mensagem (4B)][Tipo de Comando (1B)][Payload Data]`

Os códigos de comandos suportados são:

* `1` = PUT
* `2` = GET
* `3` = DELETE

## 🛠️ Implementação do Cliente (`titan_client.rb`)

Crie o arquivo abaixo para gerenciar o stream TCP bruto e a decodificação dos frames binários:

```ruby
require 'socket'

class TitanClient
  CMD_PUT, CMD_GET, CMD_DELETE = 1, 2, 3

  def initialize(host = '127.0.0.1', port = 8001)
    @host, @port = host, port
    @socket, @message_id = nil, 0
  end

  # Abre uma conexão persistente via socket TCP
  def connect
    @socket = TCPSocket.new(@host, @port)
  rescue Errno::ECONNREFUSED
    raise "Não foi possível conectar ao TitanDB em #{@host}:#{@port}"
  end

  # Operação PUT: Envia chave e valor formatados
  def put(key, value)
    payload = [key.bytesize, key, value.bytesize, value].pack("S>a*S>a*")
    send_frame(CMD_PUT, payload)
    read_frame
  end

  # Operação GET: Busca o valor de uma chave
  def get(key)
    send_frame(CMD_GET, key)
    read_frame
  end

  # Fecha o socket de forma segura
  def close
    @socket&.close
    @socket = nil
  end

  private

  # Monta o cabeçalho binário e escreve os bytes no socket
  def send_frame(cmd, payload)
    connect if @socket.nil? || @socket.closed?
    @message_id += 1
    
    # L> = 4 Bytes (Unsigned Long), C = 1 Byte. > garante Big-Endian (Network Byte Order)
    header = [payload.bytesize, @message_id, cmd].pack("L>L>C")
    @socket.write(header + payload)
  rescue Errno::EPIPE, Errno::ECONNRESET
    close
    raise "Conexão perdida com o TitanDB."
  end

  # Lê a resposta do servidor desempacotando o cabeçalho de 9 bytes fixos
  def read_frame
    header_raw = @socket.read(9)
    return nil if header_raw.nil? || header_raw.bytesize < 9
    
    payload_size, msg_id, status = header_raw.unpack("L>L>C")
    payload_data = payload_size > 0 ? @socket.read(payload_size) : ""
    
    { message_id: msg_id, status: status, data: payload_data }
  end
end

```

## 🚀 Script de Execução (`main.rb`)

Script prático de automação para testar o fluxo de escrita e leitura apontando para o nó principal (`8001`):

```ruby
require_relative 'titan_client'

puts "=== Testando Conector Ruby para TitanDB ==="
client = TitanClient.new('127.0.0.1', 8001)

begin
  # 1. Escrita (PUT)
  puts "\n[Ruby -> PUT]: Gravando chave 'infra:engineer'"
  res_put = client.put("infra:engineer", "elite")
  puts "Resposta Cluster -> MsgID: #{res_put[:message_id]}, Status: #{res_put[:status]}"

  # 2. Leitura (GET)
  puts "\n[Ruby -> GET]: Buscando chave 'infra:engineer'"
  res_get = client.get("infra:engineer")
  puts "Resposta Cluster -> Dado Recuperado: '#{res_get[:data]}'"

rescue StandardError => e
  puts "🚨 Erro: #{e.message}"
ensure
  client.close
  puts "\nConexão finalizada."
end

```

Para rodar o cliente:

```bash
ruby main.rb

```

---

## 🇺🇸 English Version

# 💎 TitanDB - Ruby Connector

This section documents the **TitanConnector**, a native client implementation written in **Ruby** for TitanDB. This connector bypasses traditional HTTP abstractions or REST overhead; it opens raw TCP streams directly to the cluster and handles binary message parsing through low-level memory layout directives (`Array#pack`), ensuring ultra-low latency.

## 🏗️ Wire Format (Binary Protocol)

The Ruby client builds frames strictly matching TitanDB’s core specification:
`[Payload Size (4B)][Message ID (4B)][Command Type (1B)][Payload Data]`

Supported command codes:

* `1` = PUT
* `2` = GET
* `3` = DELETE

## 🛠️ Client Implementation (`titan_client.rb`)

Create the following file to handle the raw TCP socket life cycle and binary frame encoding/decoding:

```ruby
require 'socket'

class TitanClient
  CMD_PUT, CMD_GET, CMD_DELETE = 1, 2, 3

  def initialize(host = '127.0.0.1', port = 8001)
    @host, @port = host, port
    @socket, @message_id = nil, 0
  end

  # Opens a persistent connection via raw TCP socket
  def connect
    @socket = TCPSocket.new(@host, @port)
  rescue Errno::ECONNREFUSED
    raise "Could not connect to TitanDB at #{@host}:#{@port}"
  end

  # PUT Operation: Sends packed key-value pairs
  def put(key, value)
    payload = [key.bytesize, key, value.bytesize, value].pack("S>a*S>a*")
    send_frame(CMD_PUT, payload)
    read_frame
  end

  # GET Operation: Fetches value from a specific key
  def get(key)
    send_frame(CMD_GET, key)
    read_frame
  end

  # Gracefully closes the active socket
  def close
    @socket&.close
    @socket = nil
  end

  private

  # Builds the binary header and writes bytes directly to the TCP stream
  def send_frame(cmd, payload)
    connect if @socket.nil? || @socket.closed?
    @message_id += 1
    
    # L> = 4 Bytes (Unsigned Long), C = 1 Byte. > enforces Big-Endian (Network Byte Order)
    header = [payload.bytesize, @message_id, cmd].pack("L>L>C")
    @socket.write(header + payload)
  rescue Errno::EPIPE, Errno::ECONNRESET
    close
    raise "Lost connection to TitanDB cluster."
  end

  # Reads the cluster response by unpacking the fixed 9-byte header first
  def read_frame
    header_raw = @socket.read(9)
    return nil if header_raw.nil? || header_raw.bytesize < 9
    
    payload_size, msg_id, status = header_raw.unpack("L>L>C")
    payload_data = payload_size > 0 ? @socket.read(payload_size) : ""
    
    { message_id: msg_id, status: status, data: payload_data }
  end
end

```

## 🚀 Execution Script (`main.rb`)

A lightweight automation script to validate mutation and query workflows against the primary cluster node (`8001`):

```ruby
require_relative 'titan_client'

puts "=== Initiating Ruby Connector for TitanDB ==="
client = TitanClient.new('127.0.0.1', 8001)

begin
  # 1. Mutation (PUT)
  puts "\n[Ruby -> PUT]: Storing key 'infra:engineer'"
  res_put = client.put("infra:engineer", "elite")
  puts "Cluster Response -> MsgID: #{res_put[:message_id]}, Status: #{res_put[:status]}"

  # 2. Query (GET)
  puts "\n[Ruby -> GET]: Fetching key 'infra:engineer'"
  res_get = client.get("infra:engineer")
  puts "Cluster Response -> Retrieved Data: '#{res_get[:data]}'"

rescue StandardError => e
  puts "🚨 Error occurred: #{e.message}"
ensure
  client.close
  puts "\nConnection closed."
end

```

To run the Ruby client:

```bash
ruby main.rb

```

---

### Integração e Ecossistema
Embora o TitanDB tenha sido construído como um banco de dados chave-valor de alta performance usando Rust, eu o projetei para ser facilmente integrado a pilhas de tecnologia modernas.

Para demonstrar sua interoperabilidade com ecossistemas focados em web, desenvolvi o [Ruby-Orchestrator](https://github.com/kauandias747474-hue/Ruby-Orchestrator). Escolhi **Ruby** para esse orquestrador devido à sua produtividade excepcional na criação de camadas de serviço de alto nível e scripts de automação complexos.

**Por que Ruby?**
Integrar o TitanDB com um orquestrador em Ruby permite a prototipagem rápida de fluxos de dados complexos. Essa estrutura aproveita a sintaxe expressiva do Ruby para gerenciar o ciclo de vida dos nós de dados e orquestrar requisições ao banco de dados, criando uma ponte entre a performance de armazenamento de baixo nível e a lógica de aplicação de alto nível.

---
### Integration & Ecosystem
While TitanDB is built as a high-performance key-value store using Rust, I've designed it to be easily integrated into modern application stacks.

To demonstrate its interoperability with web-focused ecosystems, I developed the [Ruby-Orchestrator](https://github.com/kauandias747474-hue/Ruby-Orchestrator). I chose **Ruby** for this orchestrator because of its exceptional productivity in building high-level service layers and complex automation scripts.

**Why Ruby?**
Integrating TitanDB with a Ruby-based orchestrator allows for rapid prototyping of complex data workflows. This setup leverages Ruby's expressive syntax to manage the lifecycle of data nodes and orchestrate requests to the database, bridging the gap between low-level storage performance and high-level application logic.

