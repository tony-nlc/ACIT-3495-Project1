import axios from "axios";

const STREAMING_BASE = "http://localhost:5002";

export interface Video {
  id: number;
  title: string;
}

export async function getVideoBlob(id: number): Promise<string> {
  const token = localStorage.getItem("token");

  const response = await axios.get(`${STREAMING_BASE}/view/${id}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
    responseType: "blob",
  });

  // Create a local URL for the video data
  return URL.createObjectURL(response.data);
}