import { useEffect, useRef } from "react";
import Hls from "hls.js";

interface VideoPlayerProps {
  videoId: string;
}

export default function VideoPlayer({ videoId }: VideoPlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    const hls = new Hls({
      xhrSetup: function (xhr, url) {
        // Ensure all segment requests go through the API
        if (url.includes(".ts")) {
          const segmentPath = url.split("/").pop();
          xhr.open("GET", `/api/videos/${videoId}/${segmentPath}`, true);
        }
      },
    });

    hls.loadSource(`/api/videos/${videoId}/playlist.m3u8`);
    hls.attachMedia(video);
    hls.on(Hls.Events.MANIFEST_PARSED, () => {
      video.play().catch((error) => {
        console.error("Error playing video:", error);
      });
    });

    return () => {
      hls.destroy();
    };
  }, [videoId]);

  return (
    <div className="w-full aspect-video bg-black rounded-lg overflow-hidden">
      <video ref={videoRef} controls className="w-full h-full" playsInline />
    </div>
  );
}
