import asyncio
import json
import os
import psycopg2
from opentelemetry import trace
from nats.aio.client import Client as NATS
from processors.embedding import EmbeddingProcessor
from service.milvus_client import MilvusClient

async def main():
    # Configuration
    nats_url = os.getenv("NATS_URL", "nats://localhost:4222")
    milvus_addr = os.getenv("MILVUS_ADDRESS", "localhost:19530")
    db_url = os.getenv("DATABASE_URL", "postgres://user:password@localhost:5432/ai_platform?sslmode=disable")

    # Initialize Services
    nc = NATS()
    await nc.connect(nats_url)

    processor = EmbeddingProcessor()
    milvus = MilvusClient(host=milvus_addr.split(':')[0], port=milvus_addr.split(':')[1])
    tracer = trace.get_tracer(__name__)

    # Database connection
    conn = psycopg2.connect(db_url)
    cursor = conn.cursor()

    async def content_uploaded_handler(msg):
        with tracer.start_as_current_span("process_content"):
            try:
                data = json.loads(msg.data.decode())
                content_id = data['content_id']
                file_type = data['file_type']
                s3_path = data['s3_path']

                print(f"Processing {file_type} content: {content_id}")

                # Update status to 'processing'
                cursor.execute("UPDATE content_metadata SET processing_status = 'processing' WHERE id = %s", (content_id,))
                conn.commit()

                # Embedding logic
                collection_name = "text_embeddings"
                embedding = []
                if file_type == 'document':
                    # In real app, download and read s3_path
                    embedding = processor.get_text_embedding("Sample text from " + s3_path)
                    collection_name = "text_embeddings"
                elif file_type == 'image':
                    # embedding = processor.get_image_embedding(local_path)
                    embedding = [0.1] * 512 # Placeholder
                    collection_name = "image_embeddings"
                elif file_type == 'audio':
                    # embedding = processor.get_audio_embedding(local_path)
                    embedding = [0.1] * 384 # Placeholder
                    collection_name = "audio_embeddings"

                # Store in Milvus
                milvus.insert(collection_name, content_id, embedding)

                # Update status to 'completed'
                cursor.execute("UPDATE content_metadata SET processing_status = 'completed' WHERE id = %s", (content_id,))
                conn.commit()

                print(f"Successfully processed {content_id}")
            except Exception as e:
                print(f"Error processing {content_id}: {e}")
                cursor.execute("UPDATE content_metadata SET processing_status = 'failed' WHERE id = %s", (content_id,))
                conn.commit()

    async def get_embedding_handler(msg):
        query = msg.data.decode()
        embedding = processor.get_text_embedding(query)
        await nc.publish(msg.reply, json.dumps(embedding).encode())

    await nc.subscribe("ContentUploaded", cb=content_uploaded_handler)
    await nc.subscribe("GetEmbedding", cb=get_embedding_handler)

    print("AI Service is running and listening for events...")
    while True:
        await asyncio.sleep(1)

if __name__ == '__main__':
    asyncio.run(main())
