import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import Navbar from './components/Navbar';
import ProtectedRoute from './components/ProtectedRoute';

import Login from './pages/Login';
import Register from './pages/Register';
import Listings from './pages/Listings';
import Booking from './pages/Booking';
import Profile from './pages/Profile';

function App() {
    return (
        <AuthProvider>
            <Router>
                <Navbar />
                <Routes>
                    <Route path="/login" element={<Login />} />
                    <Route path="/register" element={<Register />} />
                    <Route path="/" element={
                        <ProtectedRoute>
                            <Listings />
                        </ProtectedRoute>
                    } />
                    <Route path="/booking" element={
                        <ProtectedRoute>
                            <Booking />
                        </ProtectedRoute>
                    } />
                    <Route path="/profile" element={
                        <ProtectedRoute>
                            <Profile />
                        </ProtectedRoute>
                    } />
                    <Route path="*" element={<div style={{ padding: 20 }}>Page not found</div>} />
                </Routes>
            </Router>
        </AuthProvider>
    );
}

export default App;
