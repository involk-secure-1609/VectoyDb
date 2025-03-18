package hnsw

import (
	"cmp"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"reflect"
)

// errorEncoder is a helper type to encode multiple values

var byteOrder = binary.LittleEndian

func binaryRead(r io.Reader, data any) (int, error) {
	switch v := data.(type) {
	case *int:
		br, ok := r.(io.ByteReader)
		if !ok {
			return 0, fmt.Errorf("reader does not implement io.ByteReader")
		}

		i, err := binary.ReadVarint(br)
		if err != nil {
			return 0, err
		}

		*v = int(i)
		// TODO: this will usually overshoot size.
		return binary.MaxVarintLen64, nil

	case *string:
		var ln int
		_, err := binaryRead(r, &ln)
		if err != nil {
			return 0, err
		}

		s := make([]byte, ln)
		_, err = binaryRead(r, &s)
		*v = string(s)
		return len(s), err

	case *Embedding:
		var ln int
		_, err := binaryRead(r, &ln)
		if err != nil {
			return 0, err
		}

		*v = make([]float64, ln)
		return binary.Size(*v), binary.Read(r, byteOrder, *v)
	case io.ReaderFrom:
		n, err := v.ReadFrom(r)
		return int(n), err

	default:
		return binary.Size(data), binary.Read(r, byteOrder, data)
	}
}

func binaryWrite(w io.Writer, data any) (int, error) {
	
	switch v := data.(type) {
	case int:
		var buf [binary.MaxVarintLen64]byte
		n := binary.PutVarint(buf[:], int64(v))
		n, err := w.Write(buf[:n])
		return n, err
	case io.WriterTo:
		n, err := v.WriteTo(w)
		return int(n), err
	case string:
		n, err := binaryWrite(w, len(v))
		if err != nil {
			return n, err
		}
		n2, err := io.WriteString(w, v)
		if err != nil {
			return n + n2, err
		}

		return n + n2, nil
	case Embedding:
		n, err := binaryWrite(w, len(v))
		if err != nil {
			return n, err
		}
		log.Println("INSIDE WRITING float")
		return n + binary.Size(v), binary.Write(w, byteOrder, v)

	default:
		log.Println("Printing v",v)
		t := reflect.TypeOf(data)
		log.Println(t)
		sz := binary.Size(data)
		err := binary.Write(w, byteOrder, data)
		if err != nil {
			return 0, fmt.Errorf("encoding %T: %w", data, err)
		}
		return sz, err
	}
}

func multiBinaryWrite(w io.Writer, data ...any) (int, error) {
	var written int
	for _, d := range data {
		log.Println(d)
		n, err := binaryWrite(w, d)
		log.Println(n)
		written += n
		if err != nil {
			return written, err
		}
	}
	return written, nil
}

func multiBinaryRead(r io.Reader, data ...any) (int, error) {
	var read int
	for i, d := range data {
		n, err := binaryRead(r, d)
		read += n
		if err != nil {
			return read, fmt.Errorf("reading %T at index %v: %w", d, i, err)
		}
	}
	return read, nil
}

const encodingVersion = 1


// SavedGraph is a wrapper around a graph that persists
// changes to a file upon calls to Save. It is more convenient
// but less powerful than calling Graph.Export and Graph.Import
// directly.
type SavedGraph[K cmp.Ordered] struct {
	*HNSWGraph[K]
	Path string
}

// LoadSavedGraph opens a graph from a file, reads it, and returns it.
//
// If the file does not exist (i.e. this is a new graph),
// the equivalent of NewGraph is returned.
//
// It does not hold open a file descriptor, so SavedGraph can be forgotten
// without ever calling Save.
// func LoadSavedGraph[K cmp.Ordered](path string, distanceFunc string) (*SavedGraph[K], error) {
// 	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o600)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	info, err := f.Stat()
// 	if err != nil {
// 		return nil, err
// 	}

// 	g := NewHNSWGraph[K](distanceFunc)
// 	if info.Size() > 0 {
// 		err = g.Load(bufio.NewReader(f))
// 		if err != nil {
// 			return nil, fmt.Errorf("import: %w", err)
// 		}
// 	}

// 	return &SavedGraph[K]{HNSWGraph: g, Path: path}, nil
// }

// // Save writes the graph to the file.
// func (g *SavedGraph[K]) Save() error {
// 	tmp, err := renameio.("", g.Path)
// 	if err != nil {
// 		return err
// 	}
// 	defer tmp.Cleanup()

// 	wr := bufio.NewWriter(tmp)
// 	err = g.Export(wr)
// 	if err != nil {
// 		return fmt.Errorf("exporting: %w", err)
// 	}

// 	err = wr.Flush()
// 	if err != nil {
// 		return fmt.Errorf("flushing: %w", err)
// 	}

// 	err = tmp.CloseAtomicallyReplace()
// 	if err != nil {
// 		return fmt.Errorf("closing atomically: %w", err)
// 	}

// 	return nil
// }
