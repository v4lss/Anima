// Alert type definitions.

export type AlertType = "DISCORD" | "EMAIL";

export interface AlertConfig {
  id: string;
  monitorid: string;
  userid: string;
  type: AlertType;
  webhook: string;
  enabled: boolean;
  createdat: string;
  updatedat: string;
}
