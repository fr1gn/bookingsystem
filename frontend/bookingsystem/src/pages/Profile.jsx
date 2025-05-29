import React, { useContext, useState, useEffect } from 'react';
import axiosInstance from '../api/axiosInstance';
import { AuthContext } from '../contexts/AuthContext';

const Profile = () => {
    const { user } = useContext(AuthContext);
    const [fullName, setFullName] = useState(user?.fullName || '');
    const [email, setEmail] = useState(user?.email || '');
    const [message, setMessage] = useState(null);

    const handleUpdate = async e => {
        e.preventDefault();
        try {
            await axiosInstance.put('/auth/update', { full_name: fullName, email });
            setMessage('Profile updated successfully');
        } catch {
            setMessage('Failed to update profile');
        }
    };

    useEffect(() => {
        if (user) {
            setFullName(user.fullName || '');
            setEmail(user.email || '');
        }
    }, [user]);

    return (
        <div style={{ maxWidth: 400, margin: '40px auto' }}>
            <h2>Profile</h2>
            {message && <p>{message}</p>}
            <form onSubmit={handleUpdate}>
                <div style={{ marginBottom: 15 }}>
                    <label>Full Name</label><br />
                    <input type="text" value={fullName} onChange={e => setFullName(e.target.value)} required style={{ width: '100%', padding: 8 }} />
                </div>
                <div style={{ marginBottom: 15 }}>
                    <label>Email</label><br />
                    <input type="email" value={email} onChange={e => setEmail(e.target.value)} required style={{ width: '100%', padding: 8 }} />
                </div>
                <button type="submit" style={{ padding: '10px 20px' }}>Save</button>
            </form>
        </div>
    );
};

export default Profile;
