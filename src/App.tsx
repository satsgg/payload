import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider } from '@/components/theme-provider';
import { Toaster } from '@/components/ui/sonner';
import VideoList from '@/components/video-list';
import VideoPlayer from '@/components/video-player';
import AdminDashboard from '@/components/admin/dashboard';
import AdminLogin from '@/components/admin/login';
import ProtectedRoute from '@/components/protected-route';

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <Router>
        <Routes>
          {/* Public Routes */}
          <Route path="/" element={<VideoList />} />
          <Route path="/watch/:videoId" element={<VideoPlayer />} />
          
          {/* Admin Routes */}
          <Route path="/admin/login" element={<AdminLogin />} />
          <Route 
            path="/admin/*" 
            element={
              <ProtectedRoute>
                <AdminDashboard />
              </ProtectedRoute>
            } 
          />
        </Routes>
      </Router>
      <Toaster />
    </ThemeProvider>
  );
}

export default App;