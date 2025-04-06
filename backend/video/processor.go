package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type VideoProcessor struct {
	ffmpegPath string
}

func NewVideoProcessor() (*VideoProcessor, error) {
	// Check if ffmpeg is installed
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, fmt.Errorf("ffmpeg not found: %v", err)
	}

	return &VideoProcessor{
		ffmpegPath: ffmpegPath,
	}, nil
}

func (p *VideoProcessor) ConvertToHLS(inputPath, outputDir string) (string, error) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %v", err)
	}

	// Generate output path for HLS files
	outputPath := filepath.Join(outputDir, "playlist.m3u8")

	// ffmpeg command to convert video to HLS
	cmd := exec.Command(p.ffmpegPath,
		"-i", inputPath,
		"-profile:v", "baseline", // H.264 profile
		"-level", "3.0", // H.264 level
		"-start_number", "0", // Start segment numbers from 0
		"-hls_time", "10", // Segment duration in seconds
		"-hls_list_size", "0", // Keep all segments
		"-f", "hls", // Output format
		outputPath,
	)

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg conversion failed: %v", err)
	}

	return outputPath, nil
}

func (p *VideoProcessor) GenerateThumbnail(inputPath, outputDir string) (string, error) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %v", err)
	}

	// Generate output path for thumbnail
	outputPath := filepath.Join(outputDir, "thumbnail.jpg")

	// ffmpeg command to generate thumbnail
	cmd := exec.Command(p.ffmpegPath,
		"-i", inputPath,
		"-ss", "00:00:01", // Take frame at 1 second
		"-vframes", "1", // Take only one frame
		"-q:v", "2", // Quality level
		outputPath,
	)

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("thumbnail generation failed: %v", err)
	}

	return outputPath, nil
}

func (p *VideoProcessor) GetVideoDuration(inputPath string) (int, error) {
	// ffmpeg command to get video duration
	cmd := exec.Command(p.ffmpegPath,
		"-i", inputPath,
		"-f", "null",
		"-",
	)

	// Capture stderr for duration information
	var stderr strings.Builder
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("failed to get video duration: %v", err)
	}

	// Parse duration from stderr output
	// The duration is in the format "Duration: HH:MM:SS.mmm"
	output := stderr.String()
	durationStr := ""
	if strings.Contains(output, "Duration:") {
		parts := strings.Split(output, "Duration: ")
		if len(parts) > 1 {
			durationStr = strings.Split(parts[1], ",")[0]
		}
	}

	if durationStr == "" {
		return 0, fmt.Errorf("could not parse video duration")
	}

	// Parse duration string to seconds
	var h, m, s int
	fmt.Sscanf(durationStr, "%d:%d:%d", &h, &m, &s)
	duration := h*3600 + m*60 + s

	return duration, nil
} 