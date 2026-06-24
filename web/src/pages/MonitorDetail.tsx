import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Monitor, Check } from "../types/monitor";
import { api } from "../services/api";
import StatusBadge from "../components/StatusBadge";
import Button from "../components/Button";

export default function MonitorDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [monitor, setMonitor] = useState<Monitor | null>(null);
  const [checks, setChecks]   = useState<Check[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const [mRes, cRes] = await Promise.all([
          api.get<{ data: Monitor }>(`/api/monitors/${id}`),
          api.get<{ data: Check[] }>(`/api/monitors/${id}/history`),
        ]);
        setMonitor(mRes.data);
        setChecks(cRes.data ?? []);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [id]);

  if (loading) return <p style={{ color: "var(--subtle)" }}>Loading…</p>;
  if (!monitor) return <p style={{ color: "var(--down)" }}>Monitor not found.</p>;

  const lastCheck = checks[0];
  const upCount   = checks.filter(c => c.status === "UP").length;
  const uptime    = checks.length ? ((upCount / checks.length) * 100).toFixed(1) : "-";
  const avgMs     = checks.length
    ? Math.round(checks.reduce((s, c) => s + c.responsetime, 0) / checks.length)
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
            <h1 style={{ fontSize: "20px", fontWeight: 600 }}>{monitor.name}</h1>
            {lastCheck && <StatusBadge status={lastCheck.status} />}
          </div>
          <p style={{ color: "var(--subtle)", fontFamily: "var(--mono)", fontSize: "12px" }}>{monitor.target}</p>
        </div>
        <span style={{
          fontSize: "11px",
          color: "var(--subtle)",
          background: "var(--border)",
          padding: "3px 10px",
          borderRadius: "4px",
          fontFamily: "var(--mono)",
        }}>
          {monitor.type}
        </span>
      </div>

      {/* Stats */}
      <div style={{ display: "grid", gridTemplateColumns: "repeat(3, 1fr)", gap: "12px", marginBottom: "28px" }}>
        {[
          { label: "Uptime (last 100)", value: `${uptime}%` },
          { label: "Avg response", value: checks.length ? `${avgMs}ms` : "-" },
          { label: "Check interval", value: `${monitor.interval}s` },
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
      <p style={{ fontWeight: 500, fontSize: "13px", marginBottom: "12px" }}>Recent checks</p>

      {checks.length === 0 ? (
        <p style={{ color: "var(--subtle)", fontSize: "13px" }}>No checks recorded yet.</p>
      ) : (
        <div style={{ display: "flex", flexDirection: "column", gap: "1px" }}>
          {checks.slice(0, 50).map(c => (
            <div key={c.id} style={{
              display: "flex",
              alignItems: "center",
              gap: "16px",
              padding: "10px 14px",
              background: "var(--surface)",
              border: "1px solid var(--border)",
              borderRadius: "var(--radius)",
              fontSize: "13px",
            }}>
              <StatusBadge status={c.status} />
              <span style={{ fontFamily: "var(--mono)", color: "var(--subtle)", minWidth: "60px" }}>
                {c.responsetime}ms
              </span>
              <span style={{ color: "var(--subtle)", fontSize: "12px", flex: 1 }}>
                {new Date(c.checkedat).toLocaleString()}
              </span>
              {c.error && (
                <span style={{ fontSize: "11px", color: "var(--down)", fontFamily: "var(--mono)" }}>
                  {c.error}
                </span>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
