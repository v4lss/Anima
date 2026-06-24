import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import Input from "../components/Input";
import Button from "../components/Button";

export default function Register() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const { register, loading, error } = useAuth();
  const navigate = useNavigate();

  async function handleSubmit() {
    const ok = await register(email, password);
    if (ok) navigate("/login");
  }

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "20px" }}>
      <div>
        <h1 style={{ fontSize: "20px", fontWeight: 600, marginBottom: "4px" }}>Create account</h1>
        <p style={{ color: "var(--subtle)", fontSize: "13px" }}>Start monitoring in seconds.</p>
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
          placeholder="min. 8 characters"
        />
      </div>

      {error && (
        <p style={{ fontSize: "13px", color: "var(--down)", background: "#1a0a0a", padding: "10px 12px", borderRadius: "var(--radius)", border: "1px solid #3a1515" }}>
          {error}
        </p>
      )}

      <Button onClick={handleSubmit} loading={loading} style={{ width: "100%", justifyContent: "center" }}>
        Create account
      </Button>

      <p style={{ textAlign: "center", fontSize: "13px", color: "var(--subtle)" }}>
        Already have an account?{" "}
        <Link to="/login" style={{ color: "var(--accent)" }}>Sign in</Link>
      </p>
    </div>
  );
}
