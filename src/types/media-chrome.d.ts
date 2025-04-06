declare namespace JSX {
  interface IntrinsicElements {
    "media-controller": React.DetailedHTMLProps<
      React.HTMLAttributes<HTMLElement>,
      HTMLElement
    >;
    "hls-video": React.DetailedHTMLProps<
      React.VideoHTMLAttributes<HTMLVideoElement> & {
        src: string;
        slot?: string;
        crossorigin?: boolean;
        muted?: boolean;
      },
      HTMLVideoElement
    >;
    "media-loading-indicator": React.DetailedHTMLProps<
      React.HTMLAttributes<HTMLElement> & {
        slot?: string;
        noautohide?: boolean;
      },
      HTMLElement
    >;
    "media-control-bar": React.DetailedHTMLProps<
      React.HTMLAttributes<HTMLElement>,
      HTMLElement
    >;
    "media-play-button": React.DetailedHTMLProps<
      React.HTMLAttributes<HTMLElement>,
      HTMLElement
    >;
    "media-seek-backward-button": React.DetailedHTMLProps<
      React.HTMLAttributes<HTMLElement>,
      HTMLElement
    >;
    "media-seek-forward-button": React.DetailedHTMLProps<
      React.HTMLAttributes<HTMLElement>,
      HTMLElement
    >;
    "media-mute-button": React.DetailedHTMLProps<
      React.HTMLAttributes<HTMLElement>,
      HTMLElement
    >;
    "media-volume-range": React.DetailedHTMLProps<
      React.HTMLAttributes<HTMLElement>,
      HTMLElement
    >;
    "media-time-range": React.DetailedHTMLProps<
      React.HTMLAttributes<HTMLElement>,
      HTMLElement
    >;
    "media-time-display": React.DetailedHTMLProps<
      React.HTMLAttributes<HTMLElement> & {
        showduration?: boolean;
        remaining?: boolean;
      },
      HTMLElement
    >;
    "media-playback-rate-button": React.DetailedHTMLProps<
      React.HTMLAttributes<HTMLElement>,
      HTMLElement
    >;
    "media-fullscreen-button": React.DetailedHTMLProps<
      React.HTMLAttributes<HTMLElement>,
      HTMLElement
    >;
  }
}
