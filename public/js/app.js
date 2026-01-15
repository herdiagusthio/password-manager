function openAddModal() {
    document.getElementById('modalTitle').innerText = 'Add New Secret';
    document.getElementById('secretId').value = '';
    document.getElementById('secretForm').reset();
    document.getElementById('secretModal').classList.remove('hidden');
}

function closeModal() {
    document.getElementById('secretModal').classList.add('hidden');
}

async function saveSecret(event) {
    event.preventDefault();
    const id = document.getElementById('secretId').value;
    const title = document.getElementById('title').value;
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const url = document.getElementById('url').value;

    const payload = {
        title,
        username,
        password,
        metadata: { url }
    };

    let method = 'POST';
    let endpoint = '/api/secrets';

    if (id) {
        method = 'PUT';
        endpoint = `/api/secrets/${id}`;
    }

    try {
        const response = await fetch(endpoint, {
            method: method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        });

        if (response.ok) {
            window.location.reload();
        } else {
            const err = await response.json();
            alert('Error: ' + err.error);
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to save secret');
    }
}

async function deleteSecret(id) {
    if (!confirm('Are you sure you want to delete this secret?')) return;

    try {
        const response = await fetch(`/api/secrets/${id}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            window.location.reload();
        } else {
            alert('Failed to delete secret');
        }
    } catch (error) {
        console.error('Error:', error);
    }
}

function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        showToast('Copied to clipboard!');
    }, (err) => {
        console.error('Could not copy text: ', err);
    });
}

function showToast(message) {
    const toast = document.createElement('div');
    toast.className = 'fixed bottom-4 right-4 bg-gray-800 text-white px-4 py-2 rounded shadow-lg transition-opacity duration-300';
    toast.innerText = message;
    document.body.appendChild(toast);
    setTimeout(() => {
        toast.classList.add('opacity-0');
        setTimeout(() => toast.remove(), 300);
    }, 2000);
}

// Logic to reveal/copy would require fetching the decrypted secret first
// For simplicity in this MVP, we will fetch the secret details when "Edit" is clicked
// And we can add a "Copy" button that fetches, decrypts, and copies on the fly.
async function copyPassword(id) {
    try {
        const response = await fetch(`/api/secrets/${id}`);
        const data = await response.json();
        if (data.password) {
            copyToClipboard(data.password);
        }
    } catch (error) {
        console.error("Failed to fetch password", error);
    }
}

async function openEditModal(id) {
    try {
        const response = await fetch(`/api/secrets/${id}`);
        const data = await response.json();
        
        document.getElementById('modalTitle').innerText = 'Edit Secret';
        document.getElementById('secretId').value = data.id;
        document.getElementById('title').value = data.title;
        document.getElementById('username').value = data.username;
        document.getElementById('password').value = data.password; // This comes decrypted from GET /api/secrets/:id
        document.getElementById('url').value = data.metadata ? data.metadata.url : '';
        
        document.getElementById('secretModal').classList.remove('hidden');
    } catch (error) {
        console.error("Failed to fetch details", error);
    }
    
}
