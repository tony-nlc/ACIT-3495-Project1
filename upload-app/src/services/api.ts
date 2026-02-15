import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:5000", // change later to auth service container name
});

export default api;
