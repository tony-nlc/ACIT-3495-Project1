import { useEffect, useState } from "react";
import { getVideos, getVideoBlob, type Video } from "../services/videoService";

function Streaming() {
  const [videos, setVideos] = useState<Video[]>([]);
  const [videoSrc, setVideoSrc] = useState<string>(""); // Store the Blob URL
  const [currentTitle, setCurrentTitle] = useState<string>("");

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

  const handlePlay = async (video: Video) => {
    try {
      // If there was a previous video, revoke the old URL to save memory
      if (videoSrc) URL.revokeObjectURL(videoSrc);

      const url = await getVideoBlob(video.id);
      setVideoSrc(url);
      setCurrentTitle(video.title);
    } catch (error) {
      alert("Unauthorized or video not found");
    }
  };

  return (
    <div>
      <h2>Video Streaming Page</h2>
      <ul>
        {videos.map((v) => (
          <li key={v.id}>
            <button onClick={() => handlePlay(v)}>{v.title}</button>
          </li>
        ))}
      </ul>

      {videoSrc && (
        <div>
          <h3>Now Playing: {currentTitle}</h3>
          {/* Key fix: Use the videoSrc state which contains the authenticated Blob */}
          <video width="600" controls key={videoSrc}>
            <source src={videoSrc} type="video/mp4" />
          </video>
        </div>
      )}
    </div>
  );
}