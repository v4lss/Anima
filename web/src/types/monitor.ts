// Monitor type definitions - mirror of the Go domain structs.

export type MonitorType = "HTTP" | "HTTPS" | "TCP";
export type CheckStatus = "UP" | "DOWN";

export interface Monitor {
  ID: string;
  UserID: string;
  Name: string;
  Target: string;
  Type: MonitorType;
  Interval: number; // seconds
  Enabled: boolean;
  CreatedAt: string;
  UpdatedAt: string;
  lastStatus?: CheckStatus; // last check status from API
}

export interface Check {
  ID: string;
  MonitorID: string;
  Status: CheckStatus;
  ResponseTime: number; // ms
  Error?: string;
  CheckedAt: string;
}
