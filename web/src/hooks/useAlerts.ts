import { useState } from "react";
import { AlertConfig, AlertType } from "../types/alert";
import { api } from "../services/api";

export function useAlerts(monitorId: string) {
  const [alerts, setAlerts] = useState<AlertConfig[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function fetch() {
    try {
      setLoading(true);
      setError(null);
      const res = await api.get<{ data: AlertConfig[] }>(`/api/monitors/${monitorId}/alerts`);
      setAlerts(res.data ?? []);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  async function create(type: AlertType, webhook: string) {
    try {
      setLoading(true);
      setError(null);
      await api.post(`/api/monitors/${monitorId}/alerts`, { type, webhook });
      await fetch();
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  async function remove(id: string) {
    try {
      setLoading(true);
      setError(null);
      await api.delete(`/api/monitors/${monitorId}/alerts/${id}`);
      setAlerts(prev => prev.filter(a => a.id !== id));
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  return { alerts, loading, error, fetch, create, remove };
}
