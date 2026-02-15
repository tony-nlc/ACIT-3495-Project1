import api from "./api";

export async function uploadVideo(file: File) {
  const token = localStorage.getItem("token");

  const formData = new FormData();
  formData.append("video", file);

  /*
  const response = await api.post("/upload", formData, {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "multipart/form-data",
    },
  });

  return response.data;
  */

  // TEMP simulate upload
  return Promise.resolve({ message: "Upload successful" });
}
