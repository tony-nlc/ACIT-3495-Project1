import { streamingApi } from "./streamingApi";

export interface Video {
  id: number;
  title: string;
  path: string;
}

export async function getVideos(): Promise<Video[]> {
  const token = localStorage.getItem("token");

  
  const response = await streamingApi.get<Video[]>("/videos", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.data;
  

}
