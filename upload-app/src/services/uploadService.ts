import axios from "axios";

const UPLOAD_BASE = "http://192.168.0.85:5003";

export async function uploadVideo(file: File) {
  const token = localStorage.getItem("token");

  const formData = new FormData();
  formData.append("video", file);

  const response = await axios.post(
    `${UPLOAD_BASE}/upload`,
    formData,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  return response.data;
}
