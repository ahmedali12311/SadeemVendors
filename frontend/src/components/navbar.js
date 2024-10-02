<<<<<<< HEAD
import React, { useEffect, useRef, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode'; // Correct import for jwt-decode
import '../css/Navbar.css';
import logo from '../css/vendor.jpg';
import defaultImage from '../css/vendor.jpg';
import { useOrderUpdate } from './OrderUpdateContext'; // Import the custom hook

const Navbar = ({ initialCartItems, onCartItemsChange, refreshCart }) => {
    const [userRole, setUserRole] = useState(null);
    const [isTransparent, setIsTransparent] = useState(false);
    const [cartDropdownVisible, setCartDropdownVisible] = useState(false);
    const [cartItems, setCartItems] = useState(initialCartItems || []);
    const [cart, setCart] = useState(initialCartItems || []);
    const [cartDescription, setCartDescription] = useState('');
    const [initialCartDescription, setInitialCartDescription] = useState(''); // State to store the initial description
    const [errorMessage, setErrorMessage] = useState(null);
    const [removalErrorMessage, setRemovalErrorMessage] = useState(null);
    const [successMessage, setSuccessMessage] = useState(null); // Success message state
    const [totalPrice, setTotalPrice] = useState(0);
    const [totalQuantity, setTotalQuantity] = useState(0);
    const [loadingCartItems, setLoadingCartItems] = useState(false);
    const [loadingCheckout, setLoadingCheckout] = useState(false);
    const navigate = useNavigate();
    const dropdownRef = useRef(null);
    const { setShouldUpdateOrders } = useOrderUpdate(); // Use the context

    useEffect(() => {
        fetchCartItems();
    }, [refreshCart]);
=======
import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {jwtDecode} from 'jwt-decode'; // Correct import
import '../css/Navbar.css'; // For styling

const Navbar = () => {
    const [userRole, setUserRole] = useState(null);
    const navigate = useNavigate();
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            try {
                const decodedToken = jwtDecode(token);
                const currentTime = Date.now() / 1000;
                if (decodedToken.exp < currentTime) {
<<<<<<< HEAD
=======
                    // Token expired
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
                    localStorage.removeItem('token');
                    setUserRole(null);
                } else {
                    setUserRole(decodedToken.userRole);
<<<<<<< HEAD
                    fetchCartItems();
                }
            } catch (error) {
                console.error('Error decoding token:', error);
                localStorage.removeItem('token');
                setUserRole(null);
            }
        }
    }, []);

    useEffect(() => {
        const handleScroll = () => {
            setIsTransparent(window.scrollY > 50);
        };

        window.addEventListener('scroll', handleScroll);
        return () => window.removeEventListener('scroll', handleScroll);
    }, []);

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
                setCartDropdownVisible(false);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const fetchCartItems = async () => {
        const token = localStorage.getItem('token');
        if (!token) return;

        try {
            const [itemsResponse, cartResponse] = await Promise.all([
                fetch('http://localhost:8080/cartitems', {
                    headers: { 'Authorization': `Bearer ${token}` }
                }),
                fetch('http://localhost:8080/carts', {
                    headers: { 'Authorization': `Bearer ${token}` }
                })
            ]);

            if (!itemsResponse.ok || !cartResponse.ok) {
                throw new Error('You must have a table to checkout!');
            }

            const itemsData = await itemsResponse.json();
            setCartItems(itemsData.cart || []);

            const cartData = await cartResponse.json();
            setCart(cartData);
            setCartDescription(cartData.cart?.description || ''); // Update cart description
            setInitialCartDescription(cartData.cart?.description || ''); // Set the initial description
            setTotalPrice(cartData.cart?.total_price || 0); // Update total price
            setTotalQuantity(cartData.cart?.quantity || 0); // Update total quantity

        } catch (error) {
            console.error('Error fetching cart items:', error);
        } finally {
            setLoadingCartItems(false);
        }
    };

    const deleteCart = async (id) => {
        try {
            const token = localStorage.getItem('token');
            if (!token) {
                throw new Error('Token is missing');
            }

            const cartResponse = await fetch(`http://localhost:8080/carts/${id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            if (!cartResponse.ok) {
                const errorData = await cartResponse.json();
                console.error('Error response from server:', errorData);
                throw new Error(`Failed to delete cart: ${errorData.error.role || 'Unknown error'}`);
            }

            console.log('Cart deleted successfully');
            setCartItems([]);
            setTotalPrice(0);
            setTotalQuantity(0);

        } catch (error) {
            console.error('Error deleting cart:', error);
            setErrorMessage('Error deleting cart. Please try again.');
        } finally {
            setLoadingCartItems(false);
        }
    };

=======
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

>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
    const handleSignOut = () => {
        localStorage.removeItem('token');
        setUserRole(null);
        navigate('/signin');
    };

<<<<<<< HEAD
    const toggleCartDropdown = () => {
        setCartDropdownVisible(prev => !prev);
        setErrorMessage(null);
        setRemovalErrorMessage(null);
    };

    const handleCartDescriptionChange = (e) => {
        const value = e.target.value;
        if (value.length <= 100) {
            setCartDescription(value);
            setErrorMessage(null); // Clear error if it's within the limit
        } else {
            setErrorMessage('Description cannot exceed 100 characters.');
        }
    };

    const updateCartDescription = async () => {
        if (cartDescription.length > 100) {
            setErrorMessage('Description cannot exceed 100 characters.');
            return; // Prevent updating if the description is too long
        }

        const token = localStorage.getItem('token');
        try {
            const response = await fetch(`http://localhost:8080/carts/${cart.cart.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                    'Authorization': `Bearer ${token}`
                },
                body: new URLSearchParams({
                    description: cartDescription // Update only description
                }).toString()
            });

            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            await fetchCartItems(); // Refresh cart items after updating description
            setSuccessMessage('Cart updated successfully!'); // Set success message
            setInitialCartDescription(cartDescription); // Update the initial description to the new one
            setTimeout(() => setSuccessMessage(null), 3000); // Clear success message after 3 seconds

        } catch (error) {
            setErrorMessage('Error updating description. Please try again.');
            console.error('Error updating cart description:', error);
        }
    };

    const updateCartItemQuantity = async (itemId, newQuantity) => {
        if (newQuantity < 1) return;

        setErrorMessage(null);
        try {
            const token = localStorage.getItem('token');
            const response = await fetch(`http://localhost:8080/cartitems/${itemId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                    'Authorization': `Bearer ${token}`
                },
                body: new URLSearchParams({ item_id: itemId, quantity: newQuantity }).toString()
            });

            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            await fetchCartItems();
        } catch (error) {
            setErrorMessage('Error updating quantity. Please try again.');
            console.error('Error updating cart item:', error);
        }
    };

    const deleteCartItem = async (itemId) => {
        try {
            const token = localStorage.getItem('token');
            const response = await fetch(`http://localhost:8080/cartitems/${itemId}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            console.log('Item deleted successfully');
            await fetchCartItems();
        } catch (error) {
            setRemovalErrorMessage('Error removing item. Please try again.');
            console.error('Error deleting cart item:', error);
        }
    };

    const handleCheckout = async () => {
        setLoadingCheckout(true);
        try {
            const token = localStorage.getItem('token');
            const response = await fetch('http://localhost:8080/checkout', {
                method: 'POST',
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.errors ? JSON.stringify(errorData.errors) : 'Network response was not ok');
            }

            console.log('Checkout successful');
            setCartItems([]);
            setShouldUpdateOrders(true); // Notify context that orders need to be updated
            window.location.reload(); // Refresh the page
        } catch (error) {
            console.error('Error during checkout:', error);
            setErrorMessage(error.message);
        } finally {
            setLoadingCheckout(false);
        }
    };

    return (
        <nav className={`navbar ${isTransparent ? 'transparent' : ''}`}>
            <div className="logo">
                <img src={logo} alt="Logo" />
            </div>
=======
    return (
        <nav className="navbar">
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
            <ul className="center-links">
                <li><Link to="/">Home</Link></li>
            </ul>
            <ul className="end-links">
                {userRole === "1" && (
<<<<<<< HEAD
                    <>
                        <li><Link to="/add-vendor">Add Vendor</Link></li>
                        <li><Link to="/users">Users</Link></li>
                    </>
=======
                    <li><Link to="/add-vendor">Add Vendor</Link></li>
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
                )}
                {userRole ? (
                    <>
                        <li><Link to="/profile">Profile</Link></li>
<<<<<<< HEAD
                        <li><Link to="/orders">Orders</Link></li>
                        <li>
                            <button onClick={toggleCartDropdown} disabled={loadingCartItems}>
                                {loadingCartItems ? 'Loading Cart...' : 'Cart'}
                            </button>
                            {cartDropdownVisible && (
                                <div className="cart-dropdown" ref={dropdownRef}>
                                    {loadingCartItems ? (
                                        <p>Loading cart items...</p>
                                    ) : cartItems.length === 0 ? (
                                        <p>No items in cart</p>
                                    ) : (
                                        <>
                                            <ul className="cart-list">
                                                {cartItems.map(item => (
                                                    <li key={item.item_id} className="cart-item">
                                                        <div className="cart-item-img">
                                                            <img src={item.img || defaultImage} alt={item.name} />
                                                        </div>
                                                        <div className="cart-item-name">{item.name}</div>
                                                        <div className="cart-item-quantity">
                                                            <button onClick={() => updateCartItemQuantity(item.item_id, item.quantity - 1)}>-</button>
                                                            <span>{item.quantity}</span>
                                                            <button onClick={() => updateCartItemQuantity(item.item_id, item.quantity + 1)}>+</button>
                                                        </div>
                                                        <button onClick={() => deleteCartItem(item.item_id)} className='removebutton'>Remove</button>
                                                    </li>
                                                ))}
                                            </ul>
                                            <div className="cart-summary">
                                                <p>Total Quantity: {totalQuantity}</p>
                                                <p>Total Price: ${totalPrice.toFixed(2)}</p>
                                            </div>
                                        </>
                                    )}
                                                                    {totalQuantity > 0 && (

                                    <div className="cart-description">
                                        <textarea
                                            value={cartDescription}
                                            onChange={handleCartDescriptionChange}
                                            placeholder="Write a description for your cart..."
                                            rows="3"
                                            maxLength="100"
                                        />
                                        <p>{cartDescription.length}/100</p> {/* Character count */}
                                    </div>
                                                                    )}
                                    {errorMessage && <p className="error-message">{errorMessage}</p>}
                                    {successMessage && <p className="success-message">{successMessage}</p>}
                                    {removalErrorMessage && <p className="error-message">{removalErrorMessage}</p>}
                                    
                                    {/* Show Save button only if the cart description has changed */}
                                    {cartDescription !== initialCartDescription && (
                                        <button onClick={updateCartDescription} disabled={loadingCartItems}>
                                            Save
                                        </button>
                                    )}

                                    <button onClick={handleCheckout} disabled={loadingCheckout}>
                                        {loadingCheckout ? 'Processing Checkout...' : 'Checkout'}
                                    </button>
                                    <button onClick={() => deleteCart(cart.cart.id)} disabled={loadingCartItems} className='removebutton'>
                                        Delete Cart
                                    </button>
                                </div>
                            )}
                        </li>
=======
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
                        <li><button onClick={handleSignOut}>Sign Out</button></li>
                    </>
                ) : (
                    <li><Link to="/signin">Sign In</Link></li>
                )}
            </ul>
        </nav>
    );
};

<<<<<<< HEAD
export default Navbar;
=======
export default Navbar;
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
