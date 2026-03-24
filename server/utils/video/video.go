package video

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type VideoMetadata struct {
	DurationSeconds float64
	Width           int
	Height          int
	Resolution      string
}

type trackMetadata struct {
	Width   int
	Height  int
	IsVideo bool
}

// GetVideoMetadata 获取视频秒数和分辨率。
func GetVideoMetadata(ctx context.Context, videoURL string) (*VideoMetadata, error) {
	videoURL = strings.TrimSpace(videoURL)
	if videoURL == "" {
		return nil, fmt.Errorf("video url is empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, videoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request failed: %w", err)
	}

	client := &http.Client{Timeout: 2 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request video failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request video failed: http %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read video failed: %w", err)
	}

	meta, err := parseMP4Metadata(data)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func parseMP4Metadata(data []byte) (*VideoMetadata, error) {
	meta := &VideoMetadata{}

	if err := walkAtoms(data, func(atomType string, payload []byte) error {
		if atomType != "moov" {
			return nil
		}
		return parseMoovAtom(payload, meta)
	}); err != nil {
		return nil, err
	}

	if meta.DurationSeconds <= 0 {
		return nil, fmt.Errorf("parse duration failed")
	}
	if meta.Width <= 0 || meta.Height <= 0 {
		return nil, fmt.Errorf("parse resolution failed")
	}

	meta.Resolution = fmt.Sprintf("%dx%d", meta.Width, meta.Height)
	return meta, nil
}

func parseMoovAtom(data []byte, meta *VideoMetadata) error {
	return walkAtoms(data, func(atomType string, payload []byte) error {
		switch atomType {
		case "mvhd":
			duration, err := parseMVHD(payload)
			if err != nil {
				return err
			}
			meta.DurationSeconds = duration
		case "trak":
			track, err := parseTRAK(payload)
			if err != nil {
				return err
			}
			if track.IsVideo && track.Width*track.Height > meta.Width*meta.Height {
				meta.Width = track.Width
				meta.Height = track.Height
			}
		}
		return nil
	})
}

func parseTRAK(data []byte) (*trackMetadata, error) {
	track := &trackMetadata{}

	err := walkAtoms(data, func(atomType string, payload []byte) error {
		switch atomType {
		case "tkhd":
			width, height, err := parseTKHD(payload)
			if err != nil {
				return err
			}
			track.Width = width
			track.Height = height
		case "mdia":
			isVideo, err := parseMDIA(payload)
			if err != nil {
				return err
			}
			track.IsVideo = isVideo
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return track, nil
}

func parseMDIA(data []byte) (bool, error) {
	isVideo := false

	err := walkAtoms(data, func(atomType string, payload []byte) error {
		if atomType != "hdlr" {
			return nil
		}
		handlerType, err := parseHDLR(payload)
		if err != nil {
			return err
		}
		isVideo = handlerType == "vide"
		return nil
	})
	if err != nil {
		return false, err
	}

	return isVideo, nil
}

func parseHDLR(data []byte) (string, error) {
	if len(data) < 12 {
		return "", fmt.Errorf("hdlr too short")
	}
	return string(data[8:12]), nil
}

func parseMVHD(data []byte) (float64, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("mvhd too short")
	}

	version := data[0]
	switch version {
	case 0:
		if len(data) < 20 {
			return 0, fmt.Errorf("mvhd version 0 too short")
		}
		timescale := binary.BigEndian.Uint32(data[12:16])
		duration := binary.BigEndian.Uint32(data[16:20])
		if timescale == 0 {
			return 0, fmt.Errorf("mvhd timescale is zero")
		}
		return float64(duration) / float64(timescale), nil
	case 1:
		if len(data) < 32 {
			return 0, fmt.Errorf("mvhd version 1 too short")
		}
		timescale := binary.BigEndian.Uint32(data[20:24])
		duration := binary.BigEndian.Uint64(data[24:32])
		if timescale == 0 {
			return 0, fmt.Errorf("mvhd timescale is zero")
		}
		return float64(duration) / float64(timescale), nil
	default:
		return 0, fmt.Errorf("unsupported mvhd version: %d", version)
	}
}

func parseTKHD(data []byte) (int, int, error) {
	if len(data) < 4 {
		return 0, 0, fmt.Errorf("tkhd too short")
	}

	version := data[0]
	var widthOffset int
	var heightOffset int

	switch version {
	case 0:
		widthOffset = 76
		heightOffset = 80
	case 1:
		widthOffset = 88
		heightOffset = 92
	default:
		return 0, 0, fmt.Errorf("unsupported tkhd version: %d", version)
	}

	if len(data) < heightOffset+4 {
		return 0, 0, fmt.Errorf("tkhd data too short")
	}

	widthFixed := binary.BigEndian.Uint32(data[widthOffset : widthOffset+4])
	heightFixed := binary.BigEndian.Uint32(data[heightOffset : heightOffset+4])

	width := int(widthFixed >> 16)
	height := int(heightFixed >> 16)
	return width, height, nil
}

func walkAtoms(data []byte, visit func(atomType string, payload []byte) error) error {
	for offset := 0; offset+8 <= len(data); {
		size32 := binary.BigEndian.Uint32(data[offset : offset+4])
		atomType := string(data[offset+4 : offset+8])

		headerSize := 8
		atomSize := uint64(size32)

		if atomSize == 1 {
			if offset+16 > len(data) {
				return fmt.Errorf("atom %s has invalid large size", atomType)
			}
			atomSize = binary.BigEndian.Uint64(data[offset+8 : offset+16])
			headerSize = 16
		} else if atomSize == 0 {
			atomSize = uint64(len(data) - offset)
		}

		if atomSize < uint64(headerSize) || offset+int(atomSize) > len(data) {
			return fmt.Errorf("atom %s size is invalid", atomType)
		}

		payloadStart := offset + headerSize
		payloadEnd := offset + int(atomSize)
		if err := visit(atomType, data[payloadStart:payloadEnd]); err != nil {
			return err
		}

		offset += int(atomSize)
	}

	return nil
}
