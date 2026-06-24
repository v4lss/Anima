// Monitor type definitions - mirror of the Go domain structs.

export type MonitorType = "HTTP" | "HTTPS" | "TCP";
export type CheckStatus = "UP" | "DOWN";

export interface Monitor {
  id: string;
  userId: string;
  name: string;
  target: string;
  type: MonitorType;
  interval: number; // seconds
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Check {
  id: string;
  monitorId: string;
  status: CheckStatus;
  responseTime: number; // ms
  error?: string;
  checkedAt: string;
}
