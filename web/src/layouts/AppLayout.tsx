import { Outlet, NavLink, useNavigate } from "react-router-dom";
import * as React from "react";

export default function AppLayout() {
  const navigate = useNavigate();

  function logout() {
    localStorage.removeItem("animas_token");
    navigate("/login");
  }

  return (
    <div style={{ display: "flex", height: "100vh" }}>
      {/* Sidebar */}
      <aside style={{
        width: "220px",
        flexShrink: 0,
        background: "var(--surface)",
        borderRight: "1px solid var(--border)",
        display: "flex",
        flexDirection: "column",
        padding: "24px 16px",
      }}>
        <span style={{
          fontFamily: "var(--mono)",
          fontSize: "16px",
          fontWeight: 500,
          color: "var(--text)",
          marginBottom: "32px",
          paddingLeft: "8px",
        }}>
          ▲ animas
        </span>

        <nav style={{ display: "flex", flexDirection: "column", gap: "4px", flex: 1 }}>
          <NavLink to="/" end style={navStyle}>
            Monitors
          </NavLink>
        </nav>

        <button onClick={logout} style={{
          background: "none",
          border: "none",
          color: "var(--subtle)",
          textAlign: "left",
          padding: "8px",
          borderRadius: "var(--radius)",
          fontSize: "13px",
          transition: "color 0.15s",
        }}
          onMouseEnter={e => (e.currentTarget.style.color = "var(--text)")}
          onMouseLeave={e => (e.currentTarget.style.color = "var(--subtle)")}
        >
          Sign out
        </button>
      </aside>

      {/* Main */}
      <main style={{ flex: 1, overflow: "auto", padding: "32px" }}>
        <Outlet />
      </main>
    </div>
  );
}

function navStyle({ isActive }: { isActive: boolean }) {
  return {
    display: "block",
    padding: "8px",
    borderRadius: "var(--radius)",
    color: isActive ? "var(--text)" : "var(--subtle)",
    background: isActive ? "var(--border)" : "transparent",
    fontSize: "13px",
    fontWeight: isActive ? 500 : 400,
    transition: "all 0.15s",
  } as React.CSSProperties;
}
