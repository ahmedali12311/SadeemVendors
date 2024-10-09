import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import '../css/vendor.css';
import defaultImage from './vendor.jpg';

function Vendors() {
  const [vendors, setVendors] = useState([]);
  const [filteredVendors, setFilteredVendors] = useState([]);
  const [error, setError] = useState(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(false);
  const [isAdmin, setIsAdmin] = useState(false);
  const [userVendorIds, setUserVendorIds] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortOrder, setSortOrder] = useState('latest');
  const [showDropdown, setShowDropdown] = useState(false);
  const [userData] = useState(null);
  const navigate = useNavigate();
  const [showSearchInput] = useState(false);
  const itemsPerPage = 12;
  const visiblePages = 4;

  // Fetch vendors
  useEffect(() => {
    const fetchVendors = async () => {
      setLoading(true);
      const token = localStorage.getItem('token');
      try {
        const response = await fetch(`http://localhost:8080/vendors?page=${page}&pageSize=${itemsPerPage}&sort=${sortOrder}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          if (response.status === 401) {
            navigate('/signin', { replace: true });
            setError('Unauthorized. Please sign in.');
          }
          throw new Error(`Error! Status: ${response.status}`);
        }

        const data = await response.json();

        if (data && data.Vendors) {
          let vendorsList = data.Vendors;

          if (userVendorIds.length > 0 && !isAdmin) {
            vendorsList = vendorsList.filter(vendor => userVendorIds.includes(vendor.id));
          } else if (userVendorIds.length > 0) {
            vendorsList.sort((a, b) => {
              const isAUserVendor = userVendorIds.includes(a.id);
              const isBUserVendor = userVendorIds.includes(b.id);
              if (isAUserVendor && !isBUserVendor) return -1;
              if (!isAUserVendor && isBUserVendor) return 1;
              return 0;
            });
          }

          setVendors(vendorsList);
          setFilteredVendors(vendorsList);
          setTotalPages(Math.ceil(data.TotalCount / itemsPerPage));
        }
      } catch (error) {

        navigate('/signin', { replace: true });

        console.error('Error fetching vendors:', error);
        setError('Failed to load vendors');
      }
      setLoading(false);
    };

    fetchVendors();
  }, [page, sortOrder, userVendorIds, isAdmin, navigate]);

  // Fetch user data
  useEffect(() => {
    const fetchUser = async () => {
      const token = localStorage.getItem('token');
      if (!token) {
        setIsAdmin(false);
        setUserVendorIds([]);
        return;
      }

      try {
        const response = await fetch('http://localhost:8080/me', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          if (response.status === 401) {
            setError('Unauthorized. Please sign in.');
            setIsAdmin(false);
            setUserVendorIds([]);
          }
          throw new Error(`Error! Status: ${response.status}`);
        }

        const userData = await response.json();
        if (userData && userData.me) {
          setIsAdmin(userData.me.user_role === "1");

          if (userData.me.user_role === "2") {
            try {
              const vendorResponse = await fetch(`http://localhost:8080/uservendors/${userData.me.user_info.id}`, {
                headers: {
                  Authorization: `Bearer ${token}`,
                },
              });
              const vendorData = await vendorResponse.json();
              if (vendorData && vendorData.vendor && vendorData.vendor.length > 0) {
                setUserVendorIds(vendorData.vendor.map(v => v.id));
              } else {
                setUserVendorIds([]);
              }
            } catch (error) {
              setError('Failed to load vendor data');
            }
          }
        } else {
          setError('User data is missing or incomplete.');
        }
      } catch (error) {
        setError('Failed to load user data');
      }
    };

    fetchUser();
  }, [navigate]);

  const handlePageChange = (newPage) => {
    if (newPage > 0 && newPage <= totalPages) {
      setPage(newPage);
    }
  };

  const handleEditClick = (vendorId) => {
    if (isAdmin || (userVendorIds.length > 0 && userVendorIds.includes(vendorId))) {
      navigate(`/edit-vendor/${vendorId}`);
    } else {
      console.log('Non-admin user cannot edit vendors');
    }
  };

  const handleSearchChange = (event) => {
    const searchValue = event.target.value;
    setSearchTerm(searchValue);

    if (searchValue) {
      const results = vendors.filter(vendor => vendor.name.toLowerCase().includes(searchValue.toLowerCase()));
      setFilteredVendors(results);
      setShowDropdown(true);
    } else {
      setFilteredVendors([]);
      setShowDropdown(false);
    }
  };

  const handleSortChange = (event) => {
    setSortOrder(event.target.value);
  };

  const handleSelectVendor = (vendorId) => {
    setSearchTerm('');
    setShowDropdown(false);
    navigate(`/vendor/${vendorId}`);
  };

  const handleImageError = (e) => {
    e.target.src = defaultImage;
  };

  const handleAddVendorClick = () => {
    navigate('/add-vendor');
  };

  const handleVisitClick = (vendorId) => {
    navigate(`/vendor/${vendorId}`);
  };

  if (error) {
    return <div className="error-message">{error}</div>;
  }

  const currentPage = page;
  const startPage = Math.max(1, currentPage - visiblePages + 1);
  const endPage = Math.min(totalPages, currentPage + visiblePages - 1);

  const pages = [];
  for (let i = startPage; i <= endPage; i++) {
    pages.push(i);
  }

  return (
    <div className="page-container">
      <div className="vendor-list">
        <ul className="vendor-grid">
          <h1 className="title">Vendors</h1>

          <div className={`search-input-container ${showSearchInput ? 'show' : ''}`}>
            <input
              type="text"
              className="search-input"
              placeholder="Search vendors..."
              value={searchTerm}
              onChange={handleSearchChange}
              onFocus={() => searchTerm && setShowDropdown(true)}
              onBlur={() => setTimeout(() => setShowDropdown(false), 200)}
            />

            {showDropdown && searchTerm && filteredVendors.length > 0 && (
              <ul className="dropdown-menu">
                {filteredVendors.map((vendor) => (
                  <li
                    key={vendor.id}
                    onClick={() => handleSelectVendor(vendor.id)}
                    className="dropdown-item"
                  >
                    <img
                      src={vendor.img || defaultImage}
                      alt={vendor.name}
                      className="vendor-logo"
                      onError={handleImageError}
                    />
                    <span className="vendor-name">{vendor.name}</span>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </ul>

        <div className="sort-selection">
          <select value={sortOrder} onChange={handleSortChange}>
            <option value="latest">Latest</option>
            <option value="name_asc">Name Ascending</option>
            <option value="name_desc">Name Descending</option>
          </select>
        </div>

        {isAdmin && (
          <button className="add-vendor-button" onClick={handleAddVendorClick}>
            Add Vendor
          </button>
        )}

        <ul className="vendor-cards">
          {vendors.length > 0 ? (
            vendors.map((vendor) => (
              <li key={vendor.id} className="vendor-card">
                <div className="vendor-header">
                  <img
                    src={vendor.img || defaultImage}
                    alt={vendor.name}
                    className="vendor-logo"
                    onError={handleImageError}
                  />
                  <div className="vendor-content">
                    <h2 className="vendor-name">{vendor.name}</h2>
                    <p className="vendor-description">{vendor.description || 'No description available'}</p>
                    <button onClick={() => handleVisitClick(vendor.id)}>Visit</button>
                    {(isAdmin || (userVendorIds && vendor.id === userVendorIds && userData)) && (
                      <button onClick={() => handleEditClick(vendor.id)}>Edit</button>
                    )}
                  </div>
                </div>
              </li>
            ))
          ) : (
            <li>No vendors found</li>
          )}
        </ul>

        <div className="pagination-controls">
          <button onClick={() => handlePageChange(page - 1)} disabled={page === 1 || loading}>
            Previous
          </button>
          {pages.map((page) => (
            <button
              key={page}
              onClick={() => handlePageChange(page)}
              className={currentPage === page ? 'active' : ''}
            >
              {page}
            </button>
          ))}
          {endPage < totalPages && <button onClick={() => handlePageChange(endPage + 1)}>...</button>}
          <button onClick={() => handlePageChange(page + 1)} disabled={page === totalPages || loading}>
            Next
          </button>
          {loading && <div className="spinner"></div>}
        </div>
      </div>
    </div>
  );
}

export default Vendors;