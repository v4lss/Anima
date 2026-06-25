import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Monitor, Check } from "../types/monitor";
import { AlertConfig, AlertType } from "../types/alert";
import { api } from "../services/api";
import StatusBadge from "../components/StatusBadge";
import Button from "../components/Button";
import Input from "../components/Input";

export default function MonitorDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [monitor, setMonitor] = useState<Monitor | null>(null);
  const [checks, setChecks]   = useState<Check[]>([]);
  const [alerts, setAlerts]   = useState<AlertConfig[]>([]);
  const [loading, setLoading] = useState(true);
  const [showAlertForm, setShowAlertForm] = useState(false);
  
  // Pagination state
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [loadingMore, setLoadingMore] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  
  // Filter state
  const [statusFilter, setStatusFilter] = useState<string>("");
  const [fromDate, setFromDate] = useState<string>("");
  const [toDate, setToDate] = useState<string>("");

  async function loadChecks(pageNum: number = 1, append: boolean = false) {
    if (!id) return;
    
    try {
      if (pageNum === 1) {
        setLoading(true);
      } else {
        setLoadingMore(true);
      }
      
      const params = new URLSearchParams({
        page: pageNum.toString(),
        limit: "20",
      });
      if (statusFilter) params.append("status", statusFilter);
      if (fromDate) params.append("from", fromDate);
      if (toDate) params.append("to", toDate);
      
      const res = await api.get<{ data: { data: Check[]; total: number; page: number; limit: number } }>(
        `/api/monitors/${id}/history?${params.toString()}`
      );
      
      if (append) {
        setChecks(prev => [...prev, ...(res.data.data ?? [])]);
      } else {
        setChecks(res.data.data ?? []);
      }
      setTotal(res.data.total);
      setHasMore((res.data.data?.length ?? 0) === res.data.limit && (res.data.page * res.data.limit) < res.data.total);
    } catch (e) {
      console.error(e);
      if (pageNum === 1) {
        setChecks([]);
      }
    } finally {
      setLoading(false);
      setLoadingMore(false);
    }
  }

  async function load() {
    try {
      const [mRes, aRes] = await Promise.all([
        api.get<{ data: Monitor }>(`/api/monitors/${id}`),
        api.get<{ data: AlertConfig[] }>(`/api/monitors/${id}/alerts`),
      ]);
      setMonitor(mRes.data);
      setAlerts(aRes.data ?? []);
      await loadChecks(1, false);
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, [id]);

  useEffect(() => {
    if (!loading) {
      setPage(1);
      loadChecks(1, false);
    }
  }, [statusFilter, fromDate, toDate]);

  function loadMore() {
    const nextPage = page + 1;
    setPage(nextPage);
    loadChecks(nextPage, true);
  }

  if (loading) return <p style={{ color: "var(--subtle)" }}>Loading…</p>;
  if (!monitor) return <p style={{ color: "var(--down)" }}>Monitor not found.</p>;

  const lastCheck = checks[0];
  const upCount   = checks.filter(c => c.Status === "UP").length;
  const uptime    = checks.length ? ((upCount / checks.length) * 100).toFixed(1) : "-";
  const avgMs     = checks.length
    ? Math.round(checks.reduce((s, c) => s + c.ResponseTime, 0) / checks.length)
    : 0;

  return (
    <div style={{ maxWidth: "760px" }}>
      {/* Back */}
      <button
        onClick={() => navigate("/")}
        style={{ background: "none", border: "none", color: "var(--subtle)", fontSize: "13px", marginBottom: "20px", cursor: "pointer" }}
      >
        ← Back
      </button>

      {/* Header */}
      <div style={{ display: "flex", alignItems: "flex-start", justifyContent: "space-between", marginBottom: "28px" }}>
        <div>
          <div style={{ display: "flex", alignItems: "center", gap: "12px", marginBottom: "4px" }}>
            <h1 style={{ fontSize: "20px", fontWeight: 600 }}>{monitor.Name}</h1>
            {lastCheck && <StatusBadge status={lastCheck.Status} />}
          </div>
          <p style={{ color: "var(--subtle)", fontFamily: "var(--mono)", fontSize: "12px" }}>{monitor.Target}</p>
        </div>
        <span style={{
          fontSize: "11px",
          color: "var(--subtle)",
          background: "var(--border)",
          padding: "3px 10px",
          borderRadius: "4px",
          fontFamily: "var(--mono)",
        }}>
          {monitor.Type}
        </span>
      </div>

      {/* Stats */}
      <div style={{ display: "grid", gridTemplateColumns: "repeat(3, 1fr)", gap: "12px", marginBottom: "28px" }}>
        {[
          { label: "Uptime (last 100)", value: `${uptime}%` },
          { label: "Avg response", value: checks.length ? `${avgMs}ms` : "-" },
          { label: "Check interval", value: `${monitor.Interval}s` },
        ].map(s => (
          <div key={s.label} style={{
            background: "var(--surface)",
            border: "1px solid var(--border)",
            borderRadius: "var(--radius)",
            padding: "16px",
          }}>
            <p style={{ fontSize: "11px", color: "var(--subtle)", marginBottom: "6px", letterSpacing: "0.04em" }}>
              {s.label.toUpperCase()}
            </p>
            <p style={{ fontFamily: "var(--mono)", fontSize: "22px", fontWeight: 500 }}>{s.value}</p>
          </div>
        ))}
      </div>

      {/* Check history */}
      <div style={{ marginTop: "32px" }}>
        <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: "12px" }}>
          <p style={{ fontWeight: 500, fontSize: "13px" }}>Recent checks ({total})</p>
        </div>

        {/* Filters */}
        <div style={{ display: "flex", gap: "12px", marginBottom: "16px", flexWrap: "wrap" }}>
          <select
            value={statusFilter}
            onChange={e => setStatusFilter(e.target.value)}
            style={{
              background: "var(--surface)",
              border: "1px solid var(--border)",
              borderRadius: "var(--radius)",
              padding: "6px 12px",
              fontSize: "13px",
              color: "var(--text)",
            }}
          >
            <option value="">All Status</option>
            <option value="UP">UP</option>
            <option value="DOWN">DOWN</option>
          </select>
          <input
            type="datetime-local"
            value={fromDate}
            onChange={e => setFromDate(e.target.value)}
            style={{
              background: "var(--surface)",
              border: "1px solid var(--border)",
              borderRadius: "var(--radius)",
              padding: "6px 12px",
              fontSize: "13px",
              color: "var(--text)",
            }}
            placeholder="From date"
          />
          <input
            type="datetime-local"
            value={toDate}
            onChange={e => setToDate(e.target.value)}
            style={{
              background: "var(--surface)",
              border: "1px solid var(--border)",
              borderRadius: "var(--radius)",
              padding: "6px 12px",
              fontSize: "13px",
              color: "var(--text)",
            }}
            placeholder="To date"
          />
          {(statusFilter || fromDate || toDate) && (
            <button
              onClick={() => {
                setStatusFilter("");
                setFromDate("");
                setToDate("");
              }}
              style={{
                background: "none",
                border: "1px solid var(--border)",
                borderRadius: "var(--radius)",
                padding: "6px 12px",
                fontSize: "13px",
                color: "var(--subtle)",
                cursor: "pointer",
              }}
            >
              Clear filters
            </button>
          )}
        </div>

        {checks.length === 0 ? (
          <p style={{ color: "var(--subtle)", fontSize: "13px" }}>No checks recorded yet.</p>
        ) : (
          <>
            <div style={{ display: "flex", flexDirection: "column", gap: "1px" }}>
              {checks.map(c => (
            <div key={c.ID} style={{
              display: "flex",
              alignItems: "center",
              gap: "16px",
              padding: "10px 14px",
              background: "var(--surface)",
              border: "1px solid var(--border)",
              borderRadius: "var(--radius)",
              fontSize: "13px",
            }}>
              <StatusBadge status={c.Status} />
              <span style={{ fontFamily: "var(--mono)", color: "var(--subtle)", minWidth: "60px" }}>
                {c.ResponseTime}ms
              </span>
              <span style={{ color: "var(--subtle)", fontSize: "12px", flex: 1 }}>
                {new Date(c.CheckedAt).toLocaleString()}
              </span>
              {c.Error && (
                <span style={{ fontSize: "11px", color: "var(--down)", fontFamily: "var(--mono)" }}>
                  {c.Error}
                </span>
              )}
            </div>
          ))}
            </div>
            {hasMore && (
              <button
                onClick={loadMore}
                disabled={loadingMore}
                style={{
                  marginTop: "16px",
                  background: "var(--surface)",
                  border: "1px solid var(--border)",
                  borderRadius: "var(--radius)",
                  padding: "10px 20px",
                  fontSize: "13px",
                  color: "var(--text)",
                  cursor: loadingMore ? "not-allowed" : "pointer",
                  opacity: loadingMore ? 0.6 : 1,
                }}
              >
                {loadingMore ? "Loading..." : "Load more"}
              </button>
            )}
          </>
        )}
      </div>

      {/* Alerts */}
      <div style={{ marginTop: "32px" }}>
        <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: "12px" }}>
          <p style={{ fontWeight: 500, fontSize: "13px" }}>Alerts</p>
          <Button onClick={() => setShowAlertForm(!showAlertForm)} variant="ghost">
            {showAlertForm ? "Cancel" : "+ Add alert"}
          </Button>
        </div>

        {showAlertForm && <AlertForm monitorId={id!} onCreated={async () => { setShowAlertForm(false); await loadAlerts(); }} />}

        {alerts.length === 0 && !showAlertForm ? (
          <p style={{ color: "var(--subtle)", fontSize: "13px" }}>No alerts configured.</p>
        ) : (
          <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
            {alerts.map(a => (
              <div key={a.ID} style={{
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                padding: "12px 16px",
                background: "var(--surface)",
                border: "1px solid var(--border)",
                borderRadius: "var(--radius)",
              }}>
                <div>
                  <p style={{ fontWeight: 500, fontSize: "13px" }}>{a.Type}</p>
                  <p style={{ color: "var(--subtle)", fontSize: "12px", fontFamily: "var(--mono)" }}>
                    {a.Webhook}
                  </p>
                </div>
                <button
                  onClick={() => deleteAlert(a.ID)}
                  style={{
                    background: "none",
                    border: "none",
                    color: "var(--subtle)",
                    padding: "4px 8px",
                    borderRadius: "4px",
                    fontSize: "16px",
                    lineHeight: 1,
                    cursor: "pointer",
                  }}
                  onMouseEnter={e => (e.currentTarget.style.color = "var(--down)")}
                  onMouseLeave={e => (e.currentTarget.style.color = "var(--subtle)")}
                >
                  ×
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );

  async function loadAlerts() {
    try {
      const res = await api.get<{ data: AlertConfig[] }>(`/api/monitors/${id}/alerts`);
      setAlerts(res.data ?? []);
    } catch (e: any) {
      console.error(e);
    }
  }

  async function deleteAlert(alertId: string) {
    await api.delete(`/api/monitors/${id}/alerts/${alertId}`);
    setAlerts(prev => prev.filter(a => a.ID !== alertId));
  }
}

function AlertForm({ monitorId, onCreated }: { monitorId: string; onCreated: () => void }) {
  const [type, setType] = useState<AlertType>("DISCORD");
  const [webhook, setWebhook] = useState("");
  const [loading, setLoading] = useState(false);

  async function submit() {
    try {
      setLoading(true);
      await api.post(`/api/monitors/${monitorId}/alerts`, { type, webhook });
      onCreated();
    } catch (e: any) {
      alert(e.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div style={{
      background: "var(--surface)",
      border: "1px solid var(--accent-dim)",
      borderRadius: "var(--radius)",
      padding: "20px",
      marginBottom: "16px",
      display: "flex",
      flexDirection: "column",
      gap: "14px",
    }}>
      <p style={{ fontWeight: 500, fontSize: "13px" }}>New alert</p>

      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px" }}>
        <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
          <label style={{ fontSize: "12px", color: "var(--subtle)", fontWeight: 500, letterSpacing: "0.04em" }}>TYPE</label>
          <select
            value={type}
            onChange={e => setType(e.target.value as AlertType)}
            style={{
              background: "var(--surface)",
              border: "1px solid var(--border)",
              borderRadius: "var(--radius)",
              color: "var(--text)",
              padding: "9px 12px",
            }}
          >
            <option>DISCORD</option>
            <option>EMAIL</option>
          </select>
        </div>
        <Input label="Webhook URL" value={webhook} onChange={e => setWebhook(e.target.value)} placeholder="https://discord.com/api/webhooks/..." />
      </div>

      <div style={{ display: "flex", gap: "8px" }}>
        <Button onClick={submit} loading={loading}>Create</Button>
      </div>
    </div>
  );
}
