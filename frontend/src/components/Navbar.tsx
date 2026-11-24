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
          <span className="brand-title">سامانه اختیار معامله TSETMC</span>
          <span className="brand-tag">داده‌های تیک</span>
        </div>
        <nav className="navbar-nav">
          <NavLink
            to="/"
            className={`nav-link ${isUnderlyingsActive ? 'active' : ''}`}
          >
            دارایی‌های پایه
          </NavLink>
          <NavLink
            to="/status"
            className={({ isActive }) =>
              `nav-link ${isActive ? 'active' : ''}`
            }
          >
            وضعیت سیستم
          </NavLink>
        </nav>
      </div>
    </header>
  );
}

export default Navbar;
