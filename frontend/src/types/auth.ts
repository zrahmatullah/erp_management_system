export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  roleId: string;
}

export interface Role {
  id: string;
  name: string;
  description: string;
}

export interface Permission {
  id: string;
  module: string;
  action: string;
}

export interface LoginResponse {
  token: string;
  refreshToken: string;
  user: User;
  permissions: Permission[];
}
