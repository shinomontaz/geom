package wkb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/twpayne/gogeom/geom"
	"io"
)

type wkbReader func(io.Reader, binary.ByteOrder) (geom.T, error)

var wkbReaders map[uint32]wkbReader

func init() {
	wkbReaders = make(map[uint32]wkbReader)
	wkbReaders[wkbPoint] = pointReader
	wkbReaders[wkbPointZ] = pointZReader
	wkbReaders[wkbPointM] = pointMReader
	wkbReaders[wkbPointZM] = pointZMReader
	wkbReaders[wkbLineString] = lineStringReader
	wkbReaders[wkbLineStringZ] = lineStringZReader
	wkbReaders[wkbLineStringM] = lineStringMReader
	wkbReaders[wkbLineStringZM] = lineStringZMReader
	wkbReaders[wkbPolygon] = polygonReader
	wkbReaders[wkbPolygonZ] = polygonZReader
	wkbReaders[wkbPolygonM] = polygonMReader
	wkbReaders[wkbPolygonZM] = polygonZMReader
}

func readPoints(r io.Reader, byteOrder binary.ByteOrder) ([]geom.Point, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	points := make([]geom.Point, numPoints)
	if err := binary.Read(r, byteOrder, &points); err != nil {
		return nil, err
	}
	return points, nil
}

func readPointZs(r io.Reader, byteOrder binary.ByteOrder) ([]geom.PointZ, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointZs := make([]geom.PointZ, numPoints)
	if err := binary.Read(r, byteOrder, &pointZs); err != nil {
		return nil, err
	}
	return pointZs, nil
}

func readPointMs(r io.Reader, byteOrder binary.ByteOrder) ([]geom.PointM, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointMs := make([]geom.PointM, numPoints)
	if err := binary.Read(r, byteOrder, &pointMs); err != nil {
		return nil, err
	}
	return pointMs, nil
}

func readPointZMs(r io.Reader, byteOrder binary.ByteOrder) ([]geom.PointZM, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointZMs := make([]geom.PointZM, numPoints)
	if err := binary.Read(r, byteOrder, &pointZMs); err != nil {
		return nil, err
	}
	return pointZMs, nil
}

func pointReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	point := geom.Point{}
	if err := binary.Read(r, byteOrder, &point); err != nil {
		return nil, err
	}
	return point, nil
}

func pointZReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointZ := geom.PointZ{}
	if err := binary.Read(r, byteOrder, &pointZ); err != nil {
		return nil, err
	}
	return pointZ, nil
}

func pointMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointM := geom.PointM{}
	if err := binary.Read(r, byteOrder, &pointM); err != nil {
		return nil, err
	}
	return pointM, nil
}

func pointZMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointZM := geom.PointZM{}
	if err := binary.Read(r, byteOrder, &pointZM); err != nil {
		return nil, err
	}
	return pointZM, nil
}

func lineStringReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	points, err := readPoints(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineString{points}, nil
}

func lineStringZReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointZs, err := readPointZs(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineStringZ{pointZs}, nil
}

func lineStringMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointMs, err := readPointMs(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineStringM{pointMs}, nil
}

func lineStringZMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointZMs, err := readPointZMs(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineStringZM{pointZMs}, nil
}

func polygonReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	rings := make([][]geom.Point, numRings)
	for i := uint32(0); i < numRings; i++ {
		if points, err := readPoints(r, byteOrder); err != nil {
			return nil, err
		} else {
			rings[i] = points
		}
	}
	return geom.Polygon{rings}, nil
}

func polygonZReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	ringZs := make([][]geom.PointZ, numRings)
	for i := uint32(0); i < numRings; i++ {
		if pointZs, err := readPointZs(r, byteOrder); err != nil {
			return nil, err
		} else {
			ringZs[i] = pointZs
		}
	}
	return geom.PolygonZ{ringZs}, nil
}

func polygonMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	ringMs := make([][]geom.PointM, numRings)
	for i := uint32(0); i < numRings; i++ {
		if pointMs, err := readPointMs(r, byteOrder); err != nil {
			return nil, err
		} else {
			ringMs[i] = pointMs
		}
	}
	return geom.PolygonM{ringMs}, nil
}

func polygonZMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	ringZMs := make([][]geom.PointZM, numRings)
	for i := uint32(0); i < numRings; i++ {
		if pointZMs, err := readPointZMs(r, byteOrder); err != nil {
			return nil, err
		} else {
			ringZMs[i] = pointZMs
		}
	}
	return geom.PolygonZM{ringZMs}, nil
}

func Read(r io.Reader) (geom.T, error) {

	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return nil, err
	}
	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case wkbXDR:
		byteOrder = binary.BigEndian
	case wkbNDR:
		byteOrder = binary.LittleEndian
	default:
		return nil, fmt.Errorf("invalid byte order %u", wkbByteOrder)
	}

	var wkbGeometryType uint32
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return nil, err
	}

	if reader, ok := wkbReaders[wkbGeometryType]; ok {
		return reader(r, byteOrder)
	} else {
		return nil, fmt.Errorf("unsupported geometry type %u", wkbGeometryType)
	}

}

func Unmarshal(buf []byte) (geom.T, error) {
	return Read(bytes.NewBuffer(buf))
}
