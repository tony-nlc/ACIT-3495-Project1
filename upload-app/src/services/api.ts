import axios from "axios";

const api = axios.create({
  baseURL: "http://auth-service:5000", 
});

export default api;
