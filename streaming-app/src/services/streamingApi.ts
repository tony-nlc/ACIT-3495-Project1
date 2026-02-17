import axios from "axios";

export const streamingApi = axios.create({
  baseURL: "http://streaming-service:5002",
});
