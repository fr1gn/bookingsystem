import React, { useState } from 'react';
import axiosInstance from '../api/axiosInstance';
import { useNavigate, Link } from 'react-router-dom';

const Register = () => {
    const navigate = useNavigate();
    const [fullName, setFullName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState(null);
    const [success, setSuccess] = useState(null);

    const handleSubmit = async e => {
        e.preventDefault();
        setError(null);
        setSuccess(null);
        try {
            const res = await axiosInstance.post('/auth/register', { full_name: fullName, email, password });
            setSuccess(res.data.message || 'Registration successful. You can login now.');
            setTimeout(() => navigate('/login'), 2000);
        } catch (err) {
            setError(err.response?.data?.error || 'Registration failed');
        }
    };

    return (
        <div style={{ maxWidth: 400, margin: '40px auto' }}>
            <h2>Register</h2>
            {error && <p style={{ color: 'red' }}>{error}</p>}
            {success && <p style={{ color: 'green' }}>{success}</p>}
            <form onSubmit={handleSubmit}>
                <div style={{ marginBottom: 15 }}>
                    <label>Full Name</label><br/>
                    <input type="text" value={fullName} onChange={e => setFullName(e.target.value)} required style={{width: '100%', padding: 8}} />
                </div>
                <div style={{ marginBottom: 15 }}>
                    <label>Email</label><br/>
                    <input type="email" value={email} onChange={e => setEmail(e.target.value)} required style={{width: '100%', padding: 8}} />
                </div>
                <div style={{ marginBottom: 15 }}>
                    <label>Password</label><br/>
                    <input type="password" value={password} onChange={e => setPassword(e.target.value)} required style={{width: '100%', padding: 8}} />
                </div>
                <button type="submit" style={{ padding: '10px 20px' }}>Register</button>
            </form>
            <p style={{ marginTop: 15 }}>Have an account? <Link to="/login">Login here</Link></p>
        </div>
    );
};

export default Register;
