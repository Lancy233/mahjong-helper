package util

import (
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"sort"
)

// 计算各张待牌的和率
// 剩余为 0 则和率为 0
func CalculateAgariRateOfEachTile(waits Waits, playerInfo *model.PlayerInfo) map[int]float64 {
	if playerInfo == nil {
		playerInfo = &model.PlayerInfo{}
	}

	tileAgariRate := map[int]float64{}

	// TODO 首先根据自家舍牌计算出是否振听
	// 振听的话和率简化成和枚数相关

	// 特殊处理字牌单骑的情况
	if len(waits) == 1 {
		for tile, left := range waits {
			if tile >= 27 {
				rate := honorTileDankiAgariTable[left]
				if InInts(tile, playerInfo.DoraTiles) {
					// 调整听宝牌时的和率
					// 忽略 dora 复合的影响
					rate *= honorDoraAgariMulti
				}
				tileAgariRate[tile] = rate
				return tileAgariRate
			}
		}
	}

	for tile, left := range waits {
		var rate float64
		if tile < 27 { // 数牌
			t := tile % 9
			if t > 4 {
				t = 8 - t
			}

			// TODO: 骗筋时的和率

			rate = nonSujiAgariTable[t][left]
		} else { // 字牌，非单骑
			rate = honorTileNonDankiAgariTable[left]
		}
		if InInts(tile, playerInfo.DoraTiles) {
			// 调整听宝牌时的和率
			// 忽略 dora 复合的影响
			if tile >= 27 {
				rate *= honorDoraAgariMulti
			} else {
				rate *= numberDoraAgariMulti
			}
		}
		tileAgariRate[tile] = rate
	}

	return tileAgariRate
}

// 计算平均和率
func CalculateAvgAgariRate(waits Waits, playerInfo *model.PlayerInfo) float64 {
	if playerInfo == nil {
		playerInfo = &model.PlayerInfo{}
	}

	tileAgariRate := CalculateAgariRateOfEachTile(waits, playerInfo)
	agariRate := 0.0
	for _, rate := range tileAgariRate {
		agariRate = agariRate + rate - agariRate*rate/100
	}

	// 调整两面和牌率
	// 需要 waits 恰好是筋牌关系，不能有非筋牌
	waitTiles := []int{}
	for tile, left := range waits {
		if left > 0 {
			if tile >= 27 {
				return agariRate
			}
			waitTiles = append(waitTiles, tile)
		}
	}
	if len(waitTiles) > 1 {
		suitType := waitTiles[0] / 9
		for _, tile := range waitTiles[1:] {
			if tile/9 != suitType {
				return agariRate
			}
		}
		sort.Ints(waitTiles)
		if len(waitTiles) == 2 && waitTiles[0]+3 == waitTiles[1] ||
			len(waitTiles) == 3 && waitTiles[0]+3 == waitTiles[1] && waitTiles[1]+3 == waitTiles[2] {
			agariRate *= ryanmenAgariMulti
		}
	}

	return agariRate
}
