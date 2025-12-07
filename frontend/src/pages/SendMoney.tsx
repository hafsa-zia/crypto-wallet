import { useState } from "react";
import type { FormEvent } from "react"
import api from "../api/client";

export default function SendMoney() {
  const [receiver, setReceiver] = useState("");
  const [amount, setAmount] = useState("");
  const [note, setNote] = useState("");
  const [msg, setMsg] = useState("");

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setMsg("");
    try {
      const res = await api.post("/tx", {
        receiver_wallet: receiver,
        amount: parseFloat(amount),
        note,
      });
      setMsg(res.data.message || "Transaction created");
    } catch (err: any) {
      setMsg(err?.response?.data?.error || "Failed");
    }
  };

  return (
    <div className="p-6 max-w-lg mx-auto">
      <h1 className="text-xl font-semibold mb-4">Send Money</h1>
      <form onSubmit={onSubmit} className="space-y-4 bg-white p-4 rounded-xl shadow">
        {msg && <div className="text-sm text-blue-600">{msg}</div>}
        <div>
          <label className="block text-sm mb-1">Receiver Wallet ID</label>
          <input
            className="w-full border rounded-lg px-3 py-2"
            value={receiver}
            onChange={(e) => setReceiver(e.target.value)}
            required
          />
        </div>
        <div>
          <label className="block text-sm mb-1">Amount</label>
          <input
            type="number"
            min="0"
            step="0.0001"
            className="w-full border rounded-lg px-3 py-2"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            required
          />
        </div>
        <div>
          <label className="block text-sm mb-1">Note (optional)</label>
          <textarea
            className="w-full border rounded-lg px-3 py-2"
            value={note}
            onChange={(e) => setNote(e.target.value)}
          />
        </div>
        <button className="w-full py-2 rounded-lg bg-green-600 text-white hover:bg-green-700">
          Send
        </button>
      </form>
    </div>
  );
}
