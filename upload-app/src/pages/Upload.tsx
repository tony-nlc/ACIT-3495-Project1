import { useState } from "react";
import { uploadVideo } from "../services/uploadService";

function Upload() {
  const [file, setFile] = useState<File | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!file) {
      alert("Please select a file");
      return;
    }

    try {
      const result = await uploadVideo(file);
      alert(result.message);
    } catch (error) {
      alert("Upload failed");
    }
  };

  return (
    <div>
      <h2>Upload Video</h2>

      <form onSubmit={handleSubmit}>
        <div>
          <input
            type="file"
            accept="video/*"
            onChange={(e) => {
              if (e.target.files) {
                setFile(e.target.files[0]);
              }
            }}
          />
        </div>

        <button type="submit">Upload</button>
      </form>
    </div>
  );
}

export default Upload;
