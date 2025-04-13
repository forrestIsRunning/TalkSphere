import random
import time
from datetime import datetime, timedelta
import requests
from faker import Faker
import json

# 初始化 Faker
fake = Faker(['zh_CN'])

# API配置
BASE_URL = 'http://localhost:8989'  # 确保这是你的 TalkSphere 后端地址
MODEL_SERVICE_URL = 'http://localhost:8000'

# 管理员账号信息
ADMIN_USERNAME = 'super_admin'
ADMIN_PASSWORD = 'super_admin'
ADMIN_TOKEN = None

# 话题标签列表
TAGS = ['技术', '生活', '美食', '旅游', '电影', '音乐', '游戏', '运动', '读书', '摄影',
        '职场', '教育', '健康', '时尚', '宠物', '汽车', '财经', '科技', '艺术', '心理']

def get_admin_token():
    """获取管理员token"""
    try:
        response = requests.post(f"{BASE_URL}/api/login", json={
            "username": ADMIN_USERNAME,
            "password": ADMIN_PASSWORD
        })
        if response.status_code == 200:
            data = response.json()
            if data.get('code') == 1000:
                return data['data']['token']
        print(f"获取管理员token失败: {response.text}")
        return None
    except Exception as e:
        print(f"获取管理员token异常: {str(e)}")
        return None

def get_random_avatar():
    """获取随机头像URL"""
    # 选择使用哪种头像生成服务
    service = random.choice(['gravatar', 'robohash', 'ui-avatars'])

    # 生成随机字符串作为种子
    seed = fake.md5()

    if service == 'gravatar':
        style = random.choice(['identicon', 'monsterid', 'wavatar', 'retro', 'robohash'])
        return f"https://www.gravatar.com/avatar/{seed}?d={style}&s=200"

    elif service == 'robohash':
        set_num = random.choice(['set1', 'set2', 'set3', 'set4'])
        return f"https://robohash.org/{seed}?set={set_num}"

    else:  # ui-avatars
        name = fake.name().replace(' ', '+')
        colors = ['0096FF', 'FF4C4C', '47B39C', 'FFC54D', '8E44AD']
        bg_color = random.choice(colors)
        return f"https://ui-avatars.com/api/?name={name}&background={bg_color}&color=fff"


def generate_bio():
    """生成用户简介"""
    prompt = """请生成一个10个字以内的用户简介，要求：
    1. 必须严格控制在10个字以内
    2. 只能包含以下一种组合：
       - 城市 + 职业（如：北京程序员）
       - 城市 + 爱好（如：上海美食家）
       - 职业 + 爱好（如：设计师爱画画）
    3. 不要使用标点符号
    4. 不要出现不合理的组合
    5. 不要包含多余的解释或描述
    
    请严格按照要求生成，不要超过10个字。"""
    
    try:
        print(f"正在向 Qwen 模型服务发送请求: {MODEL_SERVICE_URL}/generate")
        response = requests.post(f"{MODEL_SERVICE_URL}/generate", json={
            "prompt": prompt,
            "max_new_tokens": 20
        })
        print(f"Qwen 模型服务响应状态码: {response.status_code}")
        print(f"Qwen 模型服务响应内容: {response.text}")
        
        if response.status_code != 200:
            raise Exception(f"请求失败: {response.status_code}")
            
        response_data = response.json()
        bio = response_data.get("generated_text", "").strip()
        
        # 如果生成的简介太长，使用备用方案
        if len(bio) > 10:
            print("生成的简介太长，使用备用方案")
            bio_templates = [
                f"{fake.city()}程序员",
                f"{fake.job()}爱好者",
                f"{fake.city()}美食家",
                f"{fake.job()}爱{random.choice(['读书', '运动', '旅行', '摄影'])}"
            ]
            bio = random.choice(bio_templates)
            
        return bio
    except Exception as e:
        print(f"生成简介失败，使用备用方案: {str(e)}")
        bio_templates = [
            f"{fake.city()}程序员",
            f"{fake.job()}爱好者",
            f"{fake.city()}美食家",
            f"{fake.job()}爱{random.choice(['读书', '运动', '旅行', '摄影'])}"
        ]
        return random.choice(bio_templates)


