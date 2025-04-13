
## python 运行
```bash
cd /Users/Zhuanz/go/src/forrest/TalkSphere/backend/scripts/

python3 -m venv venv

# 3. 激活虚拟环境
source venv/bin/activate   # 激活后命令行前面会出现 (venv)

# 4. 安装所需的包
pip install -i https://pypi.tuna.tsinghua.edu.cn/simple flask
pip install -i https://pypi.tuna.tsinghua.edu.cn/simple wordcloud jieba pillow cos-python-sdk-v5

# 5. 验证安装
pip list | grep -E "wordcloud|jieba|pillow|cos-python-sdk-v5"
```

```bash
python word_cloud_server.py
```