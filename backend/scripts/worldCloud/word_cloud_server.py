from flask import Flask, request, jsonify
import jieba
import os
from wordcloud import WordCloud
import time
import logging
import sys
from logging.handlers import RotatingFileHandler

app = Flask(__name__)

# 获取当前脚本所在目录
current_dir = os.path.dirname(os.path.abspath(__file__))

# 配置日志
log_file = os.path.join(current_dir, 'word_cloud_server.log')
formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')

# 文件处理器
file_handler = RotatingFileHandler(log_file, maxBytes=10*1024*1024, backupCount=5)
file_handler.setFormatter(formatter)

# 控制台处理器
console_handler = logging.StreamHandler(sys.stdout)
console_handler.setFormatter(formatter)

# 配置根日志记录器
logging.basicConfig(
    level=logging.INFO,
    handlers=[file_handler, console_handler]
)

# 获取当前脚本所在目录
save_dir = os.path.join(current_dir, 'wordcloud_images')

# 确保保存目录存在
if not os.path.exists(save_dir):
    os.makedirs(save_dir)
    logging.info(f"创建词云图保存目录: {save_dir}")

@app.route('/generate_wordcloud', methods=['POST'])
def generate_wordcloud():
    try:
        # 获取请求数据
        data = request.get_json()
        text = data.get('text', '')
        
        if not text:
            logging.error("接收到的文本为空")
            return jsonify({
                "success": False,
                "error": "文本内容不能为空"
            }), 400
        
        logging.info(f"接收到的文本长度: {len(text)}")
        
        # 加载停用词
        stop_words = set([
            '的', '了', '和', '是', '就', '在', '我', '有', '而', '你',
            '这', '那', '也', '还', '但', '都', '对', '与', '向', '并',
            'the', 'a', 'an', 'and', 'or', 'but', 'in', 'on', 'at', 'to',
            'for', 'of', 'with', 'by'
        ])
        
        # 使用jieba分词
        words = jieba.cut(text)
        
        # 过滤停用词并合并为字符串
        word_list = [word for word in words if word not in stop_words and len(word) > 1]
        word_space_split = ' '.join(word_list)
        
        # 创建词云对象
        wordcloud = WordCloud(
            font_path='/System/Library/Fonts/PingFang.ttc',  # 使用系统字体
            width=800,
            height=400,
            background_color='white',
            max_font_size=150,
            min_font_size=10,
            max_words=100
        ).generate(word_space_split)
        
        # 生成文件名（使用时间戳）
        filename = f'wordcloud_{int(time.time())}.png'
        save_path = os.path.join(save_dir, filename)
        
        # 保存词云图
        logging.info(f"保存词云图到: {save_path}")
        wordcloud.to_file(save_path)
        
        # 验证文件是否成功创建
        if not os.path.exists(save_path):
            logging.error(f"文件未成功创建: {save_path}")
            return jsonify({
                "success": False,
                "error": "文件未成功创建"
            }), 500
        
        # 获取文件大小
        file_size = os.path.getsize(save_path)
        logging.info(f"文件创建成功，大小: {file_size} bytes")
        
        return jsonify({
            "success": True,
            "image_path": save_path,
            "file_size": file_size
        })
        
    except Exception as e:
        logging.error(f"生成词云图失败: {str(e)}", exc_info=True)
        return jsonify({
            "success": False,
            "error": str(e)
        }), 500

if __name__ == '__main__':
    # 确保工作目录正确
    os.chdir(current_dir)
    logging.info(f"当前工作目录: {os.getcwd()}")
    
    # 创建字体目录
    font_dir = os.path.join(current_dir, 'fonts')
    if not os.path.exists(font_dir):
        os.makedirs(font_dir)
        logging.info(f"创建字体目录: {font_dir}")
    
    app.run(host='127.0.0.1', port=5008) 