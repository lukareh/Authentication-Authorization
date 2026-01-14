const API_BASE = 'http://localhost:8080/api';
let currentUsername = '';

function switchTab(tab) {
    // Hide all sections
    document.getElementById('register-section').style.display = 'none';
    document.getElementById('login-section').style.display = 'none';
    document.getElementById('welcome-section').style.display = 'none';
    
    // Show selected section
    document.getElementById(`${tab}-section`).style.display = 'block';
    
    // Hide messages and OTP section
    hideMessage();
    document.getElementById('otp-section').style.display = 'none';

    // Reset forms
    document.getElementById('register-form').reset();
    document.getElementById('login-form').reset();
}

function showMessage(message, type) {
    const messageDiv = document.getElementById('message');
    messageDiv.innerHTML = `<p><strong>${message}</strong></p>`;
}

function hideMessage() {
    document.getElementById('message').innerHTML = '';
}

async function handleRegister(event) {
    event.preventDefault();
    hideMessage();

    const username = document.getElementById('reg-username').value;
    const password = document.getElementById('reg-password').value;

    try {
        const response = await fetch(`${API_BASE}/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });

        const data = await response.json();

        if (data.success) {
            showMessage(' ' + data.message, 'success');
            document.getElementById('register-form').reset();
            setTimeout(() => switchTab('login'), 2000);
        } else {
            showMessage(' ' + data.message, 'error');
        }
    } catch (error) {
        showMessage(' Failed to connect to server', 'error');
        console.error('Error:', error);
    }
}

async function handleLogin(event) {
    event.preventDefault();
    hideMessage();

    currentUsername = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await fetch(`${API_BASE}/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username: currentUsername, password }),
        });

        const data = await response.json();

        if (data.success) {
            showMessage(' Password verified! Check the terminal for your OTP.', 'info');
            document.getElementById('otp-section').style.display = 'block';
            document.getElementById('login-form').style.display = 'none';
        } else {
            showMessage(' ' + data.message, 'error');
        }
    } catch (error) {
        showMessage('Failed to connect to server', 'error');
        console.error('Error:', error);
    }
}

async function handleVerifyOTP(event) {
    event.preventDefault();
    hideMessage();

    const otp = document.getElementById('otp-input').value;

    try {
        const response = await fetch(`${API_BASE}/verify-otp`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username: currentUsername, otp }),
        });

        const data = await response.json();

        if (data.success) {
            document.getElementById('otp-section').style.display = 'none';
            document.getElementById('welcome-username').textContent = currentUsername;
            document.getElementById('welcome-section').style.display = 'block';
        } else {
            showMessage('wrong ' + data.message, 'error');
        }
    } catch (error) {
        showMessage('Failed to connect to server', 'error');
        console.error('Error:', error);
    }
}

function resetApp() {
    currentUsername = '';
    document.getElementById('welcome-section').style.display = 'none';
    document.getElementById('login-form').style.display = 'block';
    document.getElementById('login-form').reset();
    document.getElementById('otp-form').reset();
    hideMessage();
    switchTab('login');
}
