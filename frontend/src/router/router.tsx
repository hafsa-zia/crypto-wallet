import React from "react";
import {
  Routes,
  Route,
  Navigate,
  Outlet,
} from "react-router-dom";
import { useAuth } from "../context/AuthContext";

// pages (adjust paths if needed)
import Login from "../pages/Login";
import Register from "../pages/Register";
import Dashboard from "../pages/Dashboard";
import SendMoney from "../pages/SendMoney";
import TxHistory from "../pages/TxHistory";
import BlockExplorer from "../pages/BlockExplorer";
import Reports from "../pages/Reports";
import SystemLogs from "../pages/SystemLogs";
import WalletProfile from "../pages/WalletProfile";
import Beneficiaries from "../pages/Beneficiaries";

// sidebar layout
import Sidebar from "../components/Sidebar";

// Protect private routes
function PrivateRoute({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  if (!user) return <Navigate to="/login" replace />;
  return <>{children}</>;
}

// Main layout only for logged-in area
const MainLayout: React.FC = () => {
  return (
    <div className="app-shell">
      {/* Sidebar */}
      <aside className="sidebar-shell">
        <div>
          <div className="sidebar-logo">
            <span className="sidebar-logo-mark animate-glow-pulse">◎</span>
            CRYPTO WALLET
          </div>
          <Sidebar />
        </div>
        <div className="sidebar-footer">
          <span className="text-[10px] text-slate-500">
            Built with <span className="text-emerald-400">Blockchain</span> ·{" "}
            <span className="text-sky-400">UTXO</span> ·{" "}
            <span className="text-emerald-400">PoW</span>
          </span>
        </div>
      </aside>

      {/* Main content */}
      <main className="page-wrapper">
        <div className="page-inner">
          <Outlet />
        </div>
      </main>
    </div>
  );
};

const AppRouter: React.FC = () => {
  return (
    <Routes>
      {/* Default: open website -> go to Register */}
      <Route path="/" element={<Navigate to="/register" replace />} />

      {/* Public auth routes (no sidebar, no layout) */}
      <Route path="/register" element={<Register />} />
      <Route path="/login" element={<Login />} />

      {/* Private area: wrapped in PrivateRoute + MainLayout */}
      <Route
        element={
          <PrivateRoute>
            <MainLayout />
          </PrivateRoute>
        }
      >
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/send" element={<SendMoney />} />
        <Route path="/transactions" element={<TxHistory />} />
        <Route path="/blocks" element={<BlockExplorer />} />
        <Route path="/reports" element={<Reports />} />
        <Route path="/logs" element={<SystemLogs />} />
        <Route path="/wallet" element={<WalletProfile />} />
        <Route path="/beneficiaries" element={<Beneficiaries />} />
      </Route>

      {/* Fallback -> register */}
      <Route path="*" element={<Navigate to="/register" replace />} />
    </Routes>
  );
};

export default AppRouter;
