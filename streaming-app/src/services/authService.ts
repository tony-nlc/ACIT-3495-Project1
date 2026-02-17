import axios from "axios";

const AUTH_BASE = "http://localhost:5000";

export async function login(username: string, password: string) {
  const response = await axios.post<{ token: string }>(
    `${AUTH_BASE}/login`,
    {
      User: username,
      Pass: password,
    }
  );

  return response.data.token;
}
