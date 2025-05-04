import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Card } from "@/components/ui/card";
import VideoPlayer from "@/components/video-player";

interface Video {
  id: string;
  title: string;
  description: string;
  duration: number;
  created_at: string;
}

export default function WatchPage() {
  const { videoId } = useParams();
  const [video, setVideo] = useState<Video | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchVideo = async () => {
      try {
        const response = await fetch(`/api/videos/${videoId}`);
        if (!response.ok) {
          throw new Error("Failed to fetch video");
        }
        const data = await response.json();
        setVideo(data);
      } catch (error) {
        console.error("Error fetching video:", error);
      } finally {
        setLoading(false);
      }
    };

    if (videoId) {
      fetchVideo();
    }
  }, [videoId]);

  if (loading) {
    return <div className="p-8">Loading video...</div>;
  }

  if (!video) {
    return <div className="p-8">Video not found</div>;
  }

  return (
    <div className="p-8 max-w-4xl mx-auto w-full">
      <Card className="p-6">
        <VideoPlayer videoId={video.id} />
        <div className="mt-6">
          <h1 className="text-2xl font-bold mb-2">{video.title}</h1>
          <p className="text-muted-foreground">{video.description}</p>
          <p className="text-sm text-muted-foreground mt-4">
            Uploaded on {new Date(video.created_at).toLocaleDateString()}
          </p>
        </div>
      </Card>
    </div>
  );
}
