# ML Models Directory

This directory is mounted into the ML server container for local development.

## Detector Model

The AI detector model (RoBERTa) is automatically downloaded from HuggingFace on first run.
No manual setup required.

## Feedback Model

Download the TinyLlama GGUF model for feedback generation:

```bash
mkdir -p feedback
wget -O feedback/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf \
  "https://huggingface.co/TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF/resolve/main/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf"
```

## Production

In production, models are loaded from S3 bucket (`truegul-ml-models`).
