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
