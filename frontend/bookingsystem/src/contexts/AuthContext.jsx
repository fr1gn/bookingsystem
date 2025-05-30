import React, { createContext, useState, useEffect } from 'react';
import axiosInstance from '../api/axiosInstance';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [token, setToken] = useState(localStorage.getItem('token') || null);
    const [user, setUser] = useState(null);

    useEffect(() => {
        if (token) {
            localStorage.setItem('token', token);
            // При загрузке, если токен есть, получаем данные пользователя
            axiosInstance.get('/auth/me').then(res => {
                setUser(res.data.user);
            }).catch(() => {
                setUser(null);
                setToken(null);
                localStorage.removeItem('token');
            });
        } else {
            localStorage.removeItem('token');
            setUser(null);
        }
    }, [token]);

    const login = (token, userData) => {
        setToken(token);
        setUser(userData);
    };

    const logout = () => {
        setToken(null);
        setUser(null);
        localStorage.removeItem('token');
    };

    return (
        <AuthContext.Provider value={{ token, user, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
};
