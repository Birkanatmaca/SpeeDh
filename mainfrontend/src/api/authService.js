import axios from 'axios';

const API_URL = 'http://localhost:8080/api/v1';

export const register = (firstName, lastName, email, password) => {
    return axios.post(`${API_URL}/auth/register`, {
        first_name: firstName,
        last_name: lastName,
        email,
        password,
    });
};

export const login = (email, password) => {
    return axios.post(`${API_URL}/auth/login`, {
        email,
        password,
    });
};

export const verifyEmail = (email, code) => {
    return axios.post(`${API_URL}/auth/verify-email`, {
        email,
        code,
    });
};

export const transcribeAudio = (audioFile) => {
    const formData = new FormData();
    formData.append('audio', audioFile);

    const token = localStorage.getItem('authToken');

    // HATA BURADAYDI: "API_URL" yerine "API_BASE_URL" kullanılmalı.
    // return axios.post(`${API_URL}/transcribe`, formData, {

    // DÜZELTİLMİŞ HALİ:
    return axios.post(`${API_URL}/transcribe`, formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
            'Authorization': `Bearer ${token}`
        }
    });
};