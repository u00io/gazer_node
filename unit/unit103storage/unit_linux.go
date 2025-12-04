package unit103storage

import (
	"strconv"

	"golang.org/x/sys/unix"
)

func (c *Unit103Storage) processTick() []item {
	result := make([]item, 0)
	var stat unix.Statfs_t
	var free, total uint64
	err := unix.Statfs(c.disk, &stat)
	free = uint64(stat.Bsize) * stat.Bfree
	total = uint64(stat.Bsize) * stat.Blocks

	if err != nil {
		result = append(result, item{"/", "Status", err.Error(), "error"})
	} else {
		result = append(result, item{"/", "Status", "", ""})

		result = append(result, item{"/SpaceTotal", "Space Total", strconv.FormatUint(total/1024/1024, 10), "MB"})
		result = append(result, item{"/SpaceFree", "Space Free", strconv.FormatUint(free/1024/1024, 10), "MB"})
		result = append(result, item{"/SpaceUsed", "Space Used", strconv.FormatUint((total-free)/1024/1024, 10), "MB"})

		result = append(result, item{"/BlocksTotal", "Blocks Total", strconv.FormatUint(stat.Blocks, 10), ""})
		result = append(result, item{"/BlocksFree", "Blocks Free", strconv.FormatUint(stat.Bfree, 10), ""})
		result = append(result, item{"/BlocksUsed", "Blocks Used", strconv.FormatUint(stat.Blocks-stat.Bfree, 10), ""})

		result = append(result, item{"/INodesTotal", "INodes Total", strconv.FormatUint(stat.Files, 10), ""})
		result = append(result, item{"/INodesFree", "INodes Free", strconv.FormatUint(stat.Ffree, 10), ""})
		result = append(result, item{"/INodesUsed", "INodes Used", strconv.FormatUint(stat.Files-stat.Ffree, 10), ""})

		result = append(result, item{"/SpaceUsedPercents", "Space Used Percents", strconv.FormatFloat(100*float64(total-free)/float64(total), 'f', 2, 64), "%"})
		result = append(result, item{"/BlocksUsedPercents", "Blocks Used Percents", strconv.FormatFloat(100*float64(stat.Blocks-stat.Bfree)/float64(stat.Blocks), 'f', 2, 64), "%"})
		result = append(result, item{"/INodesUsedPercents", "INodes Used Percents", strconv.FormatFloat(100*float64(stat.Files-stat.Ffree)/float64(stat.Files), 'f', 2, 64), "%"})
	}

	return nil
}

/*
	var stat unix.Statfs_t
	err = unix.Statfs(c.disk, &stat)
	free = uint64(stat.Bsize) * stat.Bfree
	total = uint64(stat.Bsize) * stat.Blocks

	if err != nil {
		c.SetString("Status", err.Error(), "error")

		c.SetString("SpaceTotal", "", "error")
		c.SetString("SpaceFree", "", "error")
		c.SetString("SpaceUsed", "", "error")

		c.SetString("BlocksTotal", "", "error")
		c.SetString("BlocksFree", "", "error")
		c.SetString("BlocksUsed", "", "error")

		c.SetString("INodesTotal", "", "error")
		c.SetString("INodesFree", "", "error")
		c.SetString("INodesUsed", "", "error")

		c.SetString("SpaceUsedPercents", "", "error")
		c.SetString("BlocksUsedPercents", "", "error")
		c.SetString("INodesUsedPercents", "", "error")

	} else {
		c.SetString("Status", "", "")

		c.SetUInt64("SpaceTotal", total/1024/1024, "MB")
		c.SetUInt64("SpaceFree", free/1024/1024, "MB")
		c.SetUInt64("SpaceUsed", (total-free)/1024/1024, "MB")

		c.SetUInt64("BlocksTotal", stat.Blocks, "")
		c.SetUInt64("BlocksFree", stat.Bfree, "")
		c.SetUInt64("BlocksUsed", stat.Blocks-stat.Bfree, "")

		c.SetUInt64("INodesTotal", stat.Files, "")
		c.SetUInt64("INodesFree", stat.Ffree, "")
		c.SetUInt64("INodesUsed", stat.Files-stat.Ffree, "")

		c.SetFloat64("SpaceUsedPercents", 100*float64(total-free)/float64(total), "%", 1)
		c.SetFloat64("BlocksUsedPercents", 100*float64(stat.Blocks-stat.Bfree)/float64(stat.Blocks), "%", 1)
		c.SetFloat64("INodesUsedPercents", 100*float64(stat.Files-stat.Ffree)/float64(stat.Files), "%", 1)
	}*/
