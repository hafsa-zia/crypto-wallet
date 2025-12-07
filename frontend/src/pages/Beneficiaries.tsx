import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import api from "../api/client";

export default function Beneficiaries() {
  const [items, setItems] = useState<string[]>([]);
  const [newItem, setNewItem] = useState("");
  const [msg, setMsg] = useState("");

  useEffect(() => {
    api.get("/wallet").then((res: { data: { beneficiaries?: string[] } }) => setItems(res.data.beneficiaries || []));
  }, []);

  const add = () => {
    if (!newItem) return;
    setItems((prev) => [...prev, newItem]);
    setNewItem("");
  };

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setMsg("");
    try {
      const res = await api.post("/wallet/beneficiaries", {
        beneficiaries: items,
      });
      setMsg(res.data.message);
    } catch (err: any) {
      setMsg(err?.response?.data?.error || "Failed");
    }
  };

  return (
    <div className="max-w-xl space-y-4">
      <h1 className="text-xl font-semibold">Beneficiary List</h1>
      {msg && <div className="text-sm text-blue-600">{msg}</div>}
      <form onSubmit={onSubmit} className="space-y-3 bg-white p-4 rounded-xl shadow">
        <div className="flex gap-2">
          <input
            className="flex-1 border rounded-lg px-3 py-2"
            placeholder="Wallet ID"
            value={newItem}
            onChange={(e) => setNewItem(e.target.value)}
          />
          <button
            type="button"
            onClick={add}
            className="px-3 py-2 bg-slate-800 text-white rounded-lg"
          >
            Add
          </button>
        </div>
        <ul className="space-y-1 text-sm">
          {items.map((b, i) => (
            <li key={i} className="flex justify-between">
              <span className="font-mono">{b}</span>
              <button
                type="button"
                onClick={() =>
                  setItems((prev) => prev.filter((_, idx) => idx !== i))
                }
                className="text-red-500 text-xs"
              >
                remove
              </button>
            </li>
          ))}
        </ul>
        <button className="w-full py-2 bg-green-600 text-white rounded-lg">
          Save
        </button>
      </form>
    </div>
  );
}
