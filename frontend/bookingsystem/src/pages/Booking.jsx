import React, { useEffect, useState } from 'react';
import axiosInstance from '../api/axiosInstance';

const Booking = () => {
    const [bookings, setBookings] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const [listings, setListings] = useState([]);
    const [selectedListing, setSelectedListing] = useState('');
    const [date, setDate] = useState('');
    const [time, setTime] = useState('');
    const [bookingError, setBookingError] = useState(null);
    const [bookingSuccess, setBookingSuccess] = useState(null);

    useEffect(() => {
        const fetchBookings = async () => {
            try {
                const res = await axiosInstance.get('/booking/list');
                setBookings(res.data.bookings || []);
            } catch {
                setError('Failed to load bookings');
            } finally {
                setLoading(false);
            }
        };

        const fetchListings = async () => {
            try {
                const res = await axiosInstance.get('/listing/list');
                setListings(res.data.listings || []);
            } catch {}
        };

        fetchBookings();
        fetchListings();
    }, []);

    const handleSubmit = async e => {
        e.preventDefault();
        setBookingError(null);
        setBookingSuccess(null);

        if (!selectedListing || !date || !time) {
            setBookingError('Fill all fields');
            return;
        }

        try {
            const res = await axiosInstance.post('/booking/create', {
                listing_id: selectedListing,
                date,
                time,
            });
            setBookingSuccess('Booking created successfully');
            setBookings(prev => [...prev, res.data.booking]);
            setSelectedListing('');
            setDate('');
            setTime('');
        } catch {
            setBookingError('Failed to create booking');
        }
    };

    if (loading) return <p>Loading bookings...</p>;
    if (error) return <p style={{ color: 'red' }}>{error}</p>;

    return (
        <div style={{ maxWidth: 600, margin: '20px auto' }}>
            <h2>My Bookings</h2>
            {bookings.length === 0 && <p>No bookings found.</p>}
            <ul>
                {bookings.map(b => (
                    <li key={b.id}>{b.listing_title} — {b.date} {b.time}</li>
                ))}
            </ul>
            <hr />
            <h3>Create New Booking</h3>
            {bookingError && <p style={{ color: 'red' }}>{bookingError}</p>}
            {bookingSuccess && <p style={{ color: 'green' }}>{bookingSuccess}</p>}
            <form onSubmit={handleSubmit}>
                <div style={{ marginBottom: 15 }}>
                    <label>Apartment</label><br />
                    <select value={selectedListing} onChange={e => setSelectedListing(e.target.value)} required style={{ width: '100%', padding: 8 }}>
                        <option value="">Select apartment</option>
                        {listings.map(l => (
                            <option key={l.id} value={l.id}>{l.title}</option>
                        ))}
                    </select>
                </div>
                <div style={{ marginBottom: 15 }}>
                    <label>Date</label><br />
                    <input type="date" value={date} onChange={e => setDate(e.target.value)} required style={{ width: '100%', padding: 8 }} />
                </div>
                <div style={{ marginBottom: 15 }}>
                    <label>Time</label><br />
                    <input type="time" value={time} onChange={e => setTime(e.target.value)} required style={{ width: '100%', padding: 8 }} />
                </div>
                <button type="submit" style={{ padding: '10px 20px' }}>Book</button>
            </form>
        </div>
    );
};

export default Booking;
