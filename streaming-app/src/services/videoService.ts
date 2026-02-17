import axios from "axios";

const STREAMING_BASE = "http://localhost:5002";

export interface Video {
  id: number;
  title: string;
}

export async function getVideos(): Promise<Video[]> {
  const token = localStorage.getItem("token");

  const response = await axios.get<Video[]>(
    `${STREAMING_BASE}/getvideos`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  return response.data;
}

export async function getVideoBlob(id: number): Promise<string> {
  const token = localStorage.getItem("token");

  const response = await axios.get(`${STREAMING_BASE}/view/${id}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
    responseType: "blob", // Critical: Tells axios to treat this as binary data
  });

  // Create a local URL for the video data
  return URL.createObjectURL(response.data);
}
