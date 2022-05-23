package wrapper

import (
	"database/sql/driver"
	"fmt"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/taosdata/driver-go/v2/common"
	"github.com/taosdata/driver-go/v2/errors"
)

// @author: xftan
// @date: 2022/4/16 15:12
// @description: test for read raw block
func TestReadBlock(t *testing.T) {
	conn, err := TaosConnect("", "root", "taosdata", "", 0)
	if err != nil {
		t.Error(err)
		return
	}

	defer TaosClose(conn)

	res := TaosQuery(conn, "create database if not exists test_block_raw")
	code := TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	TaosFreeResult(res)
	time.Sleep(time.Second)
	res = TaosQuery(conn, "drop table if exists test_block_raw.all_type")
	code = TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	TaosFreeResult(res)
	res = TaosQuery(conn, "create table if not exists test_block_raw.all_type (ts timestamp,"+
		"c1 bool,"+
		"c2 tinyint,"+
		"c3 smallint,"+
		"c4 int,"+
		"c5 bigint,"+
		"c6 tinyint unsigned,"+
		"c7 smallint unsigned,"+
		"c8 int unsigned,"+
		"c9 bigint unsigned,"+
		"c10 float,"+
		"c11 double,"+
		"c12 binary(20),"+
		"c13 nchar(20)"+
		")")
	code = TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	TaosFreeResult(res)
	now := time.Now()
	after1s := now.Add(time.Second)
	sql := fmt.Sprintf("insert into test_block_raw.all_type values('%s',1,1,1,1,1,1,1,1,1,1,1,'test_binary','test_nchar')('%s',null,null,null,null,null,null,null,null,null,null,null,null,null)", now.Format(time.RFC3339Nano), after1s.Format(time.RFC3339Nano))
	res = TaosQuery(conn, sql)
	code = TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	TaosFreeResult(res)

	sql = "select * from test_block_raw.all_type"
	res = TaosQuery(conn, sql)
	code = TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	fileCount := TaosNumFields(res)
	rh, err := ReadColumn(res, fileCount)
	if err != nil {
		t.Error(err)
		return
	}
	precision := TaosResultPrecision(res)
	var data [][]driver.Value
	for {
		blockSize, errCode, block := TaosFetchRawBlock(res)
		if errCode != int(errors.SUCCESS) {
			errStr := TaosErrorStr(res)
			err := errors.NewError(code, errStr)
			t.Error(err)
			TaosFreeResult(res)
			return
		}
		if blockSize == 0 {
			break
		}
		d := ReadBlock(block, blockSize, rh.ColTypes, precision)
		data = append(data, d...)
	}
	TaosFreeResult(res)
	assert.Equal(t, 2, len(data))
	row1 := data[0]
	assert.Equal(t, now.UnixNano()/1e6, row1[0].(time.Time).UnixNano()/1e6)
	assert.Equal(t, true, row1[1].(bool))
	assert.Equal(t, int8(1), row1[2].(int8))
	assert.Equal(t, int16(1), row1[3].(int16))
	assert.Equal(t, int32(1), row1[4].(int32))
	assert.Equal(t, int64(1), row1[5].(int64))
	assert.Equal(t, uint8(1), row1[6].(uint8))
	assert.Equal(t, uint16(1), row1[7].(uint16))
	assert.Equal(t, uint32(1), row1[8].(uint32))
	assert.Equal(t, uint64(1), row1[9].(uint64))
	assert.Equal(t, float32(1), row1[10].(float32))
	assert.Equal(t, float64(1), row1[11].(float64))
	assert.Equal(t, "test_binary", row1[12].(string))
	assert.Equal(t, "test_nchar", row1[13].(string))
	row2 := data[1]
	assert.Equal(t, after1s.UnixNano()/1e6, row2[0].(time.Time).UnixNano()/1e6)
	for i := 1; i < 14; i++ {
		assert.Nil(t, row2[i])
	}
}

