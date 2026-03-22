#!/bin/bash

MODEL_URL="https://huggingface.co/keisuke-miyako/bce-embedding-base_v1-gguf-q4_k_m/resolve/main/bce-embedding-base_v1-Q4_k_m.gguf"
MODEL_NAME="bce-embedding-base_v1-Q4_k_m.gguf"

if [ -f "$MODEL_NAME" ]; then
    echo "Model '$MODEL_NAME' already exists. Skipping download."
    exit 0
fi

echo "Downloading model from Hugging Face..."

if command -v curl >/dev/null 2>&1; then
    curl -L -o "$MODEL_NAME" "$MODEL_URL"
elif command -v wget >/dev/null 2>&1; then
    wget -O "$MODEL_NAME" "$MODEL_URL"
else
    echo "Error: Neither curl nor wget found. Please download the model manually."
    echo "URL: $MODEL_URL"
    exit 1
fi

if [ $? -eq 0 ]; then
    echo "Download successful!"
else
    echo "Download failed. Please check your internet connection."
    exit 1
fi