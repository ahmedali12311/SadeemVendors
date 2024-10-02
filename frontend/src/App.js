<<<<<<< HEAD
import React, { useState } from 'react';
=======
import React from 'react';
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
import { Routes, Route, BrowserRouter as Router, useLocation } from 'react-router-dom';
import Login from './components/login.js';
import Vendors from './components/vendors.js';
import VendorDetails from './components/vendordetails.js';
<<<<<<< HEAD
=======
import EditVendor from './components/editvendor.js';
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
import AddVendor from './components/addvendor.js';
import Navbar from './components/navbar.js';
import Sidebar from './components/sidebar.js';
import Profile from './components/editprofile.js';
import UsersPage from './components/userspage.js';
import EditUser from './components/edituser.js';
<<<<<<< HEAD
import Editvendorer from './components/editvendor.js';
import Orders from './components/orders.js';
import { OrderUpdateProvider } from './components/OrderUpdateContext'; // Import the new context provider

function App() {
  const [cartItems, setCartItems] = useState([]);
=======

function App() {
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
  const location = useLocation();

  // List of routes where Navbar and Sidebar should not be shown
  const noNavAndSidebarRoutes = ['/signin', '/signup'];
<<<<<<< HEAD
  const [refreshCart, setRefreshCart] = useState(false); // New state for triggering refresh

  const handleAddToCart = (item) => {
    setCartItems((prevCartItems) => [...prevCartItems, item]);
    setRefreshCart((prev) => !prev); // Toggle the state to trigger a refresh
};

=======
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c

  return (
    <div>
      {!noNavAndSidebarRoutes.includes(location.pathname) && (
        <>
<<<<<<< HEAD
            <Navbar cartItems={cartItems} onCartItemsChange={handleAddToCart} refreshCart={refreshCart} />
            <Sidebar />
=======
          <Navbar />
          <Sidebar />
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
        </>
      )}

      <div className="main-content">
        <Routes>
          <Route path="/signin" element={<Login />} />
          <Route path="/" element={<Vendors />} />
          <Route path="/vendors" element={<Vendors />} />
          <Route path="/profile" element={<Profile />} />
<<<<<<< HEAD
          <Route path="/orders" element={<Orders />} />

          <Route path="/vendor/:id" element={<VendorDetails onAddToCart={handleAddToCart}/>} />
          <Route path="/edit-vendor/:id" element={<Editvendorer />} />

          <Route path="/add-vendor" element={<AddVendor />} />
          <Route path="/users" element={<UsersPage />} />
          <Route path="/users/edit/:userId" element={<EditUser />} />
=======
          <Route path="/vendor/:id" element={<VendorDetails />} />
          <Route path="/edit-vendor/:id" element={<EditVendor />} />
          <Route path="/add-vendor" element={<AddVendor />} />
          <Route path="/users" element={<UsersPage />} />
          <Route path="/users/edit/:userId" element={<EditUser />} />


>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
        </Routes>
      </div>
    </div>
  );
}
<<<<<<< HEAD
export default function AppWrapper() {
  return (
    <Router>
        <OrderUpdateProvider>
        <App />
        </OrderUpdateProvider>
=======

export default function AppWrapper() {
  return (
    <Router>
      <App />
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
    </Router>
  );
}