import React, { useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { AuthContext } from '../contexts/AuthContext';

const Navbar = () => {
    const { token, logout } = useContext(AuthContext);
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    return (
        <nav style={{ padding: '10px 20px', background: '#282c34', color: '#fff', display: 'flex', justifyContent: 'space-between' }}>
            <div>
                <Link to="/" style={{ color: '#61dafb', marginRight: 15, textDecoration: 'none', fontWeight: 'bold' }}>GoBooking</Link>
            </div>
            <div>
                {token ? (
                    <>
                        <Link to="/" style={{ color: '#fff', marginRight: 15, textDecoration: 'none' }}>Listings</Link>
                        <Link to="/booking" style={{ color: '#fff', marginRight: 15, textDecoration: 'none' }}>My Bookings</Link>
                        <Link to="/profile" style={{ color: '#fff', marginRight: 15, textDecoration: 'none' }}>Profile</Link>
                        <button onClick={handleLogout} style={{ background: 'transparent', border: 'none', color: '#fff', cursor: 'pointer' }}>Logout</button>
                    </>
                ) : (
                    <>
                        <Link to="/login" style={{ color: '#fff', marginRight: 15, textDecoration: 'none' }}>Login</Link>
                        <Link to="/register" style={{ color: '#fff', textDecoration: 'none' }}>Register</Link>
                    </>
                )}
            </div>
        </nav>
    );
};

export default Navbar;
