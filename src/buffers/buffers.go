package buffers

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

const KB = 1024
const MB = 1024 * KB

const DefaultBlockSizeLimit =   1 * MB
const DefaultTotalSizeLimit = 100 * MB



const swapFileTemplate = "m_swap_*.tmp"

type dataBlock struct {
	lines          []string
}

type dataFrame struct {
	offset     int64
	firstLine  int
	noOfLines  int
	block     *dataBlock
}

type BufferedData struct {
	maxTotalSize    int64
	blockSizeLimit  int
	mutex			sync.Mutex
	frames          []dataFrame
	lruFrames      *list.List
	swapFile       *os.File
	lastBlockSize   int
}

type LineIndex struct {
	lineIndex   int
	frameIndex  int
	data       *BufferedData
}

func (buff *BufferedData)NewLineIndexer() *LineIndex {
	return &LineIndex{
		lineIndex:  0,
		frameIndex: 0,
		data:       buff,
	}
}

func (i *LineIndex)Index() int {
	return i.lineIndex
}

func (i *LineIndex)IndexBegin() bool {
	i.lineIndex = 0
	i.frameIndex = 0
	return i.lineIndex < i.data.Len()
}

func (i *LineIndex)IndexOK() bool {
	if i.frameIndex >= 0 && i.frameIndex < len(i.data.frames) {
       return i.data.frames[i.frameIndex].isLineInFrame(i.lineIndex)
	}
	return false
}

func (i *LineIndex)IndexSet(index int, force bool) bool {
	result := true
	if i.lineIndex != index {
		l := i.data.Len()
		if index < 0  || index >= l {
			if force {
				i.lineIndex = index
			}
			result = false
		}
		if l > 0 && result {
			if index == 0 {
				i.lineIndex = 0
				i.frameIndex = 0
			} else {
				var frameIndexDelta int
			    if index > i.lineIndex {
					frameIndexDelta = 1
				} else { // index < i.lineIndex
					frameIndexDelta = -1
				}
				var j int
				for j = i.frameIndex; j >= 0 && j <= len(i.data.frames) && !i.data.frames[j].isLineInFrame(index); j += frameIndexDelta {}
				if j < 0 || j >= len(i.data.frames) {
					panic(fmt.Sprintf("Internal error - wrong data frame index computed: %d", j))
				}
				i.frameIndex = j
				i.lineIndex = index
			}
		}
	}
	return result
}

func (i *LineIndex)IndexNext(delta int) bool {
	if delta == 0 {
		return false
	}
	return i.IndexSet(i.lineIndex + delta, true)
}

func (i *LineIndex)IndexIncrement() bool {
	return i.IndexNext(1)
}

func (i *LineIndex)IndexDecrement() bool {
	return i.IndexNext(-1)
}

func (i *LineIndex)IndexSetIfValid(index int) bool {
	return i.IndexSet(index, false)
}

func (i *LineIndex)GetLine() (string, error) {
	if frame, err := i.data.getFrame(i.frameIndex); err != nil {
		return "", err
	} else {
		j := i.lineIndex - frame.firstLine;
		if j < 0 || j >= frame.noOfLines {
			return "", errors.New(fmt.Sprintf("wrong line index (%s) in frame %d", j, i.frameIndex))
		} else {
			return frame.block.lines[j], nil
		}
	}
}

func NewBufferedData(blockSizeLimit int, totalSizeLimit int64) *BufferedData {
	if blockSizeLimit <= 0 {
		blockSizeLimit = DefaultBlockSizeLimit
	}

	if totalSizeLimit <= 0 {
		totalSizeLimit = DefaultTotalSizeLimit
	}

	return &BufferedData{
		maxTotalSize:   totalSizeLimit,
		blockSizeLimit: blockSizeLimit,
		lastBlockSize:  0,
		frames:         []dataFrame{},
		lruFrames:      list.New(),
		swapFile:       nil,
	}
}

func NewBufferedDataMB(blockSizeLimitMB int, totalSizeLimitMB int) *BufferedData {
	return NewBufferedData(blockSizeLimitMB * MB, int64(totalSizeLimitMB) * MB)
}

func NewBufferedDataDefault() *BufferedData {
	return NewBufferedData(-1, -1)
}

