import { useEffect, useState } from "react";
import api from "../api/client";

type SystemLog = {
  id?: string;
  _id?: string;
  event?: string;
  details?: string;
  ip?: string;
  timestamp?: string;
};

export default function SystemLogs() {
  const [logs, setLogs] = useState<SystemLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchLogs = async () => {
      try {
        setLoading(true);
        setError("");

        const res = await api.get("/logs/system");
        const list = res.data.logs ?? res.data ?? [];
        setLogs(Array.isArray(list) ? list : []);
      } catch (err: any) {
        console.error("SystemLogs error:", err);
        setError(
          err?.response?.data?.error ||
            "Failed to load system logs. Check backend /logs/system."
        );
      } finally {
        setLoading(false);
      }
    };

    fetchLogs();
  }, []);

  return (
    <div className="main-column">
      {/* Header */}
      <div className="glass-card-soft">
        <h1 className="text-xl font-semibold text-white">System Logs</h1>
        <p className="text-xs text-slate-400 mt-2 leading-relaxed">
          Authentication attempts, invalid wallet ID checks, mining events, and
          other system-level actions.
        </p>
      </div>

      {error && (
        <div className="text-sm text-red-400 bg-red-950/40 border border-red-700 px-4 py-3 rounded-xl">
          {error}
        </div>
      )}

      {/* Logs table */}
      <div className="glass-card w-full overflow-x-auto">
        <table className="min-w-full text-xs">
          <thead className="bg-slate-950/60 border-b border-slate-800">
            <tr>
              <th className="px-4 py-2 text-left font-medium text-slate-400 text-[11px]">
                Time
              </th>
              <th className="px-4 py-2 text-left font-medium text-slate-400 text-[11px]">
                Event
              </th>
              <th className="px-4 py-2 text-left font-medium text-slate-400 text-[11px]">
                Details
              </th>
              <th className="px-4 py-2 text-left font-medium text-slate-400 text-[11px]">
                IP
              </th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td
                  colSpan={4}
                  className="px-4 py-4 text-center text-slate-500"
                >
                  Loading logs…
                </td>
              </tr>
            ) : logs.length === 0 ? (
              <tr>
                <td
                  colSpan={4}
                  className="px-4 py-4 text-center text-slate-500"
                >
                  No system logs yet.
                </td>
              </tr>
            ) : (
              logs.map((l, idx) => (
                <tr
                  key={l.id || l._id || idx}
                  className="border-t border-slate-800/70 hover:bg-slate-950/40"
                >
                  <td className="px-4 py-2 text-slate-300">
                    {l.timestamp
                      ? new Date(l.timestamp).toLocaleString()
                      : "-"}
                  </td>
                  <td className="px-4 py-2 text-slate-200">{l.event || "-"}</td>
                  <td className="px-4 py-2 text-slate-400 break-all">
                    {l.details || "-"}
                  </td>
                  <td className="px-4 py-2 text-slate-400">
                    {l.ip || "unknown"}
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
