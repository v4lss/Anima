import { ButtonHTMLAttributes } from "react";

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "ghost" | "danger";
  loading?: boolean;
}

export default function Button({ variant = "primary", loading, children, style, ...props }: Props) {
  const base: React.CSSProperties = {
    display: "inline-flex",
    alignItems: "center",
    gap: "8px",
    padding: "8px 16px",
    borderRadius: "var(--radius)",
    fontWeight: 500,
    fontSize: "13px",
    border: "none",
    transition: "all 0.15s",
    opacity: props.disabled || loading ? 0.5 : 1,
    cursor: props.disabled || loading ? "not-allowed" : "pointer",
  };

  const variants = {
    primary: { background: "var(--accent)", color: "#fff" },
    ghost:   { background: "transparent", color: "var(--subtle)", border: "1px solid var(--border)" },
    danger:  { background: "transparent", color: "var(--down)", border: "1px solid var(--down)" },
  };

  return (
    <button {...props} style={{ ...base, ...variants[variant], ...style }}>
      {loading ? "Loading…" : children}
    </button>
  );
}
