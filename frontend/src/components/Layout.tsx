import { Outlet, Link, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function Layout() {
  const { logout, user } = useAuth();
  const loc = useLocation();

  const links = [
    { to: "/dashboard", label: "Dashboard" },
    { to: "/wallet", label: "Wallet" },
    { to: "/send", label: "Send Money" },
    { to: "/beneficiaries", label: "Beneficiaries" },
    { to: "/history", label: "Transactions" },
    { to: "/blocks", label: "Block Explorer" },
    { to: "/logs", label: "System Logs" },
    { to: "/reports", label: "Reports" },
  ];

  return (
    <div className="min-h-screen flex bg-slate-100">
      <aside className="w-60 bg-slate-900 text-white flex flex-col">
        <div className="p-4 font-semibold text-lg border-b border-slate-700">
          Crypto Wallet
        </div>
        <nav className="flex-1 px-2 py-4 space-y-1">
          {links.map((l) => (
            <Link
              key={l.to}
              to={l.to}
              className={`block px-3 py-2 rounded-lg text-sm ${
                loc.pathname === l.to
                  ? "bg-slate-700"
                  : "hover:bg-slate-800 text-slate-200"
              }`}
            >
              {l.label}
            </Link>
          ))}
        </nav>
        <div className="p-4 border-t border-slate-700 text-xs">
          <div className="mb-2">{user?.email}</div>
          <button
            className="w-full text-left text-red-300 hover:text-red-200"
            onClick={logout}
          >
            Logout
          </button>
        </div>
      </aside>
      <main className="flex-1">
        <div className="p-4 border-b bg-white">
          <h1 className="text-lg font-semibold">Decentralized Wallet</h1>
        </div>
        <div className="p-4">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
