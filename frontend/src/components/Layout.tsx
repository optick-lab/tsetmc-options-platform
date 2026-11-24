import { ReactNode } from 'react';
import { NavLink } from 'react-router-dom';
import './Layout.css';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  return (
    <div className="app-shell" dir="rtl">
      <header className="app-header">
        <div className="header-container">
          <div className="header-topbar">
            <div className="brand">
              <span className="brand-dot" />
              <span className="brand-name">سامانه اختیار معامله TSETMC</span>
              <span className="brand-pill">داده‌های تیک</span>
            </div>
            <nav className="header-nav">
              <NavLink
                to="/"
                className="nav-pill active"
              >
                دارایی‌های پایه
              </NavLink>
            </nav>
          </div>

          <div className="header-hero">
            <h1 className="hero-title">دارایی‌های پایه</h1>
            <p className="hero-subtitle">
              ابزارهای بازار نقدی (سهام و صندوق‌ها) مشمول قراردادهای اختیار معامله
            </p>
          </div>
        </div>
      </header>

      <main className="page-sheet">
        <div className="sheet-container">
          {children}
        </div>
      </main>
    </div>
  );
}

export default Layout;
