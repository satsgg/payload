import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { ThemeProvider } from "@/components/theme-provider";
import { Toaster } from "@/components/ui/sonner";
import VideoList from "@/components/video-list";
import AdminDashboard from "@/components/admin/dashboard";
import AdminLogin from "@/components/admin/login";
import ProtectedRoute from "@/components/protected-route";
// import Login from "./pages/login";
import WatchPage from "./pages/watch";

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <Router>
        <Routes>
          {/* Public Routes */}
          <Route path="/" element={<VideoList />} />
          <Route path="/watch/:videoId" element={<WatchPage />} />
          {/* <Route path="/login" element={<Login />} /> */}

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
