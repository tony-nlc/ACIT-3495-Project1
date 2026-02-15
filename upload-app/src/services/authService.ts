import api from "./api";

interface LoginResponse {
  token: string;
}

export async function login(username: string, password: string) {

  /*
  const response = await api.post<LoginResponse>("/login", {
    username,
    password,
  });

  return response.data.token;
  */

  // TEMP: simulate backend
  return Promise.resolve("fake-jwt-token");
}