def generate_post_content():
    """生成帖子内容"""
    prompt = """请生成一个简短的帖子内容，要求：
    1. 内容要自然，像真实用户发的帖子
    2. 可以包含以下类型的内容：
       - 分享生活经历或感受
       - 提问或求助
       - 推荐或评价
       - 讨论某个话题
    3. 字数在50-100字之间
    4. 可以适当使用表情符号
    5. 内容要真实自然，不要过于正式或生硬
    
    请生成一个随机的帖子内容。"""
    
    try:
        print(f"正在向 Qwen 模型服务发送请求: {MODEL_SERVICE_URL}/generate")
        response = requests.post(f"{MODEL_SERVICE_URL}/generate", json={
            "prompt": prompt,
            "max_new_tokens": 100
        })
        print(f"Qwen 模型服务响应状态码: {response.status_code}")
        print(f"Qwen 模型服务响应内容: {response.text}")
        
        if response.status_code != 200:
            raise Exception(f"请求失败: {response.status_code}")
            
        response_data = response.json()
        content = response_data.get("generated_text", "").strip()
        
        # 如果生成的内容太短，使用备用方案
        if len(content) < 20:
            print("生成的内容太短，使用备用方案")
            content = "今天天气真好，大家有什么有趣的事情分享吗？"
            
        return content
    except Exception as e:
        print(f"生成帖子内容失败，使用备用方案: {str(e)}")
        return "今天天气真好，大家有什么有趣的事情分享吗？"


def generate_post_title(content):
    """根据内容生成标题"""
    prompt = f"""请根据以下内容生成一个简短的标题，要求：
    1. 标题要简洁明了，不超过15个字
    2. 要能准确概括内容的主旨
    3. 不要使用标点符号
    4. 不要使用"关于"、"讨论"等词
    
    内容：
    {content}
    
    请生成一个合适的标题。"""
    
    try:
        print(f"正在向 Qwen 模型服务发送请求生成标题: {MODEL_SERVICE_URL}/generate")
        response = requests.post(f"{MODEL_SERVICE_URL}/generate", json={
            "prompt": prompt,
            "max_new_tokens": 20
        })
        print(f"Qwen 模型服务响应状态码: {response.status_code}")
        print(f"Qwen 模型服务响应内容: {response.text}")
        
        if response.status_code != 200:
            raise Exception(f"请求失败: {response.status_code}")
            
        response_data = response.json()
        title = response_data.get("generated_text", "").strip()
        
        # 如果生成的标题太长，使用备用方案
        if len(title) > 15:
            print("生成的标题太长，使用备用方案")
            title = "分享我的日常"
            
        return title
    except Exception as e:
        print(f"生成标题失败，使用备用方案: {str(e)}")
        return "分享我的日常"


def register_user(username, password, email, bio, avatar_url):
    """注册用户"""
    try:
        response = requests.post(f"{BASE_URL}/api/register", json={
            "username": username,
            "password": password,
            "email": email,
            "bio": bio,
            "avatar_url": avatar_url
        })
        if response.status_code != 200:
            print(f"注册用户失败: {response.status_code} - {response.text}")
            return None
        return response.json()
    except Exception as e:
        print(f"注册用户请求异常: {str(e)}")
        return None


def create_post(token, title, content, board_id, tags):
    """创建帖子"""
    try:
        # 确保 token 格式正确，避免重复的 "Bearer" 前缀
        if token.startswith("Bearer "):
            token = token[7:]  # 移除开头的 "Bearer "
        headers = {"Authorization": f"Bearer {token}"}
        print(f"正在创建帖子: {title}")
        print(f"请求头: {headers}")
        print(f"请求数据: {json.dumps({'title': title, 'content': content, 'board_id': board_id, 'tags': tags}, ensure_ascii=False)}")
        
        response = requests.post(f"{BASE_URL}/api/posts",
                               headers=headers,
                               json={
                                   "title": title,
                                   "content": content,
                                   "board_id": board_id,
                                   "tags": tags
                               })
        
        print(f"响应状态码: {response.status_code}")
        print(f"响应内容: {response.text}")
        
        if response.status_code != 200:
            print(f"创建帖子失败: {response.status_code} - {response.text}")
            return None
            
        try:
            response_data = response.json()
            if response_data.get('code') != 1000:
                print(f"创建帖子失败，错误码: {response_data.get('code')}")
                return None
            return response_data
        except json.JSONDecodeError as e:
            print(f"解析响应 JSON 失败: {str(e)}")
            print(f"原始响应: {response.text}")
            return None
            
    except Exception as e:
        print(f"创建帖子请求异常: {str(e)}")
        return None


