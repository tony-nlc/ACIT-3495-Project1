import { useEffect, useState } from "react";
import { getVideos, type Video } from "../services/videoService";

function Streaming() {
  const [videos, setVideos] = useState<Video[]>([]);
  const [selectedVideo, setSelectedVideo] = useState<Video | null>(null);

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
            <li key={video.id}>
              <button onClick={() => setSelectedVideo(video)}>
                {video.title}
              </button>
            </li>
          ))}
        </ul>
      )}

      {selectedVideo && (
        <div>
          <h3>Now Playing: {selectedVideo.title}</h3>
          <video width="600" controls>
            <source
              src={`http://192.168.0.85:5002/view/${selectedVideo.id}`}
              type="video/mp4"
            />

            Your browser does not support the video tag.
          </video>
        </div>
      )}
    </div>
  );
}

export default Streaming;
