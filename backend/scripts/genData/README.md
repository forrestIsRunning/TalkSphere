```bash
# 1. 清理环境（如果之前有创建过）
conda deactivate
deactivate  # 如果还有 venv 环境
rm -rf venv  # 删除旧的虚拟环境

# 2. 检查当前环境
conda env list  # 查看所有可用的环境

# 3. 创建并激活新的 conda 环境
conda create -n talksphere-gen python=3.10
conda activate talksphere-gen

# 4. 确认环境是否正确激活
# 终端提示符应该显示 (talksphere-gen)
# 可以通过以下命令确认 Python 路径
which python  # 应该显示 miniconda3/envs/talksphere-gen/bin/python

# 5. 安装依赖（分步安装以确保所有依赖都正确安装）
# 5.1 安装 FastAPI 相关
conda install -c conda-forge fastapi uvicorn starlette pydantic

# 5.2 安装 PyTorch（使用 CPU 版本）
conda install pytorch torchvision torchaudio cpuonly -c pytorch

# 5.3 安装 transformers 和 accelerate
conda install -c conda-forge transformers accelerate

# 5.4 安装其他依赖
conda install -c conda-forge faker requests

# 6. 验证安装
python -c "import fastapi; import uvicorn; import starlette; import transformers; import torch; import accelerate; print('所有依赖安装成功！')"

# 7. 下载 Qwen 模型
python -c "from transformers import AutoModelForCausalLM, AutoTokenizer; AutoModelForCausalLM.from_pretrained('Qwen/Qwen1.5-0.5B-Chat', trust_remote_code=True); AutoTokenizer.from_pretrained('Qwen/Qwen1.5-0.5B-Chat', trust_remote_code=True)"
```

## 运行步骤

1. 确保 conda 环境已激活：
```bash
conda activate talksphere-gen
# 确认环境已激活
which python  # 应该显示 miniconda3/envs/talksphere-gen/bin/python
```

2. 启动 Qwen 模型服务：
```bash
python qwen_service.py
```

3. 在另一个终端中运行数据生成脚本：
```bash
python generate-mock-data.py
```

注意：
- Qwen 模型服务默认运行在 http://localhost:8000
- 数据生成脚本默认连接到 http://localhost:8989 的 TalkSphere 后端
- 如果后端地址不同，请修改脚本中的 BASE_URL 变量
- 确保在运行数据生成脚本前设置正确的 ADMIN_TOKEN

## 常见问题解决

1. 如果遇到 "ModuleNotFoundError"：
   - 确保 conda 环境已激活：`conda activate talksphere-gen`
   - 检查 Python 路径：`which python` 应该显示 miniconda3/envs/talksphere-gen/bin/python
   - 检查是否同时激活了其他环境（不应该同时看到 (venv) 和 (talksphere-gen)）
   - 重新安装依赖：`conda install -c conda-forge <package_name>`

2. 如果遇到模型下载问题：
   - 检查网络连接
   - 尝试使用代理
   - 确保有足够的磁盘空间

3. 如果遇到 CUDA 相关错误：
   - 确保安装了正确版本的 PyTorch：`conda install pytorch torchvision torchaudio pytorch-cuda=12.1 -c pytorch -c nvidia`
   - 或者使用 CPU 版本：`conda install pytorch torchvision torchaudio cpuonly -c pytorch`

4. 如果环境混乱：
   - 完全退出所有环境：`conda deactivate && deactivate`
   - 删除旧的虚拟环境：`rm -rf venv`
   - 重新创建 conda 环境：`conda create -n talksphere-gen python=3.10`
   - 重新激活：`conda activate talksphere-gen`

5. M1/M2 Mac 用户注意事项：
   - 使用 CPU 版本的 PyTorch
   - 确保安装的是 arm64 版本的包
   - 如果遇到兼容性问题，可以尝试使用 Rosetta 2 运行

6. 关于 Sliding Window Attention 警告：
   - 当首次加载 Qwen 模型时，可能会看到 "Sliding Window Attention is enabled but not implemented for `sdpa`" 的警告
   - 这个警告不会影响模型的功能，可以安全忽略
   - 这是 transformers 库的一个已知问题，不会影响我们的使用场景

7. 关于内存使用：
   - Qwen 模型需要较大的内存，如果遇到内存不足的问题：
     - 确保系统有足够的可用内存（建议至少 8GB）
     - 可以使用 `torch_dtype=torch.float32` 来减少内存使用
     - 如果仍然遇到问题，可以尝试使用更小的模型