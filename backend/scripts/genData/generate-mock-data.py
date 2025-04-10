import random
import time
from datetime import datetime, timedelta
import requests
from faker import Faker
import json

# 初始化 Faker
fake = Faker(['zh_CN'])

# API配置
BASE_URL = 'http://localhost:8989'
ADMIN_TOKEN = ''  # 这里填入管理员token

# 头像URL列表（可以替换成你想要的头像URL列表）
AVATAR_URLS = [
    "https://api.dicebear.com/7.x/avataaars/svg?seed={}",
    "https://api.dicebear.com/7.x/bottts/svg?seed={}",
    "https://api.dicebear.com/7.x/personas/svg?seed={}"
]

# 话题标签列表
TAGS = ['技术', '生活', '美食', '旅游', '电影', '音乐', '游戏', '运动', '读书', '摄影',
        '职场', '教育', '健康', '时尚', '宠物', '汽车', '财经', '科技', '艺术', '心理']

# 帖子内容模板
POST_TEMPLATES = [
    "最近在{place}体验了一下{activity}，感觉真的很{feeling}。{detail}大家有类似经历吗？",
    "分享一下关于{topic}的心得：{content}希望对大家有帮助！",
    "求推荐靠谱的{item}！预算{price}左右，主要考虑{feature}，有经验的朋友来说说吧~",
    "今天在{place}发现了一家很不错的{store}，{review}推荐大家也去试试！",
    "{weather}天气真适合{activity}，{feeling}。{thought}",
]


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
    templates = [
        "热爱{interest}和{interest2}，{city}人",
        "{profession}，{interest}爱好者",
        "专注于{field}领域，{status}",
        "{city}上班族，喜欢{interest}",
        "{status}，{interest}达人"
    ]

    interests = ['摄影', '美食', '旅行', '读书', '音乐', '电影', '运动', '写作', '绘画', '编程']
    professions = ['程序员', '设计师', '产品经理', '自由职业者', '学生', '创业者']
    fields = ['互联网', '教育', '金融', '医疗', '文创', '科技']
    statuses = ['在学习了', '在成长了', '在进步了', '在探索了']

    template = random.choice(templates)
    return template.format(
        interest=random.choice(interests),
        interest2=random.choice(interests),
        city=fake.city(),
        profession=random.choice(professions),
        field=random.choice(fields),
        status=random.choice(statuses)
    )


def generate_post_content():
    """生成帖子内容"""
    template = random.choice(POST_TEMPLATES)

    places = ['商场', '公园', '咖啡馆', '图书馆', '健身房', '餐厅', '电影院']
    activities = ['购物', '运动', '阅读', '美食', '看电影', '听音乐会', '徒步']
    feelings = ['开心', '惊喜', '放松', '充实', '有趣', '难忘']
    items = ['手机', '电脑', '相机', '耳机', '平板', '显示器', '键盘']
    prices = ['1000', '2000', '3000', '5000', '8000', '10000']
    features = ['性价比', '质量', '外观', '性能', '便携性', '续航']
    stores = ['餐厅', '咖啡馆', '书店', '甜品店', '小店']
    weathers = ['晴朗的', '阴天的', '下雨的', '多云的', '温暖的']
    thoughts = ['真是美好的一天！', '生活就是要这样！', '感觉整个人都放松了~', '推荐大家也来试试！']
    topics = ['职场经验', '学习方法', '生活技巧', '理财心得', '健康知识']

    return template.format(
        place=random.choice(places),
        activity=random.choice(activities),
        feeling=random.choice(feelings),
        detail=fake.text(max_nb_chars=100),
        topic=random.choice(topics),
        content=fake.text(max_nb_chars=200),
        item=random.choice(items),
        price=random.choice(prices),
        feature=random.choice(features),
        store=random.choice(stores),
        review=fake.text(max_nb_chars=50),
        weather=random.choice(weathers),
        thought=random.choice(thoughts)
    )


def register_user(username, password, email, bio, avatar_url):
    """注册用户"""
    try:
        response = requests.post(f"{BASE_URL}/register", json={
            "username": username,
            "password": password,
            "email": email,
            "bio": bio,
            "avatar_url": avatar_url
        })
        return response.json()
    except Exception as e:
        print(f"注册用户失败: {str(e)}")
        return None


def create_post(token, title, content, board_id, tags):
    """创建帖子"""
    try:
        headers = {"Authorization": token}
        response = requests.post(f"{BASE_URL}/api/posts",
                                 headers=headers,
                                 json={
                                     "title": title,
                                     "content": content,
                                     "board_id": board_id,
                                     "tags": tags
                                 }
                                 )
        return response.json()
    except Exception as e:
        print(f"创建帖子失败: {str(e)}")
        return None


def get_boards(token):
    """获取所有板块"""
    try:
        headers = {"Authorization": token}
        response = requests.get(f"{BASE_URL}/api/boards", headers=headers)
        return response.json()
    except Exception as e:
        print(f"获取板块失败: {str(e)}")
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

        # 登录获取token
        try:
            login_response = requests.post(f"{BASE_URL}/login", json={
                "username": user['username'],
                "password": user['password']
            })
            login_data = login_response.json()
            if login_data.get('code') != 1000:
                print(f"用户登录失败: {user['username']}")
                continue

            token = login_data['data']['token']

            # 生成帖子数据
            title = fake.sentence()[:20]
            content = generate_post_content()
            board_id = random.choice(boards)['ID']
            post_tags = random.sample(TAGS, random.randint(1, 3))

            # 创建帖子
            result = create_post(token, title, content, board_id, post_tags)
            if result and result.get('code') == 1000:
                print(f"成功创建帖子 {i + 1}/{count}")
            else:
                print(f"创建帖子失败 {i + 1}/{count}: {result}")

            # 添加延时避免请求过快
            time.sleep(0.5)

        except Exception as e:
            print(f"处理帖子时发生错误: {str(e)}")
            continue


def main():
    """主函数"""
    # 生成用户
    user_count = 100 # 可以修改生成的用户数量
    users = generate_users(user_count)

    if users:
        # 生成帖子
        post_count = 100  # 可以修改生成的帖子数量
        generate_posts(users, post_count)

    print("数据生成完成！")


if __name__ == "__main__":
    main()