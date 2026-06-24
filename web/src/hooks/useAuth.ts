import { useState } from "react";
import { api } from "../services/api";
import { AuthResponse } from "../types/user";

export function useAuth() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function login(email: string, password: string): Promise<boolean> {
    try {
      setLoading(true);
      setError(null);
      const res = await api.post<AuthResponse>("/api/auth/login", { email, password });
      localStorage.setItem("animas_token", res.data.token);
      return true;
    } catch (e: any) {
      setError(e.message);
      return false;
    } finally {
      setLoading(false);
    }
  }

  async function register(email: string, password: string): Promise<boolean> {
    try {
      setLoading(true);
      setError(null);
      await api.post("/api/auth/register", { email, password });
      return true;
    } catch (e: any) {
      setError(e.message);
      return false;
    } finally {
      setLoading(false);
    }
  }

  return { login, register, loading, error };
}
