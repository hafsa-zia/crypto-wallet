import { useState } from "react";
import type { FormEvent } from "react";
import { useAuth } from "../context/AuthContext";

export default function Login() {
  const { login } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    try {
      await login(email, password);
      window.location.href = "/dashboard";
    } catch (err: any) {
      setError(err?.response?.data?.error || "Login failed");
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-transparent">
      <div className="w-full max-w-xs px-4">
        <div className="glass-card-soft space-y-6">
          <div className="text-center space-y-2">
            <h1 className="text-2xl font-bold text-white">Welcome Back</h1>
            <p className="text-slate-400 text-xs">
              Login to access your decentralized wallet.
            </p>
          </div>

          {error && (
            <div className="text-sm text-red-400 bg-red-950/40 border border-red-700 px-3 py-2 rounded-lg">
              {error}
            </div>
          )}

          <form onSubmit={onSubmit} className="space-y-3">
            <div>
              <label className="block text-[11px] font-medium text-slate-300 mb-1">
                Email
              </label>
              <input
                type="email"
                className="w-full rounded-lg bg-slate-900 border border-slate-700 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                placeholder="you@example.com"
              />
            </div>
            <div>
              <label className="block text-[11px] font-medium text-slate-300 mb-1">
                Password
              </label>
              <input
                type="password"
                className="w-full rounded-lg bg-slate-900 border border-slate-700 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                placeholder="••••••••"
              />
            </div>
            <button type="submit" className="btn-primary w-full py-2.5 text-sm">
              Login
            </button>
          </form>

          <p className="text-center text-[11px] text-slate-500">
            Don&apos;t have an account?{" "}
            <a
              href="/register"
              className="text-emerald-400 hover:text-emerald-300 underline underline-offset-4"
            >
              Create wallet
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}
