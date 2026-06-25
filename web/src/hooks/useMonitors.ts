import { useEffect, useState } from "react";
import { Monitor } from "../types/monitor";
import { api } from "../services/api";

export function useMonitors() {
  const [monitors, setMonitors] = useState<Monitor[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  async function fetch() {
    try {
      setLoading(true);
      const res = await api.get<{ data: Monitor[] }>("/api/monitors");
      setMonitors(res.data);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => { fetch(); }, []);

  async function deleteMonitor(id: string) {
    await api.delete(`/api/monitors/${id}`);
    setMonitors(prev => prev.filter(m => m.ID !== id));
  }

  return { monitors, loading, error, refetch: fetch, deleteMonitor };
}
