/*
键值对的存储的方式，只能用来索引静态路由
实现动态路由最常用的数据结构，被称为前缀树(Trie树)
HTTP请求的路径恰好是由/分隔的多段构成的，因此，每一段可以作为前缀树的一个节点。

有一点需要注意，/p/:lang/doc只有在第三层节点，即doc节点，pattern才会设置为/p/:lang/doc。
p和:lang节点的pattern属性皆为空。因此，当匹配结束时，我们可以使用n.pattern == ""来判断路由规则是否匹配成功。
例如，/p/python虽能成功匹配到:lang，但:lang的pattern值为空，因此匹配失败。
*/

package giny

import (
	"log"
	"strings"
)

type node struct {
	pattern  string  //表示某个完整的路由，若非完整路由，pattern为""
	part     string  //路由的一部分，比如 :lang
	children []*node //该节点的子节点
	isWild   bool    //是否为模糊匹配，含有 * 或 : 时为true
}

//第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	children := make([]*node, 0)
	for _, child := range n.children {
		if part == child.part || child.isWild {
			children = append(children, child)
		}
	}
	return children
}

//insert方法实现了根据pattern构建前缀树路由
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		//只有trie树的叶结点才有非空的pattern ✖ 因为只有"/"根节点的情况也有pattern
		//因此，对于精准匹配，只有成功匹配到叶结点，才表明匹配成功

		//检测是否有路由冲突
		if n.pattern != "" {
			// log.Printf("Route conflict! n.pattern = %s, pattern = %s\n", n.pattern, pattern)
			log.Panicf("Route conflict! 已存在pattern = %s, 当前pattern = %s\n", n.pattern, pattern)
		}

		n.pattern = pattern
		return
	}

	part := parts[height]

	child := n.matchChild(part)
	if child == nil {
		//未找到相应part的子节点，需要插入节点
		child = &node{ //！！！ 重大失误，:= 的 : 导致child成为了新的局部变量，并没有真正地对原来的child进行操作，从而导致了段错误！
			part:   part,
			isWild: part[0] == '*' || part[0] == ':',
		}

		n.children = append(n.children, child)
	}
	//递归插入下一个part节点
	child.insert(pattern, parts, height+1)
}

//search搜寻符合path的trie树叶结点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { //第二个条件不能用 n.part[0] == "*"，因为n.part可能为空，那么会发生越界
		//pattern = "" 匹配失败
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	//获取所有可能的子路径
	children := n.matchChildren(part)
	for _, child := range children {
		//递归地查找
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
