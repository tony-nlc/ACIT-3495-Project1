import axios from "axios";

export const authApi = axios.create({
  baseURL: "http://auth-service:5000",
});
