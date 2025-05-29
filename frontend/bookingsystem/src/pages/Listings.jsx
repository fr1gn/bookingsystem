import React, { useEffect, useState } from 'react';
import axiosInstance from '../api/axiosInstance';
import ListingCard from '../components/ListingCard';

const Listings = () => {
    const [listings, setListings] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchListings = async () => {
            try {
                const res = await axiosInstance.get('/listing/list');
                setListings(res.data.listings || []);
            } catch (err) {
                setError('Failed to load listings');
            } finally {
                setLoading(false);
            }
        };
        fetchListings();
    }, []);

    if (loading) return <p>Loading listings...</p>;
    if (error) return <p style={{ color: 'red' }}>{error}</p>;
    if (listings.length === 0) return <p>No listings found.</p>;

    return (
        <div style={{ maxWidth: 800, margin: '20px auto' }}>
            <h2>Available Apartments</h2>
            {listings.map(listing => (
                <ListingCard key={listing.id} listing={listing} />
            ))}
        </div>
    );
};

export default Listings;
