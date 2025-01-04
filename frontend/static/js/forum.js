document.addEventListener('DOMContentLoaded', function() {
    // 检查登录状态
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.replace('/');
        return;
    }

    // 设置默认的请求头
    const defaultHeaders = {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
    };

    // 设置全局的 fetch 拦截器
    const originalFetch = window.fetch;
    window.fetch = function(url, options = {}) {
        // 合并默认headers和自定义headers
        options.headers = {
            ...defaultHeaders,
            ...(options.headers || {})
        };
        return originalFetch(url, options);
    };

    // 获取DOM元素
    const userAvatar = document.getElementById('userAvatar');
    const logoutBtn = document.getElementById('logoutBtn');
    const newPostBtn = document.getElementById('newPostBtn');
    const newPostModal = new bootstrap.Modal(document.getElementById('newPostModal'));
    const submitPostBtn = document.getElementById('submitPost');
    const postsContainer = document.querySelector('.posts-container');

    // 加载用户信息
    loadUserProfile();

    // 加载帖子列表
    loadPosts();

    // 退出登录
    logoutBtn.addEventListener('click', function(e) {
        e.preventDefault();
        localStorage.removeItem('token');
        window.location.replace('/');
    });

    // 打开发帖模态框
    newPostBtn.addEventListener('click', function() {
        newPostModal.show();
    });

    // 提交新帖子
    submitPostBtn.addEventListener('click', async function() {
        const form = document.getElementById('newPostForm');
        const formData = new FormData(form);
        
        try {
            const response = await fetch('/api/posts/create', {
                method: 'POST',
                body: JSON.stringify({
                    title: formData.get('title'),
                    content: formData.get('content'),
                    tags: formData.get('tags').split(',').map(tag => tag.trim()).filter(tag => tag)
                })
            });

            const data = await response.json();
            if (data.code === 1000) {
                newPostModal.hide();
                form.reset();
                loadPosts(); // 重新加载帖子列表
                alert('发布成功！');
            } else {
                alert(data.msg || '发布失败，请重试');
            }
        } catch (error) {
            console.error('Error:', error);
            alert('网络错误，请重试');
        }
    });

    // 加载用户信息
    async function loadUserProfile() {
        try {
            const response = await fetch('/api/profile');
            const data = await response.json();
            if (data.code === 1000) {
                if (data.data.avatar) {
                    userAvatar.src = data.data.avatar;
                }
            } else if (data.code === CodeNeedLogin) {
                localStorage.removeItem('token');
                window.location.replace('/');
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }

    // 加载帖子列表
    async function loadPosts() {
        try {
            const response = await fetch('/api/posts');
            const data = await response.json();
            if (data.code === 1000) {
                renderPosts(data.data);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }

    // 渲染帖子列表
    function renderPosts(posts) {
        if (!Array.isArray(posts)) {
            console.error('Posts data is not an array:', posts);
            return;
        }
        
        postsContainer.innerHTML = posts.map(post => `
            <div class="card mb-3 post-card">
                <div class="card-body">
                    <div class="d-flex align-items-center mb-3">
                        <img src="${post.author.avatar || '/static/images/default-avatar.png'}" 
                             class="rounded-circle me-2" width="40" height="40">
                        <div>
                            <h6 class="mb-0">${post.author.username}</h6>
                            <small class="text-muted">${formatTime(post.created_at)}</small>
                        </div>
                    </div>
                    <h5 class="card-title">${post.title}</h5>
                    <p class="card-text">${post.content.substring(0, 100)}${post.content.length > 100 ? '...' : ''}</p>
                    <div class="d-flex justify-content-between align-items-center">
                        <div class="post-tags">
                            ${post.tags.map(tag => `
                                <span class="badge bg-light text-dark me-2">${tag}</span>
                            `).join('')}
                        </div>
                        <div class="post-stats">
                            <span class="me-3"><i class="bi bi-eye me-1"></i>${post.views}</span>
                            <span class="me-3"><i class="bi bi-chat me-1"></i>${post.comments}</span>
                            <span><i class="bi bi-heart me-1"></i>${post.likes}</span>
                        </div>
                    </div>
                </div>
            </div>
        `).join('');
    }

    // 格式化时间
    function formatTime(timestamp) {
        const now = new Date();
        const date = new Date(timestamp * 1000);
        const diff = Math.floor((now - date) / 1000);

        if (diff < 60) {
            return '刚刚';
        } else if (diff < 3600) {
            return Math.floor(diff / 60) + '分钟前';
        } else if (diff < 86400) {
            return Math.floor(diff / 3600) + '小时前';
        } else if (diff < 2592000) {
            return Math.floor(diff / 86400) + '天前';
        } else {
            return date.toLocaleDateString();
        }
    }
}); 