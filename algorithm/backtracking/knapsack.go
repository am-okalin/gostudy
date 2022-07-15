package graph

type Pack struct {
	w      int   //背包最大承重
	n      int   //物品个数
	weight []int //每个物品的重量
	maxW   int   //结果集中的最大重量
}

func NewPack(w int, n int, weight []int) *Pack {
	return &Pack{w: w, n: n, weight: weight}
}

func (p *Pack) knapsack(i, cw int) {
	if cw == p.w || i == p.n {
		if cw > p.maxW {
			p.maxW = cw
		}
		return
	}
	p.knapsack(i+1, cw)
	if cw+p.weight[i] <= p.w {
		p.knapsack(i+1, cw+p.weight[i])
	}
}
