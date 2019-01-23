package cache

import "git.parallelcoin.io/pod/pkg/util"

// CacheableBlock is a wrapper around the util.Block type which provides a
// Size method used by the cache to target certain memory usage.
type CacheableBlock struct {
	*util.Block
}

// Size returns size of this block in bytes.
func (c *CacheableBlock) Size() (uint64, error) {
	return uint64(c.Block.MsgBlock().SerializeSize()), nil
}
