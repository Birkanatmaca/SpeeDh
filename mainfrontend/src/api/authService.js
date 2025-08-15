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

