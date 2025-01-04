document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM Content Loaded'); // 调试日志

    // 检查是否已登录
    const token = localStorage.getItem('token');
    
    // 添加一个标记来防止重复重定向
    const redirecting = sessionStorage.getItem('redirecting');
    if (token && !redirecting) {
        // 设置重定向标记
        sessionStorage.setItem('redirecting', 'true');
        
        // 使用延时来避免快速重定向
        setTimeout(() => {
            window.location.href = '/forum';
            // 完成后清除标记
            sessionStorage.removeItem('redirecting');
        }, 100);
        return;
    }

    // 如果已经在重定向过程中，清除标记
    if (redirecting) {
        sessionStorage.removeItem('redirecting');
    }

    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    const switchToRegister = document.getElementById('switchToRegister');
    const switchToLogin = document.getElementById('switchToLogin');

    console.log('Forms found:', { // 调试日志
        loginForm: !!loginForm,
        registerForm: !!registerForm,
        switchToRegister: !!switchToRegister,
        switchToLogin: !!switchToLogin
    });

    // 切换到注册表单
    switchToRegister.addEventListener('click', function(e) {
        console.log('Switching to register form'); // 调试日志
        e.preventDefault();
        loginForm.classList.add('d-none');
        registerForm.classList.remove('d-none');
    });

    // 切换到登录表单
    switchToLogin.addEventListener('click', function(e) {
        console.log('Switching to login form'); // 调试日志
        e.preventDefault();
        registerForm.classList.add('d-none');
        loginForm.classList.remove('d-none');
    });

    // 处理登录表单提交
    loginForm.addEventListener('submit', async function(e) {
        console.log('Login form submitted'); // 调试日志
        e.preventDefault();
        const formData = new FormData(loginForm);
        const username = formData.get('username');
        const password = formData.get('password');

        if (!username || !password) {
            alert('请填写用户名和密码');
            return;
        }
        
        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password })
            });

            const data = await response.json();
            console.log('Login response:', data);
            
            if (data.code === 1000) {
                // 保存token
                localStorage.setItem('token', data.data.token);
                // 设置请求头
                const headers = new Headers();
                headers.append('Authorization', `Bearer ${data.data.token}`);
                // 使用 replace 进行跳转
                window.location.replace('/forum');
            } else {
                alert(data.msg || '登录失败，请重试');
            }
        } catch (error) {
            console.error('Login error:', error); // 调试日志
            alert('网络错误，请重试');
        }
    });

    // 处理注册表单提交
    registerForm.addEventListener('submit', async function(e) {
        console.log('Register form submitted'); // 调试日志
        e.preventDefault();
        const formData = new FormData(registerForm);
        const username = formData.get('username');
        const password = formData.get('password');
        const confirmPassword = formData.get('confirmPassword');
        const email = formData.get('email');

        console.log('Register form data:', { // 调试日志
            username,
            email,
            passwordLength: password?.length,
            confirmPasswordMatch: password === confirmPassword
        });

        // 表单验证
        if (!username || !password || !confirmPassword || !email) {
            alert('请填写所有必填字段');
            return;
        }

        if (password !== confirmPassword) {
            alert('两次输入的密码不一致');
            return;
        }
        
        try {
            const response = await fetch('/api/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    username,
                    password,
                    email
                })
            });

            const data = await response.json();
            console.log('Register response:', data); // 调试日志
            if (data.code === 1000) {
                alert('注册成功！请登录');
                registerForm.reset();
                registerForm.classList.add('d-none');
                loginForm.classList.remove('d-none');
            } else {
                alert(data.msg || '注册失败，请重试');
            }
        } catch (error) {
            console.error('Register error:', error); // 调试日志
            alert('网络错误，请重试');
        }
    });
});