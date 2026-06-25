import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useMonitors } from "../hooks/useMonitors";
import { api } from "../services/api";
import {Monitor, MonitorType} from "../types/monitor";
import Button from "../components/Button";
import Input from "../components/Input";
import StatusBadge from "../components/StatusBadge";

export default function Dashboard() {
  const { monitors, loading, deleteMonitor } = useMonitors();
  const [creating, setCreating] = useState(false);
  const navigate = useNavigate();

  return (
    <div style={{ maxWidth: "760px" }}>
      {/* Header */}
      <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: "28px" }}>
        <div>
          <h1 style={{ fontSize: "18px", fontWeight: 600 }}>Monitors</h1>
          <p style={{ color: "var(--subtle)", fontSize: "13px", marginTop: "2px" }}>
            {monitors.length} active
          </p>
        </div>
        <Button onClick={() => setCreating(true)}>+ Add monitor</Button>
      </div>

      {/* Create form */}
      {creating && (
        <CreateMonitorForm
          onCreated={() => setCreating(false)}
          onCancel={() => setCreating(false)}
        />
      )}

      {/* List */}
      {loading ? (
        <p style={{ color: "var(--subtle)" }}>Loading…</p>
      ) : monitors.length === 0 && !creating ? (
        <EmptyState onAdd={() => setCreating(true)} />
      ) : (
        <div style={{ display: "flex", flexDirection: "column", gap: "1px" }}>
          {monitors.map(m => (
            <div
              key={m.ID}
              onClick={() => navigate(`/monitors/${m.ID}`)}
              style={{
                display: "flex",
                alignItems: "center",
                gap: "16px",
                padding: "14px 16px",
                background: "var(--surface)",
                border: "1px solid var(--border)",
                borderRadius: "var(--radius)",
                cursor: "pointer",
                transition: "border-color 0.15s",
              }}
              onMouseEnter={e => (e.currentTarget.style.borderColor = "var(--muted)")}
              onMouseLeave={e => (e.currentTarget.style.borderColor = "var(--border)")}
            >
              <StatusBadge status={m.lastStatus || (m.Enabled ? "UP" : "DOWN")} />

              <div style={{ flex: 1, minWidth: 0 }}>
                <p style={{ fontWeight: 500, fontSize: "14px", whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>
                  {m.Name}
                </p>
                <p style={{ color: "var(--subtle)", fontSize: "12px", fontFamily: "var(--mono)", marginTop: "1px" }}>
                  {m.Target}
                </p>
              </div>

              <span style={{
                fontSize: "11px",
                color: "var(--subtle)",
                background: "var(--border)",
                padding: "2px 8px",
                borderRadius: "4px",
                fontFamily: "var(--mono)",
              }}>
                {m.Type}
              </span>

              <span style={{ fontSize: "12px", color: "var(--subtle)" }}>
                every {m.Interval}s
              </span>

              <button
                onClick={e => { e.stopPropagation(); deleteMonitor(m.ID); }}
                style={{
                  background: "none",
                  border: "none",
                  color: "var(--subtle)",
                  padding: "4px 8px",
                  borderRadius: "4px",
                  fontSize: "16px",
                  lineHeight: 1,
                  transition: "color 0.15s",
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
  );
}

// CreateMonitorForm

function CreateMonitorForm({ onCreated, onCancel }: { onCreated: () => void; onCancel: () => void }) {
  const [name, setName]         = useState("");
  const [target, setTarget]     = useState("");
  const [type, setType]         = useState<MonitorType>("HTTPS");
  const [interval, setInterval] = useState("60");
  const [loading, setLoading]   = useState(false);
  const [error, setError]       = useState<string | null>(null);

  async function submit() {
    try {
      setLoading(true);
      setError(null);
      const res = await api.post<{ data: Monitor }>("/api/monitors", { name, target, type, interval: parseInt(interval) });
      console.log("[Dashboard] Create monitor response:", res);
      console.log("[Dashboard] Monitor ID:", res.data.ID);
      onCreated();
    } catch (e: any) {
      setError(e.message);
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
      <p style={{ fontWeight: 500, fontSize: "13px" }}>New monitor</p>

      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px" }}>
        <Input label="Name" value={name} onChange={e => setName(e.target.value)} placeholder="My API" />
        <Input label="Target" value={target} onChange={e => setTarget(e.target.value)} placeholder="https://api.example.com" />
      </div>

      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px" }}>
        <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
          <label style={{ fontSize: "12px", color: "var(--subtle)", fontWeight: 500, letterSpacing: "0.04em" }}>TYPE</label>
          <select
            value={type}
            onChange={e => setType(e.target.value as MonitorType)}
            style={{
              background: "var(--surface)",
              border: "1px solid var(--border)",
              borderRadius: "var(--radius)",
              color: "var(--text)",
              padding: "9px 12px",
            }}
          >
            <option>HTTP</option>
            <option>HTTPS</option>
            <option>TCP</option>
          </select>
        </div>
        <Input label="Interval (seconds)" type="number" value={interval} onChange={e => setInterval(e.target.value)} />
      </div>

      {error && <p style={{ fontSize: "13px", color: "var(--down)" }}>{error}</p>}

      <div style={{ display: "flex", gap: "8px" }}>
        <Button onClick={submit} loading={loading}>Create</Button>
        <Button variant="ghost" onClick={onCancel}>Cancel</Button>
      </div>
    </div>
  );
}

// EmptyState

function EmptyState({ onAdd }: { onAdd: () => void }) {
  return (
    <div style={{
      textAlign: "center",
      padding: "64px 24px",
      border: "1px dashed var(--border)",
      borderRadius: "var(--radius)",
      color: "var(--subtle)",
    }}>
      <p style={{ fontSize: "32px", marginBottom: "12px" }}>◎</p>
      <p style={{ fontWeight: 500, color: "var(--text)", marginBottom: "4px" }}>No monitors yet</p>
      <p style={{ fontSize: "13px", marginBottom: "20px" }}>Add your first endpoint to start tracking uptime.</p>
      <Button onClick={onAdd}>+ Add monitor</Button>
    </div>
  );
}
