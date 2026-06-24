import { Outlet } from "react-router-dom";

export default function AuthLayout() {
  return (
    <div style={{
      minHeight: "100vh",
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
      padding: "24px",
    }}>
      <div style={{ width: "100%", maxWidth: "380px" }}>
        {/* Logo */}
        <div style={{ marginBottom: "32px", textAlign: "center" }}>
          <span style={{
            fontFamily: "var(--mono)",
            fontSize: "22px",
            fontWeight: 500,
            letterSpacing: "-0.5px",
            color: "var(--text)",
          }}>
            ▲ animas
          </span>
        </div>
        <Outlet />
      </div>
    </div>
  );
}
