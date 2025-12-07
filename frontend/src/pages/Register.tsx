import type { FormEvent } from "react";
import { useState } from "react";
import api from "../api/client";

export default function Register() {
  const [form, setForm] = useState({
    full_name: "",
    email: "",
    password: "",
    cnic: "",
  });
  const [msg, setMsg] = useState("");
  const [error, setError] = useState("");

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setMsg("");
    setError("");
    try {
      const res = await api.post("/auth/register", form);
      setMsg(res.data.message || "Registered successfully.");
    } catch (err: any) {
      const msg =
        err?.response?.data?.error ||
        "Registration failed. Check backend logs.";
      setError(msg);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-transparent">
      <div className="w-full max-w-xs px-4">
        <div className="glass-card-soft space-y-6">
          <div className="text-center space-y-2">
            <h1 className="text-2xl font-bold text-white">Create Wallet</h1>
            <p className="text-slate-400 text-xs">
              Sign up to get your decentralized wallet ID.
            </p>
          </div>

          {error && (
            <div className="text-sm text-red-400 bg-red-950/40 border border-red-700 px-3 py-2 rounded-lg">
              {error}
            </div>
          )}

          {msg && !error && (
            <div className="text-sm text-emerald-400 bg-emerald-950/40 border border-emerald-700 px-3 py-2 rounded-lg">
              {msg}
            </div>
          )}

          <form onSubmit={onSubmit} className="space-y-3">
            <div>
              <label className="block text-[11px] font-medium text-slate-300 mb-1">
                Full Name
              </label>
              <input
                name="full_name"
                placeholder="Hafsa Zia"
                className="w-full rounded-lg bg-slate-900 border border-slate-700 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                value={form.full_name}
                onChange={onChange}
                required
              />
            </div>

            <div>
              <label className="block text-[11px] font-medium text-slate-300 mb-1">
                Email
              </label>
              <input
                type="email"
                name="email"
                placeholder="you@example.com"
                className="w-full rounded-lg bg-slate-900 border border-slate-700 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                value={form.email}
                onChange={onChange}
                required
              />
            </div>

            <div>
              <label className="block text-[11px] font-medium text-slate-300 mb-1">
                Password
              </label>
              <input
                type="password"
                name="password"
                placeholder="At least 6 characters"
                className="w-full rounded-lg bg-slate-900 border border-slate-700 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                value={form.password}
                onChange={onChange}
                required
                minLength={6}
              />
            </div>

            <div>
              <label className="block text-[11px] font-medium text-slate-300 mb-1">
                CNIC / National ID
              </label>
              <input
                name="cnic"
                placeholder="12345-1234567-1"
                className="w-full rounded-lg bg-slate-900 border border-slate-700 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                value={form.cnic}
                onChange={onChange}
                required
              />
            </div>

            <button
              type="submit"
              className="btn-primary w-full py-2.5 text-sm"
            >
              Register &amp; Generate Wallet
            </button>
          </form>

          <p className="text-center text-[11px] text-slate-500">
            Already have an account?{" "}
            <a
              href="/login"
              className="text-emerald-400 hover:text-emerald-300 underline underline-offset-4"
            >
              Login
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}
