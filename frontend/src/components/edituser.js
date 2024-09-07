import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import '../css/editprofile.css';
import defaultImage from '../css/vendor.jpg';

function EditUser() {
  const { userId } = useParams();
  const navigate = useNavigate();
  const [image, setImage] = useState(null);
  const [preview, setPreview] = useState(null);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [phone, setPhone] = useState('');
  const [errorMessages, setErrorMessages] = useState({
    name: '',
    email: '',
    password: '',
    phone: '',
    image: '',
    general: ''
  });
  const [loading, setLoading] = useState(true);
  const [success, setSuccess] = useState(null);
  const [successColor, setSuccessColor] = useState(''); // Add state for success message color
  const imageRef = useRef(null);

  useEffect(() => {
    const fetchUserDetails = async () => {
      setLoading(true);
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          throw new Error('No token found');
        }

        const response = await fetch(`http://localhost:8080/users/${userId}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          const errorData = await response.json();
          console.error('Backend errors:', errorData);
      
          const newErrors = {
            name: errorData.error?.name || '',
            email: errorData.error?.email || '',
            password: errorData.error?.password || '',
            phone: errorData.error?.phone || '',
            image: '',
            general: '',
          };
      
          if (response.status === 409 && errorData.error === 'Email already exists, try something else') {
            newErrors.email = 'Email already exists, try something else.';
          }
      
          setErrorMessages(newErrors);
      
          throw new Error('Failed to update profile');
        }

        const data = await response.json();
        console.log('Fetched user data:', data);

        if (data.user) {
          const { name, email, phone, img } = data.user;
          setName(name || '');
          setEmail(email || '');
          setPhone(phone || '');
          setPreview(img || defaultImage);
        } else {
          console.error('No user data in response:', data);
          setErrorMessages(prev => ({ ...prev, general: 'User data is missing in response' }));
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
    });
  
    let hasErrors = false;
  
    // Validate Name
    if (!name.trim()) {
      setErrorMessages(prev => ({ ...prev, name: 'Name is required.' }));
      hasErrors = true;
    }
  
    // Validate Email
    if (!email.trim()) {
      setErrorMessages(prev => ({ ...prev, email: 'Email is required.' }));
      hasErrors = true;
    }
  
    // Validate Phone
    if (!phone.trim()) {
      setErrorMessages(prev => ({ ...prev, phone: 'Phone number is required.' }));
      hasErrors = true;
    }
  
    // Validate if image has errors
    if (errorMessages.image) {
      hasErrors = true;
    }
  
    if (hasErrors) {
      setSuccess(null);
      return;
    }
  
    // Create FormData object
    const formData = new FormData();
    formData.append('name', name);
    formData.append('email', email);
    formData.append('phone', phone);
  
    // Only append password if it's not empty
    if (password.trim()) {
      formData.append('password', password);
    }
  
    // Append image if the user selected one
    if (image) {
      formData.append('img', image);
    }
  
    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`http://localhost:8080/users/${userId}`, {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      });
  
      if (!response.ok) {
        const errorData = await response.json();
        console.error('Backend errors:', errorData);
  
        const newErrors = {
          name: errorData.error?.name || '',
          email: errorData.error?.email || '',
          password: errorData.error?.password || '',
          phone: errorData.error?.phone || '',
          image: '',
          general: '',
        };
  
        if (response.status === 409 && errorData.error === 'Email already exists, try something else') {
          newErrors.email = 'Email already exists, try something else.';
        }
  
        setErrorMessages(newErrors);
        throw new Error('Failed to update profile');
      }
  
      setSuccess('Profile updated successfully');
      setSuccessColor('green'); // Set success color to green
      setTimeout(() => setSuccess(null), 2000); // Hide success message after 2 seconds
      setTimeout(() => navigate(`/users`), 2000); // Optional: redirect after a short delay
    } catch (error) {
      console.error('Error updating profile:', error);
      setErrorMessages(prev => ({ ...prev, general: 'Failed to update profile' }));
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      try {
        const token = localStorage.getItem('token');
        const response = await fetch(`http://localhost:8080/users/${userId}`, {
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

        setSuccess('User deleted successfully');
        setSuccessColor('red'); // Set success color to red
        setTimeout(() => {
          setSuccess(null);
          navigate('/users'); // Redirect to users list after deletion
        }, 2000); // Optional: redirect after a short delay
      } catch (error) {
        console.error('Error deleting user:', error);
        setErrorMessages(prev => ({ ...prev, general: 'Failed to delete user' }));
      }
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

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
        style={{ display: 'none' }}
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
          />
          {errorMessages.password && (
            <div className="error-message">{errorMessages.password}</div>
          )}
        </div>
        <div className={`form-group ${errorMessages.phone ? 'error' : ''}`}>
          <label htmlFor="phone">Phone:</label>
          <input
            type="text"
            id="phone"
            value={phone}
            onChange={(e) => setPhone(e.target.value)}
          />
          {errorMessages.phone && (
            <div className="error-message">{errorMessages.phone}</div>
          )}
        </div>
        {errorMessages.general && (
          <div className="error-message">{errorMessages.general}</div>
        )}
        {success && (
          <div className="success-message" style={{ color: successColor }}>
            {success}
          </div>
        )}
        <div className="form-buttons">
          <button type="submit" className="btn btn-primary">Save</button>
          <button type="button" id="my-delete-btn" className="delete-button" onClick={handleDelete}>                        Delete
          </button>
        </div>
      </form>
    </div>
  </div>
);
}

export default EditUser;
