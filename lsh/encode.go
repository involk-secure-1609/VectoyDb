package lsh

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

// var byteOrder=binary.LittleEndian

// encode serializes the CosineLsh index to a file
func (lsh *CosineLsh) Save(storeName string) error {
	f, err := os.OpenFile(storeName+"_lsh"+".store", os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	var byteOrder = binary.LittleEndian
	writer := bufio.NewWriter(f)

	// Write scalar fields
	if err := binary.Write(writer, byteOrder, lsh.dim); err != nil {
		return err
	}
	if err := binary.Write(writer, byteOrder, lsh.l); err != nil {
		return err
	}
	if err := binary.Write(writer, byteOrder, lsh.m); err != nil {
		return err
	}
	if err := binary.Write(writer, byteOrder, lsh.h); err != nil {
		return err
	}

	// Write dFunc string
	if err := writeString(writer, lsh.dFunc); err != nil {
		return err
	}

	// Write hyperplanes
	// First write dimensions of hyperplanes

	// writing number of hyperplanes
	if err := binary.Write(writer, byteOrder, int32(len(lsh.hyperplanes))); err != nil {
		return err
	}

	// writing length of each hyperplane
	if len(lsh.hyperplanes) > 0 {
		if err := binary.Write(writer, byteOrder, int32(len(lsh.hyperplanes[0]))); err != nil {
			return err
		}
	} else {
		if err := binary.Write(writer, byteOrder, int32(0)); err != nil {
			return err
		}
	}

	// Then write the hyperplane values
	for i := range lsh.hyperplanes {
		for j := range lsh.hyperplanes[i] {
			if err := binary.Write(writer, byteOrder, lsh.hyperplanes[i][j]); err != nil {
				return err
			}
		}
	}

	// Write nextId
	if err := binary.Write(writer, byteOrder, lsh.nextID); err != nil {
		return err
	}

	// Write tables
	// First write the number of tables
	if err := binary.Write(writer, byteOrder, int32(len(lsh.tables))); err != nil {
		return err
	}

	// Write each table
	for _, table := range lsh.tables {
		// Write number of entries in this table
		if err := binary.Write(writer, byteOrder, int32(len(table))); err != nil {
			return err
		}

		// Write each key-value pair
		for key, points := range table {
			// Write the key
			if err := binary.Write(writer, byteOrder, key); err != nil {
				return err
			}

			// Write number of points for this key
			if err := binary.Write(writer, byteOrder, int32(len(points))); err != nil {
				return err
			}

			// Write each point
			for _, point := range points {
				// Write point ID
				if err := binary.Write(writer, byteOrder, point.ID); err != nil {
					return err
				}

				// Write vector length
				if err := binary.Write(writer, byteOrder, int32(len(point.Vector))); err != nil {
					return err
				}

				// Write vector values
				for _, v := range point.Vector {
					if err := binary.Write(writer, byteOrder, v); err != nil {
						return err
					}
				}

				// Write extra data
				if err := writeString(writer, point.ExtraData); err != nil {
					return err
				}
			}
		}
	}
	err=writer.Flush()
	if err!=nil{
		return err
	}
	err = f.Sync()
	if err != nil {
		return nil
	}
	return nil
}

// decode deserializes the CosineLsh index from a file
func (lsh *CosineLsh) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	var byteOrder = binary.LittleEndian
	reader := bufio.NewReader(file)

	// Read scalar fields
	if err := binary.Read(reader, byteOrder, &lsh.dim); err != nil {
		return err
	}
	if err := binary.Read(reader, byteOrder, &lsh.l); err != nil {
		return err
	}
	if err := binary.Read(reader, byteOrder, &lsh.m); err != nil {
		return err
	}
	if err := binary.Read(reader, byteOrder, &lsh.h); err != nil {
		return err
	}

	// Read dFunc string
	dFunc, err := readString(reader)
	if err != nil {
		return err
	}
	lsh.dFunc = dFunc

	// Read hyperplanes
	var rows, cols int32
	if err := binary.Read(reader, byteOrder, &rows); err != nil {
		return err
	}
	if err := binary.Read(reader, byteOrder, &cols); err != nil {
		return err
	}

	// Initialize and fill hyperplanes
	lsh.hyperplanes = make([][]float64, rows)
	for i := int32(0); i < rows; i++ {
		lsh.hyperplanes[i] = make([]float64, cols)
		for j := int32(0); j < cols; j++ {
			if err := binary.Read(reader, byteOrder, &lsh.hyperplanes[i][j]); err != nil {
				return err
			}
		}
	}

	// Read nextId
	if err := binary.Read(reader, byteOrder, &lsh.nextID); err != nil {
		return err
	}

	// Read tables
	var numTables int32
	if err := binary.Read(reader, byteOrder, &numTables); err != nil {
		return err
	}

	// Initialize tables
	lsh.tables = make([]hashTable, numTables)

	// Read each table
	for i := int32(0); i < numTables; i++ {
		var numEntries int32
		if err := binary.Read(reader, byteOrder, &numEntries); err != nil {
			return err
		}

		lsh.tables[i] = make(hashTable)

		// Read each key-value pair
		for j := int32(0); j < numEntries; j++ {
			var key uint64
			if err := binary.Read(reader, byteOrder, &key); err != nil {
				return err
			}

			var numPoints int32
			if err := binary.Read(reader, byteOrder, &numPoints); err != nil {
				return err
			}

			points := make([]Point, numPoints)

			// Read each point
			for k := int32(0); k < numPoints; k++ {
				// Read point ID
				if err := binary.Read(reader, byteOrder, &points[k].ID); err != nil {
					return err
				}

				// Read vector length
				var vectorLen int32
				if err := binary.Read(reader, byteOrder, &vectorLen); err != nil {
					return err
				}

				// Initialize and fill vector
				points[k].Vector = make([]float64, vectorLen)
				for v := int32(0); v < vectorLen; v++ {
					if err := binary.Read(reader, byteOrder, &points[k].Vector[v]); err != nil {
						return err
					}
				}

				// Read extra data
				extraData, err := readString(reader)
				if err != nil {
					return err
				}
				points[k].ExtraData = extraData
			}

			lsh.tables[i][key] = points
		}
	}

	return nil
}

// Helper function to write a string
func writeString(w io.Writer, s string) error {
	// Write the length of the string
	var byteOrder = binary.LittleEndian

	if err := binary.Write(w, byteOrder, int32(len(s))); err != nil {
		return err
	}
	// Write the string data
	if _, err := w.Write([]byte(s)); err != nil {
		return err
	}
	return nil
}

// Helper function to read a string
func readString(r io.Reader) (string, error) {
	// Read the length of the string
	var byteOrder = binary.LittleEndian
	var length int32
	if err := binary.Read(r, byteOrder, &length); err != nil {
		return "", err
	}
	if length < 0 {
		return "", errors.New("invalid string length")
	}

	// Read the string data
	bytes := make([]byte, length)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return "", err
	}

	return string(bytes), nil
}
