import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Layout } from './components/Layout';
import { UnderlyingAssetsPage } from './pages/UnderlyingAssetsPage';

export function App() {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/" element={<UnderlyingAssetsPage />} />
          <Route path="/underlyings" element={<UnderlyingAssetsPage />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
}

export default App;
