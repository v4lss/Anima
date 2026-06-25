// Alert type definitions.

export type AlertType = "DISCORD" | "EMAIL";

export interface AlertConfig {
  ID: string;
  MonitorID: string;
  UserID: string;
  Type: AlertType;
  Webhook: string;
  Enabled: boolean;
  CreatedAt: string;
  UpdatedAt: string;
}
