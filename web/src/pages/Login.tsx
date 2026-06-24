import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import Input from "../components/Input";
import Button from "../components/Button";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const { login, loading, error } = useAuth();
  const navigate = useNavigate();

  async function handleSubmit() {
    const ok = await login(email, password);
    if (ok) navigate("/");
  }

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "20px" }}>
      <div>
        <h1 style={{ fontSize: "20px", fontWeight: 600, marginBottom: "4px" }}>Sign in</h1>
        <p style={{ color: "var(--subtle)", fontSize: "13px" }}>Monitor your services from one place.</p>
      </div>

      <div style={{ display: "flex", flexDirection: "column", gap: "14px" }}>
        <Input
          label="Email"
          type="email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          placeholder="you@example.com"
        />
        <Input
          label="Password"
          type="password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          placeholder="••••••••"
        />
      </div>

      {error && (
        <p style={{ fontSize: "13px", color: "var(--down)", background: "#1a0a0a", padding: "10px 12px", borderRadius: "var(--radius)", border: "1px solid #3a1515" }}>
          {error}
        </p>
      )}

      <Button onClick={handleSubmit} loading={loading} style={{ width: "100%", justifyContent: "center" }}>
        Sign in
      </Button>

      <p style={{ textAlign: "center", fontSize: "13px", color: "var(--subtle)" }}>
        No account?{" "}
        <Link to="/register" style={{ color: "var(--accent)" }}>Register</Link>
      </p>
    </div>
  );
}
