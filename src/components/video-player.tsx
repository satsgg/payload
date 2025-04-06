interface VideoPlayerProps {
  videoId: string;
}

export default function VideoPlayer({ videoId }: VideoPlayerProps) {
  return (
    // <div className="w-full aspect-video bg-black rounded-lg overflow-hidden">
    <div className="w-full aspect-video rounded-lg overflow-hidden">
      <media-controller>
        <hls-video
          src={`/api/videos/${videoId}/playlist.m3u8`}
          slot="media"
          crossorigin
          muted
        ></hls-video>
        <media-loading-indicator
          slot="centered-chrome"
          noautohide
        ></media-loading-indicator>
        <media-control-bar>
          <media-play-button></media-play-button>
          <media-seek-backward-button></media-seek-backward-button>
          <media-seek-forward-button></media-seek-forward-button>
          <media-mute-button></media-mute-button>
          <media-volume-range></media-volume-range>
          <media-time-range></media-time-range>
          <media-time-display showduration remaining></media-time-display>
          <media-playback-rate-button></media-playback-rate-button>
          <media-fullscreen-button></media-fullscreen-button>
        </media-control-bar>
      </media-controller>
    </div>
  );
}
