import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {jwtDecode} from 'jwt-decode'; // Correct import
import '../css/Navbar.css'; // For styling

const Navbar = () => {
    const [userRole, setUserRole] = useState(null);
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            try {
                const decodedToken = jwtDecode(token);
                const currentTime = Date.now() / 1000;
                if (decodedToken.exp < currentTime) {
                    // Token expired
                    localStorage.removeItem('token');
                    setUserRole(null);
                } else {
                    setUserRole(decodedToken.userRole);
                }
            } catch (error) {
                // Invalid token
                localStorage.removeItem('token');
                setUserRole(null);
            }
        } else {
            setUserRole(null);
        }
    }, []);

    const handleSignOut = () => {
        localStorage.removeItem('token');
        setUserRole(null);
        navigate('/signin');
    };

    return (
        <nav className="navbar">
            <ul className="center-links">
                <li><Link to="/">Home</Link></li>
            </ul>
            <ul className="end-links">
                {userRole === "1" && (
                    <li><Link to="/add-vendor">Add Vendor</Link></li>
                )}
                {userRole ? (
                    <>
                        <li><Link to="/profile">Profile</Link></li>
                        <li><button onClick={handleSignOut}>Sign Out</button></li>
                    </>
                ) : (
                    <li><Link to="/signin">Sign In</Link></li>
                )}
            </ul>
        </nav>
    );
};

export default Navbar;
