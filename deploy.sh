#!/bin/bash

# Encerra o script imediatamente se qualquer comando falhar
set -e

# --- 1. Atualizar o Código-Fonte ---
echo "➡️  Puxando as últimas alterações do repositório Git..."
git pull origin main

# --- 2. Reconstruir e Reiniciar os Contêineres ---
echo "🚀  Construindo a nova imagem da API e reiniciando os serviços..."
docker compose up --build -d

# --- 3. Limpeza de Imagens Antigas ---
# Remove imagens Docker antigas e não utilizadas para economizar espaço
echo "🧹  Limpando imagens Docker antigas..."
docker image prune -f

# --- 4. Finalização ---
echo "✅  Deploy concluído com sucesso!"