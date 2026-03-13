# AI Inference and Vector Search Platform

A production-ready platform for uploading documents, images, and audio, generating semantic embeddings, and performing high-performance vector searches.

## Architecture

The system consists of several microservices and infrastructure components:

- **API Service (Golang)**: Handles file uploads, manages metadata in Postgres, publishes events to NATS, and provides search endpoints.
- **AI Processing Service (Python)**: Consumes events, processes media using state-of-the-art models, generates embeddings, and manages Milvus storage.
- **Milvus**: High-performance vector database for storing and searching embeddings.
- **Postgres**: Relational database for storing file metadata and processing status.
- **NATS**: Lightweight messaging system for event-driven communication between services.
- **Redis**: Caching layer for frequently accessed search results and metadata.
- **OpenTelemetry**: Provides observability across the entire pipeline.

### AI Pipeline

1. **Upload**: User uploads a file (Doc, Image, Audio) via the Go API.
2. **Metadata**: API stores initial metadata in Postgres.
3. **Event**: API publishes `ContentUploaded` event to NATS.
4. **Inference**: Python AI service consumes the event, downloads the file, and runs:
    - `sentence-transformers` for documents.
    - `CLIP` for images.
    - `Whisper` for audio.
5. **Storage**: Generated embeddings are stored in Milvus with a reference to the Postgres ID.
6. **Update**: AI service updates the processing status in Postgres.

### Search Flow

1. **Query**: User submits a search query (text or media).
2. **Embedding**: API requests an embedding for the query from the AI service.
3. **Vector Search**: API queries Milvus using the generated embedding.
4. **Retrieval**: Top matches are retrieved from Milvus, and corresponding metadata is fetched from Postgres/Redis.
5. **Response**: System returns the relevant content and metadata to the user.

## Directory Structure

```
.
├── api/                # Golang API Service
│   ├── cmd/            # Entry points
│   └── internal/       # Core logic
├── ai/                 # Python AI Processing Service
│   └── src/            # Embedding models and NATS consumer
├── deploy/             # Infrastructure and Deployment
│   ├── docker-compose.yml
│   └── k8s/            # Kubernetes manifests
└── scripts/            # Helper scripts for testing
```

## Deployment

### Prerequisites
- Docker & Docker Compose
- Kubernetes (optional)

### Local Development
```bash
docker-compose up -d
```

### Kubernetes
```bash
kubectl apply -f deploy/k8s/
```
