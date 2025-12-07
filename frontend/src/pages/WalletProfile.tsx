import { useEffect, useState } from "react";
import api from "../api/client";

type Profile = {
  full_name: string;
  email: string;
  cnic: string;
  wallet_id: string;
  beneficiaries: string[];
  zakat_deducted: number;
  balance: number;
};

export default function WalletProfile() {
  const [data, setData] = useState<Profile | null>(null);

  useEffect(() => {
    api.get("/wallet").then((res:any) => setData(res.data));
  }, []);

  if (!data) return <div>Loading...</div>;

  return (
    <div className="space-y-4">
      <div className="bg-white rounded-xl shadow p-4">
        <h2 className="text-lg font-semibold mb-2">Profile</h2>
        <p>Name: {data.full_name}</p>
        <p>Email: {data.email}</p>
        <p>CNIC: {data.cnic}</p>
        <p>
          Wallet: <span className="font-mono">{data.wallet_id}</span>
        </p>
        <p>Balance: {data.balance}</p>
        <p>Total Zakat Deducted: {data.zakat_deducted}</p>
      </div>
      <div className="bg-white rounded-xl shadow p-4">
        <h2 className="text-lg font-semibold mb-2">Beneficiaries</h2>
        <ul className="list-disc list-inside text-sm">
          {data.beneficiaries.map((b) => (
            <li key={b} className="font-mono">
              {b}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
