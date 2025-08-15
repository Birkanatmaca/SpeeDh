import axios from 'axios';

const API_URL = 'http://localhost:8080/api/v1/auth';

export const register = (firstName, lastName, email, password) => {
    return axios.post(`${API_URL}/register`, {
        first_name: firstName,
        last_name: lastName,
        email,
        password,
    });
};

export const login = (email, password) => {
    return axios.post(`${API_URL}/login`, {
        email,
        password,
    });
};

export const verifyEmail = (email, code) => {
    return axios.post(`${API_URL}/verify-email`, {
        email,
        code,
    });
};

export const transcribeAudio = (audioFile) => {
    const formData = new FormData();
    formData.append('audio', audioFile);

    // Token'ı localStorage'dan oku
    const token = localStorage.getItem('authToken');

    return axios.post(`${API_BASE_URL}/transcribe`, formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
            'Authorization': `Bearer ${token}` // Okunan token'ı kullan
        }
    });
};
