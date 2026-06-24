// Monitor type definitions - mirror of the Go domain structs.

export type MonitorType = "HTTP" | "HTTPS" | "TCP";
export type CheckStatus = "UP" | "DOWN";

export interface Monitor {
  id: string;
  userid: string;
  name: string;
  target: string;
  type: MonitorType;
  interval: number; // seconds
  enabled: boolean;
  createdat: string;
  updatedat: string;
}

export interface Check {
  id: string;
  monitorid: string;
  status: CheckStatus;
  responsetime: number; // ms
  error?: string;
  checkedat: string;
}