func TestReadBlock2(t *testing.T) {
	conn, err := TaosConnect("", "root", "taosdata", "", 0)
	if err != nil {
		t.Error(err)
		return
	}

	defer TaosClose(conn)

	res := TaosQuery(conn, "create database if not exists test_block_raw")
	code := TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	TaosFreeResult(res)
	time.Sleep(time.Second)
	res = TaosQuery(conn, "drop table if exists test_block_raw.all_type2")
	code = TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	TaosFreeResult(res)
	res = TaosQuery(conn, "create table if not exists test_block_raw.all_type2 (ts timestamp,"+
		"c1 bool,"+
		"c2 tinyint,"+
		"c3 smallint,"+
		"c4 int,"+
		"c5 bigint,"+
		"c6 tinyint unsigned,"+
		"c7 smallint unsigned,"+
		"c8 int unsigned,"+
		"c9 bigint unsigned,"+
		"c10 float,"+
		"c11 double,"+
		"c12 binary(20),"+
		"c13 nchar(20)"+
		")")
	code = TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	TaosFreeResult(res)
	now := time.Now()
	after1s := now.Add(time.Second)
	sql := fmt.Sprintf("insert into test_block_raw.all_type2 values('%s',1,1,1,1,1,1,1,1,1,1,1,'test_binary','test_nchar')('%s',null,null,null,null,null,null,null,null,null,null,null,null,null)", now.Format(time.RFC3339Nano), after1s.Format(time.RFC3339Nano))
	res = TaosQuery(conn, sql)
	code = TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	TaosFreeResult(res)

	sql = "select * from test_block_raw.all_type2"
	res = TaosQuery(conn, sql)
	code = TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	fileCount := TaosNumFields(res)
	rh, err := ReadColumn(res, fileCount)
	if err != nil {
		t.Error(err)
		return
	}
	precision := TaosResultPrecision(res)
	payloadOffset := uintptr(4 * fileCount)
	pHeaderList := make([]uintptr, fileCount)
	pStartList := make([]uintptr, fileCount)
	var data [][]driver.Value
	for {
		blockSize, errCode, block := TaosFetchRawBlock(res)
		if errCode != int(errors.SUCCESS) {
			errStr := TaosErrorStr(res)
			err := errors.NewError(code, errStr)
			t.Error(err)
			TaosFreeResult(res)
			return
		}
		if blockSize == 0 {
			break
		}
		nullBitMapOffset := uintptr(BitmapLen(blockSize))
		tmpPHeader := uintptr(block) + payloadOffset + 12 // length i32, group u64
		tmpPStart := tmpPHeader
		for column := 0; column < fileCount; column++ {
			colLength := *((*int32)(unsafe.Pointer(uintptr(block) + 12 + uintptr(column)*4)))
			if IsVarDataType(rh.ColTypes[column]) {
				pHeaderList[column] = tmpPHeader
				tmpPStart = tmpPHeader + uintptr(4*blockSize)
				pStartList[column] = tmpPStart
			} else {
				pHeaderList[column] = tmpPHeader
				tmpPStart = tmpPHeader + nullBitMapOffset
				pStartList[column] = tmpPStart
			}
			tmpPHeader = tmpPStart + uintptr(colLength)
		}
		for row := 0; row < blockSize; row++ {
			rowV := make([]driver.Value, fileCount)
			for column := 0; column < fileCount; column++ {
				v := ItemRawBlock(rh.ColTypes[column], pHeaderList[column], pStartList[column], row, precision, func(ts int64, precision int) driver.Value {
					return common.TimestampConvertToTime(ts, precision)
				})
				rowV[column] = v
			}
			data = append(data, rowV)
		}
	}
	TaosFreeResult(res)
	assert.Equal(t, 2, len(data))
	row1 := data[0]
	assert.Equal(t, now.UnixNano()/1e6, row1[0].(time.Time).UnixNano()/1e6)
	assert.Equal(t, true, row1[1].(bool))
	assert.Equal(t, int8(1), row1[2].(int8))
	assert.Equal(t, int16(1), row1[3].(int16))
	assert.Equal(t, int32(1), row1[4].(int32))
	assert.Equal(t, int64(1), row1[5].(int64))
	assert.Equal(t, uint8(1), row1[6].(uint8))
	assert.Equal(t, uint16(1), row1[7].(uint16))
	assert.Equal(t, uint32(1), row1[8].(uint32))
	assert.Equal(t, uint64(1), row1[9].(uint64))
	assert.Equal(t, float32(1), row1[10].(float32))
	assert.Equal(t, float64(1), row1[11].(float64))
	assert.Equal(t, "test_binary", row1[12].(string))
	assert.Equal(t, "test_nchar", row1[13].(string))
	row2 := data[1]
	assert.Equal(t, after1s.UnixNano()/1e6, row2[0].(time.Time).UnixNano()/1e6)
	for i := 1; i < 14; i++ {
		assert.Nil(t, row2[i])
	}
}