def get_boards(token):
    """获取所有板块"""
    try:
        headers = {"Authorization": f"Bearer {token}"}
        response = requests.get(f"{BASE_URL}/api/boards", headers=headers)
        if response.status_code != 200:
            print(f"获取板块失败: {response.status_code} - {response.text}")
            return None
        return response.json()
    except Exception as e:
        print(f"获取板块请求异常: {str(e)}")
        return None


def generate_users(count):
    """生成用户数据"""
    print(f"开始生成 {count} 个用户...")
    users = []
    for i in range(count):
        username = fake.user_name()
        password = "123456"  # 统一使用简单密码便于测试
        email = fake.email()
        bio = generate_bio()
        avatar_url = get_random_avatar()

        result = register_user(username, password, email, bio, avatar_url)
        if result and result.get('code') == 1000:
            users.append({
                'username': username,
                'password': password,
                'email': email
            })
            print(f"成功创建用户 {i + 1}/{count}: {username}")
        else:
            print(f"创建用户失败 {i + 1}/{count}: {result}")

        # 添加延时避免请求过快
        time.sleep(0.5)

    return users


def generate_posts(users, count):
    """生成帖子数据"""
    print(f"开始生成 {count} 个帖子...")

    # 获取板块信息
    boards_info = get_boards(ADMIN_TOKEN)
    if not boards_info or boards_info.get('code') != 1000:
        print("获取板块信息失败")
        return

    boards = boards_info.get('data', [])
    if not boards:
        print("没有可用的板块")
        return

    # 为每个用户生成帖子
    for i in range(count):
        # 随机选择一个用户
        user = random.choice(users)
        print(f"\n处理用户 {user['username']} 的帖子...")

        # 登录获取token
        try:
            print("正在登录用户...")
            login_response = requests.post(f"{BASE_URL}/api/login", json={
                "username": user['username'],
                "password": user['password']
            })
            print(f"登录响应状态码: {login_response.status_code}")
            print(f"登录响应内容: {login_response.text}")
            
            if login_response.status_code != 200:
                print(f"登录失败: {login_response.status_code}")
                continue
                
            try:
                login_data = login_response.json()
                if login_data.get('code') != 1000:
                    print(f"用户登录失败: {user['username']}")
                    continue

                token = login_data['data']['token']
                print(f"成功获取用户 token")

                # 生成帖子数据
                content = generate_post_content()
                title = generate_post_title(content)
                board_id = random.choice(boards)['ID']
                post_tags = random.sample(TAGS, random.randint(1, 3))

                print(f"正在创建帖子...")
                print(f"标题: {title}")
                print(f"内容长度: {len(content)}")
                print(f"板块ID: {board_id}")
                print(f"标签: {post_tags}")

                # 创建帖子
                result = create_post(token, title, content, board_id, post_tags)
                if result and result.get('code') == 1000:
                    print(f"成功创建帖子 {i + 1}/{count}")
                else:
                    print(f"创建帖子失败 {i + 1}/{count}: {result}")

                # 添加延时避免请求过快
                time.sleep(0.5)

            except json.JSONDecodeError as e:
                print(f"解析登录响应 JSON 失败: {str(e)}")
                print(f"原始响应: {login_response.text}")
                continue
                
        except Exception as e:
            print(f"处理帖子时发生错误: {str(e)}")
            continue


def main():
    # 获取管理员token
    global ADMIN_TOKEN
    ADMIN_TOKEN = get_admin_token()
    if not ADMIN_TOKEN:
        print("无法获取管理员token，程序退出")
        return

    # 生成用户
    user_count = 50 # 可以修改生成的用户数量
    users = generate_users(user_count)

    if users:
        # 生成帖子
        post_count = 100  # 可以修改生成的帖子数量
        generate_posts(users, post_count)

    print("数据生成完成！")


if __name__ == "__main__":
    main()