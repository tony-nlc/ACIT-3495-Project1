import api from "./api";

export interface Video {
  id: number;
  title: string;
  path: string;
}

export async function getVideos(): Promise<Video[]> {
  const token = localStorage.getItem("token");

  /*
  const response = await api.get<Video[]>("/videos", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.data;
  */

  // TEMP mock data
  return Promise.resolve([
    { id: 1, title: "Sample Video 1", path: "/videos/sample1.mp4" },
    { id: 2, title: "Sample Video 2", path: "/videos/sample2.mp4" },
  ]);
}
