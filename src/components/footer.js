// src/components/Footer.js

import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faFacebook, faTelegram, faTwitter as faX, faGoogle } from '@fortawesome/free-brands-svg-icons'; // Import the Twitter (X) icon
import '../css/footer.css'; // Import the CSS file

const Footer = () => {
  return (
    <footer className="footer">
      <div className="footer-description">
        <h2>أحمد علي</h2>
        <p>طالب في جامعة بنغازي فرع المرج ومبرمج backend</p>
      </div>
      <div className="footer-icons">
        <div className="footer-icon">
          <FontAwesomeIcon icon={faFacebook} />
          <a href="https://www.facebook.com/oGJughead" target="_blank" rel="noopener noreferrer">oGJughead</a>
        </div>
        <div className="footer-icon">
          <FontAwesomeIcon icon={faTelegram} />
          <a href="https://t.me/SaturnRings1" target="_blank" rel="noopener noreferrer">SaturnRings1</a>
        </div>
        <div className="footer-icon">
          <FontAwesomeIcon icon={faX} style={{ fontSize: '24px' }} /> {/* Set explicit size for the X icon */}
          <a href="https://twitter.com/og_jughead" target="_blank" rel="noopener noreferrer">@og_jughead</a>
        </div>
        <div className="footer-icon">
          <FontAwesomeIcon icon={faGoogle} />
          <a href="mailto:sasukesasuke755@gmail.com" target="_blank" rel="noopener noreferrer">sasukesasuke755@gmail.com</a>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