func TestBlockTag(t *testing.T) {
	conn, err := TaosConnect("", "root", "taosdata", "", 0)
	if err != nil {
		t.Error(err)
		return
	}

	defer TaosClose(conn)

	res := TaosQuery(conn, "create database if not exists test_block_abc1")
	code := TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	TaosFreeResult(res)

	res = TaosQuery(conn, "use test_block_abc1")
	code = TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	TaosFreeResult(res)

	res = TaosQuery(conn, "create table if not exists meters(ts timestamp, v int) tags(location varchar(16))")
	code = TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	TaosFreeResult(res)

	res = TaosQuery(conn, "create table if not exists tb1 using meters tags('abcd')")
	code = TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	TaosFreeResult(res)

	sql := "select tbname,location from meters;"
	res = TaosQuery(conn, sql)
	code = TaosError(res)
	if code != 0 {
		errStr := TaosErrorStr(res)
		TaosFreeResult(res)
		t.Error(errors.NewError(code, errStr))
		return
	}
	fileCount := TaosNumFields(res)
	rh, err := ReadColumn(res, fileCount)
	if err != nil {
		t.Error(err)
		return
	}
	precision := TaosResultPrecision(res)
	payloadOffset := uintptr(4 * fileCount)
	pHeaderList := make([]uintptr, fileCount)
	pStartList := make([]uintptr, fileCount)
	var data [][]driver.Value
	for {
		blockSize, errCode, block := TaosFetchRawBlock(res)
		if errCode != int(errors.SUCCESS) {
			errStr := TaosErrorStr(res)
			err := errors.NewError(code, errStr)
			t.Error(err)
			TaosFreeResult(res)
			return
		}
		if blockSize == 0 {
			break
		}
		nullBitMapOffset := uintptr(BitmapLen(blockSize))
		tmpPHeader := uintptr(block) + payloadOffset + 12 // length i32, group u64
		tmpPStart := tmpPHeader
		for column := 0; column < fileCount; column++ {
			colLength := *((*int32)(unsafe.Pointer(uintptr(block) + 12 + uintptr(column)*4)))
			if IsVarDataType(rh.ColTypes[column]) {
				pHeaderList[column] = tmpPHeader
				tmpPStart = tmpPHeader + uintptr(4*blockSize)
				pStartList[column] = tmpPStart
			} else {
				pHeaderList[column] = tmpPHeader
				tmpPStart = tmpPHeader + nullBitMapOffset
				pStartList[column] = tmpPStart
			}
			tmpPHeader = tmpPStart + uintptr(colLength)
		}
		for row := 0; row < blockSize; row++ {
			rowV := make([]driver.Value, fileCount)
			for column := 0; column < fileCount; column++ {
				v := ItemRawBlock(rh.ColTypes[column], pHeaderList[column], pStartList[column], row, precision, func(ts int64, precision int) driver.Value {
					return common.TimestampConvertToTime(ts, precision)
				})
				rowV[column] = v
			}
			data = append(data, rowV)
		}
	}
	TaosFreeResult(res)
	t.Log(data)
	t.Log(len(data[0][1].(string)))
}
