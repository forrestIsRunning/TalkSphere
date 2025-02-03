

```bash
# 删除旧的虚拟环境
rm -rf venv

# 创建新的虚拟环境
python -m venv venv

# 激活虚拟环境
source venv/bin/activate
pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple

pip install faker requests
```

```bash
python generate-mock-data.py
```