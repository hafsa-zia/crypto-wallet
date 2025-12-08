import { useEffect, useState } from "react";
import api from "../api/client";
import { useAuth } from "../context/AuthContext";

type UTXO = {
  id?: string;
  _id?: string;
  amount: number;
};

export default function Dashboard() {
  const { user } = useAuth();

  const [balance, setBalance] = useState<number | null>(null);
  const [utxos, setUtxos] = useState<UTXO[]>([]);
  const [pendingCount, setPendingCount] = useState(0);
  const [sentCount, setSentCount] = useState(0);
  const [receivedCount, setReceivedCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [mining, setMining] = useState(false);
  const [mineMsg, setMineMsg] = useState("");

  // Profile editing state
  const [fullName, setFullName] = useState(user?.full_name || "");
  const [cnic, setCnic] = useState((user as any)?.cnic || "");
  const [email, setEmail] = useState(user?.email || "");
  const [profileMsg, setProfileMsg] = useState("");
  const [profileError, setProfileError] = useState("");
  const [savingProfile, setSavingProfile] = useState(false);

  // Email change (OTP) state
  const [newEmail, setNewEmail] = useState("");
  const [emailOTP, setEmailOTP] = useState("");
  const [emailStepMsg, setEmailStepMsg] = useState("");
  const [emailBusy, setEmailBusy] = useState(false);

  useEffect(() => {
    // sync local fields when user context changes
    setFullName(user?.full_name || "");
    setEmail(user?.email || "");
    setCnic((user as any)?.cnic || "");
  }, [user]);

  // If CNIC is empty from context, try to load from backend once
  useEffect(() => {
    const fetchProfileIfNeeded = async () => {
      if (cnic) return;
      try {
        const res = await api.get("/wallet"); // your wallet/profile endpoint
        if (res.data?.cnic) {
          setCnic(res.data.cnic);
        }
        if (res.data?.full_name && !fullName) {
          setFullName(res.data.full_name);
        }
        if (res.data?.email && !email) {
          setEmail(res.data.email);
        }
      } catch {
        // silently ignore if /wallet not available or doesn't return cnic
      }
    };
    fetchProfileIfNeeded();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const loadData = async () => {
    try {
      setError("");
      setLoading(true);

      // Fetch balance, utxos and tx history so we can compute quick stats
      const [balanceRes, utxoRes, txRes] = await Promise.all([
        api.get("/wallet/balance"),
        api.get("/wallet/utxos"),
        api.get("/tx/history"),
      ]);

      setBalance(balanceRes.data.balance ?? 0);

      const u = utxoRes.data.utxos ?? [];
      setUtxos(Array.isArray(u) ? u : []);

      // Compute quick stats from tx history
      const txs = Array.isArray(txRes.data.transactions)
        ? txRes.data.transactions
        : txRes.data.transactions ?? [];
      const walletID = user?.wallet_id;
      if (walletID) {
        const pending = txs.filter(
          (t: any) => t.status === "pending" && t.sender_wallet === walletID
        ).length;
        const sent = txs.filter((t: any) => t.sender_wallet === walletID).length;
        const received = txs.filter((t: any) => t.receiver_wallet === walletID).length;
        setPendingCount(pending);
        setSentCount(sent);
        setReceivedCount(received);
      } else {
        setPendingCount(0);
        setSentCount(0);
        setReceivedCount(0);
      }
    } catch (err: any) {
      console.error("Dashboard error:", err);
      setError(
        err?.response?.data?.error ||
          "Failed to load dashboard data. Check backend."
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleMineBlock = async () => {
    try {
      setMining(true);
      setMineMsg("");
      setError("");

      const res = await api.post("/admin/mine");
      const data = res.data || {};

      setMineMsg(
        `Mined block #${data.block_index ?? "?"} | reward ${
          data.reward_amount ?? "?"
        } | txs: ${data.tx_count ?? "?"}`
      );

      await loadData();
    } catch (err: any) {
      console.error("Mine error:", err);
      setError(
        err?.response?.data?.error ||
          "Failed to mine block. Make sure backend /api/admin/mine is working."
      );
    } finally {
      setMining(false);
    }
  };

  // ---- Save name + CNIC ----
  const handleSaveProfile = async () => {
    try {
      setSavingProfile(true);
      setProfileMsg("");
      setProfileError("");

      const payload: any = {
        full_name: fullName,
        cnic,
      };

      const res = await api.put("/profile", payload);
      setProfileMsg(res.data.message || "Profile updated.");
    } catch (err: any) {
      console.error("Profile update error:", err);
      setProfileError(
        err?.response?.data?.error ||
          "Failed to update profile. Check backend /profile."
      );
    } finally {
      setSavingProfile(false);
    }
  };

  // ---- Email change (OTP flow) ----
  const handleSendEmailOTP = async () => {
    try {
      setEmailBusy(true);
      setEmailStepMsg("");
      setProfileError("");

      if (!newEmail || !newEmail.includes("@")) {
        setProfileError("Enter a valid new email first.");
        return;
      }

      await api.post("/auth/request-otp", { email: newEmail });
      setEmailStepMsg(
        "OTP sent to the new email address. Enter it below to confirm the change."
      );
    } catch (err: any) {
      console.error("Send email OTP error:", err);
      setProfileError(
        err?.response?.data?.error || "Failed to send OTP to new email."
      );
    } finally {
      setEmailBusy(false);
    }
  };

  const handleConfirmEmailChange = async () => {
    if (!newEmail || !emailOTP) {
      setProfileError("New email and OTP are required.");
      return;
    }
    try {
      setEmailBusy(true);
      setProfileError("");
      setEmailStepMsg("");

      const payload: any = {
        email: newEmail,
        email_otp: emailOTP,
      };

      const res = await api.put("/profile", payload);
      setProfileMsg(res.data.message || "Email updated.");
      setEmail(newEmail);
      setNewEmail("");
      setEmailOTP("");
    } catch (err: any) {
      console.error("Confirm email change error:", err);
      setProfileError(
        err?.response?.data?.error || "Failed to update email with OTP."
      );
    } finally {
      setEmailBusy(false);
    }
  };

  return (
    <div className="main-column">
      {/* Header */}
      <div className="glass-card-soft px-6 py-5">
        <h1 className="text-2xl font-semibold text-white">Dashboard</h1>
        <p className="text-xs text-slate-400 mt-2">
          Welcome,{" "}
          <span className="font-medium text-slate-100">{fullName}</span>{" "}
          <span className="text-[11px] text-slate-500 break-all">
            (Wallet ID: {user?.wallet_id})
          </span>
        </p>
      </div>

      {/* Messages */}
      {(profileError || profileMsg || emailStepMsg) && (
        <div className="space-y-2">
          {profileError && (
            <div className="text-sm text-red-400 bg-red-950/40 border border-red-700 px-4 py-3 rounded-xl">
              {profileError}
            </div>
          )}
          {profileMsg && !profileError && (
            <div className="text-sm text-emerald-400 bg-emerald-950/40 border border-emerald-700 px-4 py-3 rounded-xl">
              {profileMsg}
            </div>
          )}
          {emailStepMsg && (
            <div className="text-xs text-sky-300 bg-sky-950/40 border border-sky-700 px-4 py-2 rounded-xl">
              {emailStepMsg}
            </div>
          )}
        </div>
      )}

      {error && (
        <div className="text-sm text-red-400 bg-red-950/40 border border-red-700 px-4 py-3 rounded-xl">
          {error}
        </div>
      )}

    {/* 🔹 Profile card – updated to match your new CSS */}
<div className="glass-card px-6 py-5 w-full dashboard-card">
  <h2 className="text-xs font-medium text-slate-400 mb-3">
    Wallet Profile
  </h2>

  {/* Grid using your new perfect alignment class */}
  <div className="dashboard-grid text-xs mt-1">

    {/* LEFT COLUMN */}
    <div className="space-y-3">
      
      {/* Full Name */}
      <div>
        <label className="block text-[11px] text-slate-400 mb-1">
          Full Name
        </label>
        <input
          className="dashboard-input"
          value={fullName}
          onChange={(e) => setFullName(e.target.value)}
        />
      </div>

      {/* CNIC */}
      <div>
        <label className="block text-[11px] text-slate-400 mb-1">
          CNIC / National ID
        </label>
        <input
          className="dashboard-input"
          value={cnic}
          onChange={(e) => setCnic(e.target.value)}
          placeholder="12345-1234567-1"
        />
      </div>

      {/* Wallet ID */}
      <div>
        <label className="block text-[11px] text-slate-400 mb-1">
          Wallet ID
        </label>
        <div className="dashboard-box font-mono">
          {user?.wallet_id || "—"}
        </div>
      </div>

    </div>

    {/* RIGHT COLUMN */}
    <div className="space-y-3">

      {/* Current Email */}
      <div>
        <label className="block text-[11px] text-slate-400 mb-1">
          Current Email
        </label>

        <div className="dashboard-box">
          {email || "—"}
        </div>

        <p className="text-[10px] text-slate-500 mt-1">
          Email is locked. To change it, enter a new email and confirm via OTP.
        </p>
      </div>

      {/* CHANGE EMAIL BOX */}
      <div className="dashboard-email-box space-y-2">

        <p className="text-[11px] text-slate-300 font-semibold">
          Change Email (OTP Required)
        </p>

        {/* New Email */}
        <div>
          <label className="block text-[11px] text-slate-400 mb-1">
            New Email
          </label>
          <input
            type="email"
            className="dashboard-input"
            value={newEmail}
            onChange={(e) => setNewEmail(e.target.value)}
            placeholder="newemail@example.com"
          />
        </div>

        {/* Send OTP */}
        <button
          type="button"
          onClick={handleSendEmailOTP}
          disabled={emailBusy || !newEmail}
          className="dashboard-btn btn-primary w-full"
        >
          {emailBusy ? "Sending OTP..." : "Send OTP to New Email"}
        </button>

        {/* OTP Input */}
        <div>
          <label className="block text-[11px] text-slate-400 mb-1">
            Enter OTP
          </label>
          <input
            className="dashboard-input"
            value={emailOTP}
            onChange={(e) => setEmailOTP(e.target.value)}
            placeholder="6-digit code"
          />
        </div>

        {/* Verify & Update */}
        <button
          type="button"
          onClick={handleConfirmEmailChange}
          disabled={emailBusy || !newEmail || !emailOTP}
          className="dashboard-btn btn-primary w-full"
        >
          {emailBusy ? "Updating Email..." : "Verify OTP & Update Email"}
        </button>

      </div>
    </div>

  </div>

  {/* SAVE BUTTON */}
  <div className="mt-4 flex justify-end save-actions">
    <button
      type="button"
      onClick={handleSaveProfile}
      disabled={savingProfile}
      className="dashboard-btn btn-primary"
    >
      {savingProfile ? "Saving..." : "Save Name & CNIC"}
    </button>
  </div>
</div>

      {/* Quick Stats card */}
      <div className="glass-card px-6 py-5 w-full">
        <h2 className="text-xs font-medium text-slate-400 mb-3">Quick Stats</h2>
        <div className="stats-grid">
          <div className="stat-box">
            <div className="text-xs text-slate-400">Balance</div>
            <div className="text-xl font-bold text-emerald-400">
              {balance !== null ? balance.toFixed(4) : loading ? "..." : "0.0000"}
            </div>
            <div className="text-[11px] text-slate-500">CWD</div>
          </div>

          <div className="stat-box">
            <div className="text-xs text-slate-400">UTXOs</div>
            <div className="text-xl font-bold text-sky-300">{utxos.length}</div>
            <div className="text-[11px] text-slate-500">Unspent outputs</div>
          </div>

          <div className="stat-box">
            <div className="text-xs text-slate-400">Pending</div>
            <div className="text-xl font-bold text-amber-400">{pendingCount}</div>
            <div className="text-[11px] text-slate-500">Pending transactions</div>
          </div>

          <div className="stat-box">
            <div className="text-xs text-slate-400">Sent</div>
            <div className="text-xl font-bold text-rose-400">{sentCount}</div>
            <div className="text-[11px] text-slate-500">Total sent</div>
          </div>

          <div className="stat-box">
            <div className="text-xs text-slate-400">Received</div>
            <div className="text-xl font-bold text-emerald-300">{receivedCount}</div>
            <div className="text-[11px] text-slate-500">Total received</div>
          </div>
        </div>
      </div>

      {/* Balance card – full width */}
      <div className="glass-card px-6 py-5 w-full">
        <h2 className="text-xs font-medium text-slate-400 mb-3">
          Wallet Balance
        </h2>
        <div className="text-3xl font-bold text-emerald-400">
          {balance !== null ? balance.toFixed(4) : loading ? "..." : "0.0000"}{" "}
          <span className="text-sm text-slate-500 ml-1">CWD</span>
        </div>
        <p className="text-[11px] text-slate-500 mt-3 leading-relaxed">
          Balance is computed using the UTXO model from all unspent outputs
          linked to your wallet.
        </p>
      </div>

      {/* Mining card – full width, stacked under balance */}
      <div className="glass-card px-6 py-5 w-full">
        <h2 className="text-xs font-medium text-slate-400 mb-3">
          Mining (Proof-of-Work)
        </h2>
        <p className="text-xs text-slate-500 mb-4 leading-relaxed">
          Mine a new block using Proof-of-Work. You&apos;ll receive a mining
          reward directly to your wallet as a UTXO.
        </p>
        <div className="space-y-3">
          <button
            onClick={handleMineBlock}
            disabled={mining}
            className="btn-primary w-full"
          >
            {mining ? "Mining..." : "Mine Block & Get Reward"}
          </button>
          {mineMsg && (
            <p className="text-[11px] text-emerald-300 bg-emerald-950/40 border border-emerald-700 px-3 py-2 rounded-lg leading-relaxed">
              {mineMsg}
            </p>
          )}
        </div>
      </div>

      
      {/* UTXO list – full width at bottom */}
      <div className="glass-card px-6 py-5 w-full">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-sm font-semibold text-slate-200">UTXO List</h2>
          {loading && (
            <span className="text-xs text-slate-500">Loading UTXOs…</span>
          )}
        </div>

        {utxos.length === 0 && !loading ? (
          <p className="text-xs text-slate-500">
            No UTXOs found yet. Try mining a block to get a reward.
          </p>
        ) : (
          <div className="space-y-2 max-h-64 overflow-y-auto pr-1">
            {utxos.map((u, idx) => (
              <div
                key={u.id || u._id || idx}
                className="flex justify-between items-center text-xs bg-slate-950/70 border border-slate-800 rounded-xl px-4 py-2.5"
              >
                <span className="font-mono text-[10px] text-slate-400 truncate max-w-xs break-all">
                  {u.id || u._id || "reward-utxo"}
                </span>
                <span className="text-emerald-400 font-semibold">
                  {u.amount}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
