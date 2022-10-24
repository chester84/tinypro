package tools

import (
	_ "image/jpeg"

	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/types"
)

// BoxHeader 信息头
type BoxHeader struct {
	Size       uint32
	FourccType [4]byte
	Size64     uint64
}

// GetMP4Duration 获取视频时长，以秒计
func GetMP4Duration(reader io.ReaderAt) (lengthOfTime uint32, err error) {
	var info = make([]byte, 0x10)
	var boxHeader BoxHeader
	var offset int64 = 0
	// 获取moov结构偏移
	for {
		_, err = reader.ReadAt(info, offset)
		if err != nil {
			return
		}

		boxHeader = getHeaderBoxInfo(info)
		fourccType := getFourccType(boxHeader)
		if fourccType == "moov" {
			break
		}

		// 有一部分mp4 mdat尺寸过大需要特殊处理
		if fourccType == "mdat" {
			if boxHeader.Size == 1 {
				offset += int64(boxHeader.Size64)
				continue
			}
		}
		offset += int64(boxHeader.Size)
	}
	// 获取moov结构开头一部分
	moovStartBytes := make([]byte, 0x100)
	_, err = reader.ReadAt(moovStartBytes, offset)
	if err != nil {
		return
	}

	// 定义timeScale与Duration偏移
	timeScaleOffset := 0x1C
	durationOffeset := 0x20
	timeScale := binary.BigEndian.Uint32(moovStartBytes[timeScaleOffset : timeScaleOffset+4])
	Duration := binary.BigEndian.Uint32(moovStartBytes[durationOffeset : durationOffeset+4])
	lengthOfTime = Duration / timeScale

	return
}

// getHeaderBoxInfo 获取头信息
func getHeaderBoxInfo(data []byte) (boxHeader BoxHeader) {
	buf := bytes.NewBuffer(data)
	_ = binary.Read(buf, binary.BigEndian, &boxHeader)

	return
}

// getFourccType 获取信息头类型
func getFourccType(boxHeader BoxHeader) (fourccType string) {
	fourccType = string(boxHeader.FourccType[:])

	return
}

func ParseMediaBaseInfo(rid int64, d []byte) (info types.MediaBaseInfo, err error) {
	var mediaInfo types.MediaBaseInfo

	tmpFile := fmt.Sprintf(`/tmp/%d.avi`, rid)
	f, err := os.Create(tmpFile)
	if err != nil {
		logs.Warning("[ParseMediaBaseInfo] can not create file: %s, err: %v", tmpFile, err)
		return mediaInfo, err
	}

	_, _ = f.Write(d)
	_ = f.Close()

	defer func() {
		_ = Remove(tmpFile)
	}()

	cmd := exec.Command("ffmpeg", "-i", tmpFile, "-vframes", "1", "-f", "singlejpeg", "-")
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	var errBuff bytes.Buffer
	cmd.Stderr = &errBuff
	err = cmd.Run()
	if err != nil {
		logs.Error("[ParseMediaBaseInfo] could not generate frame, filename: %s, err: %v", tmpFile, err)
		return mediaInfo, err
	}

	mediaInfo.FirstFrame = buffer.Bytes()

	reader := bytes.NewReader(mediaInfo.FirstFrame)
	imgObj, _, err := image.DecodeConfig(reader)
	if err != nil {
		fmt.Print(fmt.Sprintf("can not decode img, err: %v\n", err))
	} else {
		mediaInfo.Width = imgObj.Width
		mediaInfo.Height = imgObj.Height
	}

	infoCmd := fmt.Sprintf(`ffmpeg -i %s 2>&1 | grep "Duration"`, tmpFile)
	cmd2 := exec.Command("bash", "-c", infoCmd)
	stdout, err := cmd2.CombinedOutput()
	if err != nil {
		logs.Error("[ParseMediaBaseInfo] exec info cmd exception, cmd: %s, err: %v", infoCmd, err)
		return mediaInfo, err
	}

	exp := strings.Split(string(stdout), ",")
	for _, info := range exp {
		subExp := strings.Split(info, ": ")
		if len(subExp) != 2 {
			continue
		}

		if strings.Contains(subExp[0], "Duration") {
			mediaInfo.DurationHum = subExp[1]

			var duration float64
			durationExp := strings.Split(mediaInfo.DurationHum, ".")
			if len(durationExp) == 2 {
				more, _ := strconv.ParseFloat(durationExp[1], 64)

				timeExp := strings.Split(durationExp[0], ":")
				var base float64 = 1
				for i := len(timeExp) - 1; i >= 0; i-- {
					t, _ := strconv.ParseFloat(timeExp[i], 64)
					duration += t * base
					base *= 60
				}

				duration = duration*1000 + more
			}

			mediaInfo.Duration = int64(duration)
		}

		if strings.Contains(subExp[0], "bitrate") {
			mediaInfo.BitrateHum = subExp[1]

			bitrateExp := strings.Split(mediaInfo.BitrateHum, " ")
			if len(bitrateExp) == 2 {
				bitrate, _ := strconv.ParseInt(bitrateExp[0], 10, 0)
				mediaInfo.Bitrate = int(bitrate)
			}
		}
	}

	return mediaInfo, nil
}
