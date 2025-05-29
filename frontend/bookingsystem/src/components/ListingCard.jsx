import React from 'react';

const ListingCard = ({ listing }) => {
    return (
        <div style={{ border: '1px solid #ddd', padding: 15, borderRadius: 8, marginBottom: 15 }}>
            <h3>{listing.title}</h3>
            <p>{listing.description}</p>
            <p><b>Price:</b> {listing.price} ₸ / month</p>
            <p><b>Location:</b> {listing.location}</p>
        </div>
    );
};

export default ListingCard;
