// State
let currentAuthCode = '';

// Step 1: Redirect
document.getElementById('redirect-btn').addEventListener('click', async () => {
    document.getElementById('redirect-message').hidden = false;
    
    console.log('Step 1: Click here to redirect to Identity Provider');
    
    setTimeout(() => {
        document.getElementById('step1-screen').hidden = true;
        document.getElementById('step2-screen').hidden = false;
    }, 1500);
});

// Step 2: IdP Login
document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    
    try {
        const response = await fetch('/api/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            currentAuthCode = data.auth_code;
            document.getElementById('auth-code').textContent = data.auth_code;
            document.getElementById('login-form').hidden = true;
            document.getElementById('auth-code-display').hidden = false;
            
            console.log(`Step 2: IdP Login - User: ${username}`);
            console.log(`Program prints: "Auth Code: ${data.auth_code}"`);
        } else {
            alert('Login failed: ' + (data.error || 'Unknown error'));
        }
    } catch (error) {
        alert('Error: ' + error.message);
    }
});

document.getElementById('proceed-to-exchange').addEventListener('click', () => {
    document.getElementById('step2-screen').hidden = true;
    document.getElementById('step3-screen').hidden = false;
    // Pre-filling the auth code for convenience
    document.getElementById('code-input').value = currentAuthCode;
});

// Step 3: Token Exchange
document.getElementById('exchange-btn').addEventListener('click', async () => {
    const code = document.getElementById('code-input').value.trim();
    
    if (!code) {
        alert('Please enter the authorization code');
        return;
    }
    
    try {
        const response = await fetch('/api/auth/token', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ code })
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Token exchange failed');
        }
        
        const data = await response.json();
        
        if (response.ok) {
            document.getElementById('id-token').value = data.id_token;
            document.getElementById('access-token').value = data.access_token;
            document.getElementById('token-type').textContent = data.token_type;
            document.getElementById('expires-in').textContent = data.expires_in;
            
            document.getElementById('token-display').hidden = false;
            
            console.log('Step 3: Token Exchange');
            console.log('Program returns:');
            console.log(JSON.stringify({
                id_token: data.id_token.substring(0, 50) + '...',
                access_token: data.access_token.substring(0, 50) + '...'
            }, null, 2));
        } else {
            alert('Token exchange failed: ' + (data.error || 'Unknown error'));
        }
    } catch (error) {
        alert('Error: ' + error.message);
    }
});

document.getElementById('proceed-to-verify').addEventListener('click', () => {
    const idToken = document.getElementById('id-token').value;
    const accessToken = document.getElementById('access-token').value;
    
    document.getElementById('step3-screen').hidden = true;
    document.getElementById('step4-screen').hidden = false;
    
    // Pre-fill for convenience
    document.getElementById('verify-id-token').value = idToken;
    document.getElementById('verify-access-token').value = accessToken;
});

// Step 4: Token Verification
document.getElementById('verify-btn').addEventListener('click', async () => {
    const idToken = document.getElementById('verify-id-token').value.trim();
    const accessToken = document.getElementById('verify-access-token').value.trim();
    
    if (!idToken || !accessToken) {
        alert('Please paste both tokens');
        return;
    }
    
    try {
        const response = await fetch('/api/auth/verify', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                id_token: idToken, 
                access_token: accessToken 
            })
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Verification failed');
        }
        
        const data = await response.json();
        
        // Decode tokens to show claims
        const decodedIdToken = decodeJWT(idToken);
        const decodedAccessToken = decodeJWT(accessToken);
        
        let resultHTML = '<div class="verification-result">';
        
        // Display decoded token information
        if (decodedIdToken) {
            resultHTML += '<h3>ID Token Details:</h3>';
            resultHTML += '<p><strong>Header:</strong></p>';
            resultHTML += '<pre>' + JSON.stringify(decodedIdToken.header, null, 2) + '</pre>';
            resultHTML += '<p><strong>Payload (Claims):</strong></p>';
            resultHTML += '<pre>' + JSON.stringify(decodedIdToken.payload, null, 2) + '</pre>';
            resultHTML += '<hr>';
        }
        
        if (decodedAccessToken) {
            resultHTML += '<h3>Access Token Details:</h3>';
            resultHTML += '<p><strong>Header:</strong></p>';
            resultHTML += '<pre>' + JSON.stringify(decodedAccessToken.header, null, 2) + '</pre>';
            resultHTML += '<p><strong>Payload (Claims):</strong></p>';
            resultHTML += '<pre>' + JSON.stringify(decodedAccessToken.payload, null, 2) + '</pre>';
            resultHTML += '<hr>';
        }
        
        resultHTML += '<h3>Program verifies:</h3>';
        resultHTML += '<div class="check-list">';
        
        // Display each check from the backend response
        data.checks.forEach(check => {
            const isPassed = check.status === 'passed';
            resultHTML += `<p>`;
            resultHTML += `${isPassed ? 'Passed: ' : 'Failed: '} ${check.check}`;
            if (check.message) {
                resultHTML += ` - ${check.message}`;
            }
            resultHTML += `</p>`;
        });
        
        resultHTML += '</div>';
        
        resultHTML += `<div>`;
        resultHTML += `<h2>${data.message}</h2>`;
        resultHTML += '</div>';
        resultHTML += '</div>';
        
        document.getElementById('verification-result').innerHTML = resultHTML;
        document.getElementById('verification-result').hidden = false;
        
        console.log('Step 4: Token Verification');
        console.log(data.message);
    } catch (error) {
        alert('Error: ' + error.message);
    }
});

// Helper function to decode JWT token
function decodeJWT(token) {
    try {
        const parts = token.split('.');
        if (parts.length !== 3) {
            return null;
        }
        
        // Decode header
        const header = JSON.parse(atob(parts[0]));
        
        // Decode payload with padding
        let payload = parts[1];
        const padding = 4 - (payload.length % 4);
        if (padding !== 4) {
            payload += '='.repeat(padding);
        }
        const decodedPayload = JSON.parse(atob(payload));
        
        return {
            header: header,
            payload: decodedPayload,
            signature: parts[2]
        };
    } catch (error) {
        console.error('Error decoding JWT:', error);
        return null;
    }
}

// Display decoded token details
function displayDecodedToken(token, type) {
    const decoded = decodeJWT(token);
    if (!decoded) {
        return;
    }
    
    const containerId = type === 'id' ? 'id-token-decoded' : 'access-token-decoded';
    const container = document.getElementById(containerId);
    
    let html = '<h4>Decoded Token:</h4>';
    html += '<p><strong>Header:</strong></p>';
    html += '<pre>' + JSON.stringify(decoded.header, null, 2) + '</pre>';
    html += '<p><strong>Payload (Claims):</strong></p>';
    html += '<pre>' + JSON.stringify(decoded.payload, null, 2) + '</pre>';
    html += '<p><strong>Signature:</strong></p>';
    html += '<pre>' + decoded.signature.substring(0, 50) + '...</pre>';
    
    container.innerHTML = html;
    container.hidden = false;
}

// Helper function to copy tokens
function copyToken(textareaId) {
    const textarea = document.getElementById(textareaId);
    textarea.select();
    document.execCommand('copy');
    alert('Token copied to clipboard!');
}
