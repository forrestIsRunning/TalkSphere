from flask import Flask, request, jsonify
import jieba
import os
from wordcloud import WordCloud
import time
import logging
import sys

app = Flask(__name__)

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    stream=sys.stdout
)

@app.route('/generate_wordcloud', methods=['POST'])
def generate_wordcloud():
    try:
        # 获取请求数据
        data = request.get_json()
        text = data.get('text', '')
        
        logging.info(f"接收到的文本长度: {len(text)}")
        
        # 使用绝对路径
        base_dir = "/backend/scripts"  # 修改为你的实际路径
        save_dir = os.path.join(base_dir, 'wordcloud_images')
        
        # 检查并创建保存目录
        try:
            if not os.path.exists(save_dir):
                os.makedirs(save_dir)
                logging.info(f"创建目录: {save_dir}")
            
            # 检查目录权限
            if not os.access(save_dir, os.W_OK):
                logging.error(f"没有写入权限: {save_dir}")
                return jsonify({
                    "success": False,
                    "error": f"没有写入权限: {save_dir}"
                }), 500
        except Exception as e:
            logging.error(f"创建目录失败: {str(e)}")
            return jsonify({
                "success": False,
                "error": f"创建目录失败: {str(e)}"
            }), 500
        
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
            font_path='/System/Library/Fonts/PingFang.ttc',  # 使用支持中文的字体
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
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    logging.info(f"当前工作目录: {os.getcwd()}")
    
    app.run(host='127.0.0.1', port=5000) 