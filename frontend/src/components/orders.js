import React, { useEffect, useState, useCallback } from 'react';
import '../css/orders.css'; // Import the updated CSS

const Orders = () => {
    const [orders, setOrders] = useState([]);
    const [error, setError] = useState(null);
    const [successMessage, setSuccessMessage] = useState(null);
    const [userRole, setUserRole] = useState(null);
    const [loading, setLoading] = useState(false);
    const [loadingOrderId, setLoadingOrderId] = useState(null); // Track loading for specific order

    const fetchOrders = useCallback(async () => {
        setLoading(true);
        try {
            const token = localStorage.getItem('token');
            if (!token) {
                throw new Error('No token found in localStorage');
            }
    
            const decodedToken = JSON.parse(atob(token.split('.')[1]));
            const role = decodedToken.userRole;
            const userId = decodedToken.id;

            setUserRole(role);
    
            let url = '/orders';
    
            if (role === '2') {
                const vendorResponse = await fetch(`http://localhost:8080/uservendors/${userId}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });

                if (!vendorResponse.ok) {
                    throw new Error(`Failed to fetch vendor ID: ${vendorResponse.statusText}`);
                }
    
                const vendorData = await vendorResponse.json();
                const vendor = vendorData.vendor[0];
                const vendorId = vendor?.id;
    
                if (!vendorId) {
                    throw new Error('Vendor ID is undefined');
                }
    
                url = `/vendororders/${vendorId}`;
            }
    
            const response = await fetch(`http://localhost:8080${url}`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
    
            if (!response.ok) {
                const errorData = await response.json();
                setError(errorData.error || 'An unknown error occurred');
                return;
            }
    
            const data = await response.json();
            setOrders(data.orders || []);
        } catch (error) {
            console.error('Error fetching orders:', error);
        } finally {
            setLoading(false);
        }
    }, []);
    
    useEffect(() => {
        fetchOrders(); // Call fetchOrders on component mount
    }, [fetchOrders]);
    
    const handleCompleteOrder = async (orderId) => {
        setLoadingOrderId(orderId); // Set loading for specific order
        try {
            const response = await fetch(`http://localhost:8080/orderscompleted/${orderId}`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`,
                    'Content-Type': 'application/x-www-form-urlencoded'
                },
                body: new URLSearchParams({
                    status: 'completed'
                })
            });
    
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
    
            await fetchOrders();
            setSuccessMessage('Order has been marked as completed!');
        } catch (error) {
            setError('Error updating order status. Please try again.');
            console.error('Error updating order status:', error);
        } finally {
            setLoadingOrderId(null); // Reset loading state for the order
        }
    };
    
    return (
        <div className="orders-container">
            <div className="orders-card">
                <h2 className="orders-header">Your Orders</h2>
                {error && <p className="error-message">{error}</p>}
                {successMessage && <p className="success-message">{successMessage}</p>}
                <div className="orders-list">
                    {orders.map(order => (
                        <div key={order.id} className="order-item">
                            <div className="order-summary">
                                <div className="order-details">
                                    <div className="order-detail">
                                        <strong>Total Cost : </strong>
                                        <span className="cost">${order.total_order_cost}</span>
                                    </div>
                                    <div className="order-detail">
                                        <strong>Status : </strong>
                                        <span className={`status ${order.status}`}>{order.status}</span>
                                    </div>
                                </div>
                                <div className="order-vendor">
                                    <strong>Vendor : </strong>
                                    <span>{order.vendor_name}</span>
                                </div>
                            </div>
                            {userRole === '2' && order.status !== 'completed' && (
                                <div className="user-data">
                                    <div className='user-Section'>
                                        <strong>User Info:</strong>
                                    </div>
                                    <div>📞 {order.CustomerPhone}</div>
                                    <div>👤 {order.user_name}</div>
                                </div>
                            )}
                            <div className="order-items">
                                <strong>Items:</strong>
                                <ul>
                                    {order.item_names && order.item_names.map((item, index) => (
                                        <li key={index} className="order-item-details">
                                            <span className="item-name">{item}</span>
                                            <span className="item-prices">${order.item_prices[index].toFixed(2)}</span>
                                            <span className="item-quantitys">Qty: {order.item_quantities[index]}</span>
                                            <span className="item-prices">Total Price ${order.item_prices[index] * order.item_quantities[index].toFixed(2)}</span>
                                        </li>
                                    ))}
                                </ul>
                            </div>
                            {order.description && (
                                <div className="order-description">
                                    <strong>Description : </strong>
                                    <span>{order.description}</span>
                                </div>
                            )}
                            {userRole === '2' && order.status !== 'completed' && (
                                <button onClick={() => handleCompleteOrder(order.id)} className="complete-button" disabled={loadingOrderId === order.id}>
                                    {loadingOrderId === order.id ? 'Loading...' : '✔️ Mark as Completed'}
                                </button>
                            )}
                        </div>
                    ))}
                </div>
                {loading && <div className="spinner"></div>}
            </div>
        </div>
    );
};

export default Orders;