import React, { useState, useContext } from 'react';
import axiosInstance from '../api/axiosInstance';
import { AuthContext } from '../contexts/AuthContext';
import { useNavigate, Link } from 'react-router-dom';

const Login = () => {
    const { login } = useContext(AuthContext);
    const navigate = useNavigate();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState(null);

    const handleSubmit = async e => {
        e.preventDefault();
        setError(null);
        try {
            const res = await axiosInstance.post('/auth/login', { email, password });
            console.log('Login response:', res.data);
            login(res.data.access_token, null);
            navigate('/');
        } catch (err) {
            console.error('Login error:', err);
            setError(err.response?.data?.error || 'Login failed');
        }
    };

    return (
        <div style={{ maxWidth: 400, margin: '40px auto' }}>
            <h2>Login</h2>
            {error && <p style={{ color: 'red' }}>{error}</p>}
            <form onSubmit={handleSubmit}>
                <div style={{ marginBottom: 15 }}>
                    <label>Email</label><br/>
                    <input type="email" value={email} onChange={e => setEmail(e.target.value)} required style={{width: '100%', padding: 8}} />
                </div>
                <div style={{ marginBottom: 15 }}>
                    <label>Password</label><br/>
                    <input type="password" value={password} onChange={e => setPassword(e.target.value)} required style={{width: '100%', padding: 8}} />
                </div>
                <button type="submit" style={{ padding: '10px 20px' }}>Login</button>
            </form>
            <p style={{ marginTop: 15 }}>No account? <Link to="/register">Register here</Link></p>
        </div>
    );
};

export default Login;
