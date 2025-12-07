import {
  createContext,
  useContext,
  useEffect,
  useState,
  
} from "react";
import type {
  
  ReactNode,
} from "react";
import api from "../api/client";

type User = {
  full_name: string;
  email: string;
  wallet_id: string;
};

type AuthContextType = {
  user: User | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  // Initialize from localStorage so refresh doesn't log you out
  const [user, setUser] = useState<User | null>(() => {
    const saved = localStorage.getItem("user");
    return saved ? (JSON.parse(saved) as User) : null;
  });

  useEffect(() => {
    // Optional: you can call /me here to validate token if you add such endpoint
    const token = localStorage.getItem("jwt");
    if (!token) {
      setUser(null);
    }
  }, []);

  const login = async (email: string, password: string) => {
    const res = await api.post("/auth/login", { email, password });
    const u: User = res.data.user;

    // Save token + user for later
    localStorage.setItem("jwt", res.data.token);
    localStorage.setItem("user", JSON.stringify(u));

    setUser(u);
  };

  const logout = () => {
    localStorage.removeItem("jwt");
    localStorage.removeItem("user");
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used inside AuthProvider");
  }
  return ctx;
};
