import { useEffect, useState } from "react";
import api from "../api/client";

type Block = {
  index: number;
  hash: string;
  previous_hash: string;
  timestamp: string;
  nonce?: number;
  tx_count?: number;
};

export default function BlockExplorer() {
  const [blocks, setBlocks] = useState<Block[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const loadBlocks = async () => {
    try {
      setError("");
      setLoading(true);
      const res = await api.get("/blocks"); // adjust if your endpoint differs
      const data = res.data?.blocks || res.data || [];
      setBlocks(Array.isArray(data) ? data : []);
    } catch (err: any) {
      console.error("Block explorer error:", err);
      setError(
        err?.response?.data?.error ||
          "Failed to load blocks. Check backend /api/blocks endpoint."
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadBlocks();
  }, []);

  return (
    <div className="flex flex-col gap-6 max-w-3xl w-full mx-auto">

      {/* Header */}
      <div className="glass-card-soft">
        <h1 className="text-xl font-semibold text-white mb-1">
          Block Explorer
        </h1>
        <p className="text-xs text-slate-400">
          View all mined blocks, their hashes, nonce and linkage.
        </p>
      </div>

      {error && (
        <div className="text-sm text-red-400 bg-red-950/40 border border-red-700 px-4 py-3 rounded-xl">
          {error}
        </div>
      )}

      {/* Blocks list */}
      <div className="glass-card w-full">
        {loading ? (
          <p className="text-xs text-slate-500">Loading blocks…</p>
        ) : blocks.length === 0 ? (
          <p className="text-xs text-slate-500">
            No blocks found yet. Mine a block from the dashboard.
          </p>
        ) : (
          <div className="space-y-3 max-h-[28rem] overflow-y-auto pr-1">
            {blocks
              .slice()
              .sort((a, b) => b.index - a.index)
              .map((b) => (
                <div
                  key={b.hash}
                  className="bg-slate-950/70 border border-slate-800 rounded-xl px-4 py-3 flex flex-col gap-2"
                >
                  {/* Top row: index + timestamp */}
                  <div className="flex justify-between items-center gap-3">
                    <span className="text-[11px] text-slate-400">
                      Block #{b.index}
                    </span>
                    <span className="badge badge-muted">
                      {new Date(b.timestamp).toLocaleString()}
                    </span>
                  </div>

                  {/* Hash info */}
                  <div className="space-y-1">
                    <div className="text-[11px] text-slate-500">
                      <span className="font-semibold text-slate-300">
                        Hash:
                      </span>{" "}
                      <span className="font-mono break-all block">
                        {b.hash}
                      </span>
                    </div>
                    <div className="text-[11px] text-slate-500">
                      <span className="font-semibold text-slate-300">
                        Prev:
                      </span>{" "}
                      <span className="font-mono break-all block">
                        {b.previous_hash}
                      </span>
                    </div>
                  </div>

                  {/* Nonce + Tx count row */}
                  <div className="flex justify-between items-center mt-2 text-[11px] text-slate-500">
                    <span>
                      Nonce:{" "}
                      <span className="text-slate-200 font-semibold">
                        {b.nonce ?? "-"}
                      </span>
                    </span>
                    <span>
                      Tx count:{" "}
                      <span className="text-slate-200 font-semibold">
                        {b.tx_count ?? "-"}
                      </span>
                    </span>
                  </div>
                </div>
              ))}
          </div>
        )}
      </div>
    </div>
  );
}
