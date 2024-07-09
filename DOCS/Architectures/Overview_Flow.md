# Overview Flow Diagram

```mermaid
flowchart LR
    subgraph Client
        A[HTTP Client]
    end

    subgraph "Gin Web Server"
        B[Router]
    end

    subgraph "Core Services"
        direction TB
        C[Token Service]
        D[Cloud Recording Service]
        E[Real-Time Transcription Service]
    end

    subgraph "Handlers"
        direction TB
        F[Token Handlers]
        G[Cloud Recording Handlers]
        H[RTT Handlers]
    end

    subgraph "Shared Components"
        J[Storage Config]
    end

    subgraph "Middleware"
        I[Middleware]
    end

    subgraph "External"
        K[Agora RESTful API]
    end

    A <-->|Request/Response| B
    B <-->|/token| C
    B <-->|/cloud_recording| D
    B <-->|/rtt| E
    C <--> F
    D <--> G
    E <--> H
    D & E -.->|Uses| J
    C & D & E -.->|Uses| I
    G & H <-.->|API Calls| K

    classDef request fill:#f9f,stroke:#333,stroke-width:2px;
    classDef response fill:#bbf,stroke:#333,stroke-width:2px;
```
