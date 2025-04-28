#!/bin/bash

echo "开始设置 TalkSphere 数据生成服务..."


# 2. 创建并激活新的 conda 环境
echo "创建新的 conda 环境..."
conda create -n talksphere-gen python=3.10 -y
eval "$(conda shell.bash hook)"
conda activate talksphere-gen

# 3. 验证环境
echo "验证环境..."
PYTHON_PATH=$(which python)
echo "当前 Python 路径: $PYTHON_PATH"

# 4. 安装依赖
echo "安装依赖..."
# 4.1 安装 FastAPI 相关
echo "安装 FastAPI 相关依赖..."
conda install -c conda-forge fastapi uvicorn starlette pydantic -y

# 4.2 安装 PyTorch（CPU 版本）
echo "安装 PyTorch..."
conda install pytorch torchvision torchaudio cpuonly -c pytorch -y

# 4.3 安装 transformers 和其他依赖
echo "安装 transformers 和其他依赖..."
conda install -c conda-forge transformers accelerate faker requests sympy -y

# 5. 验证安装
echo "验证安装..."
python -c "import fastapi; import uvicorn; import starlette; import transformers; import torch; import accelerate; import sympy; print('所有依赖安装成功！')"

if [ $? -ne 0 ]; then
    echo "依赖安装验证失败，请检查错误信息"
    exit 1
fi

# 6. 下载 Qwen 模型
echo "下载 Qwen 模型..."
python -c "from transformers import AutoModelForCausalLM, AutoTokenizer; AutoModelForCausalLM.from_pretrained('Qwen/Qwen1.5-0.5B-Chat', trust_remote_code=True); AutoTokenizer.from_pretrained('Qwen/Qwen1.5-0.5B-Chat', trust_remote_code=True)"

if [ $? -ne 0 ]; then
    echo "模型下载失败，请检查错误信息"
    exit 1
fi