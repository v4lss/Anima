import { CheckStatus } from "../types/monitor";

export default function StatusBadge({ status }: { status: CheckStatus }) {
  const up = status === "UP";
  return (
    <span style={{
      display: "inline-flex",
      alignItems: "center",
      gap: "6px",
      fontSize: "12px",
      fontWeight: 500,
      color: up ? "var(--up)" : "var(--down)",
      fontFamily: "var(--mono)",
    }}>
      <span style={{
        width: "6px",
        height: "6px",
        borderRadius: "50%",
        background: up ? "var(--up)" : "var(--down)",
        boxShadow: up ? "0 0 6px var(--up)" : "0 0 6px var(--down)",
      }} />
      {status}
    </span>
  );
}
