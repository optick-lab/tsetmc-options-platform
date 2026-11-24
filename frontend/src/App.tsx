import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Navbar } from './components/Navbar';
import { HomePage } from './pages/HomePage';
import { UnderlyingAssetsPage } from './pages/UnderlyingAssetsPage';

export function App() {
  return (
    <BrowserRouter>
      <Navbar />
      <Routes>
        <Route path="/" element={<UnderlyingAssetsPage />} />
        <Route path="/underlyings" element={<UnderlyingAssetsPage />} />
        <Route path="/status" element={<HomePage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
