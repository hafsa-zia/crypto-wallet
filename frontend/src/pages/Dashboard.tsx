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
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [mining, setMining] = useState(false);
  const [mineMsg, setMineMsg] = useState("");

  const loadData = async () => {
    try {
      setError("");
      setLoading(true);

      const [balanceRes, utxoRes] = await Promise.all([
        api.get("/wallet/balance"),
        api.get("/wallet/utxos"),
      ]);

      setBalance(balanceRes.data.balance ?? 0);

      const u = utxoRes.data.utxos ?? [];
      setUtxos(Array.isArray(u) ? u : []);
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

  return (
      <div className="flex flex-col gap-6 max-w-3xl w-full mx-auto">
      {/* Header */}
      <div className="glass-card-soft px-6 py-5">
        <h1 className="text-2xl font-semibold text-white">Dashboard</h1>
        <p className="text-xs text-slate-400 mt-2">
          Welcome,{" "}
          <span className="font-medium text-slate-100">{user?.full_name}</span>{" "}
          <span className="text-[11px] text-slate-500 break-all">
            ({user?.wallet_id})
          </span>
        </p>
      </div>

      {error && (
        <div className="text-sm text-red-400 bg-red-950/40 border border-red-700 px-4 py-3 rounded-xl">
          {error}
        </div>
      )}

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

      {/* Quick stats – full width below mining */}
      <div className="glass-card-soft px-6 py-5 w-full">
        <h2 className="text-xs font-medium text-slate-400 mb-3">
          Quick Stats
        </h2>
        <p className="text-xs text-slate-500 leading-relaxed">
          Later you can show total sent, total received, and zakat summary here,
          based on your on-chain reports.
        </p>
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
