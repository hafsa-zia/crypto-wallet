import { useEffect, useState } from "react";
import api from "../api/client";
import { useAuth } from "../context/AuthContext";

type ReportSummary = {
  total_sent_amount: number;
  total_sent_count: number;
  total_received_amount: number;
  total_received_count: number;
  zakat_deducted_amount: number;
  zakat_deducted_tx_count: number;
};

export default function Reports() {
  const { user } = useAuth();

  const [summary, setSummary] = useState<ReportSummary | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [zakatMsg, setZakatMsg] = useState("");
  const [zakatBusy, setZakatBusy] = useState(false);

  useEffect(() => {
    const fetchSummary = async () => {
      try {
        setLoading(true);
        setError("");
        const res = await api.get("/reports/summary");
        setSummary(res.data as ReportSummary);
      } catch (err: any) {
        console.error("Reports error:", err);
        setError(
          err?.response?.data?.error ||
            "Failed to load reports summary. Check backend /reports/summary."
        );
      } finally {
        setLoading(false);
      }
    };
    fetchSummary();
  }, []);

  const handleRunZakat = async () => {
  try {
    setZakatBusy(true);
    setZakatMsg("");
    setError("");

    const res = await api.post("/zakat/run-self");
    const msg =
      res.data?.message ||
      `Zakat run completed.`;

    setZakatMsg(msg);

    // reload summary to reflect updated zakat
    const sumRes = await api.get("/reports/summary");
    setSummary(sumRes.data as ReportSummary);
  } catch (err: any) {
    console.error("Run zakat error:", err);
    setError(
      err?.response?.data?.error ||
        "Failed to run zakat. Check backend /zakat/run-self."
    );
  } finally {
    setZakatBusy(false);
  }
};


  const formatAmount = (val?: number) =>
    typeof val === "number" ? val.toFixed(4) : "0.0000";

  return (
    <div className="main-column">
      {/* Header */}
      <div className="glass-card-soft">
        <div className="flex items-start justify-between gap-4">
          <div>
            <h1 className="text-xl font-semibold text-white">Wallet Reports</h1>
            <p className="text-xs text-slate-400 mt-2 leading-relaxed">
              Overview of your total sent, received and zakat deductions for
              wallet{" "}
              <span className="font-mono text-[11px] text-slate-300 break-all">
                {user?.wallet_id}
              </span>
              .
            </p>
          </div>

          <div className="flex flex-col items-end gap-1">
            <button
              onClick={handleRunZakat}
              disabled={zakatBusy}
              className="btn-primary text-[11px] px-3 py-1.5"
            >
              {zakatBusy ? "Running Zakat..." : "Run Zakat Now"}
            </button>
            {zakatMsg && (
              <span className="text-[11px] text-emerald-300 max-w-xs text-right">
                {zakatMsg}
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

      {/* Top summary cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <div className="glass-card">
          <h2 className="text-[11px] font-medium text-slate-400 mb-2">
            Total Sent
          </h2>
          <div className="text-xl font-semibold text-rose-400">
            {loading ? "..." : formatAmount(summary?.total_sent_amount)}{" "}
            <span className="text-[11px] text-slate-500 ml-1">CWD</span>
          </div>
          <p className="text-[11px] text-slate-500 mt-2 leading-relaxed">
            {loading
              ? ""
              : `${summary?.total_sent_count ?? 0} outgoing transaction${
                  (summary?.total_sent_count ?? 0) === 1 ? "" : "s"
                }`}
          </p>
        </div>

        <div className="glass-card">
          <h2 className="text-[11px] font-medium text-slate-400 mb-2">
            Total Received
          </h2>
          <div className="text-xl font-semibold text-emerald-400">
            {loading ? "..." : formatAmount(summary?.total_received_amount)}{" "}
            <span className="text-[11px] text-slate-500 ml-1">CWD</span>
          </div>
          <p className="text-[11px] text-slate-500 mt-2 leading-relaxed">
            {loading
              ? ""
              : `${summary?.total_received_count ?? 0} incoming transaction${
                  (summary?.total_received_count ?? 0) === 1 ? "" : "s"
                }`}
          </p>
        </div>

        <div className="glass-card">
          <h2 className="text-[11px] font-medium text-slate-400 mb-2">
            Zakat Deducted
          </h2>
          <div className="text-xl font-semibold text-sky-400">
            {loading ? "..." : formatAmount(summary?.zakat_deducted_amount)}{" "}
            <span className="text-[11px] text-slate-500 ml-1">CWD</span>
          </div>
          <p className="text-[11px] text-slate-500 mt-2 leading-relaxed">
            {loading
              ? ""
              : `${summary?.zakat_deducted_tx_count ?? 0} zakat transaction${
                  (summary?.zakat_deducted_tx_count ?? 0) === 1 ? "" : "s"
                }`}
          </p>
        </div>
      </div>

      {/* Breakdown panel */}
      <div className="glass-card-soft">
        <h2 className="text-sm font-semibold text-slate-200 mb-3">
          Breakdown &amp; Notes
        </h2>
        <div className="grid gap-4 md:grid-cols-2 text-xs text-slate-400">
          <div className="space-y-2">
            <p>
              <span className="text-slate-300">Net flow: </span>
              {loading || !summary
                ? "..."
                : (() => {
                    const net =
                      (summary?.total_received_amount || 0) -
                      (summary?.total_sent_amount || 0) -
                      (summary?.zakat_deducted_amount || 0);
                    const prefix = net >= 0 ? "+" : "";
                    return `${prefix}${net.toFixed(4)} CWD`;
                  })()}
            </p>
            <p>
              <span className="text-slate-300">Total activity: </span>
              {loading || !summary
                ? "..."
                : (summary.total_sent_count || 0) +
                  (summary.total_received_count || 0) +
                  (summary.zakat_deducted_tx_count || 0)}{" "}
              transactions
            </p>
          </div>
          <div className="space-y-1 text-[11px] text-slate-500 leading-relaxed">
            <p>
              This report is generated directly from your confirmed transactions
              stored on the custom blockchain and UTXO model, not from a cached
              balance. It&apos;s safe to show in your presentation as
              &quot;on-chain analytics&quot;.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
