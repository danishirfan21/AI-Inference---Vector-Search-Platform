from pymilvus import connections, Collection, FieldSchema, CollectionSchema, DataType

class MilvusClient:
    def __init__(self, host="milvus", port="19530"):
        connections.connect("default", host=host, port=port)
        self._init_collections()

    def _init_collections(self):
        # We'll use different collections for different media types due to different dimensions
        self._create_collection("text_embeddings", 384)
        self._create_collection("image_embeddings", 512)
        self._create_collection("audio_embeddings", 384)

    def _create_collection(self, name, dim):
        if name in [c.name for c in [Collection(n) for n in ["text_embeddings", "image_embeddings", "audio_embeddings"] if n in connections.get_connection_addr("default")]]: # Simple check
             pass # In real app, check if exists properly

        fields = [
            FieldSchema(name="id", dtype=DataType.VARCHAR, is_primary=True, auto_id=False, max_length=100),
            FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=dim),
        ]
        schema = CollectionSchema(fields, f"AI platform {name}")
        collection = Collection(name, schema)

        index_params = {
            "metric_type": "L2",
            "index_type": "IVF_FLAT",
            "params": {"nlist": 128}
        }
        collection.create_index("embedding", index_params)
        return collection

    def insert(self, collection_name, content_id, embedding):
        collection = Collection(collection_name)
        collection.insert([[content_id], [embedding]])

    def search(self, collection_name, embedding, limit=10):
        collection = Collection(collection_name)
        collection.load()
        search_params = {"metric_type": "L2", "params": {"nprobe": 10}}
        results = collection.search(
            data=[embedding],
            anns_field="embedding",
            param=search_params,
            limit=limit,
            output_fields=["id"]
        )
        return results
