import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Layout } from './components/Layout';
import { UnderlyingAssetsPage } from './pages/UnderlyingAssetsPage';
import { RawTradesPage } from './pages/RawTradesPage';

export function App() {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/" element={<UnderlyingAssetsPage />} />
          <Route path="/underlyings" element={<UnderlyingAssetsPage />} />
          <Route path="/trades" element={<RawTradesPage />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
}

export default App;
