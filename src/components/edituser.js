import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import '../css/editprofile.css';
import defaultImage from '../css/vendor.jpg';

function EditUser () {
  const { userId } = useParams();
  const navigate = useNavigate();
  const [image, setImage] = useState(null);
  const [preview, setPreview] = useState(null);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [phone, setPhone] = useState('');
  const [vendorId, setVendorId] = useState('');
  const [role, setRole] = useState('');
  const [currentRole, setCurrentRole] = useState(''); // Store current role
  const [errorMessages, setErrorMessages] = useState({
    name: '',
    email: '',
    password: '',
    phone: '',
    image: '',
    general: '',
    role: '',
    vendorId: ''
  });
  const [loading, setLoading] = useState(true);
  const [success, setSuccess] = useState(null);
  const [successColor, setSuccessColor] = useState('');
  const imageRef = useRef(null);
  const phonePrefix = '+2189';

  useEffect(() => {
    const fetchUserDetails = async () => {
      setLoading(true);
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          throw new Error('No token found');
        }

        const [userResponse, roleResponse] = await Promise.all([
          fetch(`https://backend-934694036821.europe-west1.run.app/users/${userId}`, {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }),
          fetch(`https://backend-934694036821.europe-west1.run.app/userroles/${userId}`, {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          })
        ]);

        if (!userResponse.ok || !roleResponse.ok) {
          const errorData = await userResponse.json() || {};
          console.error('Backend errors:', errorData);

          const newErrors = {
            name: errorData.error?.name || '',
            email: errorData.error?.email || '',
            password: errorData.error?.password || '',
            phone: errorData.error?.phone || '',
            image: '',
            general: '',
            role: '',
            vendorId: ''
          };

          if (userResponse.status === 409 && errorData.error === 'Email already exists, try something else') {
            newErrors.email = 'Email already exists, try something else.';
          }

          setErrorMessages(newErrors);
          throw new Error('Failed to fetch user details or role information');
        }

        const userData = await userResponse.json();
        const roleData = await roleResponse.json();

        if (userData.user && roleData.user_roles) {
          const { name, email, phone, img } = userData.user;
          const { roleID } = roleData.user_roles;

          setName(name || '');
          setEmail(email || '');
          setPhone(phone || '');
          setPreview(img || defaultImage);
          setRole(roleID);
          setCurrentRole(roleID); // Set current role for comparison
        } else {
          console.error('No user or role data in response');
          setErrorMessages(prev => ({ ...prev, general: 'User  or role data is missing in response' }));
        }
      } catch (error) {
        console.error('Error fetching user details:', error);
        setErrorMessages(prev => ({ ...prev, general: 'Failed to load user details' }));
      } finally {
        setLoading(false);
      }
    };

    fetchUserDetails();
  }, [userId]);

  const handleImageClick = (event) => {
    const file = event.target.files[0];
    if (file) {
      const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
      if (!validTypes.includes(file.type)) {
        setErrorMessages(prev => ({
          ...prev,
          image: 'Invalid image type. Please upload a JPEG, PNG, GIF, or WEBP image.'
        }));

        setImage(null);
        setPreview(prev => prev || defaultImage);

        setTimeout(() => {
          setErrorMessages(prev => ({ ...prev, image: '' }));
        }, 10000);

        return;
      }

      if (file.size > 2000000) { // 2MB
        setErrorMessages(prev => ({
          ...prev,
          image: 'Image size must be less than 2MB.'
        }));

        setImage(null);
        setPreview(prev => prev || defaultImage);

        setTimeout(() => {
          setErrorMessages(prev => ({ ...prev, image: '' }));
        }, 3000);

        return;
      }

      setErrorMessages(prev => ({ ...prev, image: '' }));
      setImage(file);

      const reader = new FileReader();
      reader.onloadend = () => {
        setPreview(reader.result);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleSave = async (event) => {
    event.preventDefault();

    setErrorMessages({
      name: '',
      email: '',
      password: '',
      phone: '',
      image: '',
      general: '',
      role: '',
      vendorId: ''
    });

    let hasErrors = false;

    if (!name.trim()) {
      setErrorMessages(prev => ({ ...prev, name: 'Name is required.' }));
      hasErrors = true;
    }

    if (!email.trim()) {
      setErrorMessages(prev => ({ ...prev, email: 'Email is required.' }));
      hasErrors = true;
    }

    if (!phone.trim()) {
      setErrorMessages(prev => ({ ...prev, phone: 'Phone number is required.' }));
      hasErrors = true;
    }

    if (errorMessages.image) {
      hasErrors = true;
    }

    if (hasErrors) {
      setSuccess(null);
      return;
    }

    const formData = new FormData();
    formData.append('name', name);
    formData.append('email', email);
    formData.append('phone', phonePrefix + phone);

    if (password.trim()) {
      formData.append('password', password);
    }

    if (image) {
      formData.append('img', image);
    }

    try {
      const token = localStorage.getItem('token');
      const userResponse = await fetch(`https://backend-934694036821.europe-west1.run.app/users/${userId}`, {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      });

      if (!userResponse.ok) {
        const errorData = await userResponse.json();
        console.error('Backend errors:', errorData);

        const newErrors = {
          name: errorData.error?.name || '',
          email: errorData.error?.email || '',
          password: errorData.error?.password || '',
          phone: errorData.error?.phone || '',
          image: '',
          general: '',
          role: '',
          vendorId: ''
        };

        if (userResponse.status === 409 && errorData.error === 'Email already exists, try something else') {
          newErrors.email = 'Email already exists, try something else.';
        }

        setErrorMessages(newErrors);
        throw new Error('Failed to update user details');
      }

      if (role !== currentRole) {
        const formDataRole = new URLSearchParams();
        formDataRole.append('role', role);
        formDataRole.append('vendorID', vendorId);

        const roleResponse = await fetch(`https://backend-934694036821.europe-west1.run.app/grantrole/${userId}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            Authorization: `Bearer ${token}`,
          },
          body: formDataRole.toString(),
        });

        if (!roleResponse.ok) {
          const errorData = await roleResponse.json();
          console.error('Backend errors:', errorData);
          setErrorMessages(prev => ({ ...prev, general: 'Failed to update user role' }));
          throw new Error(`Failed to update user role: ${errorData.error}`);
        }
      }

      setSuccess('Profile updated successfully');
      setSuccessColor('green');
      setTimeout(() => setSuccess(null), 2000);
      setTimeout(() => navigate('/users'), 1000);
    } catch (error) {
      console.error('Error updating profile or role:', error);
      setErrorMessages(prev => ({ ...prev, general: 'Failed to update profile or role' }));
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      try {
        const token = localStorage.getItem('token');
        const response = await fetch(`https://backend-934694036821.europe-west1.run.app/users/${userId}`, {
          method: 'DELETE',
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          const errorData = await response.json();
          console.error('Backend errors:', errorData);
          setErrorMessages(prev => ({ ...prev, general: 'Failed to delete user' }));
          return;
        }

        setSuccess('User  deleted successfully');
        setSuccessColor('red');
        setTimeout(() => {
          setSuccess(null);
          navigate('/users');
        }, 2000);
      } catch (error) {
        console.error('Error deleting user:', error);
        setErrorMessages(prev => ({ ...prev, general: 'Failed to delete user' }));
      }
    }
  };

  const handleImageError = (e) => {
    e.target.src = defaultImage;
  };

  return (
    <div className="profile-container">
      <div className="profile-image-container">
        <input
          type="file"
          ref={imageRef}
          onChange={handleImageClick}
          style ={{ display: 'none' }}
        />
        <img
          src={preview}
          alt={name}
          className={`profile-image ${!preview ? 'no-image' : ''}`}
          onError={handleImageError}
          onClick={() => imageRef.current?.click()}
        />
        {errorMessages.image && (
          <div className="error-message">{errorMessages.image}</div>
        )}
      </div>
      <div className="profile-info-container">
        <form onSubmit={handleSave}>
          <div className={`form-group ${errorMessages.name ? 'error' : ''}`}>
            <label htmlFor="name">Name:</label>
            <input
              type="text"
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              style={{ borderColor: errorMessages.name ? 'red' : '' }}
            />
            {errorMessages.name && (
              <div className="error-message">{errorMessages.name}</div>
            )}
          </div>
          <div className={`form-group ${errorMessages.email ? 'error' : ''}`}>
            <label htmlFor="email">Email:</label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              style={{ borderColor: errorMessages.email ? 'red' : '' }}
            />
            {errorMessages.email && (
              <div className="error-message">{errorMessages.email}</div>
            )}
          </div>
          <div className={`form-group ${errorMessages.password ? 'error' : ''}`}>
            <label htmlFor="password">Password:</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              style={{ borderColor: errorMessages.password ? 'red' : '' }}
            />
            {errorMessages.password && (
              <div className="error-message">{errorMessages.password}</div>
            )}
          </div>
          <div className={`form-group ${errorMessages.phone ? 'error' : ''}`}>
            <label htmlFor="phone">Phone:</label>
            <div className="phone-input-container">
              <span className="phone-prefix">{phonePrefix}</span>
              <input
                type="text"
                placeholder="Phone Number  21234567.."
                value={phone}
                onChange={(event) => setPhone(event.target.value)}
                maxLength={8}
              />
            </div>
            {errorMessages.phone && (
              <div className="error-message">{errorMessages.phone}</div>
            )}
          </div>

          <div className={`form-group ${errorMessages.role ? 'error' : ''}`}>
            <label htmlFor="role">Role:</label>
            <input
              type="text"
              id="role"
              value={role}
              onChange={(e) => setRole(e.target.value)}
              style={{ borderColor: errorMessages.role ? 'red' : '' }}
            />
            {errorMessages.role && (
              <div className="error-message">{errorMessages.role}</div>
            )}
          </div>
          {role === '2' && (
            <div className={`form-group ${errorMessages.vendorId ? 'error' : ''}`}>
              <label htmlFor="vendorId">Vendor ID:</label>
              <input
                type="text"
                id="vendorId"
                value={vendorId}
                onChange={(e) => setVendorId(e.target.value)}
                style={{ borderColor: errorMessages.vendorId ? 'red' : '' }}
              />
              {errorMessages.vendorId && (
                <div className="error-message">{errorMessages.vendorId}</div>
              )}
            </div>
          )}
          <button type="submit" className="save-button" disabled={loading}>Save Changes</button>
          {loading && <div className="spinner"></div>}
        </form>
        {errorMessages.general && (
          <div className="error-message">Error encountered while updating!</div>
        )}
        {success && (
          <div className="success-message" style={{ color: successColor }}>{success}</div>
        )}
        <button onClick={handleDelete} className="delete-button">Delete User</button>
      </div>

    </div>
  );
}

export default EditUser;