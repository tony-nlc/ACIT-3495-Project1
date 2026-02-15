import { useEffect, useState } from "react";
import { getVideos, type Video } from "../services/videoService";

function Streaming() {
  const [videos, setVideos] = useState<Video[]>([]);

  useEffect(() => {
    async function fetchVideos() {
      try {
        const data = await getVideos();
        setVideos(data);
      } catch (error) {
        alert("Failed to load videos");
      }
    }

    fetchVideos();
  }, []);

  return (
    <div>
      <h2>Video Streaming Page</h2>

      {videos.length === 0 ? (
        <p>No videos available.</p>
      ) : (
        <ul>
          {videos.map((video) => (
            <li key={video.id}>{video.title}</li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default Streaming;
