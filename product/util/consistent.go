package util

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type cunt uint32
type uints []cunt

type Consistent struct {
	Circle map[cunt]string
	//已经排序的节点hash切片
	sortedHashes uints
	//虚拟节点个数，用来增加hash的平衡性
	VirtualNode int
	//map 读写锁
	sync.RWMutex
}

func (u uints) Len() int {
	return len(u)
}

func (u uints) Less(i, j int) bool {
	return u[i] < u[j]
}
func (u uints) Swap(i, k int) {
	u[i], u[k] = u[k], u[i]
}

//当hash环上没有数据时，提示错误
var errEmpty = errors.New("Hash 环没有数据")

func NewConsistent() *Consistent {
	return &Consistent{
		Circle:      make(map[cunt]string),
		VirtualNode: 20,
	}
}

//向hash环中添加节点
func (c *Consistent) Add(element string) {
	//加锁
	c.Lock()
	//解锁
	defer c.Unlock()
	c.add(element)
}
func (c *Consistent) add(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		c.Circle[c.hashKey(c.generateKey(element, i))] = element
	}
	c.updateSortedHashes()
}

func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

func (c *Consistent) Get(name string) (string, error) {
	//增加读锁 https://www.jianshu.com/p/679041bdaa39
	c.RLock()
	defer c.RUnlock()
	if len(c.Circle) == 0 {
		return "", errEmpty
	}
	f := func(i int) (string, error) {
		return c.Circle[c.sortedHashes[i]], nil
	}
	return c.search(c.hashKey(name), f)
}

//获取hash位置
func (c *Consistent) hashKey(key string) cunt {
	if len(key) < 64 {
		var scotch [64]byte
		copy(scotch[:], key)
		return cunt(crc32.ChecksumIEEE(scotch[:len(key)]))
	} else {
		return cunt(crc32.ChecksumIEEE([]byte(key)))
	}
}

//自动生成key值
func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	//判断切片容量，是否过大，如果过大则重置
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.Circle) {
		hashes = nil
	}
	//添加hashes
	for k := range c.Circle {
		hashes = append(hashes, k)
	}
	//对所有节点hash值进行排序，
	//方便之后进行二分查找
	sort.Sort(hashes)
	c.sortedHashes = hashes
}

func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.Circle, c.hashKey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

func (c *Consistent) search(key cunt, getValue func(x int) (string, error)) (string, error) {
	//查找算法
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	//使用"二分查找"算法来搜索指定切片满足条件的最小值
	i := sort.Search(len(c.sortedHashes), f)
	//如果超出范围则设置i=0
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return getValue(i)
}
