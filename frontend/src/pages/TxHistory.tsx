import { useEffect, useState } from "react";
import api from "../api/client";

type Tx = {
  id?: string;
  _id?: string;
  sender_wallet?: string;
  receiver_wallet?: string;
  amount?: number;
  note?: string;
  timestamp?: string;
  type?: string;
  status?: string;
};

export default function TxHistory() {
  const [txs, setTxs] = useState<Tx[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [mining, setMining] = useState(false);
  const [mineMsg, setMineMsg] = useState("");

  const fetchTxs = async () => {
    try {
      setLoading(true);
      setError("");

      const res = await api.get("/tx/history");

      const list = res.data.transactions ?? res.data ?? [];
      setTxs(Array.isArray(list) ? list : []);
    } catch (err: any) {
      console.error("TxHistory error:", err);
      setError(
        err?.response?.data?.error ||
          "Failed to load transaction history. Check backend."
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTxs();
  }, []);

  const handleMinePending = async () => {
    try {
      setMining(true);
      setMineMsg("");
      setError("");

      const res = await api.post("/admin/mine");
      const data = res.data || {};

      setMineMsg(
        `Mined block #${data.block_index ?? "?"} | reward ${
          data.reward_amount ?? "?"
        } | txs in block: ${data.tx_count ?? "?"}`
      );

      await fetchTxs();
    } catch (err: any) {
      console.error("Mine pending error:", err);
      setError(
        err?.response?.data?.error ||
          "Failed to mine pending transactions. Check backend /admin/mine."
      );
    } finally {
      setMining(false);
    }
  };

  return (
    <div className="main-column">
      {/* Header */}
      <div className="glass-card-soft">
        <div className="flex items-start justify-between gap-3">
          <div>
            <h1 className="text-xl font-semibold text-white">
              Transaction History
            </h1>
            <p className="text-xs text-slate-400 mt-2 leading-relaxed">
              All sent, received and zakat transactions involving your wallet.
            </p>
          </div>

          <div className="flex flex-col items-end gap-1">
            <button
              onClick={handleMinePending}
              disabled={mining}
              className="btn-primary px-3 py-2 text-[11px]"
            >
              {mining ? "Mining pending..." : "Mine Pending Transactions"}
            </button>
            {mineMsg && (
              <span className="text-[11px] text-emerald-300 text-right leading-snug">
                {mineMsg}
              </span>
            )}
          </div>
        </div>
      </div>

      {error && (
        <div className="text-sm text-red-400 bg-red-950/40 border border-red-700 px-4 py-3 rounded-xl">
          {error}
        </div>
      )}

      {/* Table card */}
      <div className="glass-card w-full overflow-x-auto">
        <table className="min-w-full text-xs">
          <thead className="bg-slate-950/60 border-b border-slate-800">
            <tr>
              <th className="px-4 py-2 text-left font-medium text-slate-400 text-[11px]">
                Time
              </th>
              <th className="px-4 py-2 text-left font-medium text-slate-400 text-[11px]">
                Sender
              </th>
              <th className="px-4 py-2 text-left font-medium text-slate-400 text-[11px]">
                Receiver
              </th>
              <th className="px-4 py-2 text-right font-medium text-slate-400 text-[11px]">
                Amount
              </th>
              <th className="px-4 py-2 text-center font-medium text-slate-400 text-[11px]">
                Type
              </th>
              <th className="px-4 py-2 text-center font-medium text-slate-400 text-[11px]">
                Status
              </th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td
                  colSpan={6}
                  className="px-4 py-4 text-center text-slate-500"
                >
                  Loading transactions…
                </td>
              </tr>
            ) : txs.length === 0 ? (
              <tr>
                <td
                  colSpan={6}
                  className="px-4 py-4 text-center text-slate-500"
                >
                  No transactions found yet.
                </td>
              </tr>
            ) : (
              txs.map((t, idx) => (
                <tr
                  key={t.id || t._id || idx}
                  className="border-t border-slate-800/70 hover:bg-slate-950/40"
                >
                  <td className="px-4 py-2 text-slate-300">
                    {t.timestamp
                      ? new Date(t.timestamp).toLocaleString()
                      : "-"}
                  </td>
                  <td className="px-4 py-2 font-mono text-[10px] text-slate-400 break-all">
                    {t.sender_wallet || "-"}
                  </td>
                  <td className="px-4 py-2 font-mono text-[10px] text-slate-400 break-all">
                    {t.receiver_wallet || "-"}
                  </td>
                  <td className="px-4 py-2 text-right text-slate-100">
                    {t.amount ?? 0}
                  </td>
                  <td className="px-4 py-2 text-center">
                    <span className="inline-flex items-center rounded-full bg-slate-800/80 px-2 py-0.5 text-[10px] uppercase tracking-wide text-slate-300">
                      {t.type || "normal"}
                    </span>
                  </td>
                  <td className="px-4 py-2 text-center">
                    <span
                      className={`inline-flex items-center rounded-full px-2 py-0.5 text-[10px] uppercase tracking-wide ${
                        t.status === "confirmed"
                          ? "bg-emerald-900/60 text-emerald-300"
                          : t.status === "pending"
                          ? "bg-amber-900/60 text-amber-300"
                          : "bg-slate-800/80 text-slate-300"
                      }`}
                    >
                      {t.status || "unknown"}
                    </span>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
