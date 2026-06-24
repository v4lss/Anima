// User type definitions.

export interface User {
  id: string;
  email: string;
  createdat: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}
