import type { FormEvent } from "react";
import { useState } from "react";
import { Link } from "react-router-dom";
import api from "../api/client";

export default function Register() {
  const [form, setForm] = useState({
    full_name: "",
    email: "",
    password: "",
    cnic: "",
    otp: "",          // ✅ OTP in state
  });
  const [msg, setMsg] = useState("");
  const [error, setError] = useState("");
  const [otpStatus, setOtpStatus] = useState("");
  const [sendingOtp, setSendingOtp] = useState(false);

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSendOtp = async () => {
    if (!form.email) {
      setOtpStatus("Please enter your email first.");
      return;
    }
    try {
      setSendingOtp(true);
      setOtpStatus("");
      setError("");
      const res = await api.post("/auth/request-otp", {
        email: form.email,
      });
      setOtpStatus(
        res.data.message || "OTP sent. Check your email / backend logs."
      );
    } catch (err: any) {
      const msg =
        err?.response?.data?.error || "Failed to send OTP. Check backend.";
      setOtpStatus(msg);
    } finally {
      setSendingOtp(false);
    }
  };

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setMsg("");
    setError("");
    try {
      // ✅ sends full form including otp
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
    <div className="auth-shell bg-transparent">
      <div className="auth-card px-4">
        <div className="glass-card-soft space-y-6">
          <div className="text-center space-y-2">
            <h1 className="text-2xl font-bold text-white">Create Wallet</h1>
            <p className="text-slate-400 text-xs">
              Sign up with email + OTP to get your decentralized wallet ID.
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

          {otpStatus && (
            <div className="text-[11px] text-sky-300 bg-sky-950/40 border border-sky-800 px-3 py-2 rounded-lg">
              {otpStatus}
            </div>
          )}

          <form onSubmit={onSubmit} className="space-y-3">
            {/* Full name */}
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

            {/* Email + Send OTP button */}
            <div>
              <label className="block text-[11px] font-medium text-slate-300 mb-1">
                Email
              </label>
              <div className="flex gap-2">
                <input
                  type="email"
                  name="email"
                  placeholder="you@example.com"
                  className="flex-1 rounded-lg bg-slate-900 border border-slate-700 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                  value={form.email}
                  onChange={onChange}
                  required
                />
                <button
                  type="button"
                  onClick={handleSendOtp}
                  disabled={sendingOtp}
                  className="btn-primary text-[11px] px-3 py-2"
                >
                  {sendingOtp ? "Sending..." : "Send OTP"}
                </button>
              </div>
            </div>

            {/* Password */}
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

            {/* CNIC */}
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

            {/* OTP input */}
            <div>
              <label className="block text-[11px] font-medium text-slate-300 mb-1">
                OTP (from email)
              </label>
              <input
                name="otp"                    // ✅ must match registerRequest.OTP
                placeholder="6-digit code"
                className="w-full rounded-lg bg-slate-900 border border-slate-700 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                value={form.otp}
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
            <Link
              to="/login"
              className="text-emerald-400 hover:text-emerald-300 underline underline-offset-4"
            >
              Login
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
