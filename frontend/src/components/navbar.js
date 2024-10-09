import React, { useEffect, useRef, useState, useCallback } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode'; // Fixed import
import '../css/Navbar.css';
import logo from '../css/vendor.jpg';
import defaultImage from '../css/vendor.jpg';
import { useOrderUpdate } from './OrderUpdateContext'; // Import the context

const Navbar = ({ initialCartItems = [], onCartItemsChange, refreshCart }) => {
    const [userRole, setUserRole] = useState(null);
    const [isTransparent, setIsTransparent] = useState(false);
    const [cartDropdownVisible, setCartDropdownVisible] = useState(false);
    const [cartItems, setCartItems] = useState(initialCartItems);
    const [editedCartItems, setEditedCartItems] = useState([]);
    const [cart, setCart] = useState({});
    const [cartDescription, setCartDescription] = useState('');
    const [initialCartDescription, setInitialCartDescription] = useState('');
    const [errorMessage, setErrorMessage] = useState(null);
    const [successMessage, setSuccessMessage] = useState(null);
    const [totalPrice, setTotalPrice] = useState(0);
    const [totalQuantity, setTotalQuantity] = useState(0);
    const [loadingCartItems, setLoadingCartItems] = useState(false);
    const [loadingDeleteItems, setLoadingDeleteItems] = useState({});
    const [loadingCheckout, setLoadingCheckout] = useState(false);
    const [isQuantitySaved, setIsQuantitySaved] = useState(true);
    const navigate = useNavigate();
    const dropdownRef = useRef(null);
    const { setShouldUpdateOrders } = useOrderUpdate();
    const [isSaving, setIsSaving] = useState(false);
    const [initialQuantities, setInitialQuantities] = useState([]);
    const [initialPrices, setInitialPrices] = useState([]);

    // Timer to clear error/success message
    const clearMessageAfterDelay = (type, delay = 5000) => {
        setTimeout(() => {
            if (type === 'error') {
                setErrorMessage(null);
            } else if (type === 'success') {
                setSuccessMessage(null);
            }
        }, delay);
    };

    const fetchCartItems = useCallback(async () => {
        const token = localStorage.getItem('token');
        if (!token) return;

        try {
            setLoadingCartItems(true);
            const [itemsResponse, cartResponse] = await Promise.all([
                fetch('http://localhost:8080/cartitems', {
                    headers: { Authorization: `Bearer ${token}` }
                }),
                fetch('http://localhost:8080/carts', {
                    headers: { Authorization: `Bearer ${token}` }
                })
            ]);

            if (!itemsResponse.ok || !cartResponse.ok) {
                const errorText = await itemsResponse.text();
                throw new Error(errorText || 'You must have a table to checkout!');
            }

            const itemsData = await itemsResponse.json();
            const cartData = await cartResponse.json();

            const cartItemsData = itemsData.cart || [];
            const cartDescription = cartData.cart?.description || '';
            const cartTotalPrice = cartData.cart?.total_price || 0;
            
            setInitialQuantities(cartItemsData.map(item => item.quantity));
            setInitialPrices(cartItemsData.map(item => item.price || 0));
            setCartItems(cartItemsData);
            setEditedCartItems(cartItemsData);
            setCart(cartData);
            setCartDescription(cartDescription);
            setInitialCartDescription(cartDescription);
            setTotalPrice(cartTotalPrice);
            updateTotalQuantityAndPrice(cartItemsData);
            setIsQuantitySaved(true);
        } catch (error) {
            setErrorMessage(error.message || 'Failed to fetch cart items.');
            clearMessageAfterDelay('error');
        } finally {
            setLoadingCartItems(false);
        }
    }, []);

    useEffect(() => {
        fetchCartItems();
    }, [fetchCartItems, refreshCart]);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            try {
                const decodedToken = jwtDecode(token);
                const currentTime = Date.now() / 1000;
                if (decodedToken.exp < currentTime) {
                    localStorage.removeItem('token');
                    setUserRole(null);
                } else {
                    setUserRole(decodedToken.userRole);
                    fetchCartItems();
                }
            } catch (error) {
                console.error('Error decoding token:', error);
                localStorage.removeItem('token');
                setUserRole(null);
            }
        }
    }, [fetchCartItems]);

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

    const handleSignOut = () => {
        localStorage.removeItem('token');
        setUserRole(null);
        navigate('/signin');
    };

    const toggleCartDropdown = () => {
        setCartDropdownVisible(prev => !prev);
    };

    const handleCartDescriptionChange = (e) => {
        const value = e.target.value;
        if (value.length <= 100) {
            setCartDescription(value);
            setIsQuantitySaved(false);
            setErrorMessage(null);
        } else {
            setErrorMessage('Description cannot exceed 100 characters.');
            clearMessageAfterDelay('error');
        }
    };
    const deleteCartItem = async (itemId) => {
        if (!cart) return;
    
        // Set loading state for the item being deleted
        setLoadingDeleteItems(prevState => ({ ...prevState, [itemId]: true }));
    
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
                const errorData = await response.json();
                throw new Error(errorData.errors ? JSON.stringify(errorData.errors) : 'Network response was not ok');
            }
    
            // Successfully deleted, remove the item from the UI
            setEditedCartItems(prevItems => {
                // Filter out the deleted item
                return prevItems.filter(item => item.item_id !== itemId);
            });
            
            updateTotalQuantityAndPrice(editedCartItems.filter(item => item.item_id !== itemId));
        } catch (error) {
            console.error('Error deleting cart item:', error);
            setErrorMessage('Error removing item. Please try again.');
            clearMessageAfterDelay('error');
        } finally {
            // Remove loading state for the item
            setLoadingDeleteItems(prevState => ({ ...prevState, [itemId]: false }));
        }
    };

    const deleteCart = async (cartID) => {
        setLoadingCartItems(true);
        try {
            const token = localStorage.getItem('token');
            const response = await fetch(`http://localhost:8080/carts/${cart.cart.id}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.errors ? JSON.stringify(errorData.errors) : 'Network response was not ok');
            }

            setCartItems([]);
            setEditedCartItems([]);
            setCartDescription('');
            setInitialCartDescription('');
            setTotalPrice(0);
            setTotalQuantity(0);
            setIsQuantitySaved(true);

            console.log('Cart deleted successfully');
        } catch (error) {
            console.error('Error deleting cart:', error);
            setErrorMessage('Error removing cart. Please try again.');
            clearMessageAfterDelay('error');
        } finally {
            setLoadingCartItems(false);
        }
    };

    const updateCartItemQuantityLocally = (itemId, newQuantity) => {
        if (newQuantity < 1) return;

        const updatedItems = editedCartItems.map(item =>
            item.item_id === itemId ? { ...item, quantity: newQuantity } : item
        );
        setEditedCartItems(updatedItems);
        if (updatedItems.find(item => item.item_id === itemId).quantity === cartItems.find(item => item.item_id === itemId).quantity) {
            setIsQuantitySaved(true);
        } else {
            setIsQuantitySaved(false);
        }
        updateTotalQuantityAndPrice(updatedItems);
    };

    const updateTotalQuantityAndPrice = (items) => {
        const total = items.reduce((acc, item) => acc + item.quantity, 0);
        const price = items.reduce((acc, item) => {
            const itemPrice = item.discount || item.price || 0;
            return acc + item.quantity * parseFloat(itemPrice);
        }, 0);

        setTotalQuantity(total);
        setTotalPrice(price);
    };

    const saveCartChanges = async () => {
        setLoadingCartItems(true);
        setIsSaving(true);
        try {
            const token = localStorage.getItem('token');

            const quantityPromises = editedCartItems.map(item =>
                fetch(`http://localhost:8080/cartitems/${item.item_id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                        'Authorization': `Bearer ${token}`
                    },
                    body: new URLSearchParams({
                        item_id: item.item_id,
                        quantity: item.quantity
                    }).toString()
                })
            );
            const descriptionPromise = fetch(`http://localhost:8080/carts/${cart.cart.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                    'Authorization': `Bearer ${token}`
                },
                body: new URLSearchParams({
                    description: cartDescription
                }).toString()
            });

            const responses = await Promise.all([...quantityPromises, descriptionPromise]);

            for (const response of responses) {
                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.errors ? JSON.stringify(errorData.errors.error) : 'Error while saving cart changes');
                }
            }

            await fetchCartItems();

            setIsQuantitySaved(true);
            setSuccessMessage('Cart updated successfully!');
            setErrorMessage(null);
            clearMessageAfterDelay('success');
        } catch (error) {
            console.error('Error saving cart changes:', error);
            setErrorMessage(error.message || 'An unknown error occurred');
            setSuccessMessage(null);
            clearMessageAfterDelay('error');
        } finally {
            setLoadingCartItems(false);
            setIsSaving(false);
        }
    };

    const handleCheckout = async () => {
        if (!isQuantitySaved) {
            setErrorMessage('Please save your changes before checking out.');
            clearMessageAfterDelay('error');
            return;
        }

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
            setShouldUpdateOrders(true);
            window.location.reload();
        } catch (error) {
            console.error('Error during checkout:', error);
            setErrorMessage(error.message);
            setSuccessMessage(null);
            clearMessageAfterDelay('error');
        } finally {
            setLoadingCheckout(false);
        }
    };

    const handleDiscardChanges = () => {
        const resetItems = editedCartItems.map((item, index) => ({
            ...item,
            quantity: initialQuantities[index]
        }));
        setTotalPrice(initialPrices.reduce((acc, price) => acc + price, 0));
        setTotalQuantity(initialQuantities.reduce((acc, qty) => acc + qty, 0));
        setEditedCartItems(resetItems);
        setCartDescription(initialCartDescription);
        setIsQuantitySaved(true);
        setErrorMessage(null);
        setSuccessMessage('Changes discarded.');
        clearMessageAfterDelay('success');
    };

    return (
        <nav className={`navbar ${isTransparent ? 'transparent' : ''}`}>
            <div className="logo">
                <img src={logo} alt="Logo" />
            </div>
            <ul className="center-links">
                <li><Link to="/">Home</Link></li>
            </ul>
            <ul className="end-links">
                {userRole === "1" && (
                    <>
                        <li><Link to="/add-vendor">Add Vendor</Link></li>
                        <li><Link to="/users">Users</Link></li>
                    </>
                )}
                {userRole ? (
                    <>
                        <li><Link to="/profile">Profile</Link></li>
                        <li><Link to="/orders">Orders</Link></li>
                        <li>
                            <button onClick={toggleCartDropdown} disabled={loadingCartItems}>
                                {loadingCartItems ? 'Loading ...' : 'Cart'}
                            </button>
                            {cartDropdownVisible && (
                                <div className="cart-dropdown" ref={dropdownRef}>
                                    {loadingCartItems ? (
                                        <div className="spinner"></div>
                                    ) : cartItems.length === 0 ? (
                                        <p>No items in cart</p>
                                    ) : (
                                        <>
                                            <ul className="cart-list">
                                                {editedCartItems.map(item => (
                                                    <li key={item.item_id} className="cart-item">
                                                        <div className="cart-item-img">
                                                            <img
                                                                src={item.img || defaultImage}
                                                                alt={item.name}
                                                                onError={(e) => (e.target.src = defaultImage)}
                                                            />
                                                        </div>
                                                        <div className="cart-item-name">{item.name}</div>
                                                        <div className="cart-item-quantity">
                                                            <button
                                                                onClick={() =>
                                                                    updateCartItemQuantityLocally(item.item_id, item.quantity - 1)
                                                                }
                                                            >
                                                                -
                                                            </button>
                                                            <span>{item.quantity}</span>
                                                            <button
                                                                onClick={() =>
                                                                    updateCartItemQuantityLocally(item.item_id, item.quantity + 1)
                                                                }
                                                            >
                                                                +
                                                            </button>
                                                        </div>
                                                        <button onClick={() => deleteCartItem(item.item_id)} className="removebutton">
                                                            Remove
                                                        </button>
                                                        {loadingDeleteItems[item.item_id] && <div className="spinner"></div>}
                                                    </li>
                                                ))}
                                            </ul>
                                            <div className="cart-summary">
                                                <p>Total Quantity: {totalQuantity}</p>
                                                <p>Total Price: ${totalPrice.toFixed(2)}</p>
                                            </div>
                                        </>
                                    )}

                                    {cartItems.length !== 0 && (
                                        <div className="cart-description">
                                            <textarea
                                                value={cartDescription}
                                                onChange={handleCartDescriptionChange}
                                                placeholder="Write a note for your order..."
                                                rows="3"
                                                maxLength="100"
                                            />
                                            <p>{cartDescription.length}/100</p>
                                        </div>
                                    )}

                                    {/** Show "Delete Cart" button only if the cart has items */}
                                    {cartItems.length > 0 && cart && cart.cart?.id && (
                                        <button onClick={() => deleteCart(cart.cart.id)}>
                                            Delete Cart
                                        </button>
                                    )}

                                    {!loadingCheckout && totalQuantity > 0 && !isQuantitySaved && (
                                        <>
                                            <button onClick={saveCartChanges} disabled={isSaving}>
                                                {isSaving ? 'Saving...' : 'Save'}
                                            </button>
                                            <button onClick={handleDiscardChanges}>Discard Changes</button>
                                        </>
                                    )}

                                    {isQuantitySaved && (
                                        <button onClick={handleCheckout} disabled={loadingCheckout}>
                                            {loadingCheckout ? 'Processing Checkout...' : 'Checkout'}
                                        </button>
                                    )}
                                    {successMessage && <p style={{ color: 'green' }}>{successMessage}</p>}
                                    {errorMessage && <p style={{ color: 'red' }}>{errorMessage}</p>}
                                </div>
                            )}
                        </li>
                        <li>
                            <button onClick={handleSignOut}>Sign Out</button>
                        </li>
                    </>
                ) : (
                    <li>
                        <Link to="/signin">Sign In</Link>
                    </li>
                )}
            </ul>
        </nav>
    );
};

export default Navbar;