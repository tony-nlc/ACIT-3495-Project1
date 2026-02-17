import axios from "axios";

export const uploadApi = axios.create({
  baseURL: "http://upload-service:5003",
});
