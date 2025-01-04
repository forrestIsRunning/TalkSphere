document.addEventListener('DOMContentLoaded', function() {
    const avatarUpload = document.getElementById('avatarUpload');
    const avatarPreview = document.getElementById('avatarPreview');
    const bioTextarea = document.getElementById('bio');
    const bioLength = document.getElementById('bioLength');
    const profileForm = document.getElementById('profileForm');

    // 加载用户现有资料
    loadUserProfile();

    // 头像预览
    avatarUpload.addEventListener('change', function(e) {
        const file = e.target.files[0];
        if (file) {
            if (file.size > 5 * 1024 * 1024) { // 5MB限制
                alert('图片大小不能超过5MB');
                return;
            }
            
            const reader = new FileReader();
            reader.onload = function(e) {
                avatarPreview.src = e.target.result;
            };
            reader.readAsDataURL(file);
        }
    });

    // 字数统计
    bioTextarea.addEventListener('input', function() {
        const length = this.value.length;
        bioLength.textContent = length;
        
        if (length > 200) {
            this.value = this.value.substring(0, 200);
            bioLength.textContent = 200;
        }
    });

    // 表单提交
    profileForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const formData = new FormData();
        const avatarFile = avatarUpload.files[0];
        if (avatarFile) {
            formData.append('avatar', avatarFile);
        }
        formData.append('bio', bioTextarea.value);

        try {
            const token = localStorage.getItem('token');
            if (!token) {
                window.location.href = '/login';
                return;
            }

            const response = await fetch('/api/profile/update', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`
                },
                body: formData
            });

            const data = await response.json();
            if (data.code === 1000) {
                alert('保存成功！');
            } else {
                alert(data.msg || '保存失败，请重试');
            }
        } catch (error) {
            console.error('Error:', error);
            alert('网络错误，请重试');
        }
    });

    // 加载用户现有资料
    async function loadUserProfile() {
        try {
            const token = localStorage.getItem('token');
            if (!token) {
                window.location.href = '/login';
                return;
            }

            const response = await fetch('/api/profile', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            const data = await response.json();
            if (data.code === 1000) {
                if (data.data.avatar) {
                    avatarPreview.src = data.data.avatar;
                }
                if (data.data.bio) {
                    bioTextarea.value = data.data.bio;
                    bioLength.textContent = data.data.bio.length;
                }
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }
});