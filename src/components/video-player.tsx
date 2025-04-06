import { useState, useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';
import Hls from 'hls.js';
import { Card } from '@/components/ui/card';

interface VideoDetails {
  id: string;
  title: string;
  description: string;
  hlsUrl: string;
}

export default function VideoPlayer() {
  const { videoId } = useParams();
  const [video, setVideo] = useState<VideoDetails | null>(null);
  const videoRef = useRef<HTMLVideoElement>(null);
  const hlsRef = useRef<Hls | null>(null);

  useEffect(() => {
    fetchVideoDetails();
  }, [videoId]);

  const fetchVideoDetails = async () => {
    try {
      const response = await fetch(`/api/videos/${videoId}`);
      const data = await response.json();
      setVideo(data);
      
      if (data.hlsUrl) {
        initializeHlsPlayer(data.hlsUrl);
      }
    } catch (error) {
      console.error('Error fetching video details:', error);
    }
  };

  const initializeHlsPlayer = (hlsUrl: string) => {
    if (!videoRef.current) return;

    if (hlsRef.current) {
      hlsRef.current.destroy();
    }

    if (Hls.isSupported()) {
      const hls = new Hls({
        enableWorker: true,
        lowLatencyMode: true,
      });

      hls.loadSource(hlsUrl);
      hls.attachMedia(videoRef.current);
      hlsRef.current = hls;

      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        videoRef.current?.play();
      });
    } else if (videoRef.current.canPlayType('application/vnd.apple.mpegurl')) {
      // For Safari
      videoRef.current.src = hlsUrl;
    }
  };

  useEffect(() => {
    return () => {
      if (hlsRef.current) {
        hlsRef.current.destroy();
      }
    };
  }, []);

  if (!video) {
    return <div>Loading...</div>;
  }

  return (
    <div className="container mx-auto p-6">
      <Card className="overflow-hidden">
        <div className="aspect-video bg-black">
          <video
            ref={videoRef}
            className="w-full h-full"
            controls
            playsInline
          />
        </div>
        <div className="p-6">
          <h1 className="text-2xl font-bold mb-4">{video.title}</h1>
          <p className="text-muted-foreground">{video.description}</p>
        </div>
      </Card>
    </div>
  );
}