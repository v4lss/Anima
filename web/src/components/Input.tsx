import { InputHTMLAttributes } from "react";

interface Props extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export default function Input({ label, error, style, ...props }: Props) {
  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
      {label && (
        <label style={{ fontSize: "12px", color: "var(--subtle)", fontWeight: 500, letterSpacing: "0.04em" }}>
          {label.toUpperCase()}
        </label>
      )}
      <input
        {...props}
        style={{
          background: "var(--surface)",
          border: `1px solid ${error ? "var(--down)" : "var(--border)"}`,
          borderRadius: "var(--radius)",
          color: "var(--text)",
          padding: "9px 12px",
          outline: "none",
          transition: "border-color 0.15s",
          ...style,
        }}
        onFocus={e => { e.currentTarget.style.borderColor = "var(--accent)"; }}
        onBlur={e => { e.currentTarget.style.borderColor = error ? "var(--down)" : "var(--border)"; }}
      />
      {error && <span style={{ fontSize: "12px", color: "var(--down)" }}>{error}</span>}
    </div>
  );
}
