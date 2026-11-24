import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Layout } from './components/Layout';
import { HomePage } from './pages/HomePage';
import { UnderlyingAssetsPage } from './pages/UnderlyingAssetsPage';

export function App() {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/" element={<UnderlyingAssetsPage />} />
          <Route path="/underlyings" element={<UnderlyingAssetsPage />} />
          <Route path="/status" element={<HomePage />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
}

export default App;
