// User type definitions.

export interface User {
  id: string;
  email: string;
  createdAt: string;
}

export interface AuthResponse {
  data: {
    token: string;
    user: User;
  };
}
