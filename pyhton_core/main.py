from flask import Flask, request, jsonify
from sentence_transformers import SentenceTransformer

app = Flask(__name__)
model = SentenceTransformer("intfloat/multilingual-e5-small")


@app.get("/health")
def health():
    return jsonify({"status": "ok"})


@app.post("/embed-batch")
def embed_batch():
    data = request.get_json()
    texts = data.get("texts", [])
    if not texts:
        return jsonify({"embeddings": []})
    embeddings = model.encode(texts, normalize_embeddings=True).tolist()
    return jsonify({"embeddings": embeddings})


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
