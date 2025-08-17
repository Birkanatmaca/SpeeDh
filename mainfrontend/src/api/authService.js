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

export const getHistory = () => {
    const token = localStorage.getItem('authToken');
    return axios.get(`${API_URL}/transcripts`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
};

export const getAudioFile = (transcriptId) => {
    const token = localStorage.getItem('authToken');
    return axios.get(`${API_URL}/transcripts/${transcriptId}/audio`, {
        headers: {
            'Authorization': `Bearer ${token}`
        },
        responseType: 'blob', // Bu çok önemli! Dosyayı blob olarak almamızı sağlar.
    });
};

export const deleteTranscript = (transcriptId) => {
    const token = localStorage.getItem('authToken');
    return axios.delete(`${API_URL}/transcripts/${transcriptId}`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
};