func (buff *BufferedData) AddLine(line string) {
	defer buff.mutex.Unlock()
	buff.mutex.Lock()
	line = strings.TrimRight(line, " \t\r\n")
	lineLength := len(line)
	if len(buff.frames) == 0 || buff.lastBlockSize > 0 && buff.lastBlockSize + lineLength > buff.blockSizeLimit {
		buff.frames = append(buff.frames, *newDataFrame(buff.Len()))
		buff.lastBlockSize = 0
	}
	lastFrameIndex := len(buff.frames) - 1
	lastFrame := &(buff.frames[lastFrameIndex])
	buff.reloadFrame(lastFrame, lastFrameIndex)
	lastFrame.block.lines = append(lastFrame.block.lines, line)
	lastFrame.noOfLines += 1
	buff.lastBlockSize += lineLength
}


func (buff *BufferedData) Close() {
	if buff.swapFile != nil {
		buff.swapFile.Close()
		os.Remove(buff.swapFile.Name())
		buff.swapFile = nil
	}
	buff.lruFrames = nil
	for _, frame := range buff.frames {
		frame.block = nil
	}
	buff.frames = nil
}

func (f *dataFrame)isLineInFrame(lineIndex int) bool {
	return lineIndex >= f.firstLine  && lineIndex < f.firstLine + f.noOfLines
}

func (buff *BufferedData) Len() int {
	l := len(buff.frames)
	if l == 0 {
		return 0
	}
	f := &(buff.frames[l-1])
	return f.firstLine + f.noOfLines
}


func newDataFrame(firstLine int) *dataFrame {
	return &dataFrame{
		offset:    -1,
		firstLine: firstLine,
		noOfLines: 0,
		block:     &dataBlock{
			lines: []string{},

		},
	}
}

func (buff *BufferedData)getFrame(frameIndex int) (*dataFrame, error) {
	defer buff.mutex.Unlock()
	buff.mutex.Lock()
	length := len(buff.frames)
	if frameIndex < 0 || frameIndex >= length {
		return nil, errors.New(fmt.Sprintf("wrong index %d in getFrame()", frameIndex))
	}
	result := &(buff.frames[frameIndex])
	buff.reloadFrame(result, frameIndex)
	return result, nil
}


func (buff *BufferedData)getWorkingFile() *os.File {
	if buff.swapFile == nil {
		var err error
		buff.swapFile, err = ioutil.TempFile("", swapFileTemplate)
		if err != nil {
			panic(fmt.Sprintf("Cannot open swap file: %s", err.Error()))
		}
	}
	return buff.swapFile
}

func (buff *BufferedData)unloadFrame(frame *dataFrame) {
	if frame.block != nil {
		if frame.offset < 0 {
			f := buff.getWorkingFile()
			offset, err := f.Seek(0, 2)
			if err == nil {
				var sb strings.Builder
				for _, line := range frame.block.lines {
					sb.WriteString(line)
					sb.WriteByte('\n')
				}
				if err == nil {
					_, err = f.WriteString(sb.String())
				}
				if err == nil {
					err = f.Sync()
				}
			}
			if err != nil {
				panic(fmt.Sprintf("I/O error in unloadFrame(): %s", err.Error()))
			}
			frame.offset = offset
		}
		frame.block = nil
	}
}

func (buff *BufferedData)loadFrame(frame *dataFrame) {
	if frame.block == nil {
		f := buff.getWorkingFile()
  		if frame.offset < 0 {
			panic("internal error: loadFrame for empty block and frame offset < 0")
		}
		if _, err := f.Seek(frame.offset, 0); err != nil {
			panic(fmt.Sprintf("cannot seek in loadFrame(): %s", err.Error()))
		}
		frame.block = &dataBlock{
			lines: []string{},
		}
		reader := bufio.NewReader(f)
		eof := false
		for l := 0; l < frame.noOfLines && !eof; l++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					eof = true
				} else {
					log.Fatal(err)
				}
			} else {
				frame.block.lines = append(frame.block.lines, strings.TrimRight(line, "\n"))
			}
		}
	}
}

func (buff *BufferedData)reloadFrame(frame *dataFrame, frameIndex int) {
	if head := buff.lruFrames.Front(); head != nil && head.Value.(int) == frameIndex && frame.block != nil {
		// No need to reload
		return
	}
	buff.loadFrame(frame)
	for e := buff.lruFrames.Front(); e != nil; e = e.Next() {
		if e.Value == frameIndex {
			buff.lruFrames.MoveToFront(e)
			break
		}
	}
	if first := buff.lruFrames.Front(); first == nil || first.Value != frameIndex {
		buff.lruFrames.PushFront(frameIndex)
	}
	if len := buff.lruFrames.Len(); len > 1 && int64(len) * int64(buff.blockSizeLimit) > buff.maxTotalSize {
		last := buff.lruFrames.Back()
		buff.unloadFrame(&(buff.frames[last.Value.(int)]))
		buff.lruFrames.Remove(last)
	}
}

