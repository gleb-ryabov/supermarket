const API_BASE = localStorage.getItem('API_BASE') || 'http://localhost:8080/api';

async function apiRequest(path, options = {}) {
    const response = await fetch(`${API_BASE}${path}`, {
        headers: {
            'Content-Type': 'application/json',
            ...(options.headers || {})
        },
        ...options
    });

    let data = null;
    const text = await response.text();
    if (text) {
        data = JSON.parse(text);
    }

    if (!response.ok) {
        const message = data?.error || `Ошибка HTTP ${response.status}`;
        throw new Error(message);
    }
    return data;
}

function buildQuery(params) {
    const search = new URLSearchParams();
    Object.entries(params || {}).forEach(([key, value]) => {
        if (value !== undefined && value !== null && String(value).trim() !== '') {
            search.set(key, value);
        }
    });
    const query = search.toString();
    return query ? `?${query}` : '';
}
