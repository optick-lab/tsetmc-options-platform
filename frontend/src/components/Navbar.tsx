import { NavLink, useLocation } from 'react-router-dom';
import './Navbar.css';

export function Navbar() {
  const location = useLocation();
  const isUnderlyingsActive =
    location.pathname === '/' || location.pathname === '/underlyings';

  return (
    <header className="navbar">
      <div className="navbar-container">
        <div className="navbar-brand">
          <span className="brand-title">TSETMC Options</span>
          <span className="brand-tag">Tick Data</span>
        </div>
        <nav className="navbar-nav">
          <NavLink
            to="/"
            className={`nav-link ${isUnderlyingsActive ? 'active' : ''}`}
          >
            Underlying Assets
          </NavLink>
          <NavLink
            to="/status"
            className={({ isActive }) =>
              `nav-link ${isActive ? 'active' : ''}`
            }
          >
            System Status
          </NavLink>
        </nav>
      </div>
    </header>
  );
}

export default Navbar;
