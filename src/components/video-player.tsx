import "media-chrome";
import "hls-video-element";
// import {
//   MediaSettingsMenu,
//   MediaSettingsMenuItem,
//   MediaPlaybackRateMenu,
//   MediaRenditionMenu,
//   MediaCaptionsMenu,
// } from "media-chrome";
import "media-chrome/menu";

interface VideoPlayerProps {
  videoId: string;
}

export default function VideoPlayer({ videoId }: VideoPlayerProps) {
  return (
    // <div className="w-full aspect-video bg-black rounded-lg overflow-hidden">
    <div className="w-full aspect-video rounded-lg">
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

          {/* <media-settings-menu hidden anchor="auto">
            <media-settings-menu-item>
              Playback Speed
              <media-playback-rate-menu slot="submenu" hidden>
                <div slot="title">Playback Speed</div>
              </media-playback-rate-menu>
            </media-settings-menu-item>
            <media-settings-menu-item class="quality-settings">
              Quality
              <media-rendition-menu slot="submenu" hidden>
                <div slot="title">Quality</div>
              </media-rendition-menu>
            </media-settings-menu-item>
            <media-settings-menu-item>
              Subtitles/CC
              <media-captions-menu slot="submenu" hidden>
                <div slot="title">Subtitles/CC</div>
              </media-captions-menu>
            </media-settings-menu-item>
          </media-settings-menu> */}
        </media-control-bar>
      </media-controller>
    </div>
  );
}
