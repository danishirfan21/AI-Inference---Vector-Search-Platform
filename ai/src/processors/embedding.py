import torch
from sentence_transformers import SentenceTransformer
from PIL import Image
import clip
import whisper

class EmbeddingProcessor:
    def __init__(self):
        self.text_model = SentenceTransformer('all-MiniLM-L6-v2')
        self.clip_model, self.clip_preprocess = clip.load("ViT-B/32", device="cpu")
        self.whisper_model = whisper.load_model("base")

    def get_text_embedding(self, text):
        return self.text_model.encode(text).tolist()

    def get_image_embedding(self, image_path):
        image = self.clip_preprocess(Image.open(image_path)).unsqueeze(0)
        with torch.no_grad():
            image_features = self.clip_model.encode_image(image)
        return image_features.flatten().tolist()

    def get_audio_embedding(self, audio_path):
        # Simplified: Use whisper to get hidden states or just features
        result = self.whisper_model.transcribe(audio_path)
        # In a real app, we'd use the audio features directly
        # Here we embed the transcription for simplicity
        return self.get_text_embedding(result['text'])
