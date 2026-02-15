import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:5000", // change later for docker
});

export default api;